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

func TestDataSourcePipelineRule(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "title": "tf-e2e-rule",
  "description": "desc",
  "source": "rule \"tf-e2e-rule\"\nwhen\n  has_field(\"source\")\nthen\n  set_field(\"tf_e2e_tag\", true);\nend\n",
  "id": "5ea3e60f2ab79c00127585ac"
}`

	getRoute := flute.Route{
		Name: "get pipeline rule",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/system/pipelines/rule/5ea3e60f2ab79c00127585ac",
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
		Providers: testutil.SingleDataSourceProviders("graylog_pipeline_rule", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_pipeline_rule" "test" {
  rule_id = "5ea3e60f2ab79c00127585ac"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_pipeline_rule.test", "rule_id", "5ea3e60f2ab79c00127585ac"),
					resource.TestCheckResourceAttr("data.graylog_pipeline_rule.test", "title", "tf-e2e-rule"),
					resource.TestCheckResourceAttr("data.graylog_pipeline_rule.test", "description", "desc"),
				),
			},
		},
	})
}
