# main.tf - Main configuration for local Graylog provider testing
#
# This file demonstrates basic usage of the Graylog provider for local testing.
# It includes provider configuration, variables, and a simple example resource.

terraform {
  required_version = ">= 0.13"

  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "3.0.0"
    }
  }
}

# Provider configuration
# You can also use environment variables:
# - GRAYLOG_WEB_ENDPOINT_URI or TF_VAR_graylog_endpoint
# - GRAYLOG_AUTH_NAME or TF_VAR_graylog_username
# - GRAYLOG_AUTH_PASSWORD or TF_VAR_graylog_password
provider "graylog" {
  web_endpoint_uri = var.graylog_endpoint
  auth_name        = var.graylog_username
  auth_password    = var.graylog_password
  api_version      = var.graylog_api_version
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

variable "graylog_api_version" {
  description = "Graylog API version (v1 for older versions, v2 for newer)"
  type        = string
  default     = "v1"
}

# Example resource: Create a simple test stream
# This demonstrates the minimal configuration needed for a stream
resource "graylog_stream" "test_stream" {
  title                              = "Test Stream from Terraform"
  description                        = "A test stream created by Terraform for local provider testing"
  index_set_id                       = data.graylog_index_set.default.id
  remove_matches_from_default_stream = false
}

# Add a simple rule to the test stream (optional)
# Uncomment to test stream rules
# resource "graylog_stream_rule" "test_rule" {
#   stream_id   = graylog_stream.test_stream.id
#   field       = "source"
#   value       = "test"
#   type        = 1  # EXACT match
#   inverted    = false
#   description = "Match messages with source=test"
# }

# Output the created stream details
output "test_stream_details" {
  description = "Details of the test stream created by Terraform"
  value = {
    id          = graylog_stream.test_stream.id
    title       = graylog_stream.test_stream.title
    description = graylog_stream.test_stream.description
    index_set   = graylog_stream.test_stream.index_set_id
  }
}

# Output provider status
output "provider_status" {
  description = "Provider connection status"
  value       = "Connected to Graylog at ${var.graylog_endpoint}"
}