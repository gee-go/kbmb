package crawl

import (
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
	"github.com/apex/log"
	"github.com/nsqio/go-nsq"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// A worker downloads and parses urls from jobs it receives.
type Worker struct {
	m            *Manager
	workConsumer *nsq.Consumer
}

func (c *Worker) HandleMessage(m *nsq.Message) error {
	// unmarshall message into job
	job := &Crawl{}
	if err := msgpack.Unmarshal(m.Body, job); err != nil {
		return err
	}
	defer c.m.markDone(job)
	lg := log.WithFields(job)
	lg.Info("job")

	// Timeout of 3 seconds
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	// 1. Request the url
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, job.URL)
	if err != nil {
		lg.WithError(err).Error("http get")
		return nil
	}

	// 2. Check for a redirect from a matching host to a non-matching host.
	if resp.Request.URL.Host != job.RootHost {
		resp.Body.Close() // goquery auto closes.
		return nil
	}

	// 3. Parse response.
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		lg.WithError(err).Error("parse resp")
		return nil
	}

	parser := NewParser(job, doc)
	// 4. Iterate through links
	var emails []string
	var links []string

	parser.EachURL(func(u *url.URL) {
		if u.Scheme == "mailto" {
			emails = append(emails, u.Opaque)
		} else {
			// make sure its absolute.
			u = resp.Request.URL.ResolveReference(u)

			if u.Host == job.RootHost {
				// normalize
				links = append(links, purell.NormalizeURL(u, purell.FlagsUsuallySafeGreedy))
			}
		}
	})

	// 5. pulish results
	if err := c.m.publishEmails(job, emails); err != nil {
		return err
	}
	return c.m.publishURLs(job, links)
}

func (c *Worker) Stop() {
	if c.workConsumer != nil {
		c.workConsumer.Stop()
	}
}
