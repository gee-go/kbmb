package crawl

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gee-go/kbmb/cfg"
	"github.com/gee-go/kbmb/cfg/nsqutil"
	"github.com/nsqio/go-nsq"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const (
	RedisCrawlKey = "crawls"
	NSQTopic      = "urls"
)

type Manager struct {
	config      *cfg.NSQConfig
	nsqProducer *nsqutil.ProducerPool
	redisPool   *redis.Pool
}

func NewManager(config *cfg.NSQConfig) *Manager {
	return &Manager{
		config:      config,
		redisPool:   cfg.NewRedisPool(),
		nsqProducer: config.MustNewProducerPool(),
	}
}

func (m *Manager) EmailConsumer(c *Crawl, fn func(m string)) *nsq.Consumer {
	return m.config.MustNewConsumer(c.EmailTopic(), "emails", nsq.HandlerFunc(func(m *nsq.Message) error {
		fn(string(m.Body))
		return nil
	}), 1)
}

func (m *Manager) NewWorker(concurrency int) *Worker {
	worker := &Worker{
		m: m,
	}
	worker.workConsumer = m.config.MustNewConsumer("urls", "download", worker, concurrency)
	return worker
}

func (m *Manager) Wait(c *Crawl) error {
	// TODO - use keyspace notifications.
	for range time.Tick(100 * time.Millisecond) {
		rconn := m.redisPool.Get()
		count, err := redis.Int(rconn.Do("GET", c.WaitGroupKey()))
		rconn.Close()
		if err != nil || count == 0 {
			return err
		}
	}

	return nil
}

func (m *Manager) markDone(c *Crawl) error {
	rconn := m.redisPool.Get()
	defer rconn.Close()
	_, err := rconn.Do("DECR", c.WaitGroupKey())
	return err
}

func (m *Manager) publishEmails(parent *Crawl, emails []string) error {
	// TODO - don't need wait group for emails
	unpublished, err := redisDiffAndSet(m.redisPool, parent.EmailTopic(), "counter", emails)
	if err != nil {
		return err
	}

	if len(unpublished) == 0 {
		return nil
	}

	// make messages to send.
	out := make([][]byte, len(unpublished))
	for i, u := range unpublished {
		out[i] = []byte(u)
	}

	return m.nsqProducer.MultiPublishAsync(parent.EmailTopic(), out)
}

func (m *Manager) publishURLs(parent *Crawl, urls []string) error {
	// mark urls as visited, and return previously unvisited ones
	unvisited, err := redisDiffAndSet(m.redisPool, parent.VisitKey(), parent.WaitGroupKey(), urls)
	if err != nil {
		return err
	}

	if len(unvisited) == 0 {
		return nil
	}

	// make messages to send.
	out := make([][]byte, len(unvisited))
	for i, u := range unvisited {
		c := parent.Child(u)

		b, err := msgpack.Marshal(c)
		if err != nil {
			return err
		}

		out[i] = b
	}

	return m.nsqProducer.MultiPublishAsync(NSQTopic, out)
}

// Start a new crawl.
func (m *Manager) Start(c *Crawl) error {
	rc := m.redisPool.Get()

	b, err := msgpack.Marshal(c)
	if err != nil {
		return err
	}

	// Keep track of current jobs.
	if _, err := rc.Do("HSET", RedisCrawlKey, c.ID, b); err != nil {
		return err
	}

	// mark as visited
	if _, err := redisDiffAndSet(m.redisPool, c.VisitKey(), c.WaitGroupKey(), []string{c.URL}); err != nil {
		return err
	}

	// kick off a job
	return m.nsqProducer.PublishAsync(NSQTopic, b)
}
