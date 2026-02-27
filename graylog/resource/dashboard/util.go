package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
)

func deepCopyMap(src map[string]interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(src)
	if err != nil {
		return nil, fmt.Errorf("failed to copy map: %w", err)
	}
	var dst map[string]interface{}
	if err := json.Unmarshal(b, &dst); err != nil {
		return nil, fmt.Errorf("failed to copy map: %w", err)
	}
	return dst, nil
}

const (
	keyID            = "id"
	keyCreatedAt     = "created_at"
	keyWidgetMapping = "widget_mapping"
	keyPositions     = "positions"
	keyState         = "state"
	keyWidgets       = "widgets"
	keyConfig        = "config"
	keyTimerange     = "timerange"
	keyWidgetID      = "widget_id"
	keyTabTitle      = "tab_title"
	keyQueryString   = "query_string"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	// force type = DASHBOARD
	data["type"] = "DASHBOARD"

	// State is a list of state blocks (one per tab)
	stateList := data[keyState].([]interface{})
	stateMap := make(map[string]interface{}, len(stateList))

	for _, rawStateItem := range stateList {
		stateItem := rawStateItem.(map[string]interface{})

		// Convert *schema.Set fields in widgets to []interface{} before deep copy
		// (schema.Set contains SchemaSetFunc which cannot be JSON marshaled)
		if rawWidgets, ok := stateItem[keyWidgets].([]interface{}); ok {
			for _, w := range rawWidgets {
				if widget, ok := w.(map[string]interface{}); ok {
					if s, ok := widget["streams"]; ok {
						if set, ok := s.(*schema.Set); ok {
							widget["streams"] = set.List()
						}
					}
				}
			}
		}

		// deep copy state to avoid mutating ResourceData during API conversion
		state, err := deepCopyMap(stateItem)
		if err != nil {
			return nil, err
		}

		// Extract or generate state ID
		stateID := ""
		if v, ok := state[keyID]; ok {
			stateID, _ = v.(string)
		}
		if stateID == "" {
			stateID = uuid.New().String()
		}
		delete(state, keyID)

		// Extract tab_title and merge into titles
		tabTitle, _ := state[keyTabTitle].(string)
		delete(state, keyTabTitle)

		// Keep query_string on state for extraction in create/update
		// (will be removed before sending to API)

		if err := convert.JSONToData(state, keyWidgetMapping, keyPositions, "titles"); err != nil {
			return nil, err
		}

		// Ensure titles is a map and merge tab_title
		titles := ensureTitlesMap(state)
		if tabTitle != "" {
			titles["tab"] = map[string]interface{}{"title": tabTitle}
		}
		state["titles"] = titles

		// Process widgets
		widgets := state[keyWidgets].([]interface{})
		for i, a := range widgets {
			widget := a.(map[string]interface{})
			if wID, ok := widget[keyWidgetID]; ok {
				if s, ok := wID.(string); ok && s != "" {
					widget["id"] = s
				}
				delete(widget, keyWidgetID)
			}
			if err := convert.JSONToData(widget, keyConfig, keyTimerange); err != nil {
				return nil, err
			}
			if q, ok := widget["query"]; ok {
				switch v := q.(type) {
				case []interface{}:
					if len(v) == 0 {
						delete(widget, "query")
					} else {
						// Unwrap single-element list to object for Graylog API
						widget["query"] = v[0]
					}
				case map[string]interface{}:
					if len(v) == 0 {
						delete(widget, "query")
					}
				}
			}
			widgets[i] = widget
		}
		state[keyWidgets] = widgets

		stateMap[stateID] = state
	}

	data[keyState] = stateMap
	delete(data, keyCreatedAt)
	return data, nil
}

// ensureTitlesMap ensures state["titles"] is a map and returns it.
func ensureTitlesMap(state map[string]interface{}) map[string]interface{} {
	v, ok := state["titles"]
	if !ok || v == nil {
		return map[string]interface{}{}
	}
	switch mv := v.(type) {
	case map[string]interface{}:
		return mv
	default:
		return map[string]interface{}{}
	}
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	log.Printf("dashboard flatten input state type %T", data[keyState])
	stateMap := data[keyState].(map[string]interface{})

	if len(stateMap) == 0 {
		return errors.New("dashboard state is empty")
	}

	// Determine output order: match stored state IDs if available, else sort alphabetically
	idOrder := getStateIDOrder(d, stateMap)

	statesList := make([]interface{}, 0, len(stateMap))
	for _, stateID := range idOrder {
		sv, ok := stateMap[stateID]
		if !ok {
			continue
		}
		state := sv.(map[string]interface{})

		cleanState, err := flattenState(stateID, state)
		if err != nil {
			return err
		}
		statesList = append(statesList, cleanState)
	}

	data[keyState] = statesList

	for _, k := range []string{
		"title", "description", "summary", "type", "search_id", "owner", "created_at",
	} {
		if v, ok := data[k]; ok {
			if err := d.Set(k, v); err != nil {
				return fmt.Errorf("failed to set %s: %w", k, err)
			}
		}
	}
	if err := d.Set(keyState, data[keyState]); err != nil {
		return fmt.Errorf("failed to set state: %w", err)
	}

	a, ok := data[keyID]
	if !ok {
		return errors.New("failed to set id. 'id' isn't found")
	}
	dID, ok := a.(string)
	if !ok {
		return fmt.Errorf("'id' should be string: %v", a)
	}

	d.SetId(dID)
	return nil
}

