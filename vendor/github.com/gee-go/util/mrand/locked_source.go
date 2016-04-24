package mrand

import (
	"math/rand"
	"sync"
)

type lockedSource struct {
	mu  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.mu.Lock()
	n = r.src.Int63()
	r.mu.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.mu.Lock()
	r.src.Seed(seed)
	r.mu.Unlock()
}

// NewLockedSource returns a thread safe rand.Source seeded with the current time.
func NewLockedSource() rand.Source {
	return &lockedSource{src: NewSource()}
}
