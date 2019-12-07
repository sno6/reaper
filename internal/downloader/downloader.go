package downloader

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/sno6/reaper/internal/search"
)

const (
	nDownloaders = 100 // Image downloaders.
	nWriters     = 100 // Download consumers / writers.
)

type downloadedImage struct {
	name string
	r    io.Reader
}

type Downloader struct {
	cfg *Config
}

type Config struct {
	OutDir  string
	Verbose bool
	Timeout time.Duration
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func New(cfg *Config) *Downloader {
	return &Downloader{cfg: cfg}
}

func (d *Downloader) Download(images []*search.Result) error {
	wrkChan := make(chan string)
	errChan := make(chan error)
	imgChan := make(chan *downloadedImage)

	if err := createOutDir(d.cfg.OutDir); err != nil {
		return errors.Wrap(err, "error creating output directory")
	}

	wg := newProgWaitGroup(len(images))

	for nw := 0; nw < nWriters; nw++ {
		go func() {
			for {
				errChan <- save(<-imgChan, d.cfg.OutDir)
				wg.Done()
			}
		}()
	}

	for nd := 0; nd < nDownloaders; nd++ {
		go func() {
			for {
				img, err := d.download(<-wrkChan)
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
					if d.cfg.Verbose {
						color.Red("\n:: Error: %v", err)
					}
				}
			}
		}
	}()

	wg.Wait()
	color.Green("\n:: Successfully downloaded %d images, %d downloads failed (use -v flag for errors)",
		len(images)-len(errors), len(errors))

	return nil
}

func (d *Downloader) download(u string) (*downloadedImage, error) {
	c := &http.Client{Timeout: d.cfg.Timeout}
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	body := resp.Body.(io.Reader)
	body, ct, err := isImage(body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file content type")
	}
	return &downloadedImage{
		name: cleanFilename(u, ct),
		r:    body,
	}, nil
}

func isImage(r io.Reader) (io.Reader, string, error) {
	br := bufio.NewReader(r)
	buf, err := br.Peek(512)
	if err != nil {
		return nil, "", errors.Wrap(err, "error seeking 512 bytes to read image details")
	}
	contentType := http.DetectContentType(buf)
	if !strings.HasPrefix(contentType, "image/") {
		return nil, "", fmt.Errorf("unexpected content type %s", contentType)
	}
	return br, contentType, nil

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

func cleanFilename(u string, contentType string) string {
	entropy := genEntropy()
	base := filepath.Base(u)

	ext := filepath.Ext(base)
	if ext == "" {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && (exts != nil && len(exts) > 1) {
			return entropy + "-" + base + exts[0]
		}

		// No extension and we can't find an extension by the content type..
		return entropy + "-" + base
	}

	ind := strings.Index(base, ext)
	return entropy + "-" + base[0:ind+(len(ext))]
}

func genEntropy() string {
	const (
		low  = 1000
		high = 99999
	)
	return strconv.Itoa(low + rand.Intn(high-low))
}
