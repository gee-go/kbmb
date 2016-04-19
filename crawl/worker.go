package crawl

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

type Job struct {
	URL *url.URL
}

func (j *Job) Key() string {
	u := j.URL.String()

	if len(u) > 0 && u[len(u)-1] == '/' {
		u = u[:len(u)-1]
	}

	return u
}

type Worker struct {
	ID int

	jobChan    chan *Job
	workerChan chan chan *Job

	q         *JobQueue
	emailChan *UniqueStringChan

	Host string
}

func (w *Worker) Process(ctx context.Context, j *Job) error {
	resp, err := ctxhttp.Get(ctx, http.DefaultClient, j.URL.String())
	if err != nil {
		return err
	}

	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	doc := &Doc{d, w.Host}

	if doc.doc.Url.Host != w.Host {
		return nil
	}

	doc.EachURL(func(u *url.URL) {
		// handle mailto links
		if u.Scheme == "mailto" {
			w.emailChan.In() <- u.Opaque
		} else if u.Host == w.Host {
			w.q.Put(&Job{URL: u})
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
