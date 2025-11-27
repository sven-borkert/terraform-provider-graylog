package search

import (
	"context"
	"errors"
	"net/http"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

// Client provides methods to interact with the Graylog Views Search API.
type Client struct {
	Client httpclient.Client
}

// Get retrieves a search by ID.
func (cl Client) Get(ctx context.Context, id string) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/views/search/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

// Create creates a new search object with the given queries and search types.
func (cl Client) Create(
	ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/views/search",
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

// Delete deletes a search by ID.
func (cl Client) Delete(ctx context.Context, id string) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/views/search/" + id,
	})
	return resp, err
}
