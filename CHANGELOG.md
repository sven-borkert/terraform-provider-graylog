# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Graylog 7.0 API compatibility** - Full support for Graylog 7.0's breaking API changes
- `util.WrapEntityForCreation()` - Helper function to wrap entities in CreateEntityRequest structure
- `util.RemoveComputedFields()` - Helper function to remove read-only fields before updates
- **Migration guide** - Comprehensive guide for upgrading to Graylog 7.0 (`MIGRATION_GUIDE_V7.md`)
- **API mapping documentation** - Detailed documentation of API changes (`API_MAPPING.md`)
- **Component inventory** - Complete inventory of all provider resources and clients (`INVENTORY.md`)

### Changed
- **BREAKING:** Updated all client Create() methods to use Graylog 7.0 CreateEntityRequest wrapper
  - Affected resources: streams, dashboards, event definitions, event notifications, index sets, inputs, outputs, pipelines, pipeline rules, grok patterns, roles, users, sidecar collectors, sidecar configurations
- **BREAKING:** Requires Graylog 7.0 or later
- Updated all resource update methods to automatically remove computed fields
- Stream client Create() now wraps entity data for Graylog 7.0 compatibility
- Dashboard client Create() now wraps entity data for Graylog 7.0 compatibility
- Event Definition client Create() now wraps entity data for Graylog 7.0 compatibility
- Event Notification client Create() now wraps entity data for Graylog 7.0 compatibility
- Index Set client Create() now wraps entity data for Graylog 7.0 compatibility
- Input client Create() now wraps entity data for Graylog 7.0 compatibility
- Output client Create() now wraps entity data for Graylog 7.0 compatibility
- Pipeline client Create() now wraps entity data for Graylog 7.0 compatibility
- Pipeline Rule client Create() now wraps entity data for Graylog 7.0 compatibility
- Grok Pattern client Create() now wraps entity data for Graylog 7.0 compatibility
- Role client Create() now wraps entity data for Graylog 7.0 compatibility
- User client Create() now wraps entity data for Graylog 7.0 compatibility
- Sidecar Collector client Create() now wraps entity data for Graylog 7.0 compatibility
- Sidecar Configuration client Create() now wraps entity data for Graylog 7.0 compatibility
- All resource update methods now use `util.RemoveComputedFields()` to clean data before API calls
- Updated README with Graylog 7.0 compatibility information

### Fixed
- Fixed unknown properties validation errors in Graylog 7.0 by removing computed fields
- Fixed stream updates failing due to read-only fields (id, created_at, creator_user_id)
- Fixed dashboard updates failing due to read-only fields
- Fixed event definition updates failing due to read-only fields
- Fixed event notification updates failing due to read-only fields
- Fixed index set updates failing due to read-only fields
- Fixed all resource updates to handle Graylog 7.0's strict property validation

### Deprecated
- `graylog_alarm_callback` - Use `graylog_event_notification` instead (Event System)
- `graylog_alert_condition` - Use `graylog_event_definition` instead (Event System)

**Note:** Deprecated resources still work but are marked for future removal. Graylog recommends using the Events System for alerting.

## Migration Notes

### From v1.x to v2.0

**Prerequisites:**
- Graylog must be upgraded to version 7.0 or later first
- Backup your Terraform state before upgrading

**Breaking Changes:**
1. Provider now requires Graylog 7.0+
2. All entity creation uses new API format (handled automatically)
3. Computed fields automatically removed from updates (no configuration changes needed)

**Migration Steps:**
1. Upgrade Graylog server to 7.0+
2. Update provider version in your Terraform configuration:
   ```hcl
   terraform {
     required_providers {
       graylog = {
         source  = "terraform-provider-graylog/graylog"
         version = "~> 2.0"
       }
     }
   }
   ```
3. Run `terraform init -upgrade`
4. Run `terraform plan` to verify
5. Run `terraform apply` if changes detected

**No Terraform Configuration Changes Required:**
Your existing `.tf` files work without modification. The provider handles all API format changes internally.

See [MIGRATION_GUIDE_V7.md](MIGRATION_GUIDE_V7.md) for detailed migration instructions.

## Previous Releases

### [1.x.x] - Historical Versions

Previous versions supported Graylog 3.x through 6.x. See git history for detailed changelog of v1.x releases.

**Note:** v1.x is not compatible with Graylog 7.0 due to breaking API changes.

---

## Release Process

### Version Numbering

- **Major version (X.0.0)**: Breaking changes, requires migration
- **Minor version (0.X.0)**: New features, backward compatible
- **Patch version (0.0.X)**: Bug fixes, backward compatible

### Current Status

- **v2.0.0**: In development - Graylog 7.0 compatibility
- **v1.x.x**: Maintenance mode - Graylog 3.x - 6.x support

---

## Links

- [GitHub Repository](https://github.com/terraform-provider-graylog/terraform-provider-graylog)
- [Terraform Registry](https://registry.terraform.io/providers/terraform-provider-graylog/graylog)
- [Graylog Documentation](https://go2docs.graylog.org/current/)
- [Graylog 7.0 Upgrade Guide](https://go2docs.graylog.org/current/upgrading_graylog/upgrade_to_graylog_7.0.htm)
- [Provider Documentation](https://registry.terraform.io/providers/terraform-provider-graylog/graylog/latest/docs)

---

**Last Updated:** 2025-11-12
