package crawl

import (
	"fmt"
	"net/url"

	"github.com/eapache/channels"
)

type Job struct {
	URL  string
	Root string
}

type Crawler struct {
	root *url.URL
	vc   *VisitCache
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
	jobChan := channels.NewInfiniteChannel()

	jobChan.In() <- &Job{
		URL:  c.root.String(),
		Root: "web.mit.edu",
	}

	for c.vc.Len() > 0 {
		c.HandleURL(c.vc.Pop())
	}
	// time.Sleep(3 * time.Second)
	return nil
}
