package crawl

import "github.com/eapache/queue"

// UniqueStringChan converts a string channel to one that only outputs non-duplicate strings.
type UniqueStringChan struct {
	seen          map[string]struct{}
	length        chan int
	input, output chan string
	buffer        *queue.Queue
}

func NewUniqueStringChan() *UniqueStringChan {
	ch := &UniqueStringChan{
		seen:   make(map[string]struct{}),
		length: make(chan int),
		input:  make(chan string),
		output: make(chan string),
		buffer: queue.New(),
	}
	go ch.run()

	return ch
}

func (ch *UniqueStringChan) Close() {
	close(ch.input)
}

func (ch *UniqueStringChan) In() chan<- string {
	return ch.input
}

func (ch *UniqueStringChan) Out() <-chan string {
	return ch.output
}

func (ch *UniqueStringChan) Len() int {
	return <-ch.length
}

func (ch *UniqueStringChan) Count() int {
	return len(ch.seen)
}

func (ch *UniqueStringChan) run() {

	var input, output chan string
	var next string
	input = ch.input

	for input != nil || output != nil {
		select {
		case elem, open := <-input:
			if open {
				_, found := ch.seen[elem]

				if !found {
					ch.seen[elem] = struct{}{}
					ch.buffer.Add(elem)
				}

			} else {
				input = nil
			}
		case output <- next:
			ch.buffer.Remove()
		case ch.length <- ch.buffer.Length():
		}

		if ch.buffer.Length() > 0 {
			output = ch.output
			next = ch.buffer.Peek().(string)
		} else {
			output = nil
			next = ""
		}
	}

	close(ch.output)
	close(ch.length)
}
