package crawl

import (
	"fmt"
	"net/http"
	"net/url"
)

type Crawler struct {
	root       *url.URL
	numWorkers int
	workers    chan chan *Job
	q          *JobQueue
	emailChan  *UniqueStringChan
}

func New(root string) (*Crawler, error) {
	u, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	// Default to http
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// check for redirects
	// TODO - don't need to do 2 requests for first page.
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	size := 8
	return &Crawler{
		root:       resp.Request.URL,
		q:          NewJobQueue(),
		numWorkers: size,
		workers:    make(chan chan *Job, size),
		emailChan:  NewUniqueStringChan(),
	}, nil
}

func (c *Crawler) startWorkers() {
	for i := 0; i < c.numWorkers; i++ {
		w := &Worker{
			ID:         i,
			jobChan:    make(chan *Job),
			workerChan: c.workers,
			q:          c.q,
			emailChan:  c.emailChan,
		}
		go w.Start()
	}

	// send jobs
	for u := range c.q.Out() {
		worker := <-c.workers

		worker <- &Job{u, c.root}
	}
}

func (c *Crawler) Run() error {
	go c.startWorkers()
	c.q.Enqueue(c.root.String())
	go func() {
		for email := range c.emailChan.Out() {
			fmt.Println(email)
		}
	}()
	c.q.wg.Wait()
	return nil
}
