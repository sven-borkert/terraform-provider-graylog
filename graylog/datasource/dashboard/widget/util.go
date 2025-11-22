package widget

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	resWidget "github.com/sven-borkert/terraform-provider-graylog/graylog/resource/dashboard/widget"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}, _ *http.Response) error {
	if err := convert.DataToJSON(data, "config"); err != nil {
		return err
	}
	if err := convert.SetResourceData(d, resWidget.Resource(), data); err != nil {
		return err
	}
	// Set composite ID to match resource importer behavior
	d.SetId(d.Get("dashboard_id").(string) + "/" + d.Get("widget_id").(string))
	return nil
}
