package input

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestRefreshPreservesSecretsAndAllowsRotation(t *testing.T) {
	for _, tc := range []struct{ name, prior, remote, want string }{
		{"masked password", `{"password":"original","port":1}`, `{"password":"<password set>","port":2}`, `{"password":"original","port":2}`},
		{"empty password", `{"password":""}`, `{"password":"<password set>"}`, `{"password":""}`},
		{"import", ``, `{"password":"<password set>"}`, `{"password":"<password set>"}`},
		{"remote removal", `{"password":"original"}`, `{}`, `{}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, Resource().Schema, map[string]interface{}{"attributes": tc.prior})
			var attrs map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(tc.remote), &attrs))
			require.NoError(t, setDataToResourceData(d, map[string]interface{}{"id": "input", "attributes": attrs}))
			require.JSONEq(t, tc.want, d.Get(keyAttributes).(string))
			if tc.name == "masked password" {
				suppress := Resource().Schema[keyAttributes].DiffSuppressFunc
				require.True(t, suppress(keyAttributes, d.Get(keyAttributes).(string), tc.want, d))
				require.False(t, suppress(keyAttributes, d.Get(keyAttributes).(string), `{"password":"rotated","port":2}`, d))
			}
		})
	}
}
