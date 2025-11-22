package saved

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if v, ok := data["id"]; ok {
		_ = d.Set("saved_search_id", v)
		d.SetId(v.(string))
	}
	for _, k := range []string{"title", "summary", "description", "owner", "created_at", "search_id"} {
		if v, ok := data[k]; ok {
			_ = d.Set(k, v)
		}
	}
	return nil
}
