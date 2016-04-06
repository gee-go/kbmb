package crawl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVisitCache(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	vc := NewVisitCache()

	vc.Add("a")
	exp := []string{"a"}
	assert.Equal(exp, vc.List())

	vc.Add("a")
	assert.Equal(exp, vc.List(), "check add duplicate")

	assert.Equal([]string{"b", "c", "d"}, vc.FilterDupes([]string{"a", "b", "b", "c", "d"}))
}
