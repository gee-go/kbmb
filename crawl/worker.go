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
	quitChan   chan struct{}

	emailResultChan chan<- string
	nextVisitChan   chan<- string
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
				w.emailResultChan <- email
			}

			for _, next := range pr.Next {
				w.nextVisitChan <- next
			}

		case <-w.quitChan:
			// We have been asked to stop.
			fmt.Printf("worker%d stopping\n", w.ID)
			return
		}
	}
}

func (w *Worker) Stop() {
	w.quitChan <- struct{}{}
}

type WorkerPool struct {
	numWorkers int

	workerChan chan chan *Job
	quitChan   chan struct{}
	resultChan chan *PageResult

	emailResultChan *UniqueStringChan
	nextVisitChan   *UniqueStringChan
}

// EmailChan allows iteration over all emails seen.
func (wp *WorkerPool) EmailChan() <-chan string {
	return wp.emailResultChan.Out()
}

func (wp *WorkerPool) NextChan() <-chan string {
	return wp.nextVisitChan.Out()
}

func (wp *WorkerPool) Start(root *url.URL) {
	wp.nextVisitChan.In() <- root.String()

	for {
		select {
		case email := <-wp.EmailChan():
			fmt.Println(email)
		case work := <-wp.NextChan():

			go func() {
				worker := <-wp.workerChan
				worker <- &Job{work, root}
			}()
		}
	}

	// for wp.nextVisitChan.Len() > 0 {

	// }

	// time.Sleep(30 * time.Second)
}

func NewWorkerPool(size int) *WorkerPool {
	wp := &WorkerPool{
		numWorkers: size,
		workerChan: make(chan chan *Job, size),
		resultChan: make(chan *PageResult),

		emailResultChan: NewUniqueStringChan(),
		nextVisitChan:   NewUniqueStringChan(),
	}

	for i := 0; i < size; i++ {
		w := &Worker{
			ID:              i,
			jobChan:         make(chan *Job),
			workerChan:      wp.workerChan,
			quitChan:        make(chan struct{}),
			emailResultChan: wp.emailResultChan.In(),
			nextVisitChan:   wp.nextVisitChan.In(),
		}
		go w.Start()
	}

	return wp
}
