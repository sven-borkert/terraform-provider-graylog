package role

import (
	"context"
	"errors"
	"net/http"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client httpclient.Client
}

func (cl Client) Get(
	ctx context.Context, name string,
) (map[string]interface{}, *http.Response, error) {
	if name == "" {
		return nil, nil, errors.New("name is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/roles/" + name,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Create(
	ctx context.Context, role interface{},
) (map[string]interface{}, *http.Response, error) {
	if role == nil {
		return nil, nil, errors.New("request body is nil")
	}

	// Wrap entity for Graylog 7.0 CreateEntityRequest structure
	// See: https://go2docs.graylog.org/current/upgrading_graylog/upgrade_to_graylog_7.0.htm
	requestData := map[string]interface{}{
		"entity": role,
		"share_request": map[string]interface{}{
			"selected_grantee_capabilities": map[string]interface{}{},
		},
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/roles",
		RequestBody:  requestData,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Update(
	ctx context.Context, name string, role interface{},
) (map[string]interface{}, *http.Response, error) {
	if name == "" {
		return nil, nil, errors.New("name is required")
	}
	if role == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "PUT",
		Path:         "/roles/" + name,
		RequestBody:  role,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Delete(ctx context.Context, name string) (*http.Response, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/roles/" + name,
	})
	return resp, err
}
