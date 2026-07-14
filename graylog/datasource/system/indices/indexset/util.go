package indexset

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/resource/system/indices/indexset"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	id, ok := util.RenameKey(data, "id", "index_set_id")
	// Graylog 7 returns data_tiering and field_restrictions as JSON objects; encode
	// them to strings for the TypeString schema attributes, like the resource read.
	if err := convert.DataToJSON(data, "rotation_strategy", "retention_strategy", "data_tiering", "field_restrictions"); err != nil {
		return err
	}

	if err := convert.SetResourceData(d, indexset.Resource(), data); err != nil {
		return err
	}

	if ok {
		d.SetId(id.(string))
	}
	return nil
}
