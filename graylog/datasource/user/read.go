package user

import (
	"context"
	"errors"
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

	if id, ok := d.GetOk(keyUserID); ok {
		return readByID(ctx, d, cl, id.(string))
	}

	if username, ok := d.GetOk(keyUsername); ok {
		return readByUsername(ctx, d, cl, username.(string))
	}

	return errors.New("one of user_id or username must be set")
}

func readByID(ctx context.Context, d *schema.ResourceData, cl client.Client, id string) error {
	if _, ok := d.GetOk(keyUsername); ok {
		return errors.New("both user_id and username must not be set at the same time")
	}
	data, _, err := cl.User.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user by ID %s: %w", id, err)
	}
	return setDataToResourceData(d, data)
}

func readByUsername(ctx context.Context, d *schema.ResourceData, cl client.Client, username string) error {
	data, _, err := cl.User.Get(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user %s: %w", username, err)
	}
	return setDataToResourceData(d, data)
}
