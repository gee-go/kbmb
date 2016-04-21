package crawl

import (
	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/k0kubun/pp"
	"github.com/nsqio/go-nsq"
)

type Msg struct {
	RootHost string
	URL      string
}

// A worker downloads and parses urls from jobs it receives.
type Worker struct {
	m            *Manager
	workConsumer *nsq.Consumer
}

func (c *Worker) HandleMessage(m *nsq.Message) error {
	// unmarshall message into job
	job := &Crawl{}
	if err := msgpack.Unmarshal(m.Body, job); err != nil {
		return err
	}

	defer c.m.markDone(job)
	pp.Println(job)
	doc, err := NewDoc(job.URL, job.RootHost)
	if err != nil {
		return err
	}
	pr := doc.Result()
	return c.m.publishURLs(job, pr.Next)
}

func (c *Worker) Stop() {
	if c.workConsumer != nil {
		c.workConsumer.Stop()
	}
}
