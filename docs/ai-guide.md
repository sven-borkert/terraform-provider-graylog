# Graylog Terraform Provider - AI Agent Guide

Comprehensive reference for configuring Graylog 7.0+ with the `sven-borkert/graylog` Terraform provider. Covers provider setup, all 27 resources, 17 data sources, and common workflows.

## Provider Setup

### Using the Published Provider

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "sven-borkert/graylog"
      version = "~> 3.0"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "admin"
  auth_password    = "password"
}
```

### Provider Arguments

| Argument | Required | Env Var | Default | Description |
|----------|----------|---------|---------|-------------|
| `web_endpoint_uri` | Yes | `GRAYLOG_WEB_ENDPOINT_URI` | - | Graylog API endpoint (must end with `/api`) |
| `auth_name` | Yes | `GRAYLOG_AUTH_NAME` | - | Username or token name |
| `auth_password` | Yes | `GRAYLOG_AUTH_PASSWORD` | - | Password or token value |
| `x_requested_by` | No | `GRAYLOG_X_REQUESTED_BY` | `terraform-provider-graylog` | X-Requested-By header |
| `api_version` | No | `GRAYLOG_API_VERSION` | `v3` | API version |

### Using a Local Build

```bash
make build
cd examples/graylog7-e2e
../../bin/terraform-dev plan   # No 'terraform init' needed
../../bin/terraform-dev apply
```

## JSON String Pattern

Many resource attributes accept JSON strings instead of nested HCL blocks. Use `jsonencode()`:

```hcl
attributes = jsonencode({
  bind_address = "0.0.0.0"
  port         = 5044
})
```

This pattern is used for: `attributes`, `configuration`, `config`, `parameters`, `extractor_config`, `rotation_strategy`, `retention_strategy`, `data_tiering`, `field_spec`, `field_restrictions`, `positions`, `widget_mapping`, `titles`, `timerange`.

The JSON structure varies by type (e.g., input type, output type). Refer to the Graylog REST API browser or existing resource GET responses for the exact format.

## Resource Dependency Order

Create resources in this order to satisfy dependencies:

```
1. graylog_grok_pattern          (standalone)
2. graylog_index_set_template    (standalone, or use data source for built-in)
3. graylog_index_set             (depends on: index_set_template)
4. graylog_input                 (standalone)
5. graylog_input_static_fields   (depends on: input)
6. graylog_stream                (depends on: index_set)
7. graylog_stream_rule           (depends on: stream)
8. graylog_output                (standalone)
9. graylog_stream_output         (depends on: stream, output)
10. graylog_pipeline_rule        (standalone)
11. graylog_pipeline             (standalone, references rules in source)
12. graylog_pipeline_connection  (depends on: stream, pipeline)
13. graylog_event_notification   (standalone)
14. graylog_event_definition     (depends on: event_notification optionally)
15. graylog_role                 (standalone)
16. graylog_user                 (depends on: role optionally)
17. graylog_dashboard            (standalone)
18. graylog_saved_search         (depends on: stream optionally)
19. graylog_sidecar_collector    (standalone)
20. graylog_sidecar_configuration (depends on: sidecar_collector)
21. graylog_sidecars             (depends on: sidecar_collector, sidecar_configuration)
22. graylog_extractor            (depends on: input)
23. graylog_ldap_setting         (standalone)
```

---

## Resources

### graylog_index_set_template

Manages custom index set templates.

```hcl
resource "graylog_index_set_template" "custom" {
  title       = "Custom Template"
  description = "30-day hot retention"
  index_set_config = jsonencode({
    index_analyzer                      = "standard"
    shards                              = 1
    replicas                            = 0
    index_optimization_max_num_segments = 1
    index_optimization_disabled         = false
    field_type_refresh_interval         = 5000
    rotation_strategy_class             = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
    rotation_strategy = {
      type               = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig"
      index_lifetime_min = "P30D"
      index_lifetime_max = "P40D"
    }
    retention_strategy_class = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
    retention_strategy = {
      type                  = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig"
      max_number_of_indices = 20
    }
    data_tiering = {
      type               = "hot_only"
      index_lifetime_min = "P30D"
      index_lifetime_max = "P40D"
    }
  })
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Template title |
| `description` | No | string | Template description |
| `index_set_config` | Yes | JSON string | Index set configuration as JSON |

Import: not supported.

---

### graylog_index_set

Manages index sets for storing log data.

```hcl
data "graylog_index_set_template" "hot7" {
  title = "7 Days Hot"
}

resource "graylog_index_set" "app_logs" {
  title        = "Application Logs"
  description  = "Index set for application logs"
  index_prefix = "applogs"
  rotation_strategy_class = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
  rotation_strategy = jsonencode({
    type               = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  retention_strategy_class = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
  retention_strategy = jsonencode({
    type                  = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig"
    max_number_of_indices = 5
  })
  data_tiering = jsonencode({
    type               = "hot_only"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  index_analyzer                      = "standard"
  index_set_template_id               = data.graylog_index_set_template.hot7.id
  shards                              = 1
  replicas                            = 0
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
  writable                            = true
  use_legacy_rotation                 = false
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Index set title |
| `index_prefix` | Yes (ForceNew) | string | Index prefix, must be unique |
| `rotation_strategy_class` | Yes | string | Rotation strategy class name |
| `rotation_strategy` | Yes | JSON string | Rotation configuration |
| `retention_strategy_class` | Yes | string | Retention strategy class name |
| `retention_strategy` | Yes | JSON string | Retention configuration |
| `index_analyzer` | Yes | string | Elasticsearch analyzer (usually `"standard"`) |
| `shards` | Yes | int | Number of shards per index |
| `index_optimization_max_num_segments` | Yes | int | Max segments after optimization |
| `description` | No | string | Description |
| `replicas` | No | int | Number of replicas |
| `index_optimization_disabled` | No | bool | Disable index optimization |
| `use_legacy_rotation` | No | bool | Use legacy rotation |
| `writable` | No | bool | Whether the index set is writable |
| `default` | No | bool | Whether this is the default index set |
| `field_type_refresh_interval` | No | int | Field type refresh interval (ms) |
| `field_type_profile` | No | string | Field type profile |
| `index_set_template_id` | No | string | Index set template ID |
| `data_tiering` | No | JSON string | Data tiering config. Default: `{"type":"hot_only","index_lifetime_min":"P30D","index_lifetime_max":"P40D"}` |
| `field_restrictions` | No | JSON string | Field restrictions. Default: `{}` |

Computed: `creation_date`, `can_be_default`.

Common rotation strategy classes:
- `org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy` (recommended)
- `org.graylog2.indexer.rotation.strategies.TimeBasedRotationStrategy`
- `org.graylog2.indexer.rotation.strategies.SizeBasedRotationStrategy`
- `org.graylog2.indexer.rotation.strategies.MessageCountRotationStrategy`

Common retention strategy classes:
- `org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy`
- `org.graylog2.indexer.retention.strategies.ClosingRetentionStrategy`
- `org.graylog2.indexer.retention.strategies.NoopRetentionStrategy`

Import: `terraform import graylog_index_set.example <id>`

---

### graylog_input

Manages Graylog inputs (log receivers).

```hcl
resource "graylog_input" "syslog_udp" {
  title  = "Syslog UDP"
  type   = "org.graylog2.inputs.syslog.udp.SyslogUDPInput"
  global = true
  attributes = jsonencode({
    bind_address          = "0.0.0.0"
    port                  = 1514
    recv_buffer_size      = 262144
    number_worker_threads = 4
    override_source       = ""
    force_rdns            = false
    allow_override_date   = true
    store_full_message    = false
    expand_structured_data = false
  })
}

resource "graylog_input" "beats" {
  title  = "Beats"
  type   = "org.graylog.plugins.beats.Beats2Input"
  global = true
  attributes = jsonencode({
    bind_address          = "0.0.0.0"
    port                  = 5044
    recv_buffer_size      = 1048576
    number_worker_threads = 4
    tls_enable            = false
    tcp_keepalive         = false
    no_browser_cache      = true
  })
}

resource "graylog_input" "gelf_tcp" {
  title  = "GELF TCP"
  type   = "org.graylog2.inputs.gelf.tcp.GELFTCPInput"
  global = true
  attributes = jsonencode({
    bind_address          = "0.0.0.0"
    port                  = 12201
    recv_buffer_size      = 1048576
    number_worker_threads = 4
    tls_enable            = false
    tcp_keepalive         = false
    use_null_delimiter    = true
    max_message_size      = 2097152
  })
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Input title |
| `type` | Yes | string | Input type class name |
| `attributes` | Yes | JSON string | Input configuration (varies by type) |
| `global` | No | bool | Whether this is a global input |
| `node` | No | string | Node ID for non-global inputs |

Computed: `created_at`, `creator_user_id`.

Common input types:
- `org.graylog2.inputs.syslog.udp.SyslogUDPInput`
- `org.graylog2.inputs.syslog.tcp.SyslogTCPInput`
- `org.graylog2.inputs.gelf.tcp.GELFTCPInput`
- `org.graylog2.inputs.gelf.udp.GELFUDPInput`
- `org.graylog2.inputs.gelf.http.GELFHttpInput`
- `org.graylog.plugins.beats.Beats2Input`
- `org.graylog2.inputs.raw.tcp.RawTCPInput`
- `org.graylog2.inputs.raw.udp.RawUDPInput`

Import: `terraform import graylog_input.example <id>`

---

### graylog_input_static_fields

Manages static fields added to all messages received by an input.

```hcl
resource "graylog_input_static_fields" "syslog_fields" {
  input_id = graylog_input.syslog_udp.id
  fields = {
    environment = "production"
    source_type = "syslog"
  }
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `input_id` | Yes (ForceNew) | string | Input ID |
| `fields` | No | map(string) | Key-value pairs of static fields |

Import: `terraform import graylog_input_static_fields.example <input_id>`

---

### graylog_stream

Manages streams for routing messages.

```hcl
resource "graylog_stream" "app_stream" {
  title                              = "Application Logs"
  description                        = "All application log messages"
  index_set_id                       = graylog_index_set.app_logs.id
  matching_type                      = "AND"
  remove_matches_from_default_stream = true
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Stream title |
| `index_set_id` | Yes | string | Index set ID to route messages to |
| `description` | No | string | Description |
| `matching_type` | No | string | Rule matching: `"AND"` or `"OR"` |
| `remove_matches_from_default_stream` | No | bool | Remove matched messages from default stream |
| `disabled` | No | bool | Whether the stream is disabled |
| `is_default` | No | bool | Whether this is the default stream |

Computed: `creator_user_id`, `created_at`.

Import: `terraform import graylog_stream.example <id>`

---

### graylog_stream_rule

Manages rules that determine which messages are routed to a stream.

```hcl
resource "graylog_stream_rule" "by_source" {
  stream_id   = graylog_stream.app_stream.id
  field       = "source"
  type        = 1
  value       = "^app-.*"
  description = "Match messages from app servers"
  inverted    = false
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `stream_id` | Yes (ForceNew) | string | Stream ID |
| `field` | Yes | string | Message field to match |
| `type` | Yes | int | Rule type (see below) |
| `value` | No | string | Value to match against |
| `description` | No | string | Description |
| `inverted` | No | bool | Invert the match |

Computed: `rule_id`.

Rule types:
- `1` - match exactly
- `2` - match regular expression
- `3` - greater than
- `4` - smaller than
- `5` - field presence
- `6` - contain
- `7` - always match

Import: `terraform import graylog_stream_rule.example <stream_id>/<rule_id>`

---

### graylog_output

Manages outputs for forwarding messages.

```hcl
resource "graylog_output" "stdout" {
  title = "STDOUT"
  type  = "org.graylog2.outputs.LoggingOutput"
  configuration = jsonencode({
    prefix = "OUTPUT: "
  })
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Output title |
| `type` | Yes | string | Output type class name |
| `configuration` | Yes | JSON string | Output configuration (varies by type) |

Computed: `created_at`, `creator_user_id`, `content_pack`.

Import: `terraform import graylog_output.example <id>`

---

### graylog_stream_output

Associates outputs with a stream.

```hcl
resource "graylog_stream_output" "app_outputs" {
  stream_id  = graylog_stream.app_stream.id
  output_ids = [graylog_output.stdout.id]
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `stream_id` | Yes (ForceNew) | string | Stream ID |
| `output_ids` | Yes | set(string) | Set of output IDs |

Import: `terraform import graylog_stream_output.example <stream_id>`

---

### graylog_grok_pattern

Manages custom grok patterns for message parsing.

```hcl
resource "graylog_grok_pattern" "custom_date" {
  name    = "CUSTOM_DATE"
  pattern = "\\d{4}-\\d{2}-\\d{2}"
}

resource "graylog_grok_pattern" "custom_timestamp" {
  name    = "CUSTOM_TIMESTAMP"
  pattern = "%%{CUSTOM_DATE}[T ]%%{TIME}"
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `name` | Yes | string | Pattern name (conventionally UPPER_CASE) |
| `pattern` | Yes | string | Grok pattern definition |

When referencing other patterns with `%{PATTERN}`, escape as `%%{PATTERN}` in HCL.

Import: `terraform import graylog_grok_pattern.example <id>`

---

### graylog_pipeline_rule

Manages pipeline processing rules.

```hcl
resource "graylog_pipeline_rule" "set_source_type" {
  description = "Set source_type field"
  source      = <<-EOF
    rule "set source type"
    when
      has_field("source")
    then
      set_field("source_type", "syslog");
    end
  EOF
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `source` | Yes | string | Pipeline rule source code |
| `description` | No | string | Description |

The `title` is parsed from the `source` by the API, not set as a separate attribute.

Import: `terraform import graylog_pipeline_rule.example <id>`

---

### graylog_pipeline

Manages processing pipelines.

```hcl
resource "graylog_pipeline" "app_pipeline" {
  description = "Application log processing pipeline"
  source      = <<-EOF
    pipeline "Application Processing"
    stage 0 match either
      rule "set source type"
    end
  EOF
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `source` | Yes | string | Pipeline source code defining stages and rules |
| `description` | No | string | Description |

The `title` is parsed from the `source` by the API, not set as a separate attribute. Rules referenced in the source must exist (created via `graylog_pipeline_rule`).

Import: `terraform import graylog_pipeline.example <id>`

---

### graylog_pipeline_connection

Connects pipelines to streams.

```hcl
resource "graylog_pipeline_connection" "app_connection" {
  stream_id    = graylog_stream.app_stream.id
  pipeline_ids = [graylog_pipeline.app_pipeline.id]
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `stream_id` | Yes (ForceNew) | string | Stream ID |
| `pipeline_ids` | Yes | set(string) | Set of pipeline IDs to connect |

Each stream should have at most one `graylog_pipeline_connection` resource. The resource ID is the stream ID.

Import: `terraform import graylog_pipeline_connection.example <stream_id>`

---

### graylog_event_notification

Manages event notification destinations.

```hcl
resource "graylog_event_notification" "email_alert" {
  title       = "Email Alert"
  description = "Send alert emails"
  config = jsonencode({
    type           = "email-notification-v1"
    sender         = "graylog@example.com"
    subject        = "Graylog Alert: $${event_definition_title}"
    body_template  = "Alert: $${event_definition_description}\n\nTimestamp: $${event.timestamp}"
    email_recipients = ["ops@example.com"]
    user_recipients  = []
    html_body_template = ""
    time_zone        = "UTC"
    lookup_recipient_emails = false
    recipients_lut_name     = ""
    recipients_lut_key      = ""
  })
}

resource "graylog_event_notification" "http_webhook" {
  title       = "Webhook"
  description = "Send alerts via webhook"
  config = jsonencode({
    type        = "http-notification-v1"
    url         = "https://webhook.example.com/graylog"
    api_key     = ""
    api_secret  = ""
    basic_auth  = null
    skip_tls_verification = false
    time_zone   = "UTC"
    method      = "POST"
    content_type = "JSON"
    headers      = ""
    body_template = ""
  })
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Notification title |
| `description` | No | string | Description |
| `config` | Yes | JSON string | Notification configuration (varies by type) |

Common notification types: `email-notification-v1`, `http-notification-v1`, `slack-notification-v1`.

Import: `terraform import graylog_event_notification.example <id>`

---

### graylog_event_definition

Manages event definitions (alerts/correlations).

```hcl
resource "graylog_event_definition" "high_error_rate" {
  title       = "High Error Rate"
  description = "Triggers when error count exceeds threshold"
  priority    = 2
  alert       = true

  config = jsonencode({
    type  = "aggregation-v1"
    query = "level:ERROR"
    query_parameters = []
    streams = [graylog_stream.app_stream.id]
    group_by = []
    search_within_ms  = 300000
    execute_every_ms  = 300000
    series = [{
      id       = "count-"
      type     = "count"
      field    = null
    }]
    conditions = {
      expression = {
        expr = ">"
        left = { expr = "number-ref", ref = "count-" }
        right = { expr = "number", value = 100.0 }
      }
    }
  })

  notification_settings {
    grace_period_ms = 300000
    backlog_size    = 10
  }

  notifications {
    notification_id = graylog_event_notification.email_alert.id
  }

  key_spec = ["source"]
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Event definition title |
| `priority` | Yes | int | Priority: `1` (low), `2` (normal), `3` (high) |
| `config` | Yes | JSON string | Event definition configuration |
| `notification_settings` | Yes | block | Notification settings (see below) |
| `alert` | No | bool | Whether this is an alert |
| `description` | No | string | Description |
| `field_spec` | No | JSON string | Field specification. Default: `{}` |
| `notifications` | No | list(block) | List of notification references |
| `key_spec` | No | set(string) | Fields to group events by |

**notification_settings block:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `grace_period_ms` | No | int | Grace period in milliseconds |
| `backlog_size` | No | int | Number of messages to include |

**notifications block:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `notification_id` | Yes | string | Event notification ID |

Common config types: `aggregation-v1`, `correlation-v1`.

Import: `terraform import graylog_event_definition.example <id>`

---

### graylog_role

Manages custom roles.

```hcl
resource "graylog_role" "operator" {
  name        = "stream-operator"
  description = "Can read and manage streams"
  permissions = [
    "streams:read",
    "streams:edit",
    "streams:changestate",
  ]
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `name` | Yes | string | Role name (used as identifier) |
| `permissions` | Yes | set(string) | Set of permissions |
| `description` | No | string | Description |

Computed: `read_only`.

Built-in roles (`Admin`, `Reader`) have `read_only = true` and cannot be modified.

Import: `terraform import graylog_role.example <role_name>`

---

### graylog_user

Manages Graylog users.

```hcl
resource "graylog_user" "operator" {
  username           = "operator"
  email              = "operator@example.com"
  first_name         = "System"
  last_name          = "Operator"
  password           = "SecurePassword123"
  roles              = ["Reader", graylog_role.operator.name]
  timezone           = "Europe/Berlin"
  session_timeout_ms = 7200000
}

resource "graylog_user" "service" {
  username        = "api-service"
  email           = "api@example.com"
  first_name      = "API"
  last_name       = "Service"
  password        = "ServicePass456"
  roles           = ["Reader"]
  service_account = true
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `username` | Yes (ForceNew) | string | Username, cannot be changed after creation |
| `email` | Yes | string | Email address |
| `first_name` | Yes | string | First name |
| `last_name` | Yes | string | Last name |
| `password` | No (Sensitive) | string | Password. Required on create, optional on update |
| `roles` | No | set(string) | Roles assigned to the user |
| `permissions` | No | set(string) | Direct permissions |
| `timezone` | No | string | Timezone (e.g., `"Europe/Berlin"`) |
| `session_timeout_ms` | No | int | Session timeout in ms. Default: `3600000` (1 hour) |
| `service_account` | No | bool | Whether this is a service account |

Computed: `user_id`, `full_name`, `external`, `read_only`, `session_active`, `last_activity`, `client_address`, `account_status`.

The `password` is write-only and cannot be read back from the API.

Import: `terraform import graylog_user.example <username>`

---

### graylog_dashboard

Manages dashboards with aggregation widgets.

```hcl
resource "graylog_dashboard" "overview" {
  title       = "System Overview"
  description = "System health dashboard"

  state {
    id = "main-state"

    widgets {
      widget_id = "msg-count"
      type      = "aggregation"
      config = jsonencode({
        row_pivots = [{
          type   = "time"
          fields = ["timestamp"]
          config = { interval = { type = "auto", scaling = 1.0 } }
        }]
        column_pivots = []
        series        = [{ function = "count()", config = {} }]
        sort          = []
        rollup        = true
        visualization = "area"
        visualization_config = { interpolation = "linear" }
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    widgets {
      widget_id = "top-sources"
      type      = "aggregation"
      config = jsonencode({
        row_pivots    = [{ type = "values", fields = ["source"], config = { limit = 10 } }]
        column_pivots = []
        series        = [{ function = "count()", config = {} }]
        sort          = [{ type = "series", field = "count()", direction = "Descending" }]
        rollup        = true
        visualization = "bar"
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    positions = jsonencode({
      msg-count   = { col = 1, row = 1, height = 3, width = "Infinity" }
      top-sources = { col = 1, row = 4, height = 3, width = "Infinity" }
    })

    titles = jsonencode({
      widget = {
        msg-count   = "Message Count Over Time"
        top-sources = "Top Sources"
      }
    })

    widget_mapping = jsonencode({})
  }
}
```

**Top-level arguments:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Dashboard title |
| `description` | No | string | Description |
| `summary` | No | string | Short summary |
| `search_id` | No | string | Existing search ID. Auto-generated if omitted |

**state block:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `id` | No | string | State identifier |
| `widgets` | Yes | list(block) | Widget definitions |
| `widget_mapping` | No | JSON string | Widget-to-search-type mapping. Auto-populated |
| `positions` | No | JSON string | Widget positions: `{widget_id: {col, row, height, width}}` |
| `titles` | No | JSON string | Custom titles: `{widget: {widget_id: "Title"}}` |

**widgets block:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `widget_id` | No | string | Unique widget identifier |
| `type` | Yes | string | Widget type (usually `"aggregation"`) |
| `config` | Yes | JSON string | Widget configuration |
| `timerange` | No | JSON string | Time range, e.g., `{"type":"relative","range":900}` |
| `query` | No | block | Query with `type` and `query_string` fields |

Computed: `owner`, `created_at`, `type` (always `"DASHBOARD"`), `search_id`.

Visualization types: `area`, `bar`, `line`, `pie`, `table`, `numeric`, `heatmap`.

The provider auto-generates search objects with `search_types` for each widget.

Import: `terraform import graylog_dashboard.example <id>`

---

### graylog_saved_search

Manages saved searches.

```hcl
resource "graylog_saved_search" "error_search" {
  title           = "Application Errors"
  description     = "Search for application error messages"
  query           = "level:ERROR"
  streams         = [graylog_stream.app_stream.id]
  timerange_type  = "relative"
  timerange_range = 3600
  sort_field      = "timestamp"
  sort_order      = "Descending"
  selected_fields = ["timestamp", "source", "message", "level"]
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `title` | Yes | string | Saved search title |
| `description` | No | string | Description |
| `summary` | No | string | Summary |
| `query` | No | string | Lucene query. Default: `"*"` |
| `streams` | No | set(string) | Stream IDs to search in |
| `timerange_type` | No | string | `"relative"`, `"absolute"`, or `"keyword"`. Default: `"relative"` |
| `timerange_range` | No | int | Seconds to look back (for relative). Default: `300` |
| `selected_fields` | No | list(string) | Fields to display |
| `sort_field` | No | string | Sort field. Default: `"timestamp"` |
| `sort_order` | No | string | `"Ascending"` or `"Descending"`. Default: `"Descending"` |

Computed: `search_id`, `state_id`, `owner`, `created_at`.

Import: `terraform import graylog_saved_search.example <id>`

---

### graylog_extractor

Manages extractors that parse fields from input messages.

```hcl
resource "graylog_extractor" "kv_extractor" {
  input_id        = graylog_input.syslog_udp.id
  title           = "Key-Value Extractor"
  type            = "split_and_index"
  cursor_strategy = "copy"
  source_field    = "message"
  target_field    = "kv_data"
  condition_type  = "regex"
  condition_value = "\\w+=\\w+"
  order           = 0
  extractor_config = jsonencode({
    index         = 1
    split_by      = "="
  })

  converters {
    type   = "numeric"
    config = jsonencode({})
  }
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `input_id` | Yes (ForceNew) | string | Input ID |
| `title` | Yes | string | Extractor title |
| `type` | Yes | string | Extractor type |
| `cursor_strategy` | Yes | string | `"copy"` or `"cut"` |
| `source_field` | Yes | string | Field to extract from |
| `condition_type` | Yes | string | `"none"`, `"string"`, or `"regex"` |
| `extractor_config` | Yes | JSON string | Extractor configuration |
| `order` | No | int | Execution order |
| `condition_value` | No | string | Condition match value |
| `target_field` | No | string | Field to write result to |
| `converters` | No | list(block) | Post-extraction converters |

Computed: `extractor_id`.

**converters block:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `type` | Yes | string | Converter type |
| `config` | Yes | JSON string | Converter configuration |

Extractor types: `regex`, `substring`, `split_and_index`, `copy_input`, `grok`, `json`, `lookup_table`.

Import: `terraform import graylog_extractor.example <input_id>/<extractor_id>`

---

### graylog_alarm_callback

Manages legacy alarm callbacks on streams.

```hcl
resource "graylog_alarm_callback" "email_callback" {
  stream_id = graylog_stream.app_stream.id
  type      = "org.graylog2.alarmcallbacks.EmailAlarmCallback"
  title     = "Email Alarm"
  configuration = jsonencode({
    sender  = "graylog@example.com"
    subject = "Graylog alert for stream: $${stream.title}"
    body    = "Alert condition hit.\n\nStream: $${stream.title}"
    user_receivers  = ["admin"]
    email_receivers = ["ops@example.com"]
  })
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `stream_id` | Yes (ForceNew) | string | Stream ID |
| `type` | Yes | string | Alarm callback type class |
| `title` | Yes | string | Title |
| `configuration` | Yes | JSON string | Configuration (must be a JSON object) |

Computed: `alarmcallback_id`.

Import: `terraform import graylog_alarm_callback.example <stream_id>/<alarm_callback_id>`

---

### graylog_alert_condition

Manages legacy alert conditions on streams.

```hcl
resource "graylog_alert_condition" "message_count" {
  stream_id = graylog_stream.app_stream.id
  type      = "message_count"
  title     = "High Message Count"
  parameters = jsonencode({
    grace             = 5
    threshold         = 1000
    threshold_type    = "MORE"
    backlog           = 5
    time              = 5
    repeat_notifications = false
  })
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `stream_id` | Yes (ForceNew) | string | Stream ID |
| `type` | Yes | string | Condition type |
| `title` | Yes | string | Title |
| `parameters` | Yes | JSON string | Parameters (must be a JSON object) |
| `in_grace` | No | bool | Whether condition is in grace period |

Computed: `alert_condition_id`.

Alert condition types: `message_count`, `field_value`, `field_content_value`.

Import: `terraform import graylog_alert_condition.example <stream_id>/<alert_condition_id>`

---

### graylog_dashboard_widget

Manages individual legacy dashboard widgets (pre-Graylog 7 style).

```hcl
resource "graylog_dashboard_widget" "search_result" {
  dashboard_id = graylog_dashboard.overview.id
  type         = "SEARCH_RESULT_COUNT"
  description  = "Total Messages"
  config = jsonencode({
    timerange = { type = "relative", range = 300 }
    query     = "*"
  })
  cache_time = 10
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `dashboard_id` | Yes (ForceNew) | string | Dashboard ID |
| `type` | Yes | string | Widget type |
| `description` | Yes | string | Widget description |
| `config` | Yes | JSON string | Widget configuration |
| `cache_time` | No | int | Cache time in seconds |

Computed: `widget_id`, `creator_user_id`.

Import: `terraform import graylog_dashboard_widget.example <dashboard_id>/<widget_id>`

---

### graylog_dashboard_widget_positions

Manages widget positions on a legacy dashboard.

```hcl
resource "graylog_dashboard_widget_positions" "positions" {
  dashboard_id = graylog_dashboard.overview.id
  positions = jsonencode({
    widget_id_1 = { col = 1, row = 1, width = 2, height = 2 }
    widget_id_2 = { col = 3, row = 1, width = 2, height = 2 }
  })
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `dashboard_id` | Yes (ForceNew) | string | Dashboard ID |
| `positions` | Yes | JSON string | Map of widget IDs to position objects (must be a JSON object) |

Import: `terraform import graylog_dashboard_widget_positions.example <dashboard_id>`

---

### graylog_sidecar_collector

Manages sidecar collector definitions.

```hcl
resource "graylog_sidecar_collector" "filebeat" {
  name                  = "filebeat"
  service_type          = "exec"
  node_operating_system = "linux"
  executable_path       = "/usr/share/filebeat/bin/filebeat"
  execute_parameters    = "-c %s"
  validation_parameters = "test config -c %s"
  default_template      = ""
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `name` | Yes | string | Collector name |
| `service_type` | Yes | string | `"exec"` or `"svc"` |
| `node_operating_system` | Yes | string | `"linux"`, `"windows"`, or `"darwin"` |
| `executable_path` | Yes | string | Path to collector binary |
| `execute_parameters` | No | string | Execution parameters (`%s` = config path) |
| `validation_parameters` | No | string | Validation parameters |
| `default_template` | No | string | Default configuration template |

Import: `terraform import graylog_sidecar_collector.example <id>`

---

### graylog_sidecar_configuration

Manages sidecar configuration entries.

```hcl
resource "graylog_sidecar_configuration" "filebeat_config" {
  collector_id = graylog_sidecar_collector.filebeat.id
  name         = "filebeat-syslog"
  color        = "#FF0000"
  template     = <<-EOF
    filebeat.inputs:
      - type: log
        paths:
          - /var/log/syslog
    output.logstash:
      hosts: ["graylog.example.com:5044"]
  EOF
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `collector_id` | Yes | string | Collector ID |
| `name` | Yes | string | Configuration name |
| `color` | Yes | string | Color hex code (e.g., `"#FF0000"`) |
| `template` | Yes | string | Configuration template content |

Import: `terraform import graylog_sidecar_configuration.example <id>`

---

### graylog_sidecars

Manages sidecar-to-configuration assignments.

```hcl
resource "graylog_sidecars" "assignments" {
  sidecars {
    node_id = "node-id-from-sidecar-registration"
    assignments {
      collector_id     = graylog_sidecar_collector.filebeat.id
      configuration_id = graylog_sidecar_configuration.filebeat_config.id
    }
  }
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `sidecars` | Yes | set(block) | Sidecar assignment blocks |

**sidecars block:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `node_id` | Yes | string | Sidecar node ID |
| `assignments` | Yes | set(block) | Collector/configuration assignments |

**assignments block:**

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `collector_id` | Yes | string | Collector ID |
| `configuration_id` | Yes | string | Configuration ID |

Import: `terraform import graylog_sidecars.example <any_id>` (resource uses a fixed API endpoint)

---

### graylog_ldap_setting

Manages LDAP/Active Directory authentication settings.

```hcl
resource "graylog_ldap_setting" "main" {
  enabled                = true
  system_username        = "cn=admin,dc=example,dc=com"
  system_password        = "ldap-password"
  ldap_uri               = "ldap://ldap.example.com:389"
  use_start_tls          = true
  trust_all_certificates = false
  active_directory        = false
  search_base            = "ou=users,dc=example,dc=com"
  search_pattern         = "(uid={0})"
  display_name_attribute = "cn"
  default_group          = "Reader"
  group_search_base      = "ou=groups,dc=example,dc=com"
  group_id_attribute     = "cn"
  group_search_pattern   = "(member={0})"
  group_mapping = {
    "graylog-admins" = "Admin"
    "graylog-users"  = "Reader"
  }
}
```

| Argument | Required | Type | Description |
|----------|----------|------|-------------|
| `system_username` | Yes | string | LDAP bind DN |
| `ldap_uri` | Yes | string | LDAP URI |
| `search_base` | Yes | string | User search base DN |
| `search_pattern` | Yes | string | User search pattern (`{0}` = username) |
| `display_name_attribute` | Yes | string | LDAP attribute for display name |
| `default_group` | Yes | string | Default Graylog role for LDAP users |
| `system_password` | No (Sensitive) | string | LDAP bind password |
| `enabled` | No | bool | Enable LDAP auth |
| `use_start_tls` | No | bool | Use STARTTLS |
| `trust_all_certificates` | No | bool | Trust all TLS certs |
| `active_directory` | No | bool | Use Active Directory mode |
| `group_search_base` | No | string | Group search base DN |
| `group_id_attribute` | No | string | Group name attribute |
| `group_search_pattern` | No | string | Group search pattern |
| `group_mapping` | No | map(string) | LDAP group to Graylog role mapping |
| `additional_default_groups` | No | set(string) | Additional default groups |

Computed: `system_password_set`.

Import: `terraform import graylog_ldap_setting.example <any_id>`

---

## Data Sources

### graylog_dashboard

Look up an existing dashboard by ID or title.

```hcl
data "graylog_dashboard" "existing" {
  title = "My Dashboard"
}
# or
data "graylog_dashboard" "by_id" {
  dashboard_id = "abc123"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `dashboard_id` | string | Dashboard ID (conflicts with `title`) |
| `title` | string | Dashboard title (conflicts with `dashboard_id`) |

Attributes: `description`, `summary`, `owner`, `search_id`, `created_at`.

---

### graylog_dashboard_widget

Look up a specific widget within a dashboard.

```hcl
data "graylog_dashboard_widget" "widget" {
  dashboard_id = "abc123"
  widget_id    = "def456"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `dashboard_id` | string (Required) | Dashboard ID |
| `widget_id` | string (Required) | Widget ID |

Attributes: `type`, `description`, `config`, `cache_time`, `creator_user_id`.

---

### graylog_grok_pattern

Look up a grok pattern by name or ID.

```hcl
data "graylog_grok_pattern" "ip" {
  name = "IP"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `name` | string | Pattern name (conflicts with `pattern_id`) |
| `pattern_id` | string | Pattern ID (conflicts with `name`) |

Attributes: `pattern`.

---

### graylog_grok_patterns

List all grok patterns.

```hcl
data "graylog_grok_patterns" "all" {}
```

No arguments. Returns `patterns_json` (JSON string of all patterns).

---

### graylog_index_set

Look up an index set by ID, title, or prefix.

```hcl
data "graylog_index_set" "default" {
  title = "Default index set"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `index_set_id` | string | Index set ID |
| `title` | string | Index set title |
| `index_prefix` | string | Index prefix |

Only one of the three arguments should be specified.

Attributes: `description`, `rotation_strategy_class`, `rotation_strategy`, `retention_strategy_class`, `retention_strategy`, `index_analyzer`, `shards`, `replicas`, `index_optimization_max_num_segments`, `index_optimization_disabled`, `use_legacy_rotation`, `writable`, `default`, `can_be_default`, `field_type_refresh_interval`, `field_type_profile`, `index_template_type`, `creation_date`, `data_tiering`, `field_restrictions`.

---

### graylog_index_set_template

Look up a built-in index set template by title.

```hcl
data "graylog_index_set_template" "hot7" {
  title = "7 Days Hot"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `title` | string (Required) | Template title |

Attributes: `id`, `description`.

---

### graylog_index_set_templates

List all index set templates.

```hcl
data "graylog_index_set_templates" "all" {}
```

No arguments. Returns `templates` list with: `id`, `title`, `description`, `built_in`, `default`.

---

### graylog_input

Look up an input by ID or title.

```hcl
data "graylog_input" "syslog" {
  title = "Syslog UDP"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `input_id` | string | Input ID (conflicts with `title`) |
| `title` | string | Input title (conflicts with `input_id`) |
| `type` | string | Optional type filter |

Attributes: `attributes`, `global`, `node`, `created_at`, `creator_user_id`.

---

### graylog_output

Look up an output by ID or title.

```hcl
data "graylog_output" "stdout" {
  title = "STDOUT"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `output_id` | string | Output ID (conflicts with `title`) |
| `title` | string | Output title (conflicts with `output_id`) |
| `type` | string | Optional type filter |

Attributes: `configuration`, `created_at`, `creator_user_id`, `content_pack`.

---

### graylog_pipeline

Look up a pipeline by ID or title.

```hcl
data "graylog_pipeline" "main" {
  title = "Application Processing"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `pipeline_id` | string | Pipeline ID (conflicts with `title`) |
| `title` | string | Pipeline title (conflicts with `pipeline_id`) |

Attributes: `source`, `description`.

---

### graylog_pipeline_rule

Look up a pipeline rule by ID.

```hcl
data "graylog_pipeline_rule" "rule" {
  rule_id = "abc123"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `rule_id` | string (Required) | Rule ID |

Attributes: `source`, `description`, `title`.

---

### graylog_role

Look up a role by name.

```hcl
data "graylog_role" "admin" {
  name = "Admin"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `name` | string (Required) | Role name |

Attributes: `description`, `permissions`, `read_only`.

---

### graylog_saved_search

Look up a saved search by title.

```hcl
data "graylog_saved_search" "errors" {
  title = "Application Errors"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `title` | string (Required) | Saved search title |

Attributes: `saved_search_id`, `search_id`, `summary`, `description`, `owner`, `created_at`, `state_id`.

---

### graylog_sidecar

Look up a sidecar by node ID or name.

```hcl
data "graylog_sidecar" "server1" {
  node_name = "server1"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `node_id` | string | Node ID (conflicts with `node_name`) |
| `node_name` | string | Node name (conflicts with `node_id`) |

---

### graylog_stream

Look up a stream by title or ID.

```hcl
data "graylog_stream" "default" {
  title = "Default Stream"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `title` | string | Stream title (conflicts with `stream_id`) |
| `stream_id` | string | Stream ID (conflicts with `title`) |

Attributes: `index_set_id`, `description`, `matching_type`, `remove_matches_from_default_stream`, `creator_user_id`, `created_at`, `disabled`, `is_default`.

---

### graylog_stream_rule

Look up a stream rule.

```hcl
data "graylog_stream_rule" "rule" {
  stream_id = graylog_stream.app_stream.id
  rule_id   = "abc123"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `stream_id` | string (Required) | Stream ID |
| `rule_id` | string (Required) | Rule ID |

Attributes: `field`, `type`, `value`, `description`, `inverted`.

---

### graylog_user

Look up a user by ID or username.

```hcl
data "graylog_user" "admin" {
  username = "admin"
}
```

| Argument | Type | Description |
|----------|------|-------------|
| `user_id` | string | User ID (conflicts with `username`) |
| `username` | string | Username (conflicts with `user_id`) |

Attributes: `email`, `full_name`, `first_name`, `last_name`, `timezone`, `session_timeout_ms`, `roles`, `permissions`, `external`, `read_only`, `session_active`, `last_activity`, `client_address`, `account_status`, `service_account`.

---

## Common Patterns

### Complete Logging Pipeline

Set up an index set, stream with rules, pipeline processing, and connect them:

```hcl
# 1. Index set for storage
data "graylog_index_set_template" "hot7" {
  title = "7 Days Hot"
}

resource "graylog_index_set" "app" {
  title        = "Application Logs"
  index_prefix = "applogs"
  rotation_strategy_class = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
  rotation_strategy = jsonencode({
    type               = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  retention_strategy_class = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
  retention_strategy = jsonencode({
    type                  = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig"
    max_number_of_indices = 5
  })
  data_tiering = jsonencode({
    type               = "hot_only"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  index_analyzer                      = "standard"
  index_set_template_id               = data.graylog_index_set_template.hot7.id
  shards                              = 1
  replicas                            = 0
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
  writable                            = true
  use_legacy_rotation                 = false
}

# 2. Stream to route messages
resource "graylog_stream" "app" {
  title                              = "Application Logs"
  index_set_id                       = graylog_index_set.app.id
  matching_type                      = "OR"
  remove_matches_from_default_stream = true
}

# 3. Stream rules for routing
resource "graylog_stream_rule" "by_facility" {
  stream_id = graylog_stream.app.id
  field     = "facility"
  type      = 1
  value     = "local0"
}

# 4. Pipeline rule for enrichment
resource "graylog_pipeline_rule" "enrich_app" {
  description = "Enrich application logs"
  source      = <<-EOF
    rule "enrich app logs"
    when
      has_field("facility") AND to_string($message.facility) == "local0"
    then
      set_field("app_name", "myapp");
      set_field("environment", "production");
    end
  EOF
}

# 5. Pipeline
resource "graylog_pipeline" "app" {
  description = "Application log processing"
  source      = <<-EOF
    pipeline "App Processing"
    stage 0 match either
      rule "enrich app logs"
    end
  EOF
  depends_on = [graylog_pipeline_rule.enrich_app]
}

# 6. Connect pipeline to stream
resource "graylog_pipeline_connection" "app" {
  stream_id    = graylog_stream.app.id
  pipeline_ids = [graylog_pipeline.app.id]
}
```

### Syslog Input Setup

```hcl
resource "graylog_input" "syslog_udp" {
  title  = "Syslog UDP"
  type   = "org.graylog2.inputs.syslog.udp.SyslogUDPInput"
  global = true
  attributes = jsonencode({
    bind_address          = "0.0.0.0"
    port                  = 1514
    recv_buffer_size      = 262144
    number_worker_threads = 4
  })
}

resource "graylog_input" "syslog_tcp" {
  title  = "Syslog TCP"
  type   = "org.graylog2.inputs.syslog.tcp.SyslogTCPInput"
  global = true
  attributes = jsonencode({
    bind_address          = "0.0.0.0"
    port                  = 1514
    recv_buffer_size      = 1048576
    number_worker_threads = 4
    tls_enable            = false
    tcp_keepalive         = false
    use_null_delimiter    = false
    max_message_size      = 2097152
  })
}
```

### Role-Based Access

```hcl
resource "graylog_role" "dashboard_viewer" {
  name        = "dashboard-viewer"
  description = "Can view dashboards"
  permissions = [
    "dashboards:read",
    "searches:read",
  ]
}

resource "graylog_role" "stream_manager" {
  name        = "stream-manager"
  description = "Can manage streams"
  permissions = [
    "streams:read",
    "streams:edit",
    "streams:create",
    "streams:changestate",
  ]
}

resource "graylog_user" "ops" {
  username   = "ops-user"
  email      = "ops@example.com"
  first_name = "Ops"
  last_name  = "User"
  password   = "SecurePass123"
  roles      = [
    "Reader",
    graylog_role.dashboard_viewer.name,
    graylog_role.stream_manager.name,
  ]
}
```

---

## Import Cheat Sheet

| Resource | Import Command |
|----------|---------------|
| `graylog_index_set` | `terraform import graylog_index_set.name <id>` |
| `graylog_input` | `terraform import graylog_input.name <id>` |
| `graylog_input_static_fields` | `terraform import graylog_input_static_fields.name <input_id>` |
| `graylog_stream` | `terraform import graylog_stream.name <id>` |
| `graylog_stream_rule` | `terraform import graylog_stream_rule.name <stream_id>/<rule_id>` |
| `graylog_stream_output` | `terraform import graylog_stream_output.name <stream_id>` |
| `graylog_output` | `terraform import graylog_output.name <id>` |
| `graylog_grok_pattern` | `terraform import graylog_grok_pattern.name <id>` |
| `graylog_pipeline_rule` | `terraform import graylog_pipeline_rule.name <id>` |
| `graylog_pipeline` | `terraform import graylog_pipeline.name <id>` |
| `graylog_pipeline_connection` | `terraform import graylog_pipeline_connection.name <stream_id>` |
| `graylog_event_notification` | `terraform import graylog_event_notification.name <id>` |
| `graylog_event_definition` | `terraform import graylog_event_definition.name <id>` |
| `graylog_role` | `terraform import graylog_role.name <role_name>` |
| `graylog_user` | `terraform import graylog_user.name <username>` |
| `graylog_dashboard` | `terraform import graylog_dashboard.name <id>` |
| `graylog_dashboard_widget` | `terraform import graylog_dashboard_widget.name <dashboard_id>/<widget_id>` |
| `graylog_dashboard_widget_positions` | `terraform import graylog_dashboard_widget_positions.name <dashboard_id>` |
| `graylog_saved_search` | `terraform import graylog_saved_search.name <id>` |
| `graylog_sidecar_collector` | `terraform import graylog_sidecar_collector.name <id>` |
| `graylog_sidecar_configuration` | `terraform import graylog_sidecar_configuration.name <id>` |
| `graylog_sidecars` | `terraform import graylog_sidecars.name <any_id>` |
| `graylog_extractor` | `terraform import graylog_extractor.name <input_id>/<extractor_id>` |
| `graylog_alarm_callback` | `terraform import graylog_alarm_callback.name <stream_id>/<alarm_callback_id>` |
| `graylog_alert_condition` | `terraform import graylog_alert_condition.name <stream_id>/<alert_condition_id>` |
| `graylog_ldap_setting` | `terraform import graylog_ldap_setting.name <any_id>` |
