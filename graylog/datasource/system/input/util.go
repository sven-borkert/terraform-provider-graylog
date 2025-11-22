package input

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	rinput "github.com/sven-borkert/terraform-provider-graylog/graylog/resource/system/input"
)

// normalizeConfiguration ensures Graylog 7 responses that return "configuration" are mapped to "attributes".
func normalizeConfiguration(data map[string]interface{}) {
	if _, ok := data["attributes"]; ok {
		return
	}
	if cfg, ok := data["configuration"]; ok && cfg != nil {
		data["attributes"] = cfg
		delete(data, "configuration")
	}
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}, _ *http.Response) error {
	if err := convert.DataToJSON(data, "attributes"); err != nil {
		return err
	}
	if err := convert.SetResourceData(d, rinput.Resource(), data); err != nil {
		return err
	}
	if id, ok := data["id"]; ok {
		d.SetId(id.(string))
		_ = d.Set("input_id", id.(string))
	}
	return nil
}
