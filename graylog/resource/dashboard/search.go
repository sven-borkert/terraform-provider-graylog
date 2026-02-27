package dashboard

import (
	"fmt"

	"github.com/google/uuid"
)

// generateQueryFromWidgets creates a single search query entry with search_types derived from widget configs.
// It returns the query entry and a widget_mapping that links widget IDs to search_type IDs.
func generateQueryFromWidgets(stateID string, queryString string, widgets []interface{}, defaultTimerange map[string]interface{}) (map[string]interface{}, map[string][]string, error) {
	searchTypes := make([]interface{}, 0, len(widgets))
	widgetMapping := make(map[string][]string)

	for _, w := range widgets {
		widget := w.(map[string]interface{})
		widgetID := ""
		if id, ok := widget["id"].(string); ok && id != "" {
			widgetID = id
		} else if id, ok := widget[keyWidgetID].(string); ok && id != "" {
			widgetID = id
		}

		widgetType := ""
		if t, ok := widget["type"].(string); ok {
			widgetType = t
		}

		// Only aggregation widgets need search_types
		if widgetType != "aggregation" {
			continue
		}

		config, ok := widget[keyConfig].(map[string]interface{})
		if !ok {
			continue
		}

		// Generate a unique ID for this search_type
		searchTypeID := uuid.New().String()

		// Extract streams from widget
		var widgetStreams []interface{}
		if s, ok := widget["streams"].([]interface{}); ok {
			widgetStreams = s
		}

		// Extract per-widget query
		var widgetQuery map[string]interface{}
		if q, ok := widget["query"].(map[string]interface{}); ok {
			widgetQuery = q
		}

		// Create the search_type from widget config
		searchType := createSearchTypeFromWidgetConfig(searchTypeID, config, widget[keyTimerange], defaultTimerange, widgetStreams, widgetQuery)

		searchTypes = append(searchTypes, searchType)
		widgetMapping[widgetID] = []string{searchTypeID}
	}

	// Build persistent filter from queryString so it survives UI interactions.
	// The query.query_string maps to the editable search bar which resets on search/timeframe change,
	// so we use the "filter" field instead for persistent filtering.
	// Note: The filter UI ("Search Filters") is a Graylog commercial edition feature.
	// In open-source Graylog, the filter is applied but not visible in the dashboard UI.
	var filter interface{}
	if queryString != "" {
		filter = map[string]interface{}{
			"type": "or",
			"filters": []interface{}{
				map[string]interface{}{
					"type":    "query_string",
					"query":   queryString,
					"filters": []interface{}{},
				},
			},
		}
	}

	query := map[string]interface{}{
		"id":      stateID,
		"filters": []interface{}{},
		"filter":  filter,
		"query": map[string]interface{}{
			"type":         "elasticsearch",
			"query_string": "",
		},
		"timerange":    defaultTimerange,
		"search_types": searchTypes,
	}

	return query, widgetMapping, nil
}

// buildSearchObject assembles multiple query entries into a complete search object.
func buildSearchObject(queries []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"parameters":            []interface{}{},
		"skip_no_streams_check": false,
		"queries":               queries,
	}
}

