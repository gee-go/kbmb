package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gee-go/kbmb/crawl"
)

func check(err error) {
	if err != nil {
		log.WithError(err).Fatal("check")
	}
}

func main() {
	log.SetHandler(cli.New(os.Stderr))
	s, err := crawl.New("gee.io")
	check(err)

	check(s.Run())

}
