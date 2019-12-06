package search

type Result struct {
	URL           string
	Width, Height int
}

type Search interface {
	Images([]string) ([]*Result, error)
	Name() string
}
