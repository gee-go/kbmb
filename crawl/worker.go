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

	q *JobQueue
}

// Start processing jobs.
func (w *Worker) Start() {
	for {
		w.workerChan <- w.jobChan // worker is ready to work again.

		select {
		// got a new job.
		case job := <-w.jobChan:
			doc, err := NewDoc(job.URL, job.Root)
			if err != nil {
				fmt.Println(err)
				continue
			}
			pr := doc.Result()

			for _, email := range pr.Emails {
				fmt.Println(w.ID, email, job.URL)
			}

			for _, next := range pr.Next {
				w.q.Enqueue(next)
			}

			w.q.Complete()
		}
	}
}
