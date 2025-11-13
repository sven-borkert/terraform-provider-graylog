package setting

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
	delete(data, keySystemPasswordSet)

	// Remove computed fields for Graylog 7.0 compatibility
	util.RemoveComputedFields(data)

	if _, err := cl.LDAPSetting.Update(ctx, data); err != nil {
		return fmt.Errorf("failed to update a LDAP setting: %w", err)
	}
	return nil
}
