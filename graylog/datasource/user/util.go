package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyID             = "id"
	keyUserID         = "user_id"
	keyUsername       = "username"
	keyEmail          = "email"
	keyFullName       = "full_name"
	keyFirstName      = "first_name"
	keyLastName       = "last_name"
	keyTimezone       = "timezone"
	keySessionTimeout = "session_timeout_ms"
	keyRoles          = "roles"
	keyPermissions    = "permissions"
	keyExternal       = "external"
	keyReadOnly       = "read_only"
	keySessionActive  = "session_active"
	keyLastActivity   = "last_activity"
	keyClientAddress  = "client_address"
	keyAccountStatus  = "account_status"
	keyServiceAccount = "service_account"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	// Set ID from the "id" field or username
	if id, ok := data[keyID].(string); ok {
		d.SetId(id)
		if err := d.Set(keyUserID, id); err != nil {
			return err
		}
	} else if username, ok := data[keyUsername].(string); ok {
		d.SetId(username)
	}

	if username, ok := data[keyUsername].(string); ok {
		if err := d.Set(keyUsername, username); err != nil {
			return err
		}
	}

	if email, ok := data[keyEmail].(string); ok {
		if err := d.Set(keyEmail, email); err != nil {
			return err
		}
	}

	if fullName, ok := data[keyFullName].(string); ok {
		if err := d.Set(keyFullName, fullName); err != nil {
			return err
		}
	}

	if firstName, ok := data[keyFirstName].(string); ok {
		if err := d.Set(keyFirstName, firstName); err != nil {
			return err
		}
	}

	if lastName, ok := data[keyLastName].(string); ok {
		if err := d.Set(keyLastName, lastName); err != nil {
			return err
		}
	}

	if timezone, ok := data[keyTimezone].(string); ok {
		if err := d.Set(keyTimezone, timezone); err != nil {
			return err
		}
	}

	if sessionTimeout, ok := data[keySessionTimeout].(float64); ok {
		if err := d.Set(keySessionTimeout, int(sessionTimeout)); err != nil {
			return err
		}
	}

	if roles, ok := data[keyRoles].([]interface{}); ok {
		if err := d.Set(keyRoles, roles); err != nil {
			return err
		}
	}

	if permissions, ok := data[keyPermissions].([]interface{}); ok {
		if err := d.Set(keyPermissions, permissions); err != nil {
			return err
		}
	}

	if external, ok := data[keyExternal].(bool); ok {
		if err := d.Set(keyExternal, external); err != nil {
			return err
		}
	}

	if readOnly, ok := data[keyReadOnly].(bool); ok {
		if err := d.Set(keyReadOnly, readOnly); err != nil {
			return err
		}
	}

	if sessionActive, ok := data[keySessionActive].(bool); ok {
		if err := d.Set(keySessionActive, sessionActive); err != nil {
			return err
		}
	}

	if lastActivity, ok := data[keyLastActivity].(string); ok {
		if err := d.Set(keyLastActivity, lastActivity); err != nil {
			return err
		}
	}

	if clientAddress, ok := data[keyClientAddress].(string); ok {
		if err := d.Set(keyClientAddress, clientAddress); err != nil {
			return err
		}
	}

	if accountStatus, ok := data[keyAccountStatus].(string); ok {
		if err := d.Set(keyAccountStatus, accountStatus); err != nil {
			return err
		}
	}

	if serviceAccount, ok := data[keyServiceAccount].(bool); ok {
		if err := d.Set(keyServiceAccount, serviceAccount); err != nil {
			return err
		}
	}

	return nil
}
