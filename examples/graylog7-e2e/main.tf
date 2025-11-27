provider "graylog" {
  web_endpoint_uri = var.web_endpoint_uri
  auth_name        = var.auth_name
  auth_password    = var.auth_password
  x_requested_by   = var.x_requested_by
  api_version      = var.api_version
}

data "graylog_index_set_template" "hot30" {
  title = "30 Days Hot"
}

data "graylog_index_set_templates" "all_builtin" {}

locals {
  # Minimal rotation/retention configs for testing
  rotation_strategy = {
    type               = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig"
    index_lifetime_min = "P30D"
    index_lifetime_max = "P40D"
  }

  retention_strategy = {
    type                  = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig"
    max_number_of_indices = 5
  }

  gelf_udp_config = {
    bind_address          = "0.0.0.0"
    port                  = 12210
    recv_buffer_size      = 262144
    decompress_size_limit = 8388608
  }

  stream_rule_value = "tf-e2e-source"
  pipeline_suffix   = "tf-e2e-v2"

  pipeline_rule_source = <<-EOT
rule "${local.pipeline_suffix}-rule"
when
  has_field("source")
then
  set_field("tf_e2e_tag", true);
end
EOT

  pipeline_source = <<-EOT
pipeline "${local.pipeline_suffix}-pipeline"
stage 0 match either
  rule "${local.pipeline_suffix}-rule"
end
EOT
}

resource "graylog_index_set" "tf_e2e" {
  title                               = "tf-e2e-index-set"
  description                         = "Terraform E2E test index set"
  index_prefix                        = "tf-e2e"
  rotation_strategy_class             = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
  rotation_strategy                   = jsonencode(local.rotation_strategy)
  retention_strategy_class            = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
  retention_strategy                  = jsonencode(local.retention_strategy)
  data_tiering                        = jsonencode({ type = "hot_only", index_lifetime_min = "P30D", index_lifetime_max = "P40D" })
  field_restrictions                  = jsonencode({})
  index_analyzer                      = "standard"
  index_set_template_id               = data.graylog_index_set_template.hot30.id
  shards                              = 1
  replicas                            = 0
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
  index_optimization_disabled         = false
  writable                            = true
  use_legacy_rotation                 = false
}

resource "graylog_input" "gelf_udp" {
  title  = "tf-e2e-gelf-udp"
  type   = "org.graylog2.inputs.gelf.udp.GELFUDPInput"
  global = true
  node   = "" # let Graylog decide

  attributes = jsonencode(local.gelf_udp_config)
}

resource "graylog_input" "beats" {
  title  = "tf-e2e-beats"
  type   = "org.graylog.plugins.beats.BeatsInput"
  global = true
  node   = ""

  attributes = jsonencode({
    bind_address     = "0.0.0.0"
    port             = 5044
    recv_buffer_size = 1048576
    tls_enable       = false
  })
}

# Get the 7 Days Hot template for PacketBeat
data "graylog_index_set_template" "hot7" {
  title = "7 Days Hot"
}

