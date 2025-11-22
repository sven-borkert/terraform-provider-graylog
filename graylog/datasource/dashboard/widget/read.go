package widget

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	dID := d.Get("dashboard_id").(string)
	wID := d.Get("widget_id").(string)
	data, resp, err := cl.DashboardWidget.Get(ctx, dID, wID)
	if err != nil {
		return err
	}
	return setDataToResourceData(d, data, resp)
}