// getStateIDOrder returns state IDs in a stable order.
// It first tries to match the order from the stored Terraform state, then appends any remaining IDs sorted.
func getStateIDOrder(d *schema.ResourceData, stateMap map[string]interface{}) []string {
	// Try to get existing state IDs from stored Terraform state
	var storedIDs []string
	if stored, ok := d.GetOk(keyState); ok {
		if stateList, ok := stored.([]interface{}); ok {
			for _, s := range stateList {
				if sm, ok := s.(map[string]interface{}); ok {
					if id, ok := sm[keyID].(string); ok && id != "" {
						storedIDs = append(storedIDs, id)
					}
				}
			}
		}
	}

	if len(storedIDs) > 0 {
		// Use stored order, append any new IDs sorted
		used := make(map[string]bool, len(storedIDs))
		order := make([]string, 0, len(stateMap))
		for _, id := range storedIDs {
			if _, ok := stateMap[id]; ok {
				order = append(order, id)
				used[id] = true
			}
		}
		// Add any remaining IDs not in stored state
		var remaining []string
		for id := range stateMap {
			if !used[id] {
				remaining = append(remaining, id)
			}
		}
		sort.Strings(remaining)
		order = append(order, remaining...)
		return order
	}

	// No stored state: sort alphabetically for deterministic output
	ids := make([]string, 0, len(stateMap))
	for id := range stateMap {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// flattenState converts a single API state entry to Terraform-compatible format.
func flattenState(stateID string, state map[string]interface{}) (map[string]interface{}, error) {
	widgets := state[keyWidgets].([]interface{})
	for i, a := range widgets {
		widget := a.(map[string]interface{})
		if id, ok := widget["id"]; ok {
			widget[keyWidgetID] = id
		}
		if err := convert.DataToJSON(widget, keyConfig, keyTimerange); err != nil {
			return nil, err
		}
		// ensure timerange/config are strings even if API shape changes
		for _, k := range []string{keyConfig, keyTimerange} {
			if v, ok := widget[k]; ok {
				if _, ok := v.(string); !ok {
					b, err := json.Marshal(v)
					if err != nil {
						return nil, fmt.Errorf("failed to marshal widget %s: %w", k, err)
					}
					widget[k] = string(b)
				}
			}
		}
		// Wrap query object from API into a list for Terraform schema
		if q, ok := widget["query"]; ok {
			switch v := q.(type) {
			case map[string]interface{}:
				widget["query"] = []interface{}{v}
			}
		}
		for k := range widget {
			switch k {
			case keyWidgetID, "type", keyConfig, keyTimerange, "query", "streams":
			default:
				delete(widget, k)
			}
		}
		widgets[i] = widget
	}

	// Sort widgets by widget_id for deterministic ordering (Graylog stores them as a HashSet)
	sort.Slice(widgets, func(i, j int) bool {
		a, _ := widgets[i].(map[string]interface{})[keyWidgetID].(string)
		b, _ := widgets[j].(map[string]interface{})[keyWidgetID].(string)
		return a < b
	})

	cleanState := map[string]interface{}{
		keyWidgets: widgets,
		keyID:      stateID,
	}

	// Handle titles from API response and extract tab_title
	if v, ok := state["titles"]; ok {
		switch mv := v.(type) {
		case map[string]interface{}:
			// Extract tab title if present
			if tab, ok := mv["tab"].(map[string]interface{}); ok {
				if title, ok := tab["title"].(string); ok {
					cleanState[keyTabTitle] = title
				}
				delete(mv, "tab")
			}
			b, err := json.Marshal(mv)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal titles: %w", err)
			}
			cleanState["titles"] = string(b)
		case string:
			if mv == "" {
				cleanState["titles"] = "{}"
			} else {
				cleanState["titles"] = mv
			}
		default:
			cleanState["titles"] = "{}"
		}
	} else {
		cleanState["titles"] = "{}"
	}

	// Handle query_string from search injection
	if qs, ok := state[keyQueryString].(string); ok {
		cleanState[keyQueryString] = qs
	} else {
		cleanState[keyQueryString] = ""
	}

	if v, ok := state[keyWidgetMapping]; ok {
		switch mv := v.(type) {
		case map[string]interface{}:
			b, err := json.Marshal(mv)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal widget_mapping: %w", err)
			}
			cleanState[keyWidgetMapping] = string(b)
		case string:
			cleanState[keyWidgetMapping] = mv
		}
	}
	if v, ok := state[keyPositions]; ok {
		switch mv := v.(type) {
		case map[string]interface{}:
			b, err := json.Marshal(mv)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal positions: %w", err)
			}
			cleanState[keyPositions] = string(b)
		case string:
			cleanState[keyPositions] = mv
		}
	} else {
		cleanState[keyPositions] = "{}"
	}

	return cleanState, nil
}

