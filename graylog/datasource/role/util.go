package role

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyName        = "name"
	keyDescription = "description"
	keyPermissions = "permissions"
	keyReadOnly    = "read_only"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	// Set ID from the name
	name, ok := data[keyName].(string)
	if !ok {
		return nil
	}
	d.SetId(name)

	if err := d.Set(keyName, name); err != nil {
		return err
	}

	if description, ok := data[keyDescription].(string); ok {
		if err := d.Set(keyDescription, description); err != nil {
			return err
		}
	}

	if permissions, ok := data[keyPermissions].([]interface{}); ok {
		if err := d.Set(keyPermissions, permissions); err != nil {
			return err
		}
	}

	if readOnly, ok := data[keyReadOnly].(bool); ok {
		if err := d.Set(keyReadOnly, readOnly); err != nil {
			return err
		}
	}

	return nil
}
