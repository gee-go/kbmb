package crawl

import (
	"fmt"
	"net/http"
	"net/url"
)

type Job struct {
	URL  string
	Root string
}

type Crawler struct {
	root *url.URL

	resultChan chan *PageResult
	visitChan  *UniqueStringChan
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
		root:      resp.Request.URL,
		visitChan: NewUniqueStringChan(),
	}, nil
}

func (c *Crawler) HandleURL(u string) error {
	doc, err := NewDoc(u, c.root)
	if err != nil {
		return err
	}
	pr := doc.Result()

	for _, next := range pr.Next {
		c.visitChan.In() <- next
	}

	return nil
}

func (c *Crawler) Run() error {
	go func() {
		for u := range c.visitChan.Out() {
			fmt.Println(u)
			c.HandleURL(u)
		}
	}()

	c.HandleURL(c.root.String())

	for c.visitChan.Len() > 0 {

	}

	return nil
}
