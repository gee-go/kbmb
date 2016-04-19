package testutil

import (
	"bytes"
	"math"
	"math/rand"
)

const (
	alphaBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// AlphaBytes returns n random bytes with A-Za-z chars.
// See http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golanghttp://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func AlphaBytes(src rand.Source, n int) []byte {
	b := make([]byte, n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(alphaBytes) {
			b[i] = alphaBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

var RandAlphaSelect = RandomByteSelect([]byte(alphaBytes))

func RandomByteSelect(choices []byte) func(src rand.Source, n int) []byte {
	charCount := len(choices)
	_, bitsReq := math.Frexp(float64(charCount))
	idxBits := uint(bitsReq)
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
		u.Write(AlphaBytes(rnd, rnd.Intn(6)+1))
	}

	countParams := rnd.Intn(6)
	if countParams > 0 {
		u.WriteByte('?')
	}

	for i := 0; i < countParams; i++ {
		// Write key
		u.Write(AlphaBytes(rnd, rnd.Intn(6)+1))
		u.WriteByte('=')
		// Write value
		u.Write(AlphaBytes(rnd, rnd.Intn(6)+1))
	}

	return u.String()
}
