package role

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	name := d.Get(keyName).(string)
	data, _, err := cl.Role.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to get role %s: %w", name, err)
	}
	return setDataToResourceData(d, data)
}
