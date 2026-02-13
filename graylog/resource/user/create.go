package user

import (
	"context"

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
	if _, ok := data[keyPermissions]; !ok {
		data[keyPermissions] = []string{}
	}

	if _, err := cl.User.Create(ctx, data); err != nil {
		return err
	}
	id := data[keyUsername].(string)
	d.SetId(id)
	return util.ReadAfterCreate(d, m, id, read)
}
