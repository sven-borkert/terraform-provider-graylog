package fieldtype

import (
	"context"
	"net/http"
	"testing"

	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

// newTestClient wires the fieldtype client to a flute-mocked HTTP transport.
func newTestClient(t *testing.T, body string) Client {
	t.Helper()
	route := flute.Route{
		Name:    "get field types",
		Matcher: flute.Matcher{Method: "GET"},
		Tester:  flute.Tester{Path: "/api/system/indices/index_sets/types/idx1"},
		Response: flute.Response{
			Base:       http.Response{StatusCode: 200},
			BodyString: body,
		},
	}
	http.DefaultClient = &http.Client{
		Transport: flute.Transport{
			T: t,
			Services: []flute.Service{
				{Endpoint: "http://example.com", Routes: []flute.Route{route}},
			},
		},
	}
	return Client{Client: httpclient.New("http://example.com/api")}
}

// TestGetFieldTypeOrigin verifies that only OVERRIDDEN_* origins (custom mappings)
// are returned; index/profile defaults are treated as "no custom mapping".
func TestGetFieldTypeOrigin(t *testing.T) {
	ctx := context.Background()

	t.Run("custom override is returned", func(t *testing.T) {
		cl := newTestClient(t, `{"elements": [
			{"field_name": "http_status", "type": "long", "origin": "OVERRIDDEN_INDEX", "is_reserved": false}
		]}`)
		ft, _, err := cl.GetFieldType(ctx, "idx1", "http_status")
		if err != nil {
			t.Fatal(err)
		}
		if ft == nil {
			t.Fatal("expected a custom field type, got nil")
		}
		if ft.Type != "long" {
			t.Fatalf("expected type long, got %q", ft.Type)
		}
	})

	t.Run("index default is not a custom mapping", func(t *testing.T) {
		cl := newTestClient(t, `{"elements": [
			{"field_name": "http_status", "type": "long", "origin": "INDEX", "is_reserved": false}
		]}`)
		ft, _, err := cl.GetFieldType(ctx, "idx1", "http_status")
		if err != nil {
			t.Fatal(err)
		}
		if ft != nil {
			t.Fatalf("expected nil for non-custom origin, got %+v", ft)
		}
	})
}
