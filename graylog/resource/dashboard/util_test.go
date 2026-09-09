package dashboard

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

// Regression test: ensure we can flatten a ViewDTO state map returned by the Views API.
func TestSetDataToResourceDataFlattensPositions(t *testing.T) {
	const body = `{
		"id":"view-id",
		"state":{
			"aeb86b45-a578-484a-a165-d2693f822150":{
				"widget_mapping":{},
				"positions":{"tf-e2e-agg":{"col":1,"row":1,"height":2,"width":2}},
				"widgets":[{
					"id":"tf-e2e-agg",
					"type":"aggregation",
					"config":{"row_pivots":[],"column_pivots":[],"series":[],"sort":[],"visualization":"bar","rollup":true},
					"timerange":{"type":"relative","range":300}
				}],
				"titles":{},
				"display_mode_settings":{"positions":{}}
			}
		},
		"search_id":"search-id",
		"title":"tf-e2e-dashboard",
		"type":"DASHBOARD"
	}`
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		t.Fatalf("failed to unmarshal fixture: %v", err)
	}

	d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{})
	if err := setDataToResourceData(d, data); err != nil {
		t.Fatalf("setDataToResourceData returned error: %v", err)
	}
	state := d.Get("state").([]interface{})[0].(map[string]interface{})
	if _, ok := state["positions"].(string); !ok {
		t.Fatalf("positions should be JSON string, got %T", state["positions"])
	}
	widget := state["widgets"].([]interface{})[0].(map[string]interface{})
	if _, ok := widget["timerange"].(string); !ok {
		t.Fatalf("timerange should be JSON string, got %T", widget["timerange"])
	}
}

func TestRefreshPreservesWidgetOrder(t *testing.T) {
	for _, tc := range []struct {
		name        string
		prior, want []string
	}{
		{"configured order", []string{"z", "a"}, []string{"z", "a", "b"}},
		{"removed widget", []string{"gone", "z", "a"}, []string{"z", "a", "b"}},
		{"import", nil, []string{"a", "b", "z"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var widgets []interface{}
			for _, id := range tc.prior {
				widgets = append(widgets, map[string]interface{}{"widget_id": id, "type": "messages", "config": "{}"})
			}
			d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{"state": []interface{}{map[string]interface{}{"id": "tab", "widgets": widgets}}})
			var data map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(`{"id":"view","state":{"tab":{"widgets":[{"id":"b","type":"messages","config":{}},{"id":"a","type":"messages","config":{}},{"id":"z","type":"messages","config":{}}]}}}`), &data))
			require.NoError(t, setDataToResourceData(d, data))
			var got []string
			for _, w := range d.Get("state.0.widgets").([]interface{}) {
				got = append(got, w.(map[string]interface{})[keyWidgetID].(string))
			}
			require.Equal(t, tc.want, got)
		})
	}
}
