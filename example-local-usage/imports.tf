# imports.tf - Resource blocks for importing existing Graylog resources
#
# This file contains empty resource blocks that can be used for importing
# existing resources from your Graylog server into Terraform state.
#
# To import resources:
# 1. Uncomment the resource block you want to import
# 2. Give it a meaningful name in Terraform
# 3. Run: terraform import <resource_type>.<resource_name> <resource_id>
# 4. After import, run 'terraform plan' to see the current configuration
# 5. Update the resource block to match the imported state

# Example: Import an existing stream
# resource "graylog_stream" "existing_stream" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_stream.existing_stream <STREAM_ID>

# Example: Import an existing input
# resource "graylog_input" "existing_input" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_input.existing_input <INPUT_ID>

# Example: Import an existing index set
# resource "graylog_index_set" "existing_index_set" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_index_set.existing_index_set <INDEX_SET_ID>

# Example: Import an existing dashboard
# resource "graylog_dashboard" "existing_dashboard" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_dashboard.existing_dashboard <DASHBOARD_ID>

# Example: Import an existing user
# resource "graylog_user" "existing_user" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_user.existing_user <USERNAME>

# Example: Import an existing role
# resource "graylog_role" "existing_role" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_role.existing_role <ROLE_NAME>

# Example: Import an existing alert condition
# resource "graylog_alert_condition" "existing_alert" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_alert_condition.existing_alert <STREAM_ID>/<CONDITION_ID>

# Example: Import an existing pipeline
# resource "graylog_pipeline" "existing_pipeline" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_pipeline.existing_pipeline <PIPELINE_ID>

# Example: Import an existing pipeline rule
# resource "graylog_pipeline_rule" "existing_rule" {
#   # After import, Terraform will populate this with the actual configuration
# }
# Import command: terraform import graylog_pipeline_rule.existing_rule <RULE_ID>