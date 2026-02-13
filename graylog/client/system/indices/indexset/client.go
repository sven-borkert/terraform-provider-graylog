package indexset

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

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
		Path:         "/system/indices/index_sets/" + id,
		ResponseBody: &body,
	})
	return body, resp, err
}

type GetAllParams struct {
	Skip  int
	Limit int
	Stats bool
}

func (cl Client) Gets(
	ctx context.Context, params *GetAllParams,
) (map[string]interface{}, *http.Response, error) {
	query := url.Values{}
	if params != nil {
		if params.Skip != 0 {
			query.Add("skip", strconv.Itoa(params.Skip))
		}
		if params.Limit != 0 {
			query.Add("limit", strconv.Itoa(params.Limit))
		}
		if params.Stats {
			query.Add("stats", "true")
		}
	}
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/indices/index_sets",
		Query:        query,
		ResponseBody: &body,
	})
	return body, resp, err
}

func genCreationDate() string {
	return time.Now().In(time.FixedZone("UTC", 0)).Format(time.RFC3339Nano)
}

func (cl Client) Create(
	ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	// Set creation_date if not provided
	if v, ok := data["creation_date"]; !ok || v == "" || v == nil {
		data["creation_date"] = genCreationDate()
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/system/indices/index_sets",
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
		Path:         "/system/indices/index_sets/" + id,
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
		Path:   "/system/indices/index_sets/" + id,
	})
	return resp, err
}

// DeflectorStatus represents the deflector status for an index set
type DeflectorStatus struct {
	IsUp          bool   `json:"is_up"`
	CurrentTarget string `json:"current_target"`
}

// GetDeflectorStatus returns the deflector status for an index set
func (cl Client) GetDeflectorStatus(ctx context.Context, id string) (*DeflectorStatus, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}

	body := &DeflectorStatus{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/deflector/" + id,
		ResponseBody: body,
	})
	return body, resp, err
}

// WaitForDeflector polls the deflector status until it's up or timeout is reached
func (cl Client) WaitForDeflector(ctx context.Context, id string, timeout time.Duration, interval time.Duration) error {
	if id == "" {
		return errors.New("id is required")
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		status, _, err := cl.GetDeflectorStatus(ctx, id)
		if err == nil && status.IsUp && status.CurrentTarget != "" {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
			// continue polling
		}
	}

	return errors.New("timeout waiting for deflector to be ready")
}
