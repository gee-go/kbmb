package crawl

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/k0kubun/pp"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type Job struct {
	URL  string
	Root *url.URL
}

type Worker struct {
	ID int

	jobChan    chan *Job
	workerChan chan chan *Job

	q         *JobQueue
	emailChan *UniqueStringChan
}

func (w *Worker) Process(ctx context.Context, j *Job) error {
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, j.URL)
	if err != nil {
		return err
	}
	pp.Println(resp.Header.Get("Content-Type"))
	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	doc := &Doc{d, j.Root}

	if doc.doc.Url.Host != j.Root.Host {
		return nil
	}

	doc.EachURL(func(u *url.URL) {
		// handle mailto links
		if u.Scheme == "mailto" {
			w.emailChan.In() <- u.Opaque
		} else if u.Host == j.Root.Host {
			next := u.String()

			if len(next) > 0 && next[len(next)-1] == '/' {
				next = next[:len(next)-1]
			}

			w.q.Put(next)
		}
	})

	return nil
}

// Start processing jobs.
func (w *Worker) Start() {
	for {
		w.workerChan <- w.jobChan // worker is ready to work again.

		select {
		// got a new job.
		case job := <-w.jobChan:
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			if err := w.Process(ctx, job); err != nil {
				fmt.Println(err)
			}
			w.q.Complete()
		}
	}
}
