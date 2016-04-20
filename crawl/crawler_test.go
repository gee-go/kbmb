package crawl

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrawl(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
      case 
    }
	}))
	defer ts.Close()

	crawler, err := New(ts.URL)
	assert := require.New(t)
	assert.NoError(err)

	for e := range crawler.Run() {
		fmt.Println(e)
	}

	fmt.Println(ts.URL)
}
