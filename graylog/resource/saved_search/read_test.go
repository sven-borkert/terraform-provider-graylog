package saved_search

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/config"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestReadSortDirection(t *testing.T) {
	for _, direction := range []string{"Ascending", "Descending"} {
		t.Run(direction, func(t *testing.T) {
			original := http.DefaultClient
			t.Cleanup(func() { http.DefaultClient = original })
			http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				require.Equal(t, "/api/views/view", r.URL.Path)
				body := fmt.Sprintf(`{"id":"view","state":{"tab":{"widgets":[{"config":{"sort":[{"type":"pivot","field":"timestamp","direction":%q}]}}]}}}`, direction)
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
			})}
			d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{"title": "search", "sort_order": "stale"})
			d.SetId("view")
			require.NoError(t, read(d, config.Config{Endpoint: "http://example.invalid/api"}))
			require.Equal(t, direction, d.Get("sort_order"))
		})
	}
}
