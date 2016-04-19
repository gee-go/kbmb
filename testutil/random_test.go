package testutil

import (
	"testing"

	"github.com/gee-go/util/mrand"
)

const randSize = 10

func TestRandom(t *testing.T) {

}

func BenchmarkAlphaString(b *testing.B) {
	// a := require.New(t)
	src := mrand.NewSource()
	var s []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = RandAlphaSelect(src, randSize)
	}

	if len(s) != randSize {
		b.Fail()
	}
}

func BenchmarkRandomByteSelect(b *testing.B) {
	// a := require.New(t)
	src := mrand.NewSource()
	lsource := RandomByteSelect([]byte(alphaBytes))
	var s []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s = lsource(src, randSize)
	}

	if len(s) != randSize {
		b.Fail()
	}
}
