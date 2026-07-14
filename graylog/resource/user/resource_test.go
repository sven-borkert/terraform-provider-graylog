package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/flute/v2/flute"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/testutil"
)

const userID = "5ea23d422ab79c001251dbfa"

func userResponse(lastName string) string {
	return `{
  "id": "` + userID + `",
  "username": "test",
  "email": "test@example.com",
  "first_name": "Test",
  "last_name": "` + lastName + `",
  "full_name": "Test ` + lastName + `",
  "permissions": ["users:edit:test"],
  "preferences": {},
  "timezone": null,
  "session_timeout_ms": 3600000,
  "external": false,
  "startpage": null,
  "roles": ["Reader"],
  "read_only": false,
  "session_active": false,
  "last_activity": null,
  "client_address": null,
  "account_status": "enabled",
  "service_account": false
}`
}

func TestAccUser(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	userBody := ""

	getRoute := flute.Route{
		Name: "get a user",
		Matcher: flute.Matcher{
			Method: "GET",
		},
		Tester: flute.Tester{
			Path:         "/api/users/test",
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(userBody)),
				}, nil
			},
		},
	}

	postRoute := flute.Route{
		Name: "create a user",
		Matcher: flute.Matcher{
			Method: "POST",
		},
		Tester: flute.Tester{
			Path:         "/api/users",
			PartOfHeader: testutil.Header(),
			Test: func(t *testing.T, req *http.Request, svc flute.Service, route flute.Route) {
				body := map[string]interface{}{}
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				// Graylog 7 splits full_name into first_name/last_name; full_name
				// must NOT be sent (the server rejects unknown properties).
				require.Equal(t, "Test", body["first_name"])
				require.Equal(t, "User", body["last_name"])
				require.NotContains(t, body, "full_name")
				require.Equal(t, "password", body["password"])
				userBody = userResponse("User")
			},
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 201,
			},
		},
	}

	createStep := resource.TestStep{
		ResourceName: "graylog_user.test",
		PreConfig: func() {
			testutil.SetHTTPClient(t, getRoute, postRoute)
		},
		Config: `
resource "graylog_user" "test" {
  username   = "test"
  email      = "test@example.com"
  password   = "password"
  first_name = "Test"
  last_name  = "User"
  roles = [
    "Reader",
  ]
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("graylog_user.test", "username", "test"),
			resource.TestCheckResourceAttr("graylog_user.test", "last_name", "User"),
			// The Mongo id returned by the API is captured into user_id.
			resource.TestCheckResourceAttr("graylog_user.test", "user_id", userID),
		),
	}

	// Update by the user's Mongo id (not username).
	updateRoute := flute.Route{
		Name: "update a user",
		Matcher: flute.Matcher{
			Method: "PUT",
			Path:   "/api/users/" + userID,
		},
		Tester: flute.Tester{
			Path:         "/api/users/" + userID,
			PartOfHeader: testutil.Header(),
			Test: func(t *testing.T, req *http.Request, svc flute.Service, route flute.Route) {
				body := map[string]interface{}{}
				if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				// Password is changed via the dedicated endpoint, never the main body.
				require.NotContains(t, body, "password")
				require.Equal(t, "Userupdated", body["last_name"])
				userBody = userResponse("Userupdated")
			},
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 200,
			},
		},
	}

	// Password change goes to the dedicated endpoint, addressed by id.
	passwordRoute := flute.Route{
		Name: "change a user password",
		Matcher: flute.Matcher{
			Method: "PUT",
			Path:   "/api/users/" + userID + "/password",
		},
		Tester: flute.Tester{
			Path:           "/api/users/" + userID + "/password",
			PartOfHeader:   testutil.Header(),
			BodyJSONString: `{"password": "newpassword"}`,
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 204,
			},
		},
	}

	deleteRoute := flute.Route{
		Name: "delete a user",
		Matcher: flute.Matcher{
			Method: "DELETE",
		},
		Tester: flute.Tester{
			Path:         "/api/users/test",
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 204,
			},
		},
	}

	updateStep := resource.TestStep{
		ResourceName: "graylog_user.test",
		PreConfig: func() {
			testutil.SetHTTPClient(t, getRoute, updateRoute, passwordRoute, deleteRoute)
		},
		Config: `
resource "graylog_user" "test" {
  username   = "test"
  email      = "test@example.com"
  password   = "newpassword"
  first_name = "Test"
  last_name  = "Userupdated"
  roles = [
    "Reader",
  ]
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr("graylog_user.test", "last_name", "Userupdated"),
			resource.TestCheckResourceAttr("graylog_user.test", "password", "newpassword"),
		),
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleResourceProviders("graylog_user", Resource()),
		Steps: []resource.TestStep{
			createStep,
			updateStep,
		},
	})
}
