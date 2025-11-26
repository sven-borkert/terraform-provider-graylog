package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
)

const (
	keyUsername       = "username"
	keyPermissions    = "permissions"
	keyClientAddress  = "client_address"
	keyExternal       = "external"
	keyLastActivity   = "last_activity"
	keyUserID         = "user_id"
	keySessionActive  = "session_active"
	keyReadOnly       = "read_only"
	keyFullName       = "full_name"
	keyAccountStatus  = "account_status"
	keyServiceAccount = "service_account"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}

	// Remove computed fields that should not be sent to the API
	delete(data, keyClientAddress)
	delete(data, keyExternal)
	delete(data, keyLastActivity)
	delete(data, keySessionActive)
	delete(data, keyReadOnly)
	delete(data, keyUserID)
	delete(data, keyFullName) // full_name is computed from first_name and last_name
	delete(data, keyAccountStatus)

	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}
	d.SetId(data[keyUsername].(string))
	return nil
}
