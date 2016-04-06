package crawl

import (
	"fmt"
	"net/url"
	"sync"
)

type Crawler struct {
	root *url.URL
	vc   *VisitCache

	q  chan string
	wg sync.WaitGroup
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
		vc:   NewVisitCache(),
	}, nil
}

func (c *Crawler) HandleURL(u string) error {
	doc, err := NewDoc(u)
	if err != nil {
		return err
	}

	pr := doc.Result()

	for _, next := range pr.Next {
		c.vc.Enqueue(next)
	}
	fmt.Println(u)

	return nil
}

func (c *Crawler) Run() error {
	c.vc.Enqueue(c.root.String())
	for c.vc.Len() > 0 {
		c.HandleURL(c.vc.Pop())
	}
	// time.Sleep(3 * time.Second)
	return nil
}
