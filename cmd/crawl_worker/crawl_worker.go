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
	crawler := crawl.NewCrawler(nsqConfig)
	defer crawler.Stop()
	crawler.Start(8)

	crawler.SendURLs([]string{"http://gee.io"})

	crawler.Wait()
}
