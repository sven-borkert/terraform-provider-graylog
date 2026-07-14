package dashboard

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

func (cl Client) Get(ctx context.Context, id string) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/dashboards/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

// Gets lists all dashboards (View summaries), accumulating every page. The
// endpoint is paginated (default 50 per page), so a single request would miss
// dashboards beyond the first page. The result is normalized to {"elements": [...]}.
func (cl Client) Gets(ctx context.Context) (map[string]interface{}, *http.Response, error) {
	const perPage = 100
	elements := []interface{}{}
	var lastResp *http.Response
	for page := 1; ; page++ {
		body := map[string]interface{}{}
		query := url.Values{}
		query.Set("page", strconv.Itoa(page))
		query.Set("per_page", strconv.Itoa(perPage))
		resp, err := cl.Client.Call(ctx, httpclient.CallParams{
			Method:       "GET",
			Path:         "/dashboards",
			Query:        query,
			ResponseBody: &body,
		})
		lastResp = resp
		if err != nil {
			return nil, resp, err
		}
		// Graylog 7 returns "elements"; fall back to legacy keys for compatibility.
		pageElems, _ := body["elements"].([]interface{})
		if pageElems == nil {
			if pageElems, _ = body["dashboards"].([]interface{}); pageElems == nil {
				pageElems, _ = body["views"].([]interface{})
			}
		}
		elements = append(elements, pageElems...)
		if len(pageElems) < perPage {
			break
		}
	}
	return map[string]interface{}{"elements": elements}, lastResp, nil
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
		Path:         "/dashboards",
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Update(
	ctx context.Context, id string, data map[string]interface{},
) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	if data == nil {
		return nil, errors.New("request body is nil")
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/dashboards/" + id,
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
		Path:   "/dashboards/" + id,
	})
	return resp, err
}
