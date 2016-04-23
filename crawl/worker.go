package crawl

import (
	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/apex/log"
	"github.com/nsqio/go-nsq"
)

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
	log.WithField("url", job.URL).Info("job")
	doc, err := NewDoc(job.URL, job.RootHost)
	if err != nil {
		return err
	}
	// TODO - config
	if doc.doc.Url.Host != job.RootHost {
		return nil
	}

	pr := doc.Result()
	if err := c.m.publishEmails(job, pr.Emails); err != nil {
		return err
	}

	return c.m.publishURLs(job, pr.Next)
}

func (c *Worker) Stop() {
	if c.workConsumer != nil {
		c.workConsumer.Stop()
	}
}
