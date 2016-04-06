package crawl

import (
	"fmt"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Doc struct {
	doc *goquery.Document
}

type URLFn func(u *url.URL)

type PageResult struct {
	Emails []string
	Next   []string
}

func NewDoc(u string) (*Doc, error) {
	doc, err := goquery.NewDocument(u)
	fmt.Println(doc.Url.Host)
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
			log.Fatal(err)
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
