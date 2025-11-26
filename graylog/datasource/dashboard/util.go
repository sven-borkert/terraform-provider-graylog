package dashboard

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyID          = "id"
	keyDashboardID = "dashboard_id"
	keyElements    = "elements"
	keyDashboards  = "dashboards"
	keyViews       = "views"
	keyTitle       = "title"
	keyDescription = "description"
	keySummary     = "summary"
	keyOwner       = "owner"
	keySearchID    = "search_id"
	keyCreatedAt   = "created_at"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	id, ok := data[keyID].(string)
	if !ok {
		return nil
	}
	d.SetId(id)

	if err := d.Set(keyDashboardID, id); err != nil {
		return err
	}

	if title, ok := data[keyTitle].(string); ok {
		if err := d.Set(keyTitle, title); err != nil {
			return err
		}
	}

	if description, ok := data[keyDescription].(string); ok {
		if err := d.Set(keyDescription, description); err != nil {
			return err
		}
	}

	if summary, ok := data[keySummary].(string); ok {
		if err := d.Set(keySummary, summary); err != nil {
			return err
		}
	}

	if owner, ok := data[keyOwner].(string); ok {
		if err := d.Set(keyOwner, owner); err != nil {
			return err
		}
	}

	if searchID, ok := data[keySearchID].(string); ok {
		if err := d.Set(keySearchID, searchID); err != nil {
			return err
		}
	}

	if createdAt, ok := data[keyCreatedAt].(string); ok {
		if err := d.Set(keyCreatedAt, createdAt); err != nil {
			return err
		}
	}

	return nil
}
