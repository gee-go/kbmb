package main

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/garyburd/redigo/redis"
	"github.com/gee-go/kbmb/crawl"
	"github.com/nsqio/go-nsq"
)

var setExists = redis.NewScript(-1, `
local r = {}

for _, m in pairs(ARGV) do
  if redis.call('SISMEMBER', KEYS[1], m) == 0 then
    r[#r+1] = m 
  end
end

for _, m in pairs(ARGV) do
  redis.call('SADD', KEYS[1], m)
end

return r
 `)

func NewRedisPool() *redis.Pool {

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://127.0.0.1:6379")
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func main() {

	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
	cfg := nsq.NewConfig()
	rdis := NewRedisPool()
	c := rdis.Get()
	c.Do("DEL", "visited")
	c.Close()
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
		doc, err := crawl.NewDoc(u, "gee.io")
		if err != nil {
			return err
		}
		pr := doc.Result()
		rconn := rdis.Get()
		reply, err := redis.ByteSlices(setExists.Do(rconn, redis.Args{}.Add(1).Add("visited").AddFlat(pr.Next)...))
		rconn.Close()
		if err != nil {
			return err
		}
		if len(reply) > 0 {
			return workProducer.MultiPublishAsync("urls", reply, nil)
		}

		return nil
	}), 10)

	if err := workConsumer.ConnectToNSQD("localhost:4150"); err != nil {
		log.WithError(err).Fatal("connect")
	}

	workProducer.PublishAsync("urls", []byte("http://gee.io"), nil)

	for range time.Tick(1 * time.Second) {

	}
}