# PacketBeat index set with 7 days retention
resource "graylog_index_set" "packetbeat" {
  title                               = "PacketBeat"
  description                         = "Index set for PacketBeat network telemetry data"
  index_prefix                        = "packetbeat"
  rotation_strategy_class             = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
  rotation_strategy                   = jsonencode({
    type               = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  retention_strategy_class            = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
  retention_strategy                  = jsonencode({
    type                  = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig"
    max_number_of_indices = 5
  })
  data_tiering                        = jsonencode({
    type               = "hot_only"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  field_restrictions                  = jsonencode({})
  index_analyzer                      = "standard"
  index_set_template_id               = data.graylog_index_set_template.hot7.id
  shards                              = 1
  replicas                            = 0
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
  index_optimization_disabled         = false
  writable                            = true
  use_legacy_rotation                 = false
}

# PacketBeat stream
resource "graylog_stream" "packetbeat" {
  title                              = "PacketBeat"
  description                        = "Stream for PacketBeat network telemetry data"
  index_set_id                       = graylog_index_set.packetbeat.id
  matching_type                      = "AND"
  disabled                           = false
  remove_matches_from_default_stream = true
}

# Stream rule to match PacketBeat data by input
resource "graylog_stream_rule" "packetbeat_input" {
  stream_id   = graylog_stream.packetbeat.id
  field       = "gl2_source_input"
  type        = 1  # exact match
  value       = graylog_input.beats.id
  description = "Match messages from Beats input"
  inverted    = false
}

resource "graylog_stream" "tf_e2e" {
  title                              = "tf-e2e-stream"
  description                        = "Terraform E2E test stream"
  index_set_id                       = graylog_index_set.tf_e2e.id
  matching_type                      = "AND"
  disabled                           = false
  remove_matches_from_default_stream = true
}

resource "graylog_stream_rule" "tf_e2e_source" {
  stream_id   = graylog_stream.tf_e2e.id
  field       = "source"
  type        = 1
  value       = local.stream_rule_value
  description = "match tf-e2e source"
  inverted    = false
}

resource "graylog_output" "stdout" {
  title = "tf-e2e-stdout"
  type  = "org.graylog2.outputs.LoggingOutput"

  configuration = jsonencode({
    prefix = "tf-e2e: "
  })
}

resource "graylog_stream_output" "tf_e2e" {
  stream_id  = graylog_stream.tf_e2e.id
  output_ids = [graylog_output.stdout.id]
}

data "graylog_input" "gelf_udp" {
  input_id = graylog_input.gelf_udp.id
}

data "graylog_stream" "tf_e2e" {
  stream_id = graylog_stream.tf_e2e.id
}

data "graylog_stream_rule" "tf_e2e_source" {
  stream_id = graylog_stream_rule.tf_e2e_source.stream_id
  rule_id   = graylog_stream_rule.tf_e2e_source.rule_id
}

data "graylog_output" "stdout" {
  output_id = graylog_output.stdout.id
}

output "index_set_id" {
  value = graylog_index_set.tf_e2e.id
}

output "builtin_templates" {
  value = data.graylog_index_set_templates.all_builtin.templates
}

output "gelf_input_id" {
  value = graylog_input.gelf_udp.id
}

output "stream_id" {
  value = graylog_stream.tf_e2e.id
}

output "stream_rule_id" {
  value = graylog_stream_rule.tf_e2e_source.id
}

output "data_input_attributes" {
  value = data.graylog_input.gelf_udp.attributes
}

output "data_stream_matching_type" {
  value = data.graylog_stream.tf_e2e.matching_type
}

output "data_stream_rule_value" {
  value = data.graylog_stream_rule.tf_e2e_source.value
}

output "output_id" {
  value = graylog_output.stdout.id
}

output "data_output_type" {
  value = data.graylog_output.stdout.type
}

data "graylog_saved_search" "first" {
  title = "My Saved Search"
}

data "graylog_grok_patterns" "all" {}

resource "graylog_pipeline_rule" "tf_e2e" {
  source      = local.pipeline_rule_source
  description = "E2E pipeline rule adds tf_e2e_tag"
}

resource "graylog_pipeline" "tf_e2e" {
  source      = local.pipeline_source
  description = "E2E pipeline using tf-e2e-rule"
}

resource "graylog_pipeline_connection" "tf_e2e" {
  stream_id    = graylog_stream.tf_e2e.id
  pipeline_ids = [graylog_pipeline.tf_e2e.id]
}

output "pipeline_id" {
  value = graylog_pipeline.tf_e2e.id
}

output "pipeline_rule_id" {
  value = graylog_pipeline_rule.tf_e2e.id
}

output "saved_search_id" {
  value = data.graylog_saved_search.first.saved_search_id
}

output "grok_patterns_json" {
  value = data.graylog_grok_patterns.all.patterns_json
}

resource "graylog_dashboard" "tf_e2e" {
  title       = "tf-e2e-dashboard"
  description = "Terraform E2E dashboard"
  summary     = "Dashboard via views API"
  search_id   = data.graylog_saved_search.first.search_id

  state {
    id = data.graylog_saved_search.first.state_id
    widgets {
      widget_id = "tf-e2e-agg"
      type      = "aggregation"
      config = jsonencode({
        row_pivots = [
          {
            type   = "time"
            fields = ["timestamp"]
            config = {
              interval = { type = "auto", scaling = 1.0 }
            }
          }
        ]
        column_pivots = []
        series = [
          {
            function = "count()"
            config   = {}
          }
        ]
        sort          = []
        rollup        = true
        visualization = "bar"
      })
      timerange = jsonencode({
        type  = "relative"
        range = 300
      })
    }
    widget_mapping = jsonencode({})
    positions = jsonencode({
      tf-e2e-agg = {
        col    = 1
        row    = 1
        height = 2
        width  = 2
      }
    })
  }
}

# Dashboard data source - lookup by ID
data "graylog_dashboard" "tf_e2e_by_id" {
  dashboard_id = graylog_dashboard.tf_e2e.id
}

# Dashboard data source - lookup by title
data "graylog_dashboard" "tf_e2e_by_title" {
  title = graylog_dashboard.tf_e2e.title
}

output "dashboard_id" {
  value = graylog_dashboard.tf_e2e.id
}

output "data_dashboard_by_id_title" {
  value = data.graylog_dashboard.tf_e2e_by_id.title
}

output "data_dashboard_by_title_owner" {
  value = data.graylog_dashboard.tf_e2e_by_title.owner
}

# User resource
resource "graylog_user" "tf_e2e" {
  username           = "tf-e2e-user"
  email              = "tf-e2e@example.com"
  first_name         = "Terraform"
  last_name          = "E2E Test User"
  password           = "TestPassword123"
  roles              = ["Reader"]
  timezone           = "UTC"
  session_timeout_ms = 3600000
}

# User data source - lookup by username
data "graylog_user" "tf_e2e_by_username" {
  username = graylog_user.tf_e2e.username
}

# User data source - lookup existing admin user
data "graylog_user" "admin" {
  username = "admin"
}

output "user_id" {
  value = graylog_user.tf_e2e.user_id
}

output "data_user_by_username_email" {
  value = data.graylog_user.tf_e2e_by_username.email
}

output "data_admin_user_roles" {
  value = data.graylog_user.admin.roles
}

# Role resource
resource "graylog_role" "tf_e2e" {
  name        = "tf-e2e-role"
  description = "Terraform E2E test role"
  permissions = ["streams:read"]
}

# Role data source - lookup by name
data "graylog_role" "tf_e2e_by_name" {
  name = graylog_role.tf_e2e.name
}

# Role data source - lookup existing Reader role
data "graylog_role" "reader" {
  name = "Reader"
}

output "role_name" {
  value = graylog_role.tf_e2e.name
}

output "data_role_by_name_permissions" {
  value = data.graylog_role.tf_e2e_by_name.permissions
}

output "data_reader_role_read_only" {
  value = data.graylog_role.reader.read_only
}

# Grok pattern resource
resource "graylog_grok_pattern" "tf_e2e" {
  name    = "TF_E2E_TEST"
  pattern = "\\d{4}-\\d{2}-\\d{2}"
}

# Grok pattern data source - lookup by name
data "graylog_grok_pattern" "tf_e2e_by_name" {
  name = graylog_grok_pattern.tf_e2e.name
}

# Grok pattern data source - lookup existing pattern
data "graylog_grok_pattern" "year" {
  name = "YEAR"
}

output "grok_pattern_name" {
  value = graylog_grok_pattern.tf_e2e.name
}

output "data_grok_pattern_by_name" {
  value = data.graylog_grok_pattern.tf_e2e_by_name.pattern
}

output "data_grok_year_pattern" {
  value = data.graylog_grok_pattern.year.pattern
}

# Input static fields resource
resource "graylog_input_static_fields" "tf_e2e" {
  input_id = graylog_input.gelf_udp.id
  fields = {
    tf_e2e_env    = "test"
    tf_e2e_source = "terraform"
  }
}

output "input_static_fields_input_id" {
  value = graylog_input_static_fields.tf_e2e.input_id
}

# PacketBeat Dashboard
resource "graylog_dashboard" "packetbeat" {
  title       = "PacketBeat Network Overview"
  description = "Network traffic analysis from PacketBeat"
  summary     = "Simple dashboard showing network flows"
  search_id   = data.graylog_saved_search.first.search_id

  state {
    id = data.graylog_saved_search.first.state_id

    # Widget 1: Message count over time
    widgets {
      widget_id = "pb-traffic-over-time"
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

    # Widget 2: Traffic by transport protocol
    widgets {
      widget_id = "pb-by-transport"
      type      = "aggregation"
      config = jsonencode({
        row_pivots    = [{ type = "values", fields = ["packetbeat_network_transport"], config = { limit = 10 } }]
        column_pivots = []
        series        = [{ function = "count()", config = {} }]
        sort          = []
        rollup        = true
        visualization = "pie"
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    # Widget 3: Top destination ports
    widgets {
      widget_id = "pb-top-ports"
      type      = "aggregation"
      config = jsonencode({
        row_pivots    = [{ type = "values", fields = ["packetbeat_destination_port"], config = { limit = 10 } }]
        column_pivots = []
        series = [
          { function = "count()", config = {} },
          { function = "sum(packetbeat_network_bytes)", config = {} }
        ]
        sort          = [{ type = "series", field = "count()", direction = "Descending" }]
        rollup        = true
        visualization = "table"
      })
      timerange = jsonencode({ type = "relative", range = 900 })
    }

    # Widget 4: Top source IPs
    widgets {
      widget_id = "pb-top-sources"
      type      = "aggregation"
      config = jsonencode({
        row_pivots    = [{ type = "values", fields = ["packetbeat_source_ip"], config = { limit = 10 } }]
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
      pb-traffic-over-time = { col = 1, row = 1, height = 3, width = "Infinity" }
      pb-by-transport      = { col = 1, row = 4, height = 3, width = 3 }
      pb-top-ports         = { col = 4, row = 4, height = 3, width = 3 }
      pb-top-sources       = { col = 1, row = 7, height = 3, width = "Infinity" }
    })
    titles = jsonencode({
      widget = {
        pb-traffic-over-time = "Network Traffic Over Time"
        pb-by-transport      = "Traffic by Protocol"
        pb-top-ports         = "Top Destination Ports"
        pb-top-sources       = "Top Source IPs"
      }
    })
  }
}

output "packetbeat_dashboard_id" {
  value = graylog_dashboard.packetbeat.id
}
