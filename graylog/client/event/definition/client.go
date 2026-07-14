package definition

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"

	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

// scheduleQuery builds the ?schedule=<bool> query that controls whether Graylog
// enables (schedules) the event definition on create/update.
func scheduleQuery(schedule bool) url.Values {
	q := url.Values{}
	q.Set("schedule", strconv.FormatBool(schedule))
	return q
}

type Client struct {
	Client httpclient.Client
}

func (cl Client) Get(
	ctx context.Context, id string,
) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/events/definitions/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Create(
	ctx context.Context, data map[string]interface{}, schedule bool,
) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/events/definitions",
		Query:        scheduleQuery(schedule),
		RequestBody:  util.WrapEntityForCreation(data),
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Update(
	ctx context.Context, id string, data map[string]interface{}, schedule bool,
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
		Path:         "/events/definitions/" + id,
		Query:        scheduleQuery(schedule),
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
		Path:   "/events/definitions/" + id,
	})
	return resp, err
}
