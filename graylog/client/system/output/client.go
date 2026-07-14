package output

import (
	"context"
	"errors"
	"net/http"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client httpclient.Client
}

func (cl Client) Get(ctx context.Context, id string) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/outputs/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

// Gets returns all outputs. The GET /system/outputs endpoint takes no
// pagination or stats parameters in Graylog 7.
func (cl Client) Gets(ctx context.Context) (map[string]interface{}, *http.Response, error) {
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/outputs",
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Create(
	ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/system/outputs",
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Update(
	ctx context.Context, id string, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "PUT",
		Path:         "/system/outputs/" + id,
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Delete(ctx context.Context, id string) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/system/outputs/" + id,
	})
	return resp, err
}
