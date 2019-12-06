package reaper

import (
	"github.com/fatih/color"
	"github.com/sno6/reaper/internal/download"
	"github.com/sno6/reaper/internal/search/duck"
)

type Config struct {
	SearchTerms []string
	OutDir      string
}

type Reaper struct {
	cfg *Config
}

func New(cfg *Config) *Reaper {
	return &Reaper{cfg: cfg}
}

func (r *Reaper) Run() error {
	color.Green(":: Searching search engines for the following terms: %v", r.cfg.SearchTerms)
	ddg := duck.New(nil)
	res, err := ddg.Images(r.cfg.SearchTerms)
	if err != nil {
		return err
	}

	color.Green(":: Downloading %d images sourced from %s", len(res), ddg.Name())
	dl := download.New()
	if err := dl.Download(res, r.cfg.OutDir); err != nil {
		return err
	}

	return err
}
