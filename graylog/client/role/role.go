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

func (cl Client) Gets(
	ctx context.Context,
) ([]map[string]interface{}, *http.Response, error) {
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/roles",
		ResponseBody: &body,
	})
	if err != nil {
		return nil, resp, err
	}

	// Response contains "roles" array
	rolesRaw, ok := body["roles"]
	if !ok {
		return nil, resp, errors.New("response does not contain 'roles' key")
	}

	rolesArray, ok := rolesRaw.([]interface{})
	if !ok {
		return nil, resp, errors.New("roles is not an array")
	}

	roles := make([]map[string]interface{}, len(rolesArray))
	for i, r := range rolesArray {
		roleMap, ok := r.(map[string]interface{})
		if !ok {
			return nil, resp, errors.New("role is not a map")
		}
		roles[i] = roleMap
	}

	return roles, resp, nil
}

func (cl Client) Create(
	ctx context.Context, role interface{},
) (map[string]interface{}, *http.Response, error) {
	if role == nil {
		return nil, nil, errors.New("request body is nil")
	}

	// Note: Role API does NOT use entity wrapping - takes RoleResponse directly
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/roles",
		RequestBody:  role,
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
