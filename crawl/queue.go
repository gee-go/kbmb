package crawl

import (
	"sync"

	"github.com/Workiva/go-datastructures/queue"
)

type JobQueue struct {
	seen  map[string]struct{}
	queue *queue.Queue

	mu sync.Mutex
	wg sync.WaitGroup
}

func NewJobQueue() *JobQueue {
	return &JobQueue{
		queue: queue.New(1000),
		seen:  make(map[string]struct{}),
	}
}

// Put url in Job queue.
// Ignore if already seen, otherwise append to queue.
func (ch *JobQueue) Put(j *Job) error {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	u := j.Key()
	_, found := ch.seen[u]
	if found {
		return nil
	}
	ch.seen[u] = struct{}{}

	ch.wg.Add(1)
	return ch.queue.Put(j)
}

func (ch *JobQueue) Complete() {
	ch.wg.Done()
}

func (ch *JobQueue) Wait() {
	ch.wg.Wait()
}

func (ch *JobQueue) Poll() (*Job, error) {
	out, err := ch.queue.Get(1)
	if err != nil {
		return nil, err
	}

	return out[0].(*Job), nil
}
