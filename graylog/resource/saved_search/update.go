package saved_search

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func update(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	// Get existing state ID or generate new one
	stateID := d.Get("state_id").(string)
	if stateID == "" {
		stateID = uuid.New().String()
	}
	widgetID := uuid.New().String()
	searchTypeID := uuid.New().String()

	// Build the timerange
	timerange := map[string]interface{}{
		"type":  d.Get("timerange_type").(string),
		"range": d.Get("timerange_range").(int),
	}

	// Get query
	query := d.Get("query").(string)

	// Get streams
	var streams []interface{}
	if v, ok := d.GetOk("streams"); ok {
		for _, s := range v.(*schema.Set).List() {
			streams = append(streams, s.(string))
		}
	}

	// Get selected fields
	var selectedFields []interface{}
	if v, ok := d.GetOk("selected_fields"); ok {
		for _, f := range v.([]interface{}) {
			selectedFields = append(selectedFields, f.(string))
		}
	}

	// Build message list widget config
	// Widget sort needs "type" field for SortConfigDTO
	widgetSortOrder := d.Get("sort_order").(string)
	widgetConfig := map[string]interface{}{
		"decorators":       []interface{}{},
		"fields":           selectedFields,
		"show_message_row": true,
		"show_summary":     false,
		"sort":             []interface{}{map[string]interface{}{"type": "pivot", "field": d.Get("sort_field").(string), "direction": widgetSortOrder}},
	}

	// Build the widget
	widget := map[string]interface{}{
		"id":        widgetID,
		"type":      "messages",
		"timerange": timerange,
		"streams":   streams,
		"query": map[string]interface{}{
			"type":         "elasticsearch",
			"query_string": query,
		},
		"config": widgetConfig,
	}

	// Convert sort order to API format (DESC/ASC)
	sortOrder := d.Get("sort_order").(string)
	if sortOrder == "Descending" {
		sortOrder = "DESC"
	} else if sortOrder == "Ascending" {
		sortOrder = "ASC"
	}

	// Build the search type for the messages widget
	searchType := map[string]interface{}{
		"id":        searchTypeID,
		"name":      "messages",
		"type":      "messages",
		"timerange": timerange,
		"streams":   streams,
		"query": map[string]interface{}{
			"type":         "elasticsearch",
			"query_string": query,
		},
		"filter":     nil,
		"filters":    []interface{}{},
		"decorators": []interface{}{},
		"limit":      150,
		"offset":     0,
		"sort":       []interface{}{map[string]interface{}{"field": d.Get("sort_field").(string), "order": sortOrder}},
	}

	// Build the search object
	searchData := map[string]interface{}{
		"parameters":            []interface{}{},
		"skip_no_streams_check": false,
		"queries": []interface{}{
			map[string]interface{}{
				"id":      stateID,
				"filters": []interface{}{},
				"filter":  nil,
				"query": map[string]interface{}{
					"type":         "elasticsearch",
					"query_string": query,
				},
				"timerange":    timerange,
				"search_types": []interface{}{searchType},
			},
		},
	}

	// Create new search
	searchResp, _, err := cl.ViewSearch.Create(ctx, searchData)
	if err != nil {
		return fmt.Errorf("failed to create search for saved search update: %w", err)
	}

	searchID, ok := searchResp["id"].(string)
	if !ok || searchID == "" {
		return fmt.Errorf("failed to get search ID from response")
	}
	log.Printf("[DEBUG] Created new search %s for saved search update", searchID)

	// Build widget mapping
	widgetMapping := map[string]interface{}{
		widgetID: []interface{}{searchTypeID},
	}

	// Build widget positions
	positions := map[string]interface{}{
		widgetID: map[string]interface{}{
			"col":    1,
			"row":    1,
			"height": 6,
			"width":  "Infinity",
		},
	}

	// Build the state
	state := map[string]interface{}{
		stateID: map[string]interface{}{
			"selected_fields": selectedFields,
			"titles": map[string]interface{}{
				"widget": map[string]interface{}{},
			},
			"widgets":        []interface{}{widget},
			"widget_mapping": widgetMapping,
			"positions":      positions,
		},
	}

	// Build the view data
	viewData := map[string]interface{}{
		"type":        "SEARCH",
		"title":       d.Get("title").(string),
		"description": d.Get("description").(string),
		"summary":     d.Get("summary").(string),
		"search_id":   searchID,
		"state":       state,
		"properties":  []interface{}{},
	}

	// Update the view
	_, _, err = cl.View.Update(ctx, d.Id(), viewData)
	if err != nil {
		// Clean up the new search if view update fails
		_, _ = cl.ViewSearch.Delete(ctx, searchID)
		return fmt.Errorf("failed to update saved search view: %w", err)
	}

	if err := d.Set("search_id", searchID); err != nil {
		return err
	}
	if err := d.Set("state_id", stateID); err != nil {
		return err
	}

	log.Printf("[DEBUG] Updated saved search %s with search %s", d.Id(), searchID)
	return nil
}
