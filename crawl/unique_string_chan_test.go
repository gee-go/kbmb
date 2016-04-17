package crawl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniqueStringChan(t *testing.T) {
	t.Parallel()
	a := require.New(t)
	ch := NewUniqueStringChan()

	go func() {
		for _, v := range []string{"a", "b", "a", "a", "b", "d"} {
			ch.In() <- v
		}
		ch.Close()
	}()

	var out []string
	for v := range ch.Out() {
		out = append(out, v)
	}

	a.Equal([]string{"a", "b", "d"}, out)
}
