package crawl

// UniqueStringChan converts a string channel to one that only outputs non-duplicate strings.
type UniqueStringChan struct {
	seen          map[string]struct{}
	input, output chan string
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

func (ch *UniqueStringChan) run() {
	defer close(ch.output)
	for iv := range ch.input {
		_, found := ch.seen[iv]
		if !found {
			ch.seen[iv] = struct{}{}
			ch.output <- iv
		}
	}
}