// injectStreamsFromSearch maps streams from search_types back to widgets via widget_mapping.
func injectStreamsFromSearch(viewData, searchData map[string]interface{}) {
	// Build search_type_id → streams mapping from search queries
	searchTypeStreams := map[string][]interface{}{}
	if queries, ok := searchData["queries"].([]interface{}); ok {
		for _, q := range queries {
			query, ok := q.(map[string]interface{})
			if !ok {
				continue
			}
			searchTypes, ok := query["search_types"].([]interface{})
			if !ok {
				continue
			}
			for _, st := range searchTypes {
				stMap, ok := st.(map[string]interface{})
				if !ok {
					continue
				}
				stID, ok := stMap["id"].(string)
				if !ok || stID == "" {
					continue
				}
				if streams, ok := stMap["streams"].([]interface{}); ok && len(streams) > 0 {
					searchTypeStreams[stID] = streams
				}
			}
		}
	}

	if len(searchTypeStreams) == 0 {
		return
	}

	// Walk through the view state and inject streams into widgets
	stateMap, ok := viewData[keyState].(map[string]interface{})
	if !ok {
		return
	}
	for _, sv := range stateMap {
		state, ok := sv.(map[string]interface{})
		if !ok {
			continue
		}
		widgetMapping, _ := state[keyWidgetMapping].(map[string]interface{})
		widgets, ok := state[keyWidgets].([]interface{})
		if !ok {
			continue
		}
		for _, w := range widgets {
			widget, ok := w.(map[string]interface{})
			if !ok {
				continue
			}
			widgetID, _ := widget["id"].(string)
			if widgetID == "" {
				continue
			}
			// Look up search_type IDs for this widget via widget_mapping
			if stIDs, ok := widgetMapping[widgetID]; ok {
				if idList, ok := stIDs.([]interface{}); ok {
					for _, stID := range idList {
						if id, ok := stID.(string); ok {
							if streams, ok := searchTypeStreams[id]; ok {
								widget["streams"] = streams
								break
							}
						}
					}
				}
			}
		}
	}
}

// injectQueryFromSearch maps per-tab query strings from search queries back to view state entries.
// It reads from the persistent filter field (preferred) or falls back to query.query_string.
func injectQueryFromSearch(viewData, searchData map[string]interface{}) {
	queryStrings := map[string]string{}
	if queries, ok := searchData["queries"].([]interface{}); ok {
		for _, q := range queries {
			query, ok := q.(map[string]interface{})
			if !ok {
				continue
			}
			qID, ok := query["id"].(string)
			if !ok || qID == "" {
				continue
			}
			// Prefer persistent filter over query.query_string
			if filterObj, ok := query["filter"].(map[string]interface{}); ok {
				// Direct query_string filter
				if fq, ok := filterObj["query"].(string); ok && fq != "" {
					queryStrings[qID] = fq
					continue
				}
				// Compound "or" filter wrapping a query_string filter
				if filters, ok := filterObj["filters"].([]interface{}); ok {
					for _, f := range filters {
						if fm, ok := f.(map[string]interface{}); ok {
							if ft, _ := fm["type"].(string); ft == "query_string" {
								if fq, ok := fm["query"].(string); ok && fq != "" {
									queryStrings[qID] = fq
									break
								}
							}
						}
					}
					if _, found := queryStrings[qID]; found {
						continue
					}
				}
			}
			if qObj, ok := query["query"].(map[string]interface{}); ok {
				if qs, ok := qObj["query_string"].(string); ok {
					queryStrings[qID] = qs
				}
			}
		}
	}

	stateMap, ok := viewData[keyState].(map[string]interface{})
	if !ok {
		return
	}
	for stateID, sv := range stateMap {
		state, ok := sv.(map[string]interface{})
		if !ok {
			continue
		}
		if qs, ok := queryStrings[stateID]; ok {
			state[keyQueryString] = qs
		}
	}
}
