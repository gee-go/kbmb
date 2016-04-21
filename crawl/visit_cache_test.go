package crawl

import (
	"testing"

	"github.com/gee-go/kbmb/testutil"
	"github.com/gee-go/util/mrand"
	"github.com/stretchr/testify/require"
)

// func assertDiffResp(urls ...string) {
// 	b := make([][]byte, len(urls))

// }

func TestRedisVisitCache(t *testing.T) {
	t.Parallel()
	key := string(testutil.RandAlphaSelect(mrand.NewSource(), 10))
	vc := NewRedisVisitCache(key)
	assert := require.New(t)
	assert.NoError(vc.Clear())

	_, err := vc.DiffAndSet([]string{"a", "b", "c"})
	assert.NoError(err)
	// assert.Equal([][]byte{[]byte{"a"}})
}
