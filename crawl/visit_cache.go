package crawl

import "sync"

type VisitCache struct {
	mu sync.RWMutex

	data map[string]struct{} // use empty struct value because it uses no space.
	q    []string
}

func NewVisitCache() *VisitCache {
	return &VisitCache{
		data: make(map[string]struct{}),
	}
}

// List returns a list of all urls in the cache.
func (v *VisitCache) List() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	var out []string
	for u := range v.data {
		out = append(out, u)
	}

	return out
}

// Add adds a url to the cache
func (v *VisitCache) Add(u string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.data[u] = struct{}{}
}

// Has returns true if url in cache
func (v *VisitCache) Has(u string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	_, ok := v.data[u]
	return ok
}

func (v *VisitCache) Enqueue(u string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, found := v.data[u]; found {
		return
	}

	v.data[u] = struct{}{}
	v.q = append(v.q, u)
}

func (v *VisitCache) Pop() string {
	v.mu.Lock()
	defer v.mu.Unlock()
	var x string
	x, v.q = v.q[0], v.q[1:]
	return x
}

func (v *VisitCache) Len() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return len(v.q)
}

// FilterDupes removes all url's in the slice that already are in visit cache.
// Also removes dupes in given list.
func (v *VisitCache) FilterDupes(urls []string) []string {
	// Dedupe input
	urlSet := make(map[string]struct{})
	for _, u := range urls {
		urlSet[u] = struct{}{}
	}

	v.mu.RLock()
	defer v.mu.RUnlock()

	var out []string
	for u := range urlSet {
		if _, has := v.data[u]; !has {
			out = append(out, u)
		}
	}

	return out
}
