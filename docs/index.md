---
page_title: "Provider: Graylog"
description: |-
  Terraform provider for managing Graylog resources
---

# Graylog Provider

The Graylog provider enables Terraform to manage [Graylog](https://www.graylog.org/) infrastructure as code. It supports Graylog 7.0+ and provides comprehensive resource management capabilities.

## Provider Features

- ✅ **Graylog 7.0+ Compatible** - Full support for Graylog 7.0 API changes
- ✅ **Comprehensive Resource Coverage** - Manage streams, inputs, outputs, pipelines, users, and more
- ✅ **Data Sources** - Query existing Graylog resources
- ✅ **Import Support** - Import existing resources into Terraform state
- ✅ **Type Safety** - Strong typing with proper validation
- ✅ **Documentation** - Extensive examples and guides

## Quick Start

### Basic Configuration

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "~> 3.0"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "admin"
  auth_password    = var.graylog_password
  api_version      = "v1"
}
```

### Example Resource

```hcl
# Create a stream for application logs
resource "graylog_stream" "app_logs" {
  title                              = "Application Logs"
  description                        = "Stream for application log messages"
  index_set_id                       = data.graylog_index_set.default.id
  remove_matches_from_default_stream = true
}

# Add a routing rule
resource "graylog_stream_rule" "app_filter" {
  stream_id   = graylog_stream.app_logs.id
  field       = "application"
  value       = "myapp"
  type        = 1  # EXACT match
  description = "Route messages from myapp"
}

# Create an input
resource "graylog_input" "gelf_http" {
  title  = "GELF HTTP"
  type   = "org.graylog2.inputs.gelf.http.GELFHttpInput"
  global = true

  attributes = jsonencode({
    bind_address = "0.0.0.0"
    port         = 12201
  })
}

# Query existing resources
data "graylog_index_set" "default" {
  index_set_id = "default"
}
```

## Provider Configuration

### Arguments

* `web_endpoint_uri` - (Required) Graylog API endpoint URL. Must include `/api` path (e.g., `https://graylog.example.com/api`). Can be set via `GRAYLOG_WEB_ENDPOINT_URI` environment variable.

* `auth_name` - (Required) Authentication username, API token, or session token. Can be set via `GRAYLOG_AUTH_NAME` environment variable.

* `auth_password` - (Required) Authentication password, or literal `"token"`/`"session"` when using tokens. Can be set via `GRAYLOG_AUTH_PASSWORD` environment variable.

* `api_version` - (Optional) Graylog API version. Defaults to `v1`. Can be set via `GRAYLOG_API_VERSION` environment variable.

* `x_requested_by` - (Optional) Value for the `X-Requested-By` header. Defaults to `terraform-provider-graylog`. Can be set via `GRAYLOG_X_REQUESTED_BY` environment variable.

### Authentication Methods

#### Password Authentication

```hcl
provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "admin"
  auth_password    = "your-password"
}
```

#### API Token Authentication

```hcl
provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "your-api-token-here"
  auth_password    = "token"
}
```

#### Session Token Authentication

```hcl
provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "your-session-token-here"
  auth_password    = "session"
}
```

### Environment Variables

```bash
export GRAYLOG_WEB_ENDPOINT_URI="https://graylog.example.com/api"
export GRAYLOG_AUTH_NAME="admin"
export GRAYLOG_AUTH_PASSWORD="your-password"
export GRAYLOG_API_VERSION="v1"
```

With environment variables set, the provider configuration can be minimal:

```hcl
provider "graylog" {
  # Configuration loaded from environment variables
}
```

## Compatibility

### Version Matrix

| Provider Version | Graylog Version | Status |
|-----------------|-----------------|--------|
| v3.0.0+ | Graylog 7.0+ | ✅ Supported |
| v1.x.x | Graylog 3.x - 6.x | ⚠️ Legacy (not compatible with 7.0) |

### Graylog 7.0 Changes

Version 3.0.0+ of this provider includes full support for Graylog 7.0's API changes:

- **Automatic API Format Handling** - Entity creation requests are automatically wrapped
- **Computed Field Management** - Read-only fields are automatically removed from updates
- **Zero Configuration Changes** - Your existing Terraform configurations work without modification

See the [Migration Guide](guides/migration_guide_v7) for upgrading from earlier versions.

## Available Resources

### Log Management
- **[graylog_index_set](resources/index_set)** - Manage Elasticsearch index sets
- **[graylog_stream](resources/stream)** - Create and configure log streams
- **[graylog_stream_rule](resources/stream_rule)** - Define stream routing rules
- **[graylog_stream_output](resources/stream_output)** - Connect streams to outputs

