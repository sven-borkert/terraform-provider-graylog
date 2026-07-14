package fieldtype

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

// isCustomOrigin reports whether a field type's origin indicates a custom
// mapping managed by this provider. Graylog exposes OVERRIDDEN_INDEX and
// OVERRIDDEN_PROFILE for fields whose type overrides the index/profile default;
// INDEX and PROFILE origins are not custom mappings.
func isCustomOrigin(origin string) bool {
	return strings.HasPrefix(origin, "OVERRIDDEN")
}

// GetFieldType retrieves the custom field type mapping for a specific field in an
// index set. It returns nil when the field has no custom override (origin is not
// OVERRIDDEN_*), so callers can clear state after an out-of-band removal. All
// result pages are scanned so matches beyond the first page are not missed.
func (cl Client) GetFieldType(ctx context.Context, indexSetID, fieldName string) (*IndexSetFieldType, *http.Response, error) {
	if indexSetID == "" {
		return nil, nil, errors.New("index_set_id is required")
	}
	if fieldName == "" {
		return nil, nil, errors.New("field_name is required")
	}

	const perPage = 50
	var lastResp *http.Response
	for page := 1; ; page++ {
		body := struct {
			Elements []IndexSetFieldType `json:"elements"`
		}{}

		query := url.Values{}
		query.Set("query", fieldName)
		query.Set("per_page", strconv.Itoa(perPage))
		query.Set("page", strconv.Itoa(page))

		resp, err := cl.Client.Call(ctx, httpclient.CallParams{
			Method:       "GET",
			Path:         "/system/indices/index_sets/types/" + indexSetID,
			Query:        query,
			ResponseBody: &body,
		})
		lastResp = resp
		if err != nil {
			return nil, resp, err
		}

		for _, ft := range body.Elements {
			if ft.FieldName == fieldName && isCustomOrigin(ft.Origin) {
				ft := ft
				return &ft, resp, nil
			}
		}

		if len(body.Elements) < perPage {
			break
		}
	}

	return nil, lastResp, nil
}
