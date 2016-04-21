package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
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

	cfg := nsq.NewConfig()
	workProducer, err := crawl.NewProducerPool("localhost:4150")
	if err != nil {
		log.WithError(err).Fatal("work producer")
	}

	workConsumer, err := nsq.NewConsumer("urls", "download", cfg)
	if err != nil {
		log.WithError(err).Fatal("work consumer")
	}

	visitCache := crawl.NewRedisVisitCache("visited")
	if err := visitCache.Clear(); err != nil {
		log.WithError(err).Fatal("Clear visitCache")
	}

	workConsumer.AddConcurrentHandlers(crawl.MessageHandler(workProducer, visitCache), 10)
	if err := workConsumer.ConnectToNSQD("localhost:4150"); err != nil {
		log.WithError(err).Fatal("connect")
	}
	visitCache.DiffAndSet([]string{"http://gee.io"})
	workProducer.MultiPublishAsync("urls", [][]byte{[]byte("http://gee.io")})

	visitCache.Wait()
}
