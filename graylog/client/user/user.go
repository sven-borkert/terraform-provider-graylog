package user

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
		return nil, nil, errors.New("username is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/users/" + name,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) GetByID(
	ctx context.Context, id string,
) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("user id is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/users/id/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Gets(ctx context.Context) (map[string]interface{}, *http.Response, error) {
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/users",
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Create(ctx context.Context, user map[string]interface{}) (*http.Response, error) {
	if user == nil {
		return nil, errors.New("request body is nil")
	}

	// Note: User API does NOT use entity wrapping like other Graylog 7.0 APIs
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "POST",
		Path:        "/users",
		RequestBody: user,
	})
	return resp, err
}

func (cl Client) Update(ctx context.Context, name string, user map[string]interface{}) (*http.Response, error) {
	if name == "" {
		return nil, errors.New("username is required")
	}
	if user == nil {
		return nil, errors.New("request body is nil")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/users/" + name,
		RequestBody: user,
	})
	return resp, err
}

func (cl Client) Delete(ctx context.Context, name string) (*http.Response, error) {
	if name == "" {
		return nil, errors.New("username is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/users/" + name,
	})
	return resp, err
}
