package crawl

import (
	"testing"

	"github.com/gee-go/util/mrand"
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

func BenchmarkUniqueStringChan(b *testing.B) {
	q := NewUniqueStringChan()

	src := mrand.NewSource()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.In() <- string(mrand.AlphaBytes(src, 32))
	}

	for i := 0; i < b.N; i++ {
		<-q.Out()
	}
}
