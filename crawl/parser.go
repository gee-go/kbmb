package crawl

import (
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
)

type Doc struct {
	doc *goquery.Document
}

type URLFn func(u *url.URL)

type PageResult struct {
	// Emails contains all emails found on page
	Emails []string

	// Next contains links with the correct host.
	// Not guarenteed to be unique or unvisited
	Next []string
}

func NewDoc(u string) (*Doc, error) {
	doc, err := goquery.NewDocument(u)

	if err != nil {
		return nil, err
	}
	return &Doc{doc}, nil
}

func (d *Doc) EachURL(fn URLFn) {
	sel := d.doc.Find("a[href]")

	for i := range sel.Nodes {
		l, ok := sel.Eq(i).Attr("href")
		if !ok {
			continue
		}

		u, err := url.Parse(l)
		if err != nil {
			log.WithError(err).Info("url")
		}

		fn(d.doc.Url.ResolveReference(u))
	}

}

func (d *Doc) Result() *PageResult {
	out := &PageResult{}
	// TODO - Dedupe on a page basis.

	d.EachURL(func(u *url.URL) {
		// handle mailto links
		if u.Scheme == "mailto" {
			out.Emails = append(out.Emails, u.Opaque)
			return
		}

		if u.Host == d.doc.Url.Host {
			out.Next = append(out.Next, u.String())
		}
	})

	return out
}
