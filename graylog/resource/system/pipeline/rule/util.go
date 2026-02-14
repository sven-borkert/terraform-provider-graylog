package rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

const (
	keyID     = "id"
	keySource = "source"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}

	// Compute content_hash from the rule source for cache invalidation workaround.
	// This allows users to use lifecycle { replace_triggered_by } on pipeline
	// connections to force Graylog to reload rules when content changes.
	if source, ok := data[keySource].(string); ok {
		if err := d.Set("content_hash", util.ComputeSHA256(source)); err != nil {
			return err
		}
	}

	d.SetId(data[keyID].(string))
	return nil
}
