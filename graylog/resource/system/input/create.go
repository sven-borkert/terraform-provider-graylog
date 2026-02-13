package input

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func create(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	data, err := getDataFromResourceData(d)
	if err != nil {
		return err
	}

	input, _, err := cl.Input.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to create a input: %w", err)
	}
	id := input[keyID].(string)
	d.SetId(id)
	return util.ReadAfterCreate(d, m, id, read)
}
