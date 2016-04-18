package main

import (
	"flag"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gee-go/kbmb/crawl"
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
	conf := parseFlags()

	log.SetHandler(cli.New(os.Stderr))
	s, err := crawl.New(conf.Host)
	check(err)

	check(s.Run())

}
