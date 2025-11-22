package rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	ruleRes "github.com/sven-borkert/terraform-provider-graylog/graylog/resource/system/pipeline/rule"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.SetResourceData(d, ruleRes.Resource(), data); err != nil {
		return err
	}
	if id, ok := data["id"]; ok {
		d.SetId(id.(string))
		_ = d.Set("rule_id", id.(string))
	}
	if t, ok := data["title"]; ok {
		_ = d.Set("title", t)
	}
	return nil
}
