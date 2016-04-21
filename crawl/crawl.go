package crawl

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/satori/go.uuid"
)

// Crawl represents a single crawl job.
type Crawl struct {
	ID       string
	URL      string // starting url
	RootHost string // the host used
}

func (c *Crawl) EmailTopic() string {
	return fmt.Sprintf("email_%s", c.ID)
}

func (c *Crawl) VisitKey() string {
	return fmt.Sprintf("visit:%s", c.ID)
}

func (c *Crawl) WaitGroupKey() string {
	return fmt.Sprintf("wait:%s", c.ID)
}

// Child returns a copy with the given url
func (c *Crawl) Child(url string) *Crawl {
	return &Crawl{
		ID:       c.ID,
		URL:      url,
		RootHost: c.RootHost,
	}
}

func NewCrawl(start string) (*Crawl, error) {
	u, err := url.Parse(start)
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
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &Crawl{
		ID:       uuid.NewV4().String(),
		URL:      resp.Request.URL.String(),
		RootHost: resp.Request.URL.Host,
	}, nil
}
