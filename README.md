# terraform-provider-graylog

> **Note:** This project has been significantly refactored with AI assistance to add Graylog 7.0 support. While functional, not all features have been thoroughly tested. Use with caution in production environments and please report any issues you encounter.

[![Build Status](https://github.com/sven-borkert/terraform-provider-graylog/workflows/CI/badge.svg)](https://github.com/sven-borkert/terraform-provider-graylog/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/sven-borkert/terraform-provider-graylog)](https://goreportcard.com/report/github.com/sven-borkert/terraform-provider-graylog)
[![GitHub last commit](https://img.shields.io/github/last-commit/sven-borkert/terraform-provider-graylog.svg)](https://github.com/sven-borkert/terraform-provider-graylog)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/sven-borkert/terraform-provider-graylog/master/LICENSE)

Terraform Provider for [Graylog](https://docs.graylog.org/)

**Now with Graylog 7.0 support!** 🎉

## Compatibility

| Provider Version | Graylog Version | Status |
|-----------------|-----------------|--------|
| v3.0.0+ | Graylog 7.0+ | ✅ Fully supported |
| v1.x.x | Graylog 3.x - 6.x | ⚠️ Not compatible with Graylog 7.0 |

**Upgrading to Graylog 7.0?** See the [Migration Guide](MIGRATION_GUIDE_V7.md)

## Quick Start

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "sven-borkert/graylog"
      version = "~> 2.0"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "admin"
  auth_password    = var.graylog_password
}
```

## Documentation

- [Provider Documentation](https://registry.terraform.io/providers/sven-borkert/graylog/latest/docs)
- [Terraform Registry](https://registry.terraform.io/providers/sven-borkert/graylog/latest)
- [Migration Guide to Graylog 7.0](MIGRATION_GUIDE_V7.md)
- [API Mapping for Graylog 7.0](API_MAPPING.md)

## What's New in v2.0

### Graylog 7.0 API Compatibility

The provider has been fully updated to support Graylog 7.0's breaking API changes:

- ✅ **CreateEntityRequest wrapper** - All entity creation automatically uses the new Graylog 7.0 request format
- ✅ **Unknown properties validation** - Computed fields are automatically removed from update requests
- ✅ **Backward compatible** - Existing Terraform configurations work without changes

### Supported Resources (25)

**Streams & Alerting:**
- `graylog_stream` - Stream management
- `graylog_stream_rule` - Stream routing rules
- `graylog_stream_output` - Stream output associations
- `graylog_alarm_callback` - Legacy alarm callbacks (deprecated)
- `graylog_alert_condition` - Legacy alert conditions (deprecated)

**Events System:**
- `graylog_event_definition` - Modern event definitions
- `graylog_event_notification` - Modern notifications

**Inputs & Processing:**
- `graylog_input` - Input configuration
- `graylog_input_static_fields` - Static field enrichment
- `graylog_extractor` - Message extractors
- `graylog_grok_pattern` - Grok pattern library

**Pipelines:**
- `graylog_pipeline` - Processing pipelines
- `graylog_pipeline_rule` - Pipeline rules
- `graylog_pipeline_connection` - Pipeline to stream connections

**Dashboards:**
- `graylog_dashboard` - Dashboard management
- `graylog_dashboard_widget` - Dashboard widgets
- `graylog_dashboard_widget_positions` - Widget layout

**System:**
- `graylog_index_set` - Index set configuration
- `graylog_output` - Output destinations
- `graylog_ldap_setting` - LDAP authentication

**Security:**
- `graylog_user` - User management
- `graylog_role` - Role-based access control

**Sidecars:**
- `graylog_sidecars` - Sidecar registration
- `graylog_sidecar_collector` - Collector configuration
- `graylog_sidecar_configuration` - Sidecar configs

### Supported Data Sources (4)

- `graylog_stream` - Query streams
- `graylog_dashboard` - Query dashboards
- `graylog_index_set` - Query index sets
- `graylog_sidecar` - Query sidecars
 
## Development

- Build locally:
  ```bash
  make build
  ```
- Optional local mirror install (for workflows without dev_overrides):
  ```bash
  make dev-install
  ```
- When using `~/.terraformrc` dev_overrides (recommended for development), Terraform uses whatever binary is in `bin/` at plan/apply time. Re-run `make build` to pick up changes; no re-init needed. If you do run `terraform init`, clear `.terraform/` and `.terraform.lock.hcl` to avoid stale binaries.

## License

[MIT](LICENSE)
