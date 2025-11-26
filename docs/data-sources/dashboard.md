# graylog_dashboard Data Source

Use this data source to retrieve information about an existing Graylog dashboard.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/datasource/dashboard)

## Example Usage

### Lookup by ID

```tf
data "graylog_dashboard" "by_id" {
  dashboard_id = "6926e4342562186bc3ea1a26"
}

output "dashboard_title" {
  value = data.graylog_dashboard.by_id.title
}
```

### Lookup by Title

```tf
data "graylog_dashboard" "by_title" {
  title = "My Dashboard"
}

output "dashboard_owner" {
  value = data.graylog_dashboard.by_title.owner
}
```

### Reference from a Dashboard Resource

```tf
resource "graylog_dashboard" "example" {
  title       = "Example Dashboard"
  description = "Created by Terraform"
  # ... other configuration
}

data "graylog_dashboard" "example" {
  dashboard_id = graylog_dashboard.example.id
}

output "dashboard_search_id" {
  value = data.graylog_dashboard.example.search_id
}
```

## Argument Reference

One of `dashboard_id` or `title` must be set. Both cannot be set at the same time.

* `dashboard_id` - (Optional) The ID of the dashboard.
* `title` - (Optional) The title of the dashboard. Must be unique.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the dashboard (same as `dashboard_id`).
* `title` - The title of the dashboard.
* `description` - The description of the dashboard.
* `summary` - The summary of the dashboard.
* `owner` - The owner (username) of the dashboard.
* `search_id` - The ID of the search associated with the dashboard.
* `created_at` - The timestamp when the dashboard was created.
