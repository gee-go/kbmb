package nsqutil

import (
	"sync/atomic"

	"github.com/nsqio/go-nsq"
)

// ProducerPool provides a unified interface to a set of nsqd queues
type ProducerPool struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	rrCount uint64

	producers []*nsq.Producer
}

func (pp *ProducerPool) getProducer() *nsq.Producer {
	// round robin selection
	rrCount := atomic.AddUint64(&pp.rrCount, 1)
	pidx := rrCount % uint64(len(pp.producers))

	return pp.producers[pidx]
}

func (pp *ProducerPool) PublishAsync(topic string, msgs []byte) error {
	return pp.getProducer().PublishAsync(topic, msgs, nil, nil)
}

func (pp *ProducerPool) MultiPublishAsync(topic string, msgs [][]byte) error {
	return pp.getProducer().MultiPublishAsync(topic, msgs, nil, nil)
}

func NewProducerPool(hosts []string) (*ProducerPool, error) {
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
