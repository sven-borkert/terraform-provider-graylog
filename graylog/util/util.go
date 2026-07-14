package util

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/suzuki-shunsuke/go-dataeq/dataeq"
)

func HandleGetResourceError(
	d *schema.ResourceData, resp *http.Response, err error, codes ...int,
) error {
	if resp == nil {
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	for _, code := range codes {
		if resp.StatusCode == code {
			d.SetId("")
			return nil
		}
	}
	return err
}

var (
	readAfterCreateAttempts = 10
	readAfterCreateDelay    = 500 * time.Millisecond
)

// ReadAfterCreate retries a read call to absorb eventual consistency immediately
// after resource creation. readFunc is expected to clear ID when the resource is
// not found.
func ReadAfterCreate(
	d *schema.ResourceData,
	m interface{},
	id string,
	readFunc func(*schema.ResourceData, interface{}) error,
) error {
	if readFunc == nil {
		return errors.New("read function is required")
	}
	if id == "" {
		return errors.New("id is required")
	}

	for i := 0; i < readAfterCreateAttempts; i++ {
		d.SetId(id)
		if err := readFunc(d, m); err != nil {
			return err
		}
		if d.Id() != "" {
			return nil
		}
		if i < readAfterCreateAttempts-1 {
			time.Sleep(readAfterCreateDelay)
		}
	}

	return fmt.Errorf("resource %s not found after create", id)
}

func SchemaDiffSuppressJSONString(k, oldV, newV string, d *schema.ResourceData) bool {
	b, err := dataeq.JSON.Equal([]byte(oldV), []byte(newV))
	if err != nil {
		return false
	}
	return b
}

// SchemaDiffSuppressJSONSubset suppresses diffs when the planned (new) JSON value
// is a subset of the state (old) JSON value. This handles cases where the API
// enriches the config with additional default fields not specified by the user.
func SchemaDiffSuppressJSONSubset(k, oldV, newV string, d *schema.ResourceData) bool {
	// First try exact equality
	b, err := dataeq.JSON.Equal([]byte(oldV), []byte(newV))
	if err == nil && b {
		return true
	}

	var oldData, newData interface{}
	if err := json.Unmarshal([]byte(oldV), &oldData); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newV), &newData); err != nil {
		return false
	}

	return jsonIsSubset(newData, oldData)
}

// graylogPasswordMask is the placeholder Graylog returns in place of password
// input attributes (AbstractInputsResource.maskPasswordsInConfiguration).
const graylogPasswordMask = "<password set>"

// SchemaDiffSuppressJSONMaskedSecret suppresses diffs on a JSON attributes blob
// when the only differences are fields the server masked as "<password set>".
// Inputs with secret configuration otherwise show a perpetual diff because the
// read-back value never matches the configured secret. Because the real secret is
// never returned, a change to only a masked field cannot be detected here; recreate
// the resource (or change another attribute) to rotate such secrets.
func SchemaDiffSuppressJSONMaskedSecret(k, oldV, newV string, d *schema.ResourceData) bool {
	if b, err := dataeq.JSON.Equal([]byte(oldV), []byte(newV)); err == nil && b {
		return true
	}

	var oldData, newData map[string]interface{}
	if err := json.Unmarshal([]byte(oldV), &oldData); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(newV), &newData); err != nil {
		return false
	}
	if len(oldData) != len(newData) {
		return false
	}

	for key, nv := range newData {
		ov, ok := oldData[key]
		if !ok {
			return false
		}
		if s, isStr := ov.(string); isStr && s == graylogPasswordMask {
			// Masked secret in state; cannot compare, treat as unchanged.
			continue
		}
		ob, _ := json.Marshal(ov)
		nb, _ := json.Marshal(nv)
		if string(ob) != string(nb) {
			return false
		}
	}
	return true
}

