package crawl

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Job struct {
	URL  string
	Root *url.URL
}

type Crawler struct {
	root *url.URL

	numWorkers int
	visitChan  *UniqueStringChan
	workerChan chan chan *Job
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

	return &Crawler{
		root:       resp.Request.URL,
		numWorkers: 10,
		visitChan:  NewUniqueStringChan(),
		workerChan: make(chan chan *Job, 10),
	}, nil
}

func (c *Crawler) newWorker(i int) *Worker {
	w := &Worker{
		ID:          i,
		Work:        make(chan *Job),
		WorkerQueue: c.workerChan,
		QuitChan:    make(chan bool),
		VisitChan:   c.visitChan.In(),
	}
	w.Start()
	return w
}

func (c *Crawler) HandleURL(u string) error {
	doc, err := NewDoc(u, c.root)
	if err != nil {
		return err
	}
	pr := doc.Result()

	for _, next := range pr.Next {
		c.visitChan.In() <- next
	}

	return nil
}

func (c *Crawler) Run() error {
	for i := 0; i < c.numWorkers; i++ {
		c.newWorker(i)
	}

	go func() {
		for {
			select {
			case work := <-c.visitChan.Out():
				fmt.Println(work)
				go func() {
					worker := <-c.workerChan
					worker <- &Job{work, c.root}
				}()
			}
		}

	}()

	c.HandleURL(c.root.String())

	for c.visitChan.Len() > 0 {

	}

	time.Sleep(30 * time.Second)

	return nil
}
