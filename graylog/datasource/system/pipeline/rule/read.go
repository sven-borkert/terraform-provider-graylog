package rule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	rID := d.Get("rule_id").(string)
	data, resp, err := cl.PipelineRule.Get(ctx, rID)
	if err != nil {
		return util.HandleGetResourceError(d, resp, err)
	}
	if err := setDataToResourceData(d, data); err != nil {
		return err
	}
	return nil
}
