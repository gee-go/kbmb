package crawl

import (
	"fmt"
	"sync/atomic"

	"github.com/nsqio/go-nsq"
)

type ProducerPool struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	rrCount uint64

	producers []*nsq.Producer
}

func (pp *ProducerPool) MultiPublishAsync(topic string, msgs [][]byte) error {
	// round robin selection
	rrCount := atomic.AddUint64(&pp.rrCount, 1)
	pidx := rrCount % uint64(len(pp.producers))

	producer := pp.producers[pidx]
	return producer.MultiPublishAsync(topic, msgs, nil, nil)
}

func NewProducerPool(hosts ...string) (*ProducerPool, error) {
	pp := &ProducerPool{
		producers: make([]*nsq.Producer, len(hosts)),
	}
	var err error
	nsqConfig := nsq.NewConfig()

	// Create a nsq producer for each host
	for i, h := range hosts {
		pp.producers[i], err = nsq.NewProducer(h, nsqConfig)
		if err != nil {
			return nil, err
		}
	}
	return pp, err
}

func MessageHandler(producer *ProducerPool, visitCache VisitCache) nsq.HandlerFunc {

	return nsq.HandlerFunc(func(m *nsq.Message) error {
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