// createSearchTypeFromWidgetConfig converts a widget config to a search_type (pivot query).
func createSearchTypeFromWidgetConfig(id string, config map[string]interface{}, widgetTimerange interface{}, defaultTimerange map[string]interface{}, streams []interface{}, widgetQuery map[string]interface{}) map[string]interface{} {
	if streams == nil {
		streams = []interface{}{}
	}

	// Include widget-level query in the search_type so it's applied on initial dashboard load
	var query interface{}
	if widgetQuery != nil {
		if qs, ok := widgetQuery["query_string"].(string); ok && qs != "" {
			query = map[string]interface{}{
				"type":         "elasticsearch",
				"query_string": qs,
			}
		}
	}

	searchType := map[string]interface{}{
		"id":                id,
		"type":              "pivot",
		"name":              "chart",
		"filter":            nil,
		"filters":           []interface{}{},
		"query":             query,
		"streams":           streams,
		"stream_categories": []interface{}{},
		"rollup":            true,
		"column_groups":     []interface{}{},
	}

	// Convert row_pivots to row_groups
	if rowPivots, ok := config["row_pivots"].([]interface{}); ok {
		rowGroups := make([]interface{}, 0, len(rowPivots))
		for _, rp := range rowPivots {
			pivot := rp.(map[string]interface{})
			rowGroup := convertPivotToGroup(pivot)
			rowGroups = append(rowGroups, rowGroup)
		}
		searchType["row_groups"] = rowGroups
	} else {
		searchType["row_groups"] = []interface{}{}
	}

	// Convert column_pivots to column_groups
	if colPivots, ok := config["column_pivots"].([]interface{}); ok {
		colGroups := make([]interface{}, 0, len(colPivots))
		for _, cp := range colPivots {
			pivot := cp.(map[string]interface{})
			colGroup := convertPivotToGroup(pivot)
			colGroups = append(colGroups, colGroup)
		}
		searchType["column_groups"] = colGroups
	}

	// Convert series
	if series, ok := config["series"].([]interface{}); ok {
		searchSeries := make([]interface{}, 0, len(series))
		for _, s := range series {
			seriesItem := s.(map[string]interface{})
			searchSeriesItem := convertSeries(seriesItem)
			searchSeries = append(searchSeries, searchSeriesItem)
		}
		searchType["series"] = searchSeries
	} else {
		searchType["series"] = []interface{}{}
	}

	// Copy sort
	if sort, ok := config["sort"].([]interface{}); ok {
		searchType["sort"] = sort
	} else {
		searchType["sort"] = []interface{}{}
	}

	// Copy rollup if specified
	if rollup, ok := config["rollup"].(bool); ok {
		searchType["rollup"] = rollup
	}

	// Set timerange from widget or use default
	if widgetTimerange != nil {
		if tr, ok := widgetTimerange.(map[string]interface{}); ok && len(tr) > 0 {
			searchType["timerange"] = tr
		} else {
			searchType["timerange"] = defaultTimerange
		}
	} else {
		searchType["timerange"] = defaultTimerange
	}

	return searchType
}

// convertPivotToGroup converts a row_pivot/column_pivot to a row_group/column_group.
func convertPivotToGroup(pivot map[string]interface{}) map[string]interface{} {
	group := map[string]interface{}{}

	// Copy type
	if t, ok := pivot["type"].(string); ok {
		group["type"] = t
	}

	// Copy fields
	if fields, ok := pivot["fields"].([]interface{}); ok {
		group["fields"] = fields
	}

	// Copy config options based on type
	if cfg, ok := pivot["config"].(map[string]interface{}); ok {
		// For "values" type, extract limit
		if limit, ok := cfg["limit"]; ok {
			group["limit"] = limit
		}
		// For "time" type, extract interval
		if interval, ok := cfg["interval"]; ok {
			group["interval"] = interval
		}
	}

	// Add skip_empty_values for values type
	if t, ok := pivot["type"].(string); ok && t == "values" {
		group["skip_empty_values"] = false
	}

	return group
}

// convertSeries converts a widget series to a search_type series.
func convertSeries(series map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	// Parse the function string like "count()" or "sum(packetbeat_network_bytes)"
	if fn, ok := series["function"].(string); ok {
		fnType, field := parseFunction(fn)
		result["type"] = fnType
		result["id"] = fn
		if field != "" {
			result["field"] = field
		} else {
			result["field"] = nil
		}
	}

	return result
}

// parseFunction parses a function string like "count()" or "sum(field_name)".
func parseFunction(fn string) (fnType string, field string) {
	// Find the opening parenthesis
	for i, c := range fn {
		if c == '(' {
			fnType = fn[:i]
			// Extract field name between parentheses
			if i+1 < len(fn) && fn[len(fn)-1] == ')' {
				field = fn[i+1 : len(fn)-1]
			}
			return
		}
	}
	// No parenthesis found, treat whole string as type
	return fn, ""
}

// applyWidgetMappingToState updates the state with the widget_mapping.
func applyWidgetMappingToState(state map[string]interface{}, widgetMapping map[string][]string) {
	state[keyWidgetMapping] = widgetMapping
}

// getDefaultTimerange returns a default timerange for the search.
func getDefaultTimerange() map[string]interface{} {
	return map[string]interface{}{
		"type": "relative",
		"from": 300,
	}
}

// validateSearchCreation checks that search was created successfully.
func validateSearchCreation(searchResponse map[string]interface{}) (string, error) {
	id, ok := searchResponse["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("search creation did not return an id")
	}
	return id, nil
}
