package cfg

import (
	"flag"
	"strings"
)

type Cfg struct {
	Redis
	NSQConfig
}

var (
	nsqdHosts string
)

func FromFlags() *Cfg {
	c := &Cfg{}
	flag.StringVar(&c.Redis.URL, "redis-url", "redis://127.0.0.1:6379", "redis url")
	flag.StringVar(&nsqdHosts, "nsqd", "localhost:4150", "nsqd HTTP address (may be given multiple times)")

	flag.Parse()
	c.NSQDHosts = strings.Fields(nsqdHosts)
	return c
}
