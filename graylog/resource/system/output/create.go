package output

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

	output, _, err := cl.Output.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to create a output: %w", err)
	}
	id := output[keyID].(string)
	d.SetId(id)
	return util.ReadAfterCreate(d, m, id, read)
}
