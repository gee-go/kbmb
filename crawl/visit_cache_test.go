package crawl

// func assertDiffResp(t *testing.T, actual [][]byte, urls ...string) {
// 	b := make([][]byte, len(urls))
// 	for i, u := range urls {
// 		b[i] = []byte(u)
// 	}

// 	require.Equal(t, b, actual)
// }

// func TestRedisVisitCache(t *testing.T) {
// 	t.Parallel()
// 	key := string(testutil.RandAlphaSelect(mrand.NewSource(), 10))
// 	vc := NewRedisVisitCache(key)
// 	assert := require.New(t)
// 	assert.NoError(vc.Clear())

// 	// add 3 new urls
// 	next, err := vc.DiffAndSet([]string{"a", "b", "c"})
// 	assert.NoError(err)
// 	assertDiffResp(t, next, "a", "b", "c")
// 	count, err := vc.Count()
// 	assert.NoError(err)
// 	assert.Equal(count, 3)

// 	// add the same urls, none should be added
// 	next, err = vc.DiffAndSet([]string{"a", "b", "c"})
// 	assert.NoError(err)
// 	assertDiffResp(t, next)
// 	count, err = vc.Count()
// 	assert.NoError(err)
// 	assert.Equal(count, 3)

// 	// add a different url
// 	next, err = vc.DiffAndSet([]string{"a", "b", "c", "d"})
// 	assert.NoError(err)
// 	assertDiffResp(t, next, "d")
// 	count, err = vc.Count()
// 	assert.NoError(err)
// 	assert.Equal(count, 4)
// }
