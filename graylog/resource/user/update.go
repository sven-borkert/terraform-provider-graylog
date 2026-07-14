package user

import (
	"context"
	"errors"
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
	if _, ok := data[keyPermissions]; !ok {
		data[keyPermissions] = []string{}
	}

	// Graylog 7 addresses the update endpoint by the user's Mongo ID, not the
	// username. Prefer the captured user_id; fall back to resolving it by
	// username for state written before user_id was captured.
	userID, _ := d.Get(keyUserID).(string)
	if userID == "" {
		u, _, gErr := cl.User.Get(ctx, d.Id())
		if gErr != nil {
			return fmt.Errorf("resolving user id for update: %w", gErr)
		}
		userID, _ = u["id"].(string)
	}
	if userID == "" {
		return errors.New("could not determine user id for update")
	}

	// The main update endpoint ignores password; change it via the dedicated
	// endpoint when it has changed.
	if d.HasChange(keyPassword) {
		if pw, ok := data[keyPassword].(string); ok && pw != "" {
			if _, err := cl.User.ChangePassword(ctx, userID, pw); err != nil {
				return err
			}
		}
	}
	delete(data, keyPassword)

	// Remove computed fields for Graylog 7.0 compatibility
	util.RemoveComputedFields(data)

	if _, err := cl.User.Update(ctx, userID, data); err != nil {
		return err
	}
	return nil
}
