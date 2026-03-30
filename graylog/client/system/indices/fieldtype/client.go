package fieldtype

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client httpclient.Client
}

// FieldTypeChangeRequest represents a request to change a field type.
type FieldTypeChangeRequest struct {
	Field     string   `json:"field"`
	Type      string   `json:"type"`
	IndexSets []string `json:"index_sets"`
	Rotate    bool     `json:"rotate"`
}

// CustomFieldMappingRemovalRequest represents a request to remove custom field mappings.
type CustomFieldMappingRemovalRequest struct {
	Fields    []string `json:"fields"`
	IndexSets []string `json:"index_sets"`
	Rotate    bool     `json:"rotate"`
}

// IndexSetFieldType represents a field type entry from the API response.
type IndexSetFieldType struct {
	FieldName  string `json:"field_name"`
	Type       string `json:"type"`
	Origin     string `json:"origin"`
	IsReserved bool   `json:"is_reserved"`
}

// ChangeFieldType sets a custom field type mapping for the given field on the given index sets.
func (cl Client) ChangeFieldType(ctx context.Context, req FieldTypeChangeRequest) (*http.Response, error) {
	if req.Field == "" {
		return nil, errors.New("field is required")
	}
	if req.Type == "" {
		return nil, errors.New("type is required")
	}
	if len(req.IndexSets) == 0 {
		return nil, errors.New("at least one index set is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/system/indices/mappings",
		RequestBody: &req,
	})
	return resp, err
}

// RemoveCustomMapping removes a custom field type mapping for the given fields on the given index sets.
func (cl Client) RemoveCustomMapping(ctx context.Context, req CustomFieldMappingRemovalRequest) (*http.Response, error) {
	if len(req.Fields) == 0 {
		return nil, errors.New("at least one field is required")
	}
	if len(req.IndexSets) == 0 {
		return nil, errors.New("at least one index set is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/system/indices/mappings/remove_mapping",
		RequestBody: &req,
	})
	return resp, err
}

// GetFieldType retrieves the field type for a specific field in an index set.
// Returns nil if the field is not found or has no custom override.
func (cl Client) GetFieldType(ctx context.Context, indexSetID, fieldName string) (*IndexSetFieldType, *http.Response, error) {
	if indexSetID == "" {
		return nil, nil, errors.New("index_set_id is required")
	}
	if fieldName == "" {
		return nil, nil, errors.New("field_name is required")
	}

	body := struct {
		Elements []IndexSetFieldType `json:"elements"`
	}{}

	query := url.Values{}
	query.Set("query", fieldName)
	query.Set("per_page", "50")

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/indices/index_sets/types/" + indexSetID,
		Query:        query,
		ResponseBody: &body,
	})
	if err != nil {
		return nil, resp, err
	}

	for _, ft := range body.Elements {
		if ft.FieldName == fieldName {
			return &ft, resp, nil
		}
	}

	return nil, resp, nil
}