### Data Inputs
- **[graylog_input](resources/input)** - Configure log inputs (Syslog, GELF, Beats, etc.)
- **[graylog_input_static_fields](resources/input_static_fields)** - Add static fields to inputs
- **[graylog_extractor](resources/extractor)** - Extract data from log messages

### Processing
- **[graylog_pipeline](resources/pipeline)** - Create message processing pipelines
- **[graylog_pipeline_rule](resources/pipeline_rule)** - Define pipeline processing rules
- **[graylog_pipeline_connection](resources/pipeline_connection)** - Connect pipelines to streams
- **[graylog_grok_pattern](resources/grok_pattern)** - Manage Grok patterns

### Outputs
- **[graylog_output](resources/output)** - Configure outputs for forwarding messages

### Alerting & Events
- **[graylog_event_definition](resources/event_definition)** - Define event conditions
- **[graylog_event_notification](resources/event_notification)** - Configure notifications
- **[graylog_alert_condition](resources/alert_condition)** - Legacy alert conditions (deprecated)
- **[graylog_alarm_callback](resources/alarm_callback)** - Legacy alarm callbacks (deprecated)

### Visualization
- **[graylog_dashboard](resources/dashboard)** - Create dashboards
- **[graylog_dashboard_widget](resources/dashboard_widget)** - Add widgets to dashboards
- **[graylog_dashboard_widget_positions](resources/dashboard_widget_positions)** - Manage widget layouts

### Access Control
- **[graylog_user](resources/user)** - Manage users
- **[graylog_role](resources/role)** - Define roles and permissions
- **[graylog_ldap_setting](resources/ldap_setting)** - Configure LDAP authentication

### Sidecar Management
- **[graylog_sidecar_configuration](resources/sidecar_configuration)** - Configure sidecar collectors
- **[graylog_sidecar_collector](resources/sidecar_collector)** - Manage collector definitions
- **[graylog_sidecars](resources/sidecars)** - Manage sidecar instances

## Available Data Sources

- **[graylog_index_set](data-sources/index_set)** - Query index set information
- **[graylog_stream](data-sources/stream)** - Query stream details
- **[graylog_dashboard](data-sources/dashboard)** - Query dashboard configuration
- **[graylog_sidecar](data-sources/sidecar)** - Query sidecar information

## Documentation

### Getting Started
- **[Local Testing Guide](guides/local_usage)** - Test the provider locally
- **[Migration Guide](guides/migration_guide_v7)** - Upgrade to Graylog 7.0
- **[JSON String Attributes](guides/json-string-attribute)** - Working with JSON configurations

### Reference
- **[API Mapping](reference/api_mapping)** - Graylog API endpoint documentation
- **[Resource Inventory](reference/inventory)** - Complete component inventory
- **[Changelog](changelog)** - Version history and changes
- **[Provider Naming](guides/provider_naming)** - Naming conventions

### Examples
- **[Local Usage Examples](../../example-local-usage/)** - Complete working examples with multiple testing methods

## Development

### Local Testing

See the comprehensive [example-local-usage](../../example-local-usage/) directory for:
- Multiple testing methods (standard, developer override, Makefile)
- Build scripts and helper tools
- Import examples and data source queries
- Complete documentation

Quick start:
```bash
cd example-local-usage/
make setup    # Build provider and configure
make plan     # Test with terraform plan
```

### Contributing

Contributions are welcome! Please see the [GitHub repository](https://github.com/terraform-provider-graylog/terraform-provider-graylog) for:
- Issue reporting
- Feature requests
- Pull requests
- Development guidelines

## Support

### Resources
- **[GitHub Issues](https://github.com/terraform-provider-graylog/terraform-provider-graylog/issues)** - Bug reports and feature requests
- **[Graylog Community](https://community.graylog.org/)** - Community support and discussions
- **[Graylog Documentation](https://go2docs.graylog.org/)** - Official Graylog documentation
- **[Terraform Registry](https://registry.terraform.io/providers/terraform-provider-graylog/graylog)** - Provider documentation

### Graylog 7.0 Resources
- **[Upgrade Guide](https://go2docs.graylog.org/current/upgrading_graylog/upgrade_to_graylog_7.0.htm)** - Official upgrade guide
- **[API Documentation](https://go2docs.graylog.org/current/setting_up_graylog/rest_api.html)** - REST API reference
- **[Release Notes](https://www.graylog.org/post/announcing-graylog-v7-0)** - What's new in Graylog 7.0