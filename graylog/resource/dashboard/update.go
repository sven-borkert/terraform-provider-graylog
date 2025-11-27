package dashboard

import (
	"context"
	"errors"
	"fmt"
	"log"

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

	// Create a new search (Graylog creates a new search on each dashboard update)
	searchResp, _, err := cl.ViewSearch.Create(ctx, searchData)
	if err != nil {
		return fmt.Errorf("failed to create search for dashboard update: %w", err)
	}

	searchID, err := validateSearchCreation(searchResp)
	if err != nil {
		return fmt.Errorf("search creation failed: %w", err)
	}
	log.Printf("[DEBUG] Created search %s for dashboard update", searchID)

	// Update the dashboard data with the new search_id and widget_mapping
	data["search_id"] = searchID
	applyWidgetMappingToState(state, widgetMapping)

	// Remove computed fields for Graylog 7.0 compatibility
	util.RemoveComputedFields(data)

	if _, _, err := cl.View.Update(ctx, d.Id(), data); err != nil {
		return fmt.Errorf("failed to update a dashboard %s: %w", d.Id(), err)
	}
	log.Printf("[DEBUG] Updated dashboard %s with search %s", d.Id(), searchID)
	return nil
}
