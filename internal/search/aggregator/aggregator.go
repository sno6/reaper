package aggregator

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sno6/reaper/internal/search"
	"github.com/sno6/reaper/internal/search/duck"
)

type Aggregator struct {
	cfg       *Config
	searchers []search.Searcher
}

type Config struct {
	SearchTerms         []string
	MinWidth, MinHeight int
}

func New(cfg *Config) *Aggregator {
	return &Aggregator{
		cfg:       cfg,
		searchers: []search.Searcher{&duck.Duck{}},
	}
}

func (a *Aggregator) GetAllImages() ([]*search.Result, error) {
	var results []*search.Result

	for _, searcher := range a.searchers {
		res, err := searcher.Images(a.cfg.SearchTerms)
		if err != nil {
			return nil, errors.Wrap(err, "search aggregator error")
		}
		results = append(results, res...)
	}

	results = filterDuplicates(results)
	results = filterDimensions(results, a.cfg.MinWidth, a.cfg.MinHeight)

	return results, nil
}

func filterDuplicates(res []*search.Result) []*search.Result {
	var noDups []*search.Result
	m := make(map[string]struct{})

	for _, r := range res {
		if _, has := m[r.URL]; has {
			continue
		}
		m[r.URL] = struct{}{}
		noDups = append(noDups, r)
	}

	return noDups
}

func filterDimensions(res []*search.Result, mw, mh int) []*search.Result {
	var newRes []*search.Result
	for _, v := range res {
		if v.Width >= mw && v.Height >= mh {
			newRes = append(newRes, v)
			fmt.Printf("Passing image with dims: %d %d, confg: %dx%d\n", v.Width, v.Height, mw, mh)
		}
	}
	return newRes
}
