package cfg

import (
	"github.com/apex/log"
	"github.com/gee-go/kbmb/cfg/nsqutil"
	"github.com/nsqio/go-nsq"
)

type NSQConfig struct {
	NSQDHosts      []string
	NSQLookupHosts []string
}

func (ncfg *NSQConfig) MustNewProducerPool() *nsqutil.ProducerPool {
	nsqProducer, err := ncfg.NewProducerPool()
	if err != nil {
		log.WithError(err).Fatal("NewProducerPool")
	}

	return nsqProducer
}

func (ncfg *NSQConfig) NewProducerPool() (*nsqutil.ProducerPool, error) {
	return nsqutil.NewProducerPool(ncfg.NSQDHosts)
}

func (ncfg *NSQConfig) MustNewConsumer(topic, channel string, handler nsq.Handler, concurrency int) *nsq.Consumer {
	consumer, err := ncfg.NewConsumer(topic, channel, handler, concurrency)
	if err != nil {
		log.WithError(err).Fatal("NewConsumer")
	}

	return consumer
}

func (ncfg *NSQConfig) NewConsumer(topic, channel string, handler nsq.Handler, concurrency int) (*nsq.Consumer, error) {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	consumer.AddConcurrentHandlers(handler, concurrency)
	if len(ncfg.NSQLookupHosts) > 0 {
		return consumer, consumer.ConnectToNSQLookupds(ncfg.NSQLookupHosts)
	}

	return consumer, consumer.ConnectToNSQDs(ncfg.NSQDHosts)
}
