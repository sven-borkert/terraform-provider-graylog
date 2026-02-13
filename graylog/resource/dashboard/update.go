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

	// Create a new search (Graylog creates a new search on each dashboard update)
	searchResp, _, err := cl.ViewSearch.Create(ctx, searchData)
	if err != nil {
		return fmt.Errorf("failed to create search for dashboard update: %w", err)
	}

	searchID, err := validateSearchCreation(searchResp)
	if err != nil {
		return fmt.Errorf("search creation failed: %w", err)
	}
	log.Printf("[DEBUG] Created search %s for dashboard update with %d queries", searchID, len(queries))

	// Update the dashboard data with the new search_id
	data["search_id"] = searchID

	// Remove computed fields for Graylog 7.0 compatibility
	util.RemoveComputedFields(data)

	if _, _, err := cl.View.Update(ctx, d.Id(), data); err != nil {
		return fmt.Errorf("failed to update a dashboard %s: %w", d.Id(), err)
	}
	log.Printf("[DEBUG] Updated dashboard %s with search %s", d.Id(), searchID)
	return nil
}
