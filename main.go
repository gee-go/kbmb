package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gee-go/kbmb/cfg"
	"github.com/gee-go/kbmb/crawl"
)

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
