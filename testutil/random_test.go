package testutil

import (
	"testing"

	"github.com/gee-go/util/mrand"
	"github.com/stretchr/testify/require"
)

const randSize = 10

func TestRandom(t *testing.T) {
	assert := require.New(t)

	assert.Equal(nextPow2Exp(52), uint(6))
	assert.Equal(nextPow2Exp(4), uint(2))
	assert.Equal(nextPow2Exp(2), uint(1))
	assert.Equal(nextPow2Exp(1), uint(1))
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
