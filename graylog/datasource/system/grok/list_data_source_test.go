package grok

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourceGrokList(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "patterns": [
    {"name": "MYNUMBER", "pattern": "\\\\d+", "id": "abc123"},
    {"name": "IPV4", "pattern": "%{IPV4}", "id": "def456"}
  ]
}`

	route := flute.Route{
		Name: "list grok",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/system/grok",
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
		Providers: testutil.SingleDataSourceProviders("graylog_grok_patterns", DataSourceList()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, route) },
				Config: `
data "graylog_grok_patterns" "all" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.graylog_grok_patterns.all", "patterns_json"),
				),
			},
		},
	})
}
