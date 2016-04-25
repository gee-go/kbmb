package crawl

import (
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"
	"github.com/apex/log"
)

var skipExtensions = map[string]struct{}{
	".pdf": struct{}{},
	".zip": struct{}{},
	".js":  struct{}{}, // Could be useful to download
	".css": struct{}{},
	".gif": struct{}{},
	".jpg": struct{}{},
	".png": struct{}{},
}

type URLFn func(u *url.URL)

type Parser struct {
	Job *Crawl
	Doc *goquery.Document
}

func NewParser(job *Crawl, Doc *goquery.Document) *Parser {
	return &Parser{
		Job: job,
		Doc: Doc,
	}
}

func (p *Parser) ShouldSkip(u *url.URL) bool {
	// skip certain file endings like .zip, .pdf, .gif
	_, skip := skipExtensions[path.Ext(u.Path)]
	return skip
}

func (p *Parser) EachURL(fn URLFn) {
	sel := p.Doc.Find("[href]")

	for i := range sel.Nodes {
		l, ok := sel.Eq(i).Attr("href")
		if !ok {
			continue
		}

		u, err := url.Parse(l)
		if err != nil {
			log.WithFields(p.Job).WithError(err).Info("url")
			continue
		}

		if !p.ShouldSkip(u) {
			fn(u)
		}
	}
}
