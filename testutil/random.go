package testutil

import (
	"bytes"
	"math/rand"
)

const (
	alphaBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var RandAlphaSelect = RandomByteSelect([]byte(alphaBytes))

func nextPow2Exp(n int64) uint {
	// TODO - optimize
	exp := uint(0)
	isPow2 := n != 1 && (n&(n-1) == 0)

	for n != 0 {
		n >>= 1
		exp++
	}

	if isPow2 {
		exp--
	}

	return exp
}

func RandomByteSelect(choices []byte) func(src rand.Source, n int) []byte {
	charCount := len(choices)

	idxBits := nextPow2Exp(int64(charCount))
	idxMask := int64(1)<<idxBits - 1
	idxMax := int64(63 / idxBits)

	return func(src rand.Source, n int) []byte {
		b := make([]byte, n)

		// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
		for i, cache, remain := n-1, src.Int63(), idxMask; i >= 0; {
			if remain == 0 {
				cache, remain = src.Int63(), idxMax
			}
			if idx := int(cache & idxMask); idx < charCount {
				b[i] = choices[idx]
				i--
			}
			cache >>= idxBits
			remain--
		}

		return b
	}
}

// Return a random relative url path
func RandomRelativeURL(rnd *rand.Rand) string {
	var u bytes.Buffer

	// 0 to 6 path components
	for i := 0; i < rnd.Intn(6); i++ {
		u.WriteByte('/')
		u.Write(RandAlphaSelect(rnd, rnd.Intn(6)+1))
	}

	countParams := rnd.Intn(6)
	if countParams > 0 {
		u.WriteByte('?')
	}

	for i := 0; i < countParams; i++ {
		// Write key
		u.Write(RandAlphaSelect(rnd, rnd.Intn(6)+1))
		u.WriteByte('=')
		// Write value
		u.Write(RandAlphaSelect(rnd, rnd.Intn(6)+1))
	}

	return u.String()
}
