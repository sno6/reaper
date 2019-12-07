package duck

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/sno6/reaper/internal/search"
)

type response struct {
	Next    string `json:"next"`
	Query   string `json:"query"`
	Results []struct {
		Title     string `json:"title"`
		Height    int    `json:"height"`
		Thumbnail string `json:"thumbnail"`
		URL       string `json:"url"`
		Source    string `json:"source"`
		Image     string `json:"image"`
		Width     int    `json:"width"`
	} `json:"results"`
	QueryEncoded string `json:"queryEncoded"`
	ResponseType string `json:"response_type"`
}

const (
	baseURL   = "https://duckduckgo.com"
	searchURL = baseURL + "/" + "i.js"
	maxIter   = 20
	step      = 50

	defTimeout = time.Second * 10
)

func GetImages(searchTerm string) ([]*search.Result, error) {
	fmt.Println("getImages with term:", searchTerm)

	var results []*search.Result
	searchTerm = url.QueryEscape(searchTerm)

	token, err := getToken(searchTerm)
	if err != nil {
		return nil, err
	}

	for i := 0; i < maxIter; i++ {
		images, next, err := getImages(searchTerm, token, (i*step)+step)
		if err != nil {
			return nil, err
		}
		results = append(results, images...)
		if next == "" {
			break
		}
	}
	return results, nil
}

func getImages(searchTerm string, token string, step int) ([]*search.Result, string, error) {
	u, err := url.Parse(searchURL)
	if err != nil {
		return nil, "", err
	}

	q := u.Query()
	q.Add("o", "json")
	q.Add("u", "yahoo")
	q.Add("l", "wt-wt")
	q.Add("p", "1")
	q.Add("s", strconv.Itoa(step))
	q.Add("q", searchTerm)
	q.Add("vqd", token)
	u.RawQuery = q.Encode()

	rsp, err := http.Get(u.String())
	if err != nil {
		return nil, "", err
	}
	defer rsp.Body.Close()

	var res response
	if err := json.NewDecoder(rsp.Body).Decode(&res); err != nil {
		return nil, "", err
	}
	return convertResponse(res), res.Next, nil
}

func getToken(searchTerm string) (string, error) {
	r, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return "", nil
	}

	q := r.URL.Query()
	q.Add("q", searchTerm)
	r.URL.RawQuery = q.Encode()

	c := &http.Client{Timeout: defTimeout}
	rsp, err := c.Do(r)
	if err != nil {
		return "", nil
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", nil
	}

	param := "vqd="
	rg := regexp.MustCompile(fmt.Sprintf("%s([\\d-]+)", param))
	token := rg.FindString(string(body))
	if token == "" || len(token) <= len(param) {
		return "", errors.New("unable to find DDG request token")
	}

	return token[len(param):], nil
}

func convertResponse(r response) []*search.Result {
	var results []*search.Result
	for _, res := range r.Results {
		results = append(results, &search.Result{
			URL:    res.Image,
			Width:  res.Width,
			Height: res.Height,
		})
	}
	return results
}
