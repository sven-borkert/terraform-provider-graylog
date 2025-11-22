package saved

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

func TestDataSourceSavedSearch(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	body := `{
  "total": 1,
  "elements": [
    {
      "id": "abc123",
      "title": "My Saved Search",
      "summary": "summary",
      "description": "desc",
      "owner": "admin",
      "created_at": "2025-11-22T00:00:00.000Z"
    }
  ]
}`

	listRoute := flute.Route{
		Name: "list saved searches",
		Matcher: flute.Matcher{
			Method: "GET",
			Path:   "/api/search/saved",
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
		Providers: testutil.SingleDataSourceProviders("graylog_saved_search", DataSource()),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testutil.SetHTTPClient(t, listRoute) },
				Config: `
data "graylog_saved_search" "example" {
  title = "My Saved Search"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.graylog_saved_search.example", "saved_search_id", "abc123"),
					resource.TestCheckResourceAttr("data.graylog_saved_search.example", "owner", "admin"),
					resource.TestCheckResourceAttr("data.graylog_saved_search.example", "summary", "summary"),
				),
			},
		},
	})
}
