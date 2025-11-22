package stream

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourceStreamByID(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	streamBody := `{
  "id": "5ea26bb42ab79c0012521287",
  "creator_user_id": "admin",
  "outputs": [],
  "matching_type": "AND",
  "description": "test",
  "created_at": "2020-04-24T04:31:48.481Z",
  "disabled": true,
  "rules": [],
  "title": "test",
  "remove_matches_from_default_stream": false,
  "index_set_id": "5e9861442ab79c0012e7d1c4",
  "is_default": false
}`

	getRoute := flute.Route{
		Name: "get stream by id",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/streams/5ea26bb42ab79c0012521287",
		},
		Tester: flute.Tester{
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(streamBody)),
				}, nil
			},
		},
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleDataSourceProviders("graylog_stream", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getRoute) },
				Config: `
data "graylog_stream" "by_id" {
  stream_id = "5ea26bb42ab79c0012521287"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_stream.by_id", "stream_id", "5ea26bb42ab79c0012521287"),
					resource.TestCheckResourceAttr("data.graylog_stream.by_id", "title", "test"),
					resource.TestCheckResourceAttr("data.graylog_stream.by_id", "matching_type", "AND"),
					resource.TestCheckResourceAttr("data.graylog_stream.by_id", "index_set_id", "5e9861442ab79c0012e7d1c4"),
					resource.TestCheckResourceAttr("data.graylog_stream.by_id", "disabled", "true"),
				),
			},
		},
	})
}

func TestDataSourceStreamByTitle(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	listBody := `{
  "total": 1,
  "streams": [
    {
      "id": "5ea26bb42ab79c0012521287",
      "creator_user_id": "admin",
      "outputs": [],
      "matching_type": "AND",
      "description": "test",
      "created_at": "2020-04-24T04:31:48.481Z",
      "disabled": false,
      "rules": [],
      "title": "test",
      "remove_matches_from_default_stream": false,
      "index_set_id": "5e9861442ab79c0012e7d1c4",
      "is_default": false
    }
  ]
}`

	getList := flute.Route{
		Name: "list streams",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/streams",
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
		Providers: testutil.SingleDataSourceProviders("graylog_stream", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, getList) },
				Config: `
data "graylog_stream" "by_title" {
  title = "test"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_stream.by_title", "stream_id", "5ea26bb42ab79c0012521287"),
					resource.TestCheckResourceAttr("data.graylog_stream.by_title", "disabled", "false"),
					resource.TestCheckResourceAttr("data.graylog_stream.by_title", "matching_type", "AND"),
				),
			},
		},
	})
}
