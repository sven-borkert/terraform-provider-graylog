package input

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourceInputByID(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	inputBody := `{
  "id": "5ea252212ab79c001251f682",
  "title": "gelf udp",
  "global": true,
  "type": "org.graylog2.inputs.gelf.udp.GELFUDPInput",
  "node": null,
  "created_at": "2020-04-24T02:42:41.927Z",
  "creator_user_id": "admin",
  "attributes": {
    "recv_buffer_size": 262144,
    "decompress_size_limit": 8388608,
    "bind_address": "0.0.0.0",
    "port": 12201
  }
}`

	getRoute := flute.Route{
		Name: "get input by id",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/system/inputs/5ea252212ab79c001251f682",
		},
		Tester: flute.Tester{
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(inputBody)),
				}, nil
			},
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleDataSourceProviders("graylog_input", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_input" "by_id" {
  input_id = "5ea252212ab79c001251f682"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_input.by_id", "title", "gelf udp"),
					resource.TestCheckResourceAttr("data.graylog_input.by_id", "type", "org.graylog2.inputs.gelf.udp.GELFUDPInput"),
					resource.TestCheckResourceAttr("data.graylog_input.by_id", "global", "true"),
					resource.TestCheckResourceAttr("data.graylog_input.by_id", "input_id", "5ea252212ab79c001251f682"),
				),
			},
		},
	})
}

func TestDataSourceInputByTitle(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	listBody := `{
  "total": 1,
  "inputs": [
    {
      "id": "5ea252212ab79c001251f682",
      "title": "gelf udp",
      "global": true,
      "type": "org.graylog2.inputs.gelf.udp.GELFUDPInput",
      "node": null,
      "created_at": "2020-04-24T02:42:41.927Z",
      "creator_user_id": "admin",
      "configuration": {
        "recv_buffer_size": 262144,
        "decompress_size_limit": 8388608,
        "bind_address": "0.0.0.0",
        "port": 12201
      }
    }
  ]
}`

	listRoute := flute.Route{
		Name: "list inputs",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/system/inputs",
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
		Providers: testutil.SingleDataSourceProviders("graylog_input", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, listRoute) },
				Config: `
data "graylog_input" "by_title" {
  title = "gelf udp"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_input.by_title", "input_id", "5ea252212ab79c001251f682"),
					resource.TestCheckResourceAttr("data.graylog_input.by_title", "type", "org.graylog2.inputs.gelf.udp.GELFUDPInput"),
					resource.TestCheckResourceAttr("data.graylog_input.by_title", "global", "true"),
					resource.TestCheckResourceAttr("data.graylog_input.by_title", "attributes", `{"bind_address":"0.0.0.0","decompress_size_limit":8388608,"port":12201,"recv_buffer_size":262144}`),
				),
			},
		},
	})
}
