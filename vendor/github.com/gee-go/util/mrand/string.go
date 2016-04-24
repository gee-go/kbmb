package mrand

import "math/rand"

const (
	alphaBytes    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
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
