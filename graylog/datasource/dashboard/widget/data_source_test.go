package widget

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourceDashboardWidget(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "id": "wid123",
  "type": "search_result_chart",
  "description": "desc",
  "cache_time": 10,
  "creator_user_id": "admin",
  "config": {
    "interval": "minute",
    "query": "",
    "timerange": {
      "type": "relative",
      "range": 300
    }
  }
}`

	route := flute.Route{
		Name: "get widget",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/dashboards/dash123/widgets/wid123",
		},
		Tester: flute.Tester{
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(body)),
				}, nil
			},
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleDataSourceProviders("graylog_dashboard_widget", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, route) },
				Config: `
data "graylog_dashboard_widget" "w" {
  dashboard_id = "dash123"
  widget_id    = "wid123"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_dashboard_widget.w", "type", "search_result_chart"),
					resource.TestCheckResourceAttr("data.graylog_dashboard_widget.w", "dashboard_id", "dash123"),
					resource.TestCheckResourceAttr("data.graylog_dashboard_widget.w", "widget_id", "wid123"),
				),
			},
		},
	})
}
