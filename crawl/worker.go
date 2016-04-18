package crawl

import "fmt"

type Worker struct {
	ID          int
	Work        chan *Job
	WorkerQueue chan chan *Job
	QuitChan    chan bool
	VisitChan   chan<- string
}

func (w *Worker) Start() {
	for {
		w.WorkerQueue <- w.Work

		select {
		case job := <-w.Work:
			doc, err := NewDoc(job.URL, job.Root)
			if err != nil {
				fmt.Println(err)
			}
			pr := doc.Result()

			for _, next := range pr.Next {
				w.VisitChan <- next
			}

		case <-w.QuitChan:
			// We have been asked to stop.
			fmt.Printf("worker%d stopping\n", w.ID)
			return
		}
	}
}

func (w *Worker) Stop() {
	w.QuitChan <- true
}

type WorkerPool struct {
	// number of workers.
	Size int

	WorkerChan chan chan *Job
}
