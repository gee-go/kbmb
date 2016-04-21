package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gee-go/kbmb/cfg"
	"github.com/gee-go/kbmb/crawl"
	"github.com/nsqio/go-nsq"
)

type QueueHandler struct {
	// 64bit atomic vars need to be first for proper alignment on 32bit platforms
	counter uint64 // for round robin producer selection.

	nsqConfig *nsq.Config
	producers []*nsq.Producer
}

func NewQueueHandler() *QueueHandler {
	return &QueueHandler{
		nsqConfig: nsq.NewConfig(),
	}
}

// Add a producer to the set. Not thread safe.
func (qh *QueueHandler) AddProducer(h string) error {
	producer, err := nsq.NewProducer(h, qh.nsqConfig)
	if err != nil {
		return err
	}
	qh.producers = append(qh.producers, producer)
	return nil
}

func main() {
	// Setup
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)

	visitCache := crawl.NewRedisVisitCache("visited")
	if err := visitCache.Clear(); err != nil {
		log.WithError(err).Fatal("Clear visitCache")
	}

	nsqConfig := &cfg.NSQConfig{
		NSQDHosts: []string{"localhost:4150"},
	}

	nsqProducer := nsqConfig.MustNewProducerPool()
	nsqConsumer := nsqConfig.MustNewConsumer("urls", "download", crawl.MessageHandler(nsqProducer, visitCache), 8)

	defer nsqConsumer.Stop()
	visitCache.DiffAndSet([]string{"http://gee.io"})
	nsqProducer.PublishAsync("urls", []byte("http://gee.io"))

	visitCache.Wait()
}
