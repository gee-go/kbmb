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

	nsqConfig := &cfg.NSQConfig{
		NSQDHosts: []string{"localhost:4150"},
	}

	manager := crawl.NewManager(nsqConfig)

	crawl, err := crawl.NewCrawl("gee.io")
	if err != nil {
		panic(err)
	}

	if err := manager.Start(crawl); err != nil {
		panic(err)
	}

	worker := manager.NewWorker(8)
	defer worker.Stop()

	if err := manager.Wait(crawl); err != nil {
		panic(err)
	}
}
