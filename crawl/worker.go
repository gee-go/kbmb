package crawl

import (
	"fmt"
	"net/url"
)

type Job struct {
	URL  string
	Root *url.URL
}

type CrawlState struct {
	Root *url.URL

	WorkerChan chan chan *Job
}

type Worker struct {
	ID int

	jobChan    chan *Job
	workerChan chan chan *Job
	quitChan   chan struct{}

	emailResultChan chan<- string
	nextVisitChan   chan<- string
	q               *JobQueue
}

// Start processing jobs.
func (w *Worker) Start() {
	for {
		w.workerChan <- w.jobChan // worker is ready to work again.

		select {
		// got a new job.
		case job := <-w.jobChan:
			fmt.Printf("%d job\n", w.ID)
			doc, err := NewDoc(job.URL, job.Root)
			if err != nil {
				fmt.Println(err)
				continue
			}
			pr := doc.Result()

			for _, email := range pr.Emails {
				fmt.Println(email)
			}

			for _, next := range pr.Next {
				w.q.Add(next)
			}

			w.q.complete <- true

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
	workers    []*Worker
	workerChan chan chan *Job
	quitChan   chan struct{}
	resultChan chan *PageResult

	emailResultChan *UniqueStringChan
	nextVisitChan   *UniqueStringChan

	q *JobQueue
}

// EmailChan allows iteration over all emails seen.
func (wp *WorkerPool) EmailChan() <-chan string {
	return wp.emailResultChan.Out()
}

func (wp *WorkerPool) NextChan() <-chan string {
	return wp.nextVisitChan.Out()
}

func (wp *WorkerPool) Start(root *url.URL) {
	wp.q.Add(root.String())
	go func() {
		for {
			select {
			case email := <-wp.EmailChan():
				fmt.Println(email)
			case work := <-wp.q.output:

				go func() {
					worker := <-wp.workerChan
					worker <- &Job{work, root}

				}()
			}
		}
	}()
	wp.q.wg.Wait()

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
		q:               NewJobQueue(),
	}

	for i := 0; i < size; i++ {
		w := &Worker{
			ID:              i,
			jobChan:         make(chan *Job),
			workerChan:      wp.workerChan,
			quitChan:        make(chan struct{}),
			emailResultChan: wp.emailResultChan.In(),
			nextVisitChan:   wp.nextVisitChan.In(),
			q:               wp.q,
		}
		wp.workers = append(wp.workers, w)
		go w.Start()
	}

	return wp
}
