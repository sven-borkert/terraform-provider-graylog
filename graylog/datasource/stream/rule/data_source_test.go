package rule

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourceStreamRule(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "field": "source",
  "stream_id": "5ea26bb42ab79c0012521287",
  "description": "test",
  "id": "5ea26bb42ab79c0012521299",
  "type": 1,
  "inverted": false,
  "value": "foo"
}`

	getRoute := flute.Route{
		Name: "get stream rule",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/streams/5ea26bb42ab79c0012521287/rules/5ea26bb42ab79c0012521299",
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
		Providers: testutil.SingleDataSourceProviders("graylog_stream_rule", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_stream_rule" "test" {
  stream_id = "5ea26bb42ab79c0012521287"
  rule_id   = "5ea26bb42ab79c0012521299"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_stream_rule.test", "stream_id", "5ea26bb42ab79c0012521287"),
					resource.TestCheckResourceAttr("data.graylog_stream_rule.test", "rule_id", "5ea26bb42ab79c0012521299"),
					resource.TestCheckResourceAttr("data.graylog_stream_rule.test", "field", "source"),
					resource.TestCheckResourceAttr("data.graylog_stream_rule.test", "value", "foo"),
					resource.TestCheckResourceAttr("data.graylog_stream_rule.test", "type", "1"),
					resource.TestCheckResourceAttr("data.graylog_stream_rule.test", "inverted", "false"),
				),
			},
		},
	})
}
