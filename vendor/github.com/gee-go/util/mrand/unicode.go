package mrand

import (
	"math/rand"
	"unicode"
)

func countR16(rt *unicode.RangeTable) int {
	count := 0
	for _, r := range rt.R16 {
		count += int((r.Hi-r.Lo)/r.Stride) + 1
	}
	return count
}

func countR32(rt *unicode.RangeTable) int {
	count := 0
	for _, r := range rt.R32 {
		count += int((r.Hi-r.Lo)/r.Stride) + 1
	}
	return count
}

func selectR16(rt *unicode.RangeTable, i int) rune {
	count := 0
	for _, r := range rt.R16 {

		ri := i - count
		count += int((r.Hi-r.Lo)/r.Stride) + 1

		if i < count {
			return rune(r.Lo + uint16(ri)*r.Stride)
		}
	}

	return -1
}

func selectR32(rt *unicode.RangeTable, i int) rune {
	count := 0
	for _, r := range rt.R32 {
		ri := i - count
		count += int((r.Hi-r.Lo)/r.Stride) + 1

		if i < count {
			return rune(r.Lo + uint32(ri)*r.Stride)
		}
	}

	return -1
}

// Unicode returns a random set of runes from a given range table.
func Unicode(rnd *rand.Rand, rt *unicode.RangeTable, n int) []rune {
	c16 := countR16(rt)
	c32 := countR32(rt)

	out := make([]rune, n)
	for i := 0; i < n; i++ {
		ri := rnd.Intn(c16 + c32)
		if ri < c16 {
			out[i] = selectR16(rt, ri)
		} else {
			out[i] = selectR32(rt, ri-c16)
		}

	}

	return out
}
