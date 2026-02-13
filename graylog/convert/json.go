package convert

import (
	"encoding/json"
	"fmt"

	"github.com/suzuki-shunsuke/go-dataeq/dataeq"
)

func OneSizeListToJSON(data map[string]interface{}, keys ...string) error {
	for _, key := range keys {
		raw, ok := data[key]
		if !ok {
			return fmt.Errorf("key '%s' not found in data", key)
		}
		list, ok := raw.([]interface{})
		if !ok || len(list) == 0 {
			return fmt.Errorf("field '%s' must be a non-empty list", key)
		}
		v, ok := list[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("first element of '%s' must be a map", key)
		}
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal attributes '%s' as JSON: %w", key, err)
		}
		data[key] = string(b)
	}
	return nil
}

func DataToJSON(data map[string]interface{}, keys ...string) error {
	if len(keys) == 0 {
		// all keys
		for key, a := range data {
			s, err := json.Marshal(a)
			if err != nil {
				return fmt.Errorf("failed to marshal the '%s' as JSON: %w", key, err)
			}
			data[key] = string(s)
		}
		return nil
	}
	for _, key := range keys {
		s, err := json.Marshal(data[key])
		if err != nil {
			return fmt.Errorf("failed to marshal the '%s' as JSON: %w", key, err)
		}
		data[key] = string(s)
	}
	return nil
}

func JSONToData(data map[string]interface{}, keys ...string) error {
	if len(keys) == 0 {
		// all keys
		for key, v := range data {
			attr, err := dataeq.JSON.ConvertByte([]byte(v.(string)))
			if err != nil {
				return fmt.Errorf("failed to parse the '%s'. '%s' must be a JSON string: %w", key, key, err)
			}
			data[key] = attr
		}
		return nil
	}
	for _, key := range keys {
		v, ok := data[key]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok || s == "" {
			// Skip non-string or empty string values
			continue
		}
		attr, err := dataeq.JSON.ConvertByte([]byte(s))
		if err != nil {
			return fmt.Errorf("failed to parse the '%s'. '%s' must be a JSON string: %w", key, key, err)
		}
		data[key] = attr
	}
	return nil
}

// StringJSONToData parses a raw JSON string into map[string]interface{}
func StringJSONToData(s string) (map[string]interface{}, error) {
	attr, err := dataeq.JSON.ConvertByte([]byte(s))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	m, ok := attr.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("parsed JSON is not a map object")
	}
	return m, nil
}
