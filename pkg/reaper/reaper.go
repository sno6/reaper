package reaper

import (
	"time"

	"github.com/fatih/color"
	"github.com/sno6/reaper/internal/downloader"
	"github.com/sno6/reaper/internal/search/aggregator"
)

type Config struct {
	SearchTerms         []string
	OutDir              string
	Verbose             bool
	Timeout             time.Duration
	MinWidth, MinHeight int
}

type Reaper struct {
	cfg *Config
}

func New(cfg *Config) *Reaper {
	return &Reaper{cfg: cfg}
}

func (r *Reaper) Run() error {
	color.Green(":: Searching search engines for the following terms: %v", r.cfg.SearchTerms)
	s := aggregator.New(&aggregator.Config{
		SearchTerms: r.cfg.SearchTerms,
		MinWidth:    r.cfg.MinWidth,
		MinHeight:   r.cfg.MinHeight,
	})
	res, err := s.GetAllImages()
	if err != nil {
		return err
	}

	color.Green(":: Found %d images, downloading...", len(res))
	dl := downloader.New(&downloader.Config{
		OutDir:  r.cfg.OutDir,
		Verbose: r.cfg.Verbose,
	})
	return dl.Download(res)
}
