package pipeline

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourcePipelineByID(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "id": "5ea3e4122ab79c001275832c",
  "title": "tf",
  "description": "desc",
  "source": "pipeline \"tf\"\nstage 0 match either\nend\n"
}`

	getRoute := flute.Route{
		Name: "get pipeline by id",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/system/pipelines/pipeline/5ea3e4122ab79c001275832c",
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
		Providers: testutil.SingleDataSourceProviders("graylog_pipeline", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_pipeline" "by_id" {
  pipeline_id = "5ea3e4122ab79c001275832c"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_pipeline.by_id", "pipeline_id", "5ea3e4122ab79c001275832c"),
					resource.TestCheckResourceAttr("data.graylog_pipeline.by_id", "description", "desc"),
				),
			},
		},
	})
}

func TestDataSourcePipelineByTitle(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	listBody := `{
  "pipelines": [
    {
      "id": "5ea3e4122ab79c001275832c",
      "title": "tf",
      "description": "desc",
      "source": "pipeline \"tf\"\nstage 0 match either\nend\n"
    }
  ]
}`

	getRoute := flute.Route{
		Name: "list pipelines",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/system/pipelines/pipeline",
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
		Providers: testutil.SingleDataSourceProviders("graylog_pipeline", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_pipeline" "by_title" {
  title = "tf"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_pipeline.by_title", "pipeline_id", "5ea3e4122ab79c001275832c"),
					resource.TestCheckResourceAttr("data.graylog_pipeline.by_title", "description", "desc"),
				),
			},
		},
	})
}
