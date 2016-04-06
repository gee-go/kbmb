package crawl

import (
	"net/url"

	"github.com/apex/log"
	"github.com/k0kubun/pp"
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

	return &Crawler{
		root: u,
	}, nil
}

func (c *Crawler) HandleURL(u string) error {
	log.WithField("url", u).Info("start")
	doc, err := NewDoc(u)
	if err != nil {
		return err
	}

	pr := doc.Result()
	pp.Print(pr)
	return nil
}

func (c *Crawler) Run() error {
	c.HandleURL(c.root.String())

	return nil
}
