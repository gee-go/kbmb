package mrand

import (
	"math/rand"
	"time"
)

// SeedDefault seeds the default rand source.
func SeedDefault() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// NewSource returns a source that is already seeded.
func NewSource() rand.Source {
	return rand.NewSource(time.Now().UnixNano())
}

// New returns a preseeded rand.Rand
func New() *rand.Rand {
	return rand.New(NewSource())
}
