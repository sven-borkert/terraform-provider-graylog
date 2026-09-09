package indexset

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/config"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestDefaultIndexSetLifecycle(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		create, wantDefault, fail bool
	}{
		{name: "create default", create: true, wantDefault: true},
		{name: "update default", wantDefault: true},
		{name: "create ordinary", create: true},
		{name: "update ordinary"},
		{name: "create default fails", create: true, wantDefault: true, fail: true},
		{name: "update default fails", wantDefault: true, fail: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := http.DefaultClient
			t.Cleanup(func() { http.DefaultClient = original })
			promoted := false
			writes := 0
			http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
				status := 200
				body := `{"id":"set-id","default":false}`
				switch r.Method + " " + r.URL.Path {
				case "POST /api/system/indices/index_sets", "PUT /api/system/indices/index_sets/set-id":
					var data map[string]interface{}
					require.NoError(t, json.NewDecoder(r.Body).Decode(&data))
					require.NotContains(t, data, keyDefault)
					writes++
				case "PUT /api/system/indices/index_sets/set-id/default":
					require.Equal(t, 1, writes, "default selection must follow the index-set write")
					if tc.fail {
						status = 403
						body = `{"message":"forbidden"}`
					} else {
						promoted = true
						body = `{"id":"set-id","default":true}`
					}
				case "GET /api/system/deflector/set-id":
					body = `{"is_up":true,"current_target":"logs_0"}`
				case "GET /api/system/indices/index_sets/set-id":
					if promoted {
						body = `{"id":"set-id","default":true}`
					}
				default:
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
			})}
			d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{"title": "logs", "index_prefix": "logs", "default": tc.wantDefault})
			cfg := config.Config{Endpoint: "http://example.invalid/api"}
			var err error
			if tc.create {
				err = create(d, cfg)
			} else {
				d.SetId("set-id")
				err = update(d, cfg)
			}
			if tc.fail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantDefault, promoted)
				require.Equal(t, tc.wantDefault, d.Get(keyDefault))
			}
			require.Equal(t, "set-id", d.Id(), "retain the created object when promotion fails")
		})
	}
}

func TestCannotUnsetDefaultWithoutReplacement(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Error("must reject demotion before sending a request")
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}
	_, diags := Resource().Apply(context.Background(),
		&terraform.InstanceState{ID: "set-id", Attributes: map[string]string{"default": "true"}},
		&terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{"default": {Old: "true", New: "false"}}},
		config.Config{Endpoint: "http://example.invalid/api"})
	require.True(t, diags.HasError())
	require.Contains(t, diags[0].Summary, "another index set")
}
