package crawl

import (
	"fmt"
	"net/url"
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

func (w *Worker) Process(j *Job) error {
	doc, err := NewDoc(j.URL, j.Root)
	if err != nil {
		return err
	}

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
			if err := w.Process(job); err != nil {
				fmt.Println(err)
			}
			w.q.Complete()
		}
	}
}
