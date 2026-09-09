package sidecar

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

func TestReadThenDestroyOnlyManagedSidecars(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	var deleted map[string]interface{}
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := `{"sidecars":[{"node_id":"managed","assignments":[{"collector_id":"c","configuration_id":"changed"}]},{"node_id":"unmanaged","assignments":[{"collector_id":"c","configuration_id":"other"}]}]}`
		if r.Method != http.MethodGet {
			require.NoError(t, json.NewDecoder(r.Body).Decode(&deleted))
			body = `{}`
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}
	d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{"sidecars": []interface{}{map[string]interface{}{"node_id": "managed", "assignments": []interface{}{map[string]interface{}{"collector_id": "c", "configuration_id": "original"}}}}})
	cfg := config.Config{Endpoint: "http://example.invalid/api"}
	require.NoError(t, create(d, cfg))
	require.NotContains(t, deleted, "assignment_scope_version")
	require.NoError(t, read(d, cfg))
	nodes := d.Get(keySidecars).(*schema.Set).List()
	require.Len(t, nodes, 1)
	require.Equal(t, "changed", nodes[0].(map[string]interface{})[keyAssignments].(*schema.Set).List()[0].(map[string]interface{})["configuration_id"])
	require.NoError(t, destroy(d, cfg))
	require.Equal(t, map[string]interface{}{"nodes": []interface{}{map[string]interface{}{"node_id": "managed", "assignments": []interface{}{}}}}, deleted)
}

func TestImportSeedsSidecarMembership(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodGet, r.Method)
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"sidecars":[{"node_id":"imported","assignments":[]}]}`))}, nil
	})}
	d := schema.TestResourceDataRaw(t, Resource().Schema, nil)
	d.SetId("all")
	cfg := config.Config{Endpoint: "http://example.invalid/api"}
	imported, err := Resource().Importer.StateContext(context.Background(), d, cfg)
	require.NoError(t, err)
	require.Len(t, imported, 1)
	require.NoError(t, read(imported[0], cfg))
	require.Equal(t, 1, imported[0].Get(keySidecars).(*schema.Set).Len())
}

func TestUpdateClearsRemovedSidecarAssignments(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	var sent map[string]interface{}
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		require.Equal(t, "PUT", r.Method)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&sent))
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}
	node := func(id string) interface{} {
		return map[string]interface{}{"node_id": id, "assignments": []interface{}{map[string]interface{}{"collector_id": "c", "configuration_id": "config"}}}
	}
	r := Resource()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"sidecars": []interface{}{node("removed"), node("retained")}})
	d.SetId(systemID)
	require.NoError(t, d.Set("assignment_scope_version", 1))
	state := d.State()
	cfg := config.Config{Endpoint: "http://example.invalid/api"}
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(map[string]interface{}{"sidecars": []interface{}{node("retained")}}), cfg)
	require.NoError(t, err)
	_, diags := r.Apply(context.Background(), state, diff, cfg)
	require.False(t, diags.HasError(), "%v", diags)
	require.ElementsMatch(t, []interface{}{node("retained"), map[string]interface{}{"node_id": "removed", "assignments": []interface{}{}}}, sent[keyNodes])
}

func TestLegacySidecarStateCannotModifyAssignments(t *testing.T) {
	original := http.DefaultClient
	t.Cleanup(func() { http.DefaultClient = original })
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		t.Error("legacy ownership must be rejected before API requests")
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})}
	for _, tc := range []struct {
		name string
		fn   func(*schema.ResourceData, interface{}) error
	}{
		{"refresh", read}, {"update", update}, {"destroy", destroy},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{"sidecars": []interface{}{map[string]interface{}{"node_id": "possibly-unmanaged", "assignments": []interface{}{}}}})
			d.SetId(systemID)
			require.ErrorContains(t, tc.fn(d, config.Config{Endpoint: "http://example.invalid/api"}), "legacy sidecar state")
		})
	}
}
