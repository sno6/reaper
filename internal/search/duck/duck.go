package duck

import (
	"github.com/pkg/errors"
	"github.com/sno6/reaper/internal/search"
)

type Duck struct{}

func (d *Duck) Images(terms []string) ([]*search.Result, error) {
	var results []*search.Result
	for _, term := range terms {
		res, err := GetImages(term)
		if err != nil {
			return nil, errors.Wrap(err, "error getting images from DDG")
		}
		results = append(results, res...)
	}
	return results, nil
}

func (d *Duck) Name() string {
	return "Duck Duck Go"
}
