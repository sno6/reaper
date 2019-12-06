package download

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sno6/reaper/internal/search"
)

const (
	nDownloaders = 100 // Image downloaders.
	nWriters     = 100 // Download consumers / writers.

	defaultTimeout = time.Second * 15
)

type downloadedImage struct {
	name string
	r    io.Reader
}

type Downloader struct{}

func New() *Downloader {
	return &Downloader{}
}

func (d *Downloader) Download(images []*search.Result, outDir string) error {
	wrkChan := make(chan string)
	errChan := make(chan error)
	imgChan := make(chan *downloadedImage)

	if err := createOutDir(outDir); err != nil {
		return errors.Wrap(err, "error creating output directory")
	}

	wg := newProgWaitGroup(len(images))

	for nw := 0; nw < nWriters; nw++ {
		go func() {
			for {
				errChan <- save(<-imgChan, outDir)
				wg.Done()
			}
		}()
	}

	for nd := 0; nd < nDownloaders; nd++ {
		go func() {
			for {
				img, err := download(<-wrkChan)
				if err != nil {
					wg.Done()
					errChan <- err
					continue
				}
				imgChan <- img
			}
		}()
	}

	go func() {
		for _, work := range images {
			wrkChan <- work.URL
		}
	}()

	var errors []error
	go func() {
		for {
			select {
			case err := <-errChan:
				if err != nil {
					errors = append(errors, err)
				}
			}
		}
	}()

	wg.Wait()
	color.Green("\n:: Successfully downloaded %d images, %d downloads failed (use -v flag for errors)", len(images)-len(errors), len(errors))

	return nil
}

func download(u string) (*downloadedImage, error) {
	c := &http.Client{Timeout: defaultTimeout}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	body := resp.Body.(io.Reader)
	if body, err = isImage(body); err != nil {
		return nil, errors.Wrap(err, "error downloading image")
	}
	return &downloadedImage{
		name: filepath.Base(u),
		r:    body,
	}, nil
}

func isImage(r io.Reader) (io.Reader, error) {
	br := bufio.NewReader(r)
	buf, err := br.Peek(512)
	if err != nil {
		return nil, errors.Wrap(err, "error seeking 512 bytes to read image details")
	}
	contentType := http.DetectContentType(buf)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, fmt.Errorf("unexpected content type %s", contentType)
	}
	return br, nil

}

func save(img *downloadedImage, outDir string) error {
	f, err := os.Create(filepath.Join(outDir, img.name))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, img.r)
	if err != nil {
		// An error occured while copying the image. There could have been something
		// wrong with the file or an HTTP timeout occured while reading from the underlying buffer.
		// In either case, remove the partially written file.
		f.Close()
		return os.Remove(f.Name())
	}

	return nil
}

func createOutDir(dir string) error {
	info, err := os.Stat(dir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.Mkdir(dir, 0777)
}
