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

// Update changes a user. In Graylog 7 the update endpoint is addressed by the
// user's ID (Mongo ObjectId), not by username.
func (cl Client) Update(ctx context.Context, id string, user map[string]interface{}) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}
	if user == nil {
		return nil, errors.New("request body is nil")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/users/" + id,
		RequestBody: user,
	})
	return resp, err
}

// ChangePassword sets a user's password. In Graylog 7 the main update endpoint
// ignores the password field; passwords are changed via a dedicated endpoint
// addressed by the user's ID.
func (cl Client) ChangePassword(ctx context.Context, id, password string) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/users/" + id + "/password",
		RequestBody: map[string]interface{}{"password": password},
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
