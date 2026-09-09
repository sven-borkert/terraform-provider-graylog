package input

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/suzuki-shunsuke/go-dataeq/dataeq"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
)

const (
	keyID            = "id"
	keyAttributes    = "attributes"
	keyCreatedAt     = "created_at"
	keyCreatorUserID = "creator_user_id"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}

	attrS := d.Get(keyAttributes).(string)
	attr, err := dataeq.JSON.ConvertByte([]byte(attrS))
	if err != nil {
		return nil, fmt.Errorf("failed to parse the 'attributes'. 'attributes' must be a JSON string: %w", err)
	}
	data[keyAttributes] = attr

	delete(data, keyCreatedAt)
	delete(data, keyCreatorUserID)

	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	attrVal, ok := data[keyAttributes]
	if !ok || attrVal == nil {
		if cfg, ok := data["configuration"]; ok && cfg != nil {
			attrVal = cfg
		}
	}
	if attrVal == nil {
		attrVal = map[string]interface{}{}
	}
	// Graylog masks password attributes on read. Preserve the last applied
	// value so an unchanged secret stays stable and a new value still plans
	// an update. Imports have no prior secret and retain the server mask.
	var prior map[string]interface{}
	if err := json.Unmarshal([]byte(d.Get(keyAttributes).(string)), &prior); err == nil {
		if attrs, ok := attrVal.(map[string]interface{}); ok {
			for key, value := range attrs {
				if value == "<password set>" {
					if secret, exists := prior[key]; exists {
						attrs[key] = secret
					}
				}
			}
		}
	}
	data[keyAttributes] = attrVal
	delete(data, "configuration")

	attrS, err := json.Marshal(attrVal)
	if err != nil {
		return fmt.Errorf("failed to marshal the 'attributes' as JSON: %w", err)
	}
	data[keyAttributes] = string(attrS)

	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}

	d.SetId(data[keyID].(string))
	return nil
}
