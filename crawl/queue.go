package crawl

import "sync"

type JobQueue struct {
	seen     map[string]struct{}
	input    chan string
	complete chan bool
	output   chan string
	wg       sync.WaitGroup
}

func (q *JobQueue) Add(u string) {
	q.wg.Add(1)
	q.input <- u
}

func (q *JobQueue) run() {
	// Dedupe input
	dedupedIn := make(chan string)
	go func() {
		defer close(dedupedIn)

		for v := range q.input {
			_, found := q.seen[v]
			if found {
				q.wg.Done()
			} else {
				q.seen[v] = struct{}{}
				dedupedIn <- v
			}
		}
	}()

	// buffered in
	bufferedIn := make(chan string)
	go SliceIQ(dedupedIn, bufferedIn)
}

func NewJobQueue() *JobQueue {
	q := &JobQueue{
		seen:     make(map[string]struct{}),
		input:    make(chan string),
		output:   make(chan string),
		complete: make(chan bool),
	}

	go func() {
		for range q.complete {
			q.wg.Done()
		}
	}()

	go q.run()

	return q
}
