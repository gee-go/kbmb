package crawl

import (
	"fmt"

	"github.com/apex/log"
	"github.com/gee-go/kbmb/cfg"
	"github.com/gee-go/kbmb/cfg/nsqutil"
	"github.com/nsqio/go-nsq"
)

type Crawler struct {
	config       *cfg.NSQConfig
	producerPool *nsqutil.ProducerPool
	workConsumer *nsq.Consumer
	visitCache   VisitCache
}

func (c *Crawler) HandleMessage(m *nsq.Message) error {
	defer c.visitCache.Done()

	u := string(m.Body)
	fmt.Println(u)
	doc, err := NewDoc(u, "gee.io")
	if err != nil {
		return err
	}
	pr := doc.Result()
	return c.SendURLs(pr.Next)
}

func NewCrawler(config *cfg.NSQConfig) *Crawler {
	return &Crawler{
		config:       config,
		visitCache:   NewRedisVisitCache("visited"),
		producerPool: config.MustNewProducerPool(),
	}
}

func (c *Crawler) SendURLs(urls []string) error {
	unvisited, err := c.visitCache.DiffAndSet(urls)
	if err != nil {
		return err
	}

	if len(unvisited) > 0 {
		return c.producerPool.MultiPublishAsync("urls", unvisited)
	}
	return nil
}

func (c *Crawler) Wait() error {
	return c.visitCache.Wait()
}

func (c *Crawler) Stop() {
	if err := c.visitCache.Clear(); err != nil {
		log.WithError(err).Fatal("Visit Cache Clear")
	}
	if c.workConsumer != nil {
		c.workConsumer.Stop()
	}
}

func (c *Crawler) Start(concurrency int) {
	c.workConsumer = c.config.MustNewConsumer("urls", "download", c, concurrency)
}
