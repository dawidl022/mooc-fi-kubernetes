package processor

import (
	"bytes"
	"io"
	"net/http"
	netUrl "net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ExampleUrl = "https://example.com/"

func TestScrape_ReturnsError_WhenUrlIsMalformed(t *testing.T) {
	scraper := NewScraper()
	_, err := scraper.Scrape("not a url")
	assert.Error(t, err)
}

func TestScraper_CallsGivenGetter(t *testing.T) {
	getter := &getterSpy{}
	scraper := scraper{get: getter.get}
	scraper.Scrape(ExampleUrl)
	assert.True(t, getter.wasCalled)
}

type getterSpy struct {
	wasCalled bool
}

func (g *getterSpy) get(url string) (resp *http.Response, err error) {
	g.wasCalled = true
	return
}

func TestScrape_ReturnsError_WhenCannotConnect(t *testing.T) {
	getter := &getterNoConnectionStub{}
	scraper := scraper{get: getter.get}
	_, err := scraper.Scrape(ExampleUrl)
	assert.Error(t, err)
}

type getterNoConnectionStub struct{}

func (g *getterNoConnectionStub) get(url string) (*http.Response, error) {
	return nil, &netUrl.Error{}
}

func TestScrape_ReturnsError_WhenNotHttpOK(t *testing.T) {
	statusCodes := []int{201, 204, 400, 401, 403, 404, 500}
	for _, status := range statusCodes {
		getter := &getterHttpStatusStub{status: status}
		scraper := scraper{get: getter.get}
		_, err := scraper.Scrape(ExampleUrl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, ExampleUrl)
		assert.ErrorContains(t, err, strconv.Itoa(status))
		_ = err.(responseNotOKException)
	}
}

type getterHttpStatusStub struct {
	status int
}

func (g *getterHttpStatusStub) get(url string) (*http.Response, error) {
	return &http.Response{StatusCode: g.status}, nil
}

func TestScrape_ReturnsBody_WhenHttpOK(t *testing.T) {
	messages := []string{"Hello, world!", "I need scissors! 61!"}
	for _, message := range messages {
		getter := &getterStub{status: http.StatusOK, body: []byte(message)}
		scraper := scraper{get: getter.get}
		body, err := scraper.Scrape(ExampleUrl)
		assert.NoError(t, err)
		assert.Equal(t, message, string(body))
	}
}

type getterStub struct {
	status int
	body   []byte
}

func (g *getterStub) get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(g.body)),
	}, nil
}
