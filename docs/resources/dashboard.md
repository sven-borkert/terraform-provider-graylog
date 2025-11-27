# Resource: graylog_dashboard

Manages Graylog dashboards with widgets for visualizing log data.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/dashboard/resource.go)

## Example Usage

```hcl
resource "graylog_dashboard" "network_overview" {
  title       = "Network Overview"
  description = "Network traffic analysis dashboard"
  summary     = "Dashboard showing network flows"

  state {
    id = "dashboard-state-id"

    # Area chart: Traffic over time
    widgets {
      widget_id = "traffic-over-time"
      type      = "aggregation"
      config = jsonencode({
        row_pivots = [
          {
            type   = "time"
            fields = ["timestamp"]
            config = { interval = { type = "auto", scaling = 1.0 } }
          }
        ]
        column_pivots = []
        series        = [{ function = "count()", config = {} }]
        sort          = []
        rollup        = true
        visualization = "area"
        visualization_config = {
          interpolation = "linear"
          axis_type     = "linear"
        }
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    # Pie chart: Traffic by protocol
    widgets {
      widget_id = "by-protocol"
      type      = "aggregation"
      config = jsonencode({
        row_pivots    = [{ type = "values", fields = ["protocol"], config = { limit = 10 } }]
        column_pivots = []
        series        = [{ function = "count()", config = {} }]
        sort          = []
        rollup        = true
        visualization = "pie"
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    # Table: Top destinations
    widgets {
      widget_id = "top-destinations"
      type      = "aggregation"
      config = jsonencode({
        row_pivots    = [{ type = "values", fields = ["destination"], config = { limit = 10 } }]
        column_pivots = []
        series = [
          { function = "count()", config = {} },
          { function = "sum(bytes)", config = {} }
        ]
        sort          = [{ type = "series", field = "count()", direction = "Descending" }]
        rollup        = true
        visualization = "table"
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    # Bar chart: Top sources
    widgets {
      widget_id = "top-sources"
      type      = "aggregation"
      config = jsonencode({
        row_pivots    = [{ type = "values", fields = ["source_ip"], config = { limit = 10 } }]
        column_pivots = []
        series        = [{ function = "count()", config = {} }]
        sort          = [{ type = "series", field = "count()", direction = "Descending" }]
        rollup        = true
        visualization = "bar"
        visualization_config = {
          axis_type = "linear"
          barmode   = "group"
        }
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    widget_mapping = jsonencode({})

    positions = jsonencode({
      traffic-over-time = { col = 1, row = 1, height = 3, width = "Infinity" }
      by-protocol       = { col = 1, row = 4, height = 3, width = 3 }
      top-destinations  = { col = 4, row = 4, height = 3, width = 3 }
      top-sources       = { col = 1, row = 7, height = 3, width = "Infinity" }
    })

    titles = jsonencode({
      widget = {
        traffic-over-time = "Network Traffic Over Time"
        by-protocol       = "Traffic by Protocol"
        top-destinations  = "Top Destinations"
        top-sources       = "Top Source IPs"
      }
    })
  }
}
```

## Argument Reference

### Top-level Arguments

* `title` - (Required) Dashboard title. The data type is `string`.
* `description` - (Required) Dashboard description. The data type is `string`.
* `summary` - (Optional) Short summary of the dashboard. The data type is `string`.
* `search_id` - (Optional) ID of an existing search. If not provided, the provider automatically creates a search with proper search_types for the dashboard widgets.

### State Block

The `state` block defines the dashboard layout and widgets:

* `id` - (Required) Unique identifier for the state. This is used internally by Graylog.
* `widgets` - (Required) One or more widget blocks defining the dashboard widgets.
* `widget_mapping` - (Optional) JSON string mapping widget IDs to search type IDs. The provider populates this automatically.
* `positions` - (Optional) JSON string defining widget positions. Each widget ID maps to an object with `col`, `row`, `height`, and `width`. Use `"Infinity"` for full-width widgets.
* `titles` - (Optional) JSON string defining custom widget titles. Format: `{"widget": {"widget_id": "Title"}}`.

### Widget Block

Each `widgets` block supports:

* `widget_id` - (Required) Unique identifier for the widget.
* `type` - (Required) Widget type. Use `"aggregation"` for charts and tables.
* `config` - (Required) JSON string containing widget configuration.
* `timerange` - (Optional) JSON string defining the widget's time range. Example: `{"type": "relative", "range": 900}` for 15 minutes.
* `query` - (Optional) JSON string with query configuration.

### Widget Config Options

The `config` JSON object supports:

* `row_pivots` - Array of row groupings. Each has `type` ("time" or "values"), `fields`, and `config`.
* `column_pivots` - Array of column groupings (for pivot tables).
* `series` - Array of metrics. Each has `function` (e.g., "count()", "sum(field)") and `config`.
* `sort` - Array of sort definitions with `type`, `field`, and `direction`.
* `rollup` - Boolean, whether to include rollup totals.
* `visualization` - Chart type: "area", "bar", "line", "pie", "table", "numeric", "heatmap".
* `visualization_config` - Visualization-specific settings (see below).

### Visualization Config by Type

**Area/Line Charts:**
* `interpolation` - Line interpolation: "linear", "step-after", "spline".
* `axis_type` - Y-axis type: "linear", "logarithmic".

**Bar Charts:**
* `axis_type` - Y-axis type: "linear", "logarithmic".
* `barmode` - Bar arrangement: "group", "stack", "relative", "overlay".

**Pie Charts:**
No additional configuration required.

**Tables:**
No additional configuration required.

## Attributes Reference

* `id` - The dashboard ID.
* `created_at` - The date time when the Dashboard was created.
* `owner` - The owner of the dashboard.
* `type` - The view type (always "DASHBOARD").
* `search_id` - The ID of the associated search object.

## How It Works

When creating or updating a dashboard with aggregation widgets, the provider:

1. Generates a search object with `search_types` (pivot queries) for each widget
2. Creates the search via the `/views/search` API
3. Populates the `widget_mapping` to link widget IDs to search_type IDs
4. Creates/updates the dashboard via the `/views` API

This ensures widgets display data immediately without manual configuration.

## Import

`graylog_dashboard` can be imported using the Dashboard ID:

```console
$ terraform import graylog_dashboard.example 5c4acaefc9e77bbbbbbbbbbb
```
