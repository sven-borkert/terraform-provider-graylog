# Example Terraform configuration for using the local Graylog provider

terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "2.0.0"
    }
  }
}

# Provider configuration
# You can also use environment variables:
# - GRAYLOG_WEB_ENDPOINT_URI
# - GRAYLOG_AUTH_NAME
# - GRAYLOG_AUTH_PASSWORD
provider "graylog" {
  web_endpoint_uri = var.graylog_endpoint
  auth_name        = var.graylog_username
  auth_password    = var.graylog_password
  api_version      = "v1"
}

# Input variables
variable "graylog_endpoint" {
  description = "Graylog API endpoint (e.g., https://graylog.example.com/api)"
  type        = string
}

variable "graylog_username" {
  description = "Graylog admin username"
  type        = string
  default     = "admin"
}

variable "graylog_password" {
  description = "Graylog admin password"
  type        = string
  sensitive   = true
}

# Example resource: Create a simple stream
resource "graylog_stream" "test_stream" {
  title                              = "Test Stream from Terraform"
  description                        = "A test stream created by Terraform"
  index_set_id                       = data.graylog_index_set.default.id
  remove_matches_from_default_stream = false
}

# Data source to get the default index set
data "graylog_index_set" "default" {
  index_set_id = "default"
}

# Output the created stream ID
output "stream_id" {
  value       = graylog_stream.test_stream.id
  description = "ID of the created stream"
}