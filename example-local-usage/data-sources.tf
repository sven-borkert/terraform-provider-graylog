# data-sources.tf - Examples of using data sources to query existing Graylog resources
#
# Data sources allow you to fetch information about existing resources
# without managing them in Terraform. This is useful for:
# - Referencing existing resources in new configurations
# - Discovering resource IDs for imports
# - Validation and testing

# Example: Query the default index set
data "graylog_index_set" "default" {
  index_set_id = "default"
}

# Output the default index set details
output "default_index_set" {
  description = "Details of the default index set"
  value = {
    id          = data.graylog_index_set.default.id
    title       = data.graylog_index_set.default.title
    description = data.graylog_index_set.default.description
    index_prefix = data.graylog_index_set.default.index_prefix
  }
}

# Example: Query a specific stream (uncomment and provide actual ID)
# data "graylog_stream" "example" {
#   stream_id = "your-stream-id-here"
# }
#
# output "stream_details" {
#   description = "Stream configuration"
#   value       = data.graylog_stream.example
# }

# Example: Query a specific dashboard (uncomment and provide actual ID)
# data "graylog_dashboard" "example" {
#   dashboard_id = "your-dashboard-id-here"
# }
#
# output "dashboard_details" {
#   description = "Dashboard configuration"
#   value       = data.graylog_dashboard.example
# }

# Example: Query a specific input (uncomment and provide actual ID)
# data "graylog_input" "example" {
#   input_id = "your-input-id-here"
# }
#
# output "input_details" {
#   description = "Input configuration"
#   value = {
#     id    = data.graylog_input.example.id
#     title = data.graylog_input.example.title
#     type  = data.graylog_input.example.type
#   }
# }

# Example: Query system information
# Note: This might require admin privileges
# data "graylog_system" "info" {}
#
# output "system_info" {
#   description = "Graylog system information"
#   value       = data.graylog_system.info
# }