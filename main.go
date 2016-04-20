package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/gee-go/kbmb/crawl"
	"github.com/nsqio/go-nsq"
)

type Config struct {
	Host string
}

func parseFlags() *Config {
	c := &Config{}

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("Need a file")
	}

	c.Host = flag.Arg(0)
	return c
}

func check(err error) {
	if err != nil {
		log.WithError(err).Fatal("check")
	}
}

func main() {
	cfg := nsq.NewConfig()

	workProducer, err := nsq.NewProducer("localhost:4150", cfg)
	if err != nil {
		log.WithError(err).Fatal("work producer")
	}

	workConsumer, err := nsq.NewConsumer("urls", "download", cfg)
	if err != nil {
		log.WithError(err).Fatal("work consumer")
	}
	workConsumer.AddConcurrentHandlers(nsq.HandlerFunc(func(m *nsq.Message) error {
		u := string(m.Body)
		fmt.Println(u)
		doc, err := crawl.NewDoc(u, "web.mit.edu")
		if err != nil {
			return err
		}
		pr := doc.Result()

		next := [][]byte{}
		for _, n := range pr.Next {
			next = append(next, []byte(n))
		}

		workProducer.MultiPublishAsync("urls", next, nil)

		return nil
	}), 6)

	if err := workConsumer.ConnectToNSQLookupd("localhost:4161"); err != nil {
		log.WithError(err).Fatal("connect")
	}

	workProducer.PublishAsync("urls", []byte("http://mit.edu"), nil)

	for range time.Tick(1 * time.Second) {

	}
}
