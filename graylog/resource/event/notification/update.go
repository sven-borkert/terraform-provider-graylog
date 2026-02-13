package notification

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
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
	// Remove computed fields for Graylog 7.0 compatibility
	util.RemoveComputedFields(data)

	// Graylog 7 Update requires id in body
	data[keyID] = d.Id()

	if _, _, err := cl.EventNotification.Update(ctx, d.Id(), data); err != nil {
		return fmt.Errorf("failed to update a event notification %s: %w", d.Id(), err)
	}
	return nil
}
