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

func TestDataSourceGrokByID(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "name": "MYNUMBER",
  "pattern": "\\\\d+",
  "id": "abc123"
}`

	getRoute := flute.Route{
		Name: "get grok by id",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/system/grok/abc123",
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
		Providers: testutil.SingleDataSourceProviders("graylog_grok_pattern", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_grok_pattern" "by_id" {
  pattern_id = "abc123"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_grok_pattern.by_id", "pattern_id", "abc123"),
					resource.TestCheckResourceAttr("data.graylog_grok_pattern.by_id", "name", "MYNUMBER"),
					resource.TestCheckResourceAttr("data.graylog_grok_pattern.by_id", "pattern", "\\d+"),
				),
			},
		},
	})
}

func TestDataSourceGrokByName(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	listBody := `{
  "patterns": [
    {"name": "MYNUMBER", "pattern": "\\\\d+", "id": "abc123"}
  ]
}`

	getRoute := flute.Route{
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
					Body:       ioutil.NopCloser(strings.NewReader(listBody)),
				}, nil
			},
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleDataSourceProviders("graylog_grok_pattern", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_grok_pattern" "by_name" {
  name = "MYNUMBER"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_grok_pattern.by_name", "pattern_id", "abc123"),
					resource.TestCheckResourceAttr("data.graylog_grok_pattern.by_name", "pattern", "\\d+"),
				),
			},
		},
	})
}
