package dashboard

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

	if id, ok := d.GetOk(keyDashboardID); ok {
		return readByID(ctx, d, cl, id.(string))
	}

	if t, ok := d.GetOk(keyTitle); ok {
		return readByTitle(ctx, d, cl, t.(string))
	}
	return errors.New("one of dashboard_id or title must be set")
}

func readByID(ctx context.Context, d *schema.ResourceData, cl client.Client, id string) error {
	if _, ok := d.GetOk(keyTitle); ok {
		return errors.New("both dashboard_id and title must not be set at the same time")
	}
	// Use the View client for Graylog 7.0+ which uses /views/{id} endpoint
	data, _, err := cl.View.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get dashboard %s: %w", id, err)
	}
	return setDataToResourceData(d, data)
}

func readByTitle(ctx context.Context, d *schema.ResourceData, cl client.Client, title string) error {
	dashboards, _, err := cl.Dashboard.Gets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list dashboards: %w", err)
	}

	// Graylog 7.0+ returns paginated response with "elements" key
	elements, ok := dashboards[keyElements]
	if !ok {
		// Fallback for older versions that might use "dashboards" or "views"
		elements, ok = dashboards[keyDashboards]
		if !ok {
			elements, ok = dashboards[keyViews]
			if !ok {
				return errors.New(`the response of Graylog API GET /api/dashboards is unexpected: no "elements", "dashboards", or "views" field found`)
			}
		}
	}

	var data map[string]interface{}
	for _, a := range elements.([]interface{}) {
		dashboard := a.(map[string]interface{})
		if dashboard[keyTitle].(string) == title {
			if data != nil {
				return fmt.Errorf("title %q is not unique, multiple dashboards found", title)
			}
			data = dashboard
		}
	}
	if data == nil {
		return fmt.Errorf("no dashboard found with title %q", title)
	}
	return setDataToResourceData(d, data)
}
