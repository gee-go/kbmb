package crawl

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
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

	size := runtime.NumCPU() - 1
	if size < 1 {
		size = 1
	}

	fmt.Println(size)
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

			Host: c.root.Host,
		}
		go w.Start()
	}

	// send jobs
	for {
		u, err := c.q.Poll()

		if err != nil {
			fmt.Println(err)
			continue
		}
		worker := <-c.workers
		worker <- u
	}
}

func (c *Crawler) Run() <-chan string {
	go c.startWorkers()
	c.q.Put(&Job{URL: c.root})

	go func() {
		c.q.wg.Wait()
		c.emailChan.Close()
	}()

	return c.emailChan.Out()
}
