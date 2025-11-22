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
    type              = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig"
    index_lifetime_min = "P30D"
    index_lifetime_max = "P40D"
  }

  retention_strategy = {
    type                 = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig"
    max_number_of_indices = 5
  }
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

output "index_set_id" {
  value = graylog_index_set.tf_e2e.id
}

output "builtin_templates" {
  value = data.graylog_index_set_templates.all_builtin.templates
}
