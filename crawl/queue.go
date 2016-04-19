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
func (ch *JobQueue) Put(u string) error {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	//
	_, found := ch.seen[u]
	if found {
		return nil
	}
	ch.seen[u] = struct{}{}

	ch.wg.Add(1)
	return ch.queue.Put(u)
}

func (ch *JobQueue) Complete() {
	ch.wg.Done()
}

func (ch *JobQueue) Wait() {
	ch.wg.Wait()
}

func (ch *JobQueue) Poll() (string, error) {
	out, err := ch.queue.Get(1)
	if err != nil {
		return "", err
	}

	return out[0].(string), nil
}
