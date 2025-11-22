package output

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	routput "github.com/sven-borkert/terraform-provider-graylog/graylog/resource/system/output"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}, _ *http.Response) error {
	if err := convert.DataToJSON(data, "configuration"); err != nil {
		return err
	}

	if err := convert.SetResourceData(d, routput.Resource(), data); err != nil {
		return err
	}

	if id, ok := data["id"]; ok {
		d.SetId(id.(string))
		_ = d.Set("output_id", id.(string))
	}

	return nil
}
