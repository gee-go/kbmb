package crawl

// UniqueStringChan converts a string channel to one that only outputs non-duplicate strings.
type UniqueStringChan struct {
	seen          map[string]struct{}
	input, output chan string
	buffer        []string
}

func NewUniqueStringChan() *UniqueStringChan {
	ch := &UniqueStringChan{
		seen:   make(map[string]struct{}),
		input:  make(chan string),
		output: make(chan string),
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

func (ch *UniqueStringChan) Count() int {
	return len(ch.seen)
}

func (ch *UniqueStringChan) send(v string) {
	_, found := ch.seen[v]
	if !found {
		ch.seen[v] = struct{}{}
		ch.buffer = append(ch.buffer, v)
	}
}

func (ch *UniqueStringChan) flush() {
	for _, v := range ch.buffer {
		ch.output <- v
	}
	close(ch.output)
}

func (ch *UniqueStringChan) run() {
	defer ch.flush()

	for {
		// always append to buffer if empty
		if len(ch.buffer) == 0 {
			v, open := <-ch.input
			if !open {
				return
			}
			ch.send(v)
		} else {
			select {
			case v, open := <-ch.input:
				if !open {
					return
				}
				ch.send(v)
			case ch.output <- ch.buffer[0]:
				ch.buffer = ch.buffer[1:]
			}
		}
	}
}