// jsonIsSubset checks if 'subset' is contained within 'superset'.
// For maps: all keys in subset must exist in superset with matching values.
// For slices: must have same length with each element matching.
// For scalars: must be equal.
func jsonIsSubset(subset, superset interface{}) bool {
	switch s := subset.(type) {
	case map[string]interface{}:
		sup, ok := superset.(map[string]interface{})
		if !ok {
			return false
		}
		for k, v := range s {
			sv, ok := sup[k]
			if !ok {
				return false
			}
			if !jsonIsSubset(v, sv) {
				return false
			}
		}
		return true
	case []interface{}:
		sup, ok := superset.([]interface{})
		if !ok {
			return false
		}
		if len(s) != len(sup) {
			return false
		}
		for i := range s {
			if !jsonIsSubset(s[i], sup[i]) {
				return false
			}
		}
		return true
	default:
		// Compare scalars (string, float64, bool, nil)
		return fmt.Sprintf("%v", subset) == fmt.Sprintf("%v", superset)
	}
}

func GenStateFunc(keys ...string) schema.StateContextFunc {
	return func(_ context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		a := strings.Split(d.Id(), "/")
		size := len(keys)
		if len(a) != size {
			return nil, errors.New("format of import argument should be " + strings.Join(keys, "/"))
		}
		for i, k := range keys {
			if err := d.Set(k, a[i]); err != nil {
				return nil, err
			}
		}
		return []*schema.ResourceData{d}, nil
	}
}

func RenameKey(data map[string]interface{}, oldKey, newKey string) (interface{}, bool) {
	v, ok := data[oldKey]
	if !ok {
		return nil, false
	}
	delete(data, oldKey)
	data[newKey] = v
	return v, true
}

var ValidateIsJSON = WrapValidateFunc(func(value interface{}, key string) error {
	var a interface{}
	if err := json.Unmarshal([]byte(value.(string)), &a); err != nil {
		return fmt.Errorf("the value of the field '%s' must be JSON string: %w", key, err)
	}
	return nil
})

var ValidateIsMapJSON = WrapValidateFunc(func(value interface{}, key string) error {
	var a interface{}
	if err := json.Unmarshal([]byte(value.(string)), &a); err != nil {
		return fmt.Errorf("the value of the field '%s' must be JSON string: %w", key, err)
	}
	if _, ok := a.(map[string]interface{}); !ok {
		return errors.New("the value of the field '" + key + "' must be JSON string of map")
	}
	return nil
})

func WrapValidateFunc(f func(v interface{}, k string) error) schema.SchemaValidateFunc { //nolint:staticcheck
	return func(v interface{}, k string) (s []string, es []error) {
		if err := f(v, k); err != nil {
			es = append(es, err)
		}
		return
	}
}

func SetDefaultValue(data map[string]interface{}, key string, value interface{}) {
	if _, ok := data[key]; !ok {
		data[key] = value
	}
}

// WrapEntityForCreation wraps entity data in CreateEntityRequest structure
// required by Graylog 7.0+ for entity creation endpoints (streams, dashboards,
// event definitions, event notifications, etc.)
//
// Example transformation:
//
//	Input:  {"title": "My Stream", "index_set_id": "123"}
//	Output: {"entity": {"title": "My Stream", "index_set_id": "123"}, "share_request": {"selected_grantee_capabilities": {}}}
func WrapEntityForCreation(entityData map[string]interface{}) map[string]interface{} {
	if entityData == nil {
		return map[string]interface{}{
			"entity": map[string]interface{}{},
			"share_request": map[string]interface{}{
				"selected_grantee_capabilities": map[string]interface{}{},
			},
		}
	}
	return map[string]interface{}{
		"entity": entityData,
		"share_request": map[string]interface{}{
			"selected_grantee_capabilities": map[string]interface{}{},
		},
	}
}

// RemoveComputedFields removes read-only/computed fields from data that should not
// be sent in update requests. Graylog 7.0+ rejects unknown/read-only properties.
//
// Note: Many Graylog 7 PUT endpoints require the id field in the request body.
// Callers should set data["id"] = d.Id() AFTER calling this function.
//
// Computed fields removed:
//   - id: Resource identifier (must be re-added after this call for most endpoints)
//   - created_at: Creation timestamp (server-generated)
//   - creator_user_id: Creator identifier (server-generated)
//   - last_modified: Last modification timestamp (server-generated)
func RemoveComputedFields(data map[string]interface{}) {
	delete(data, "id")
	delete(data, "created_at")
	delete(data, "creator_user_id")
	delete(data, "last_modified")
}

// ComputeSHA256 computes the SHA256 hash of a string and returns it as a hex string.
// This is used to create content hashes for cache invalidation workarounds.
func ComputeSHA256(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}
