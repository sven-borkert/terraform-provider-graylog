package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	// force type = DASHBOARD
	data["type"] = "DASHBOARD"
	// deep copy state to avoid mutating ResourceData during API conversion
	state, err := deepCopyMap(data[keyState].(map[string]interface{}))
	if err != nil {
		return nil, err
	}
	stateID := ""
	if v, ok := state[keyID]; ok {
		stateID, _ = v.(string)
	}
	if stateID == "" {
		if v, ok := data["search_id"]; ok {
			stateID, _ = v.(string)
		}
	}
	if stateID == "" {
		return nil, errors.New("state.id or search_id must be set")
	}
	delete(state, keyID)
	if err := convert.JSONToData(state, keyWidgetMapping, keyPositions, "titles"); err != nil {
		return nil, err
	}
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
	data[keyState] = map[string]interface{}{
		stateID: state,
	}
	delete(data, keyCreatedAt)
	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	log.Printf("dashboard flatten input state type %T", data[keyState])
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
	// log.Printf("dashboard flatten raw widget_mapping type %T value %v", state[keyWidgetMapping], state[keyWidgetMapping])

	widgets := state[keyWidgets].([]interface{})
	for i, a := range widgets {
		widget := a.(map[string]interface{})
		if id, ok := widget["id"]; ok {
			widget[keyWidgetID] = id
		}
		if err := convert.DataToJSON(widget, keyConfig, keyTimerange); err != nil {
			return err
		}
		// ensure timerange/config are strings even if API shape changes
		for _, k := range []string{keyConfig, keyTimerange} {
			if v, ok := widget[k]; ok {
				if _, ok := v.(string); !ok {
					b, err := json.Marshal(v)
					if err != nil {
						return fmt.Errorf("failed to marshal widget %s: %w", k, err)
					}
					widget[k] = string(b)
				}
			}
		}
		for k := range widget {
			switch k {
			case keyWidgetID, "type", keyConfig, keyTimerange, "query":
			default:
				delete(widget, k)
			}
		}
		widgets[i] = widget
	}

	cleanState := map[string]interface{}{
		keyWidgets: widgets,
		keyID:      stateID,
	}
	// Handle titles from API response
	if v, ok := state["titles"]; ok {
		switch mv := v.(type) {
		case map[string]interface{}:
			b, err := json.Marshal(mv)
			if err != nil {
				return fmt.Errorf("failed to marshal titles: %w", err)
			}
			cleanState["titles"] = string(b)
		case string:
			cleanState["titles"] = mv
		default:
			cleanState["titles"] = "{}"
		}
	} else {
		cleanState["titles"] = "{}"
	}
	if v, ok := state[keyWidgetMapping]; ok {
		switch mv := v.(type) {
		case map[string]interface{}:
			b, err := json.Marshal(mv)
			if err != nil {
				return fmt.Errorf("failed to marshal widget_mapping: %w", err)
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
				return fmt.Errorf("failed to marshal positions: %w", err)
			}
			cleanState[keyPositions] = string(b)
		case string:
			cleanState[keyPositions] = mv
		}
	} else {
		cleanState[keyPositions] = "{}"
	}
	// log.Printf("dashboard flatten state: widget_mapping type %T value %v positions type %T value %v", cleanState[keyWidgetMapping], cleanState[keyWidgetMapping], cleanState[keyPositions], cleanState[keyPositions])
	// if len(widgets) > 0 {
	// 	if w, ok := widgets[0].(map[string]interface{}); ok {
	// 		log.Printf("dashboard flatten first widget timerange type %T value %v", w[keyTimerange], w[keyTimerange])
	// 	}
	// }
	data[keyState] = []interface{}{cleanState}

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
