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

	// Iterate all states (tabs) and generate search queries
	stateMap := data[keyState].(map[string]interface{})
	if len(stateMap) == 0 {
		return errors.New("dashboard state is empty")
	}

	defaultTimerange := getDefaultTimerange()
	queries := make([]interface{}, 0, len(stateMap))

	for stateID, sv := range stateMap {
		state := sv.(map[string]interface{})

		widgets, ok := state[keyWidgets].([]interface{})
		if !ok {
			widgets = []interface{}{}
		}

		// Extract per-tab query string
		queryString, _ := state[keyQueryString].(string)
		delete(state, keyQueryString)

		query, widgetMapping, err := generateQueryFromWidgets(stateID, queryString, widgets, defaultTimerange)
		if err != nil {
			return fmt.Errorf("failed to generate search query for tab %s: %w", stateID, err)
		}
		queries = append(queries, query)
		applyWidgetMappingToState(state, widgetMapping)
	}

	searchData := buildSearchObject(queries)

	// Create the search first
	searchResp, _, err := cl.ViewSearch.Create(ctx, searchData)
	if err != nil {
		return fmt.Errorf("failed to create search for dashboard: %w", err)
	}

	searchID, err := validateSearchCreation(searchResp)
	if err != nil {
		return fmt.Errorf("search creation failed: %w", err)
	}
	log.Printf("[DEBUG] Created search %s for dashboard with %d queries", searchID, len(queries))

	// Update the dashboard data with the new search_id
	data["search_id"] = searchID

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
	return util.ReadAfterCreate(d, m, dID, read)
}
