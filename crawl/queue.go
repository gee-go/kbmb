package crawl

import "sync"

type JobQueue struct {
	seen          map[string]struct{}
	input, output chan string
	buffer        []string

	wg sync.WaitGroup
}

func NewJobQueue() *JobQueue {
	ch := &JobQueue{
		seen:   make(map[string]struct{}),
		input:  make(chan string),
		output: make(chan string),
	}
	go ch.run()

	return ch
}

func (ch *JobQueue) Close() {
	close(ch.input)
}

func (ch *JobQueue) Out() <-chan string {
	return ch.output
}

func (ch *JobQueue) Enqueue(u string) {
	ch.wg.Add(1)
	ch.input <- u
}

func (ch *JobQueue) Count() int {
	return len(ch.seen)
}

func (ch *JobQueue) Complete() {
	ch.wg.Done()
}

func (ch *JobQueue) send(v string) {
	_, found := ch.seen[v]
	if !found {
		ch.seen[v] = struct{}{}
		ch.buffer = append(ch.buffer, v)
	} else {
		ch.wg.Done()
	}
}

func (ch *JobQueue) flush() {
	for _, v := range ch.buffer {
		ch.output <- v
	}
	close(ch.output)
}

func (ch *JobQueue) Wait() {
	ch.wg.Wait()
}

func (ch *JobQueue) run() {
	defer ch.flush()

	for {
		// always append to buffer if empty
		if len(ch.buffer) == 0 {
			v, open := <-ch.input
			if !open {
				return
			}
			ch.send(v)
		} else {
			select {
			case v, open := <-ch.input:
				if !open {
					return
				}
				ch.send(v)
			case ch.output <- ch.buffer[0]:
				ch.buffer = ch.buffer[1:]
			}
		}
	}
}
