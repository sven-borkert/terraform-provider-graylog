package saved

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client httpclient.Client
}

// Gets lists saved searches (View summaries), accumulating every page. The
// endpoint is paginated (default 50 per page), so a single request would miss
// saved searches beyond the first page. The result is normalized to {"elements": [...]}.
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
			Path:         "/search/saved",
			Query:        query,
			ResponseBody: &body,
		})
		lastResp = resp
		if err != nil {
			return nil, resp, err
		}
		pageElems, _ := body["elements"].([]interface{})
		elements = append(elements, pageElems...)
		if len(pageElems) < perPage {
			break
		}
	}
	return map[string]interface{}{"elements": elements}, lastResp, nil
}
