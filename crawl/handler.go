package crawl

import (
	"fmt"

	"github.com/gee-go/kbmb/cfg/nsqutil"
	"github.com/nsqio/go-nsq"
)

func MessageHandler(producer *nsqutil.ProducerPool, visitCache VisitCache) nsq.HandlerFunc {

	return nsq.HandlerFunc(func(m *nsq.Message) error {
		defer visitCache.Done()

		u := string(m.Body)
		fmt.Println(u)
		doc, err := NewDoc(u, "gee.io")
		if err != nil {
			return err
		}
		pr := doc.Result()
		next, err := visitCache.DiffAndSet(pr.Next)
		if err != nil {
			return err
		}

		if len(next) > 0 {
			return producer.MultiPublishAsync("urls", next)
		}

		return nil
	})
}

// type Crawler struct {
// 	producerPool *ProducerPool
// 	visitCache   VisitCache
// 	consumer     *nsq.Consumer

// 	nsqds []string
// }

// func NewCrawler(nsqds []string) *Crawler {
// 	return &Crawler{
// 		nsqds:        nsqds,
// 		producerPool: NewProducerPool(nsqds...),
// 	}
// }

// func (c *Crawler) Start(root string, concurrency int) error {
// 	var err error
// 	c.consumer, err = nsq.NewConsumer("urls", "download", nsq.NewConfig())
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
