package output

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourceOutputByID(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "id": "abc123",
  "title": "stdout",
  "type": "org.graylog2.outputs.LoggingOutput",
  "configuration": {"prefix": "Writing message: "},
  "creator_user_id": "admin",
  "created_at": "2020-04-24T08:33:08.136Z"
}`

	getRoute := flute.Route{
		Name:    "get output by id",
		Matcher: flute.Matcher{Method: "GET", Path: "/api/system/outputs/abc123"},
		Tester:  flute.Tester{PartOfHeader: testutil.Header()},
		Response: flute.Response{Response: func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(body))}, nil
		}},
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleDataSourceProviders("graylog_output", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_output" "by_id" {
  output_id = "abc123"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_output.by_id", "output_id", "abc123"),
					resource.TestCheckResourceAttr("data.graylog_output.by_id", "title", "stdout"),
					resource.TestCheckResourceAttr("data.graylog_output.by_id", "type", "org.graylog2.outputs.LoggingOutput"),
					resource.TestCheckResourceAttr("data.graylog_output.by_id", "configuration", "{\"prefix\":\"Writing message: \"}"),
				),
			},
		},
	})
}

func TestDataSourceOutputByTitle(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	listBody := `{
  "outputs": [
    {
      "id": "abc123",
      "title": "stdout",
      "type": "org.graylog2.outputs.LoggingOutput",
      "configuration": {"prefix": "Writing message: "}
    },
    {
      "id": "def456",
      "title": "stdout",
      "type": "org.graylog2.outputs.GelfOutput",
      "configuration": {"host": "example"}
    }
  ],
  "total": 2
}`

	getRoute := flute.Route{
		Name:    "list outputs",
		Matcher: flute.Matcher{Method: "GET", Path: "/api/system/outputs"},
		Tester:  flute.Tester{PartOfHeader: testutil.Header()},
		Response: flute.Response{Response: func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(listBody))}, nil
		}},
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleDataSourceProviders("graylog_output", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_output" "by_title" {
  title = "stdout"
  type  = "org.graylog2.outputs.LoggingOutput"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_output.by_title", "output_id", "abc123"),
					resource.TestCheckResourceAttr("data.graylog_output.by_title", "type", "org.graylog2.outputs.LoggingOutput"),
					resource.TestCheckResourceAttr("data.graylog_output.by_title", "configuration", "{\"prefix\":\"Writing message: \"}"),
				),
			},
		},
	})
}
