package indexset

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
)

const (
	keyID                = "id"
	keyRotationStrategy  = "rotation_strategy"
	keyRetentionStrategy = "retention_strategy"
	keyDefault           = "default"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}
	// Parse rotation/retention JSON fields
	if err := convert.JSONToData(data, keyRotationStrategy, keyRetentionStrategy); err != nil {
		return nil, err
	}

	// When data_tiering is active, rotation/retention fields may be empty.
	// Remove them so they are not sent as empty values to the API.
	for _, key := range []string{
		keyRotationStrategy, "rotation_strategy_class",
		keyRetentionStrategy, "retention_strategy_class",
	} {
		if v, ok := data[key].(string); ok && v == "" {
			delete(data, key)
		}
	}

	// data_tiering: parse JSON string into map
	if s, ok := data["data_tiering"].(string); ok && s != "" {
		m, err := convert.StringJSONToData(s)
		if err != nil {
			return nil, err
		}
		data["data_tiering"] = m
	} else {
		delete(data, "data_tiering")
	}

	// field_restrictions: optional JSON string
	if v, ok := data["field_restrictions"].(string); ok && v != "" {
		m, err := convert.StringJSONToData(v)
		if err != nil {
			return nil, err
		}
		data["field_restrictions"] = m
	} else {
		delete(data, "field_restrictions")
	}

	// Remove computed/unsupported fields from request (id is added back only for Update)
	delete(data, "can_be_default")
	delete(data, "creation_date")
	delete(data, "field_type_profile")
	if v, ok := data["index_template_type"].(string); ok && v == "" {
		data["index_template_type"] = "default"
	}
	if _, ok := data["index_template_type"]; ok {
		delete(data, "index_template_type")
	}
	if v, ok := data["index_set_template_id"].(string); ok && v == "" {
		delete(data, "index_set_template_id")
	}

	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	if err := convert.DataToJSON(data, keyRotationStrategy, keyRetentionStrategy, "data_tiering", "field_restrictions"); err != nil {
		return err
	}

	// When data_tiering is active, Graylog returns null for rotation/retention
	// fields. Remove them so SetResourceData skips them, preserving config values.
	for _, key := range []string{
		keyRotationStrategy, "rotation_strategy_class",
		keyRetentionStrategy, "retention_strategy_class",
	} {
		v, ok := data[key]
		if !ok {
			continue
		}
		if v == nil {
			delete(data, key)
		} else if s, isStr := v.(string); isStr && s == "null" {
			delete(data, key)
		}
	}

	// field_restrictions: normalize null to empty object to match schema default
	if v, ok := data["field_restrictions"]; ok {
		if v == nil {
			data["field_restrictions"] = "{}"
		} else if s, isStr := v.(string); isStr && s == "null" {
			data["field_restrictions"] = "{}"
		}
	}

	// index_set_template_id: built-in template IDs regenerate on restart.
	// Remove so state retains the data source reference from config.
	delete(data, "index_set_template_id")

	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}

	d.SetId(data[keyID].(string))
	return nil
}
