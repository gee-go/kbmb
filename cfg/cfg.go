package cfg

import (
	"flag"
	"strings"

	"github.com/k0kubun/pp"
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
	flag.StringVar(&c.Redis.URL, "redis", "redis://127.0.0.1:6379", "redis url")
	flag.StringVar(&nsqdHosts, "nsqd", "localhost:4150", "nsqd HTTP address (may be given multiple times)")

	flag.Parse()
	c.NSQDHosts = strings.Fields(nsqdHosts)
	pool := c.NewRedisPool()
	pp.Println(pool.Get().Do("PING"))
	return c
}
