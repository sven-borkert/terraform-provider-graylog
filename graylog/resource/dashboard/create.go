package dashboard

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func create(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	// First, get the dashboard data from resource data
	data, err := getDataFromResourceData(d)
	if err != nil {
		return err
	}

	// Extract state and widgets to generate search_types
	stateMap := data[keyState].(map[string]interface{})
	var stateID string
	var state map[string]interface{}
	for k, v := range stateMap {
		stateID = k
		state = v.(map[string]interface{})
		break
	}

	if state == nil {
		return errors.New("dashboard state is empty")
	}

	widgets, ok := state[keyWidgets].([]interface{})
	if !ok {
		widgets = []interface{}{}
	}

	// Generate search with search_types for each widget
	defaultTimerange := getDefaultTimerange()
	searchData, widgetMapping, err := generateSearchFromWidgets(stateID, widgets, defaultTimerange)
	if err != nil {
		return fmt.Errorf("failed to generate search from widgets: %w", err)
	}

	// Create the search first
	searchResp, _, err := cl.ViewSearch.Create(ctx, searchData)
	if err != nil {
		return fmt.Errorf("failed to create search for dashboard: %w", err)
	}

	searchID, err := validateSearchCreation(searchResp)
	if err != nil {
		return fmt.Errorf("search creation failed: %w", err)
	}
	log.Printf("[DEBUG] Created search %s for dashboard", searchID)

	// Update the dashboard data with the new search_id and widget_mapping
	data["search_id"] = searchID
	applyWidgetMappingToState(state, widgetMapping)

	// Create the dashboard view
	ds, _, err := cl.View.Create(ctx, data)
	if err != nil {
		// Try to clean up the search if dashboard creation fails
		_, _ = cl.ViewSearch.Delete(ctx, searchID)
		return fmt.Errorf("failed to create a dashboard view: %w", err)
	}

	a, ok := ds[keyID]
	if !ok {
		return errors.New("response body of Graylog API is unexpected. 'id' isn't found")
	}
	dID, ok := a.(string)
	if !ok {
		return fmt.Errorf(
			"response body of Graylog API is unexpected. 'id' should be string: %v", a)
	}

	d.SetId(dID)
	log.Printf("[DEBUG] Created dashboard %s with search %s", dID, searchID)
	return nil
}
