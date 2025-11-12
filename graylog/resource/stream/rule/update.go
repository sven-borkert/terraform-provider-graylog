package rule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-graylog/terraform-provider-graylog/graylog/client"
	"github.com/terraform-provider-graylog/terraform-provider-graylog/graylog/util"
)

func update(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	data, err := getDataFromResourceData(d)
	if err != nil {
		return err
	}
	delete(data, keyStreamID)
	delete(data, keyRuleID)

	// Remove computed fields for Graylog 7.0 compatibility
	util.RemoveComputedFields(data)

	sID := d.Get(keyStreamID).(string)
	rID := d.Get(keyRuleID).(string)
	if _, _, err := cl.StreamRule.Update(ctx, sID, rID, data); err != nil {
		return fmt.Errorf(
			"failed to update a stream rule (stream id: %s, rule id: %s): %w", sID, rID, err)
	}
	return nil
}
