package api

import (
	"fmt"
	"net/url"
	"runtime"

	"github.com/garyburd/redigo/redis"
	"github.com/gee-go/kbmb/crawl"
	"github.com/nsqio/go-nsq"
)

const (
	RedisVisitKey = "visited"
	RedisEmailKey = "email"
	NSQChannel    = "urls"
	NSQTopic      = "download"
)

type Crawler struct {
	rpool *redis.Pool

	nsqProducer *nsq.Producer
	NsqConsumer *nsq.Consumer
}

func NewCrawler() (*Crawler, error) {
	cfg := nsq.NewConfig()
	// cfg.MaxInFlight = 200
	nsqProducer, err := nsq.NewProducer("localhost:4151", cfg)
	if err != nil {
		return nil, err
	}

	nsqConsumer, err := nsq.NewConsumer(NSQChannel, NSQTopic, cfg)
	if err != nil {
		nsqProducer.Stop()
		return nil, err
	}

	c := &Crawler{
		rpool:       NewRedisPool(),
		nsqProducer: nsqProducer,
		NsqConsumer: nsqConsumer,
	}

	return c, nil
}

func (c *Crawler) Start(root string) error {
	u, err := url.Parse(root)
	if err != nil {
		return err
	}

	// Default to http
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// check for redirects
	// TODO - don't need to do 2 requests for first page.
	// resp, err := http.Get(u.String())
	// if err != nil {
	// 	return err
	// }

	size := runtime.NumCPU() - 1
	if size < 1 {
		size = 1
	}

	c.NsqConsumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
		fmt.Println(string(m.Body))
		doc, err := crawl.NewDoc(string(m.Body), "web.mit.edu")
		if err != nil {
			return err
		}
		pr := doc.Result()
		next, err := SetDiffExists(c.rpool, RedisVisitKey, pr.Next)
		if err != nil {
			return err
		}

		return c.nsqProducer.MultiPublishAsync(NSQChannel, next, nil)
	}), size)

	if err := c.NsqConsumer.ConnectToNSQD("localhost:4151"); err != nil {
		return err
	}

	return c.nsqProducer.PublishAsync(NSQChannel, []byte(u.String()), nil)
}
