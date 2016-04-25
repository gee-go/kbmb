package crawl

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

var parserExcludeTestCase = []struct {
	url  string
	skip bool
}{
	{"gee.io", false},
	{"gee.io/pdf", false},
	{"gee.io/a.pdf", true},
	{"gee.io/a.pdf?a=b", true},
	{"me@gmail.com", false},
}

func TestParserExclude(t *testing.T) {
	job := &Crawl{
		ID:       "a",
		RootHost: "gee.io",
	}

	parser := NewParser(job, nil)
	assert := require.New(t)

	for _, tc := range parserExcludeTestCase {
		u, err := url.Parse(tc.url)
		assert.NoError(err)
		assert.Equal(tc.skip, parser.ShouldSkip(u))
	}

}
