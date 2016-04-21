package crawl

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

func MessageHandler(producer *ProducerPool, visitCache VisitCache) nsq.HandlerFunc {

	return nsq.HandlerFunc(func(m *nsq.Message) error {
		defer visitCache.Done()

		u := string(m.Body)
		fmt.Println(u)
		doc, err := NewDoc(u, "gee.io")
		if err != nil {
			return err
		}
		pr := doc.Result()
		next, err := visitCache.DiffAndSet(pr.Next)
		if err != nil {
			return err
		}

		if len(next) > 0 {
			return producer.MultiPublishAsync("urls", next)
		}

		return nil
	})
}
