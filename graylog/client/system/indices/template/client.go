package template

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client httpclient.Client
}

func (cl Client) BuiltIns(ctx context.Context, warmTierEnabled *bool) ([]map[string]interface{}, *http.Response, error) {
	body := []map[string]interface{}{}
	var query url.Values
	if warmTierEnabled != nil {
		query = url.Values{}
		query.Add("warm_tier_enabled", strconv.FormatBool(*warmTierEnabled))
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/indices/index_sets/templates/built-in",
		Query:        query,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Get(ctx context.Context, id string) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/indices/index_sets/templates/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Create(ctx context.Context, data map[string]interface{}) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/system/indices/index_sets/templates",
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Update(ctx context.Context, id string, data map[string]interface{}) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if data == nil {
		return nil, errors.New("request body is nil")
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/system/indices/index_sets/templates/" + id,
		RequestBody: data,
	})
	return resp, err
}

func (cl Client) Delete(ctx context.Context, id string) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/system/indices/index_sets/templates/" + id,
	})
	return resp, err
}
