package rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	ruleRes "github.com/sven-borkert/terraform-provider-graylog/graylog/resource/stream/rule"
)

const (
	keyID       = "id"
	keyRuleID   = "rule_id"
	keyStreamID = "stream_id"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.SetResourceData(d, ruleRes.Resource(), data); err != nil {
		return err
	}
	if id, ok := data[keyID]; ok {
		_ = d.Set(keyRuleID, id)
	}
	d.SetId(d.Get(keyStreamID).(string) + "/" + d.Get(keyRuleID).(string))
	return nil
}
