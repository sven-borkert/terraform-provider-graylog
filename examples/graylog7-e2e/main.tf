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
