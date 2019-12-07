package search

type Searcher interface {
	Images([]string) ([]*Result, error)
	Name() string
}

type Result struct {
	URL           string
	Width, Height int
}
