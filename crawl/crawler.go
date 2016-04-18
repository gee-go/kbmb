package crawl

import (
	"net/http"
	"net/url"
)

type Crawler struct {
	root *url.URL
}

func New(root string) (*Crawler, error) {
	u, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	// Default to http
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// check for redirects
	// TODO - don't need to do 2 requests for first page.
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	return &Crawler{
		root: resp.Request.URL,
	}, nil
}

func (c *Crawler) Run() error {
	wp := NewWorkerPool(4)
	wp.Start(c.root)

	return nil
}
