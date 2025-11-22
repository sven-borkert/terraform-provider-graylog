package grok

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	resourceGrok "github.com/sven-borkert/terraform-provider-graylog/graylog/resource/system/grok"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}, _ *http.Response) error {
	if err := convert.SetResourceData(d, resourceGrok.Resource(), data); err != nil {
		return err
	}
	if id, ok := data["id"]; ok {
		d.SetId(id.(string))
		_ = d.Set("pattern_id", id.(string))
	}
	if n, ok := data["name"]; ok {
		_ = d.Set("name", n)
	}
	return nil
}
