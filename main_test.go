package main

// func TestNewSpider(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("skip integration test in short mode.")
// 	}
// 	t.Parallel()

// 	cases := []struct {
// 		root string
// 		host string
// 	}{
// 		{"mit.edu", "web.mit.edu"},
// 		{"jana.com", "jana.com"},
// 		{"gee.io", "gee.io"},
// 	}
// 	assert := require.New(t)

// 	for _, tc := range cases {
// 		s, err := NewSpider(tc.root)
// 		assert.NoError(err)
// 		s.Run()
// 		assert.Equal(tc.host, s.root.Host)
// 	}
// }
