# Example: Managed Graylog Resources
# Rename this file to managed-resources.tf to use these examples

# Create a new index set for application logs
resource "graylog_index_set" "application_logs" {
  title                = "Application Logs"
  description          = "Index set for application logs"
  index_prefix         = "app"
  shards               = 2
  replicas             = 0
  rotation_strategy    = "count"

  rotation_strategy_details = jsonencode({
    type               = "org.graylog2.indexer.rotation.strategies.MessageCountRotationStrategy"
    max_docs_per_index = 20000000
  })

  retention_strategy = "delete"
  retention_strategy_details = jsonencode({
    type                  = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
    max_number_of_indices = 20
  })

  index_analyzer                      = "standard"
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
}

# Create a stream for application logs
resource "graylog_stream" "application_stream" {
  title                              = "Application Stream"
  description                        = "Stream for application logs"
  index_set_id                       = graylog_index_set.application_logs.id
  remove_matches_from_default_stream = true
}

# Add a rule to the stream
resource "graylog_stream_rule" "app_source_rule" {
  stream_id   = graylog_stream.application_stream.id
  field       = "source"
  value       = "application"
  type        = 1  # EXACT
  inverted    = false
  description = "Match logs with source=application"
}

# Create a Syslog UDP input
resource "graylog_input" "syslog_udp" {
  title  = "Syslog UDP"
  type   = "org.graylog2.inputs.syslog.udp.SyslogUDPInput"
  global = true

  configuration = jsonencode({
    bind_address          = "0.0.0.0"
    port                  = 1514
    recv_buffer_size      = 262144
    number_worker_threads = 8
    override_source       = null
    force_rdns            = false
    allow_override_date   = true
    store_full_message    = true
  })
}

# Create a GELF HTTP input
resource "graylog_input" "gelf_http" {
  title  = "GELF HTTP"
  type   = "org.graylog2.inputs.gelf.http.GELFHttpInput"
  global = true

  configuration = jsonencode({
    bind_address                 = "0.0.0.0"
    port                         = 12201
    recv_buffer_size             = 1048576
    max_message_size             = 2097152
    number_worker_threads        = 8
    tls_cert_file                = ""
    tls_key_file                 = ""
    tls_key_password             = ""
    tls_enable                   = false
    tls_client_auth              = "disabled"
    tls_client_auth_cert_file    = ""
    tcp_keepalive                = false
    use_null_delimiter           = true
    override_source              = null
    idle_writer_timeout          = 60
    enable_cors                  = true
    decompress_size_limit        = 8388608
  })
}

# Create a pipeline for processing logs
resource "graylog_pipeline" "log_processing" {
  source = <<EOF
pipeline "Log Processing"
stage 0 match all
  rule "parse application logs"
  rule "set timestamp"
stage 1 match all
  rule "route to streams"
end
EOF

  description = "Main log processing pipeline"
}

# Create pipeline rules
resource "graylog_pipeline_rule" "parse_app_logs" {
  source = <<EOF
rule "parse application logs"
when
  has_field("application")
then
  let parsed = parse_json(to_string($message.message));
  set_fields(parsed);
end
EOF

  description = "Parse JSON application logs"
}

resource "graylog_pipeline_rule" "set_timestamp" {
  source = <<EOF
rule "set timestamp"
when
  has_field("timestamp")
then
  let ts = parse_date(value: to_string($message.timestamp), pattern: "yyyy-MM-dd'T'HH:mm:ss.SSSZ");
  set_field("timestamp", ts);
end
EOF

  description = "Set correct timestamp from log"
}

# Create a user for application access
resource "graylog_user" "app_user" {
  username  = "app_reader"
  email     = "app@example.com"
  full_name = "Application Reader"
  password  = "ChangeMeImmediately123!"  # Change this password!
  roles = [
    "Reader"
  ]
}

# Create a custom role
resource "graylog_role" "log_analyst" {
  name        = "Log Analyst"
  description = "Role for log analysis with limited permissions"

  permissions = [
    "streams:read",
    "dashboards:read",
    "searches:absolute",
    "searches:relative",
    "searches:keyword"
  ]
}

# Create an alert event definition
resource "graylog_event_definition" "high_error_rate" {
  title       = "High Error Rate"
  description = "Alert when error rate is too high"
  priority    = 2
  alert       = true

  config = jsonencode({
    type        = "aggregation-v1"
    query       = "level:ERROR"
    streams     = []
    group_by    = []
    series      = []
    conditions  = {
      expression = {
        expr = "count() > 100"
        ref  = ""
      }
    }
    search_within_ms         = 300000  # 5 minutes
    execute_every_ms         = 60000   # 1 minute
  })

  notification_settings = jsonencode({
    grace_period_ms = 300000
    backlog_size    = 10
  })

  field_spec = jsonencode({})
}

# Create an email notification
resource "graylog_event_notification" "email_alert" {
  title       = "Email Alert"
  description = "Send email notifications for critical alerts"

  config = jsonencode({
    type              = "email-notification-v1"
    sender            = "graylog@example.com"
    subject           = "Graylog Alert: $${event_definition_title}"
    body_template     = "Alert triggered: $${event_definition_description}\n\nBacklog:\n$${backlog}"
    html_body_template = ""
    email_recipients  = ["ops-team@example.com"]
    user_recipients   = []
    time_zone         = "UTC"
  })
}

# Output useful information
output "application_stream_id" {
  description = "ID of the created application stream"
  value       = graylog_stream.application_stream.id
}

output "syslog_input_port" {
  description = "Port for the Syslog UDP input"
  value       = 1514
}

output "gelf_http_endpoint" {
  description = "GELF HTTP input endpoint"
  value       = "http://${var.graylog_endpoint}:12201/gelf"
}