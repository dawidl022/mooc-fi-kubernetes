package processor

import (
	"fmt"
	"io"
	"net/http"
)

type scraper struct {
	get func(url string) (resp *http.Response, err error)
}

func NewScraper() *scraper {
	return &scraper{get: http.Get}
}

func (s *scraper) Scrape(url string) ([]byte, error) {
	r, err := s.get(url)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, noResponseError{}
	}
	if r.StatusCode != http.StatusOK {
		return nil, responseNotOKException{url: url, status: r.StatusCode}
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

type noResponseError struct{}

func (e noResponseError) Error() string {
	return "getter did not return response"
}

type responseNotOKException struct {
	url    string
	status int
}

func (e responseNotOKException) Error() string {
	return fmt.Sprintf("server returned %d status code for url: %s", e.status, e.url)
}
