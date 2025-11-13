# Migration Guide: Graylog 7.0 Compatibility

**Version:** terraform-provider-graylog v3.0.0+
**Target:** Graylog 7.0+
**Date:** 2025-11-12

---

## Overview

This guide helps you migrate from previous versions of terraform-provider-graylog to the version compatible with Graylog 7.0. The provider has been updated to support the breaking changes introduced in Graylog 7.0's REST API.

---

## Breaking Changes in Graylog 7.0

### 1. CreateEntityRequest Wrapper (CRITICAL)

**What Changed:**
Graylog 7.0 now requires all entity creation requests to use a new `CreateEntityRequest` structure that separates entity data from sharing information.

**Impact:**
- The terraform provider's internal API calls have been updated
- **No changes required in your Terraform configurations**
- The provider handles the wrapper automatically

**Affected Resources:**
- `graylog_stream`
- `graylog_dashboard`
- `graylog_event_definition`
- `graylog_event_notification`
- `graylog_index_set`
- `graylog_input`
- `graylog_output`
- `graylog_pipeline`
- `graylog_pipeline_rule`
- `graylog_grok_pattern`
- `graylog_role`
- `graylog_user`
- `graylog_sidecar_collector`
- `graylog_sidecar_configuration`

### 2. Unknown Properties Validation

**What Changed:**
Graylog 7.0 now rejects API requests containing unknown or read-only JSON properties.

**Impact:**
- The provider now automatically removes computed/read-only fields before updates
- **No changes required in your Terraform configurations**

**Fixed Automatically:**
- `id`, `created_at`, `creator_user_id`, `last_modified` fields are now automatically removed from update requests

---

## Compatibility Matrix

| Provider Version | Graylog Version | Status |
|-----------------|-----------------|--------|
| v1.x.x | Graylog 3.x - 6.x | ⚠️ Not compatible with Graylog 7.0 |
| v3.0.0+ | Graylog 7.0+ | ✅ Fully compatible |

---

## Migration Steps

### Step 1: Upgrade Graylog Server

1. Backup your Graylog configuration and data
2. Upgrade Graylog to version 7.0 or later
3. Follow the [official Graylog 7.0 upgrade guide](https://go2docs.graylog.org/current/upgrading_graylog/upgrade_to_graylog_7.0.htm)
4. Verify Graylog is running correctly

### Step 2: Upgrade Provider Version

**In your Terraform configuration:**

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "~> 3.0"  # Update to v3.0.0+
    }
  }
}
```

**Run:**
```bash
terraform init -upgrade
```

### Step 3: Validate Configuration

Run a plan to ensure everything is working:

```bash
terraform plan
```

**Expected behavior:**
- No changes should be detected if your infrastructure matches your configuration
- The provider should successfully communicate with Graylog 7.0

### Step 4: Apply if Needed

If changes are detected, review them carefully:

```bash
terraform apply
```

---

## Common Migration Scenarios

### Scenario 1: No Infrastructure Changes

**Situation:** Your Graylog configuration hasn't changed, you're just upgrading versions.

**Steps:**
1. Upgrade Graylog server to 7.0
2. Upgrade provider to v3.0.0+
3. Run `terraform plan`
4. Verify no changes detected

**Expected Result:** ✅ No changes, state is in sync

---

### Scenario 2: Resource Recreation

**Situation:** Terraform wants to recreate resources after upgrade.

**Troubleshooting:**
```bash
# Check Terraform state
terraform show

# If resources appear different, refresh state
terraform refresh

# Review what changed
terraform plan

# If only computed fields changed, this is expected
# Apply the changes
terraform apply
```

**Common causes:**
- Computed fields now handled correctly
- State format updates from v1 to v2

---

### Scenario 3: API Permission Errors

**Situation:** Getting 403 errors when running Terraform.

**Solution:**
Ensure your Graylog user has the new `api_browser:read` permission:

```bash
# In Graylog UI:
# System → Users & Teams → [Your User] → Roles
# Assign "API Browser Reader" role
```

---

## Testing Your Migration

### Basic Connectivity Test

```hcl
# test.tf
provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com/api"
  auth_name        = "admin"
  auth_password    = "your-password"
}

data "graylog_stream" "default" {
  id = "000000000000000000000001"
}

output "test" {
  value = data.graylog_stream.default.title
}
```

Run:
```bash
terraform plan
```

**Expected:** Successfully fetches the default stream.

---

### Create Test Resource

```hcl
resource "graylog_stream" "test" {
  title        = "Test Migration Stream"
  index_set_id = var.index_set_id
  description  = "Testing Graylog 7.0 compatibility"
}
```

Run:
```bash
terraform apply
```

**Expected:** Stream created successfully in Graylog 7.0.

---

## Deprecated Resources

The following resources are deprecated (since Graylog 3.0) but still functional:

| Resource | Replacement | Status |
|----------|-------------|--------|
| `graylog_alarm_callback` | `graylog_event_notification` | ⚠️ Deprecated - use Events System |
| `graylog_alert_condition` | `graylog_event_definition` | ⚠️ Deprecated - use Events System |

**Recommendation:** Migrate to the Events System resources for better support and features.

---

## Rollback Plan

If you encounter issues, you can rollback:

### Option 1: Rollback Graylog Server

1. Restore Graylog from backup (pre-7.0 version)
2. Downgrade provider version in Terraform:
   ```hcl
   version = "~> 1.0"
   ```
3. Run `terraform init -upgrade`

### Option 2: Pin Provider Version

If Graylog 7.0 is stable but provider has issues:

```hcl
terraform {
  required_providers {
    graylog = {
      source  = "terraform-provider-graylog/graylog"
      version = "= 1.15.0"  # Last v1.x version
    }
  }
}
```

---

## API Changes Reference

### CreateEntityRequest Format

**Old (Pre-7.0):**
```json
POST /api/streams
{
  "title": "My Stream",
  "index_set_id": "123",
  "description": "Test stream"
}
```

**New (7.0+):**
```json
POST /api/streams
{
  "entity": {
    "title": "My Stream",
    "index_set_id": "123",
    "description": "Test stream"
  },
  "share_request": {
    "selected_grantee_capabilities": {}
  }
}
```

**Provider Handling:** The provider automatically wraps your data in the correct format.

---

## Troubleshooting

### Error: "Unknown properties in request body"

**Cause:** Graylog 7.0 rejects unknown fields.

**Solution:** Upgrade to provider v3.0.0+ which automatically removes computed fields.

---

### Error: "Invalid request format"

**Cause:** Using old provider with Graylog 7.0.

**Solution:** Upgrade provider to v3.0.0+:
```bash
terraform init -upgrade
```

---

### State Drift Detected

**Cause:** Computed fields now handled differently.

**Solution:** Run refresh and apply:
```bash
terraform refresh
terraform apply
```

---

### Permission Denied on API Browser

**Cause:** New permission required in Graylog 7.0.

**Solution:** This only affects manual API browser access, not Terraform. Assign "API Browser Reader" role in Graylog UI if needed.

---

## Validation Checklist

Before completing migration:

- [ ] Graylog server upgraded to 7.0+
- [ ] Graylog server accessible and healthy
- [ ] Provider upgraded to v3.0.0+
- [ ] `terraform init -upgrade` completed successfully
- [ ] `terraform plan` runs without errors
- [ ] Test resource creation works
- [ ] Test resource updates work
- [ ] Test resource deletion works
- [ ] All existing resources still managed correctly
- [ ] State is in sync with infrastructure

---

## Getting Help

### Provider Issues

Report issues at: https://github.com/terraform-provider-graylog/terraform-provider-graylog/issues

### Graylog Issues

- Documentation: https://go2docs.graylog.org/current/
- Community: https://community.graylog.org/
- Upgrade Guide: https://go2docs.graylog.org/current/upgrading_graylog/upgrade_to_graylog_7.0.htm

---

## Additional Resources

- [Graylog 7.0 Release Notes](https://go2docs.graylog.org/current/changelogs/changelog.html)
- [Graylog 7.0 Upgrade Guide](https://go2docs.graylog.org/current/upgrading_graylog/upgrade_to_graylog_7.0.htm)
- [Provider Documentation](https://registry.terraform.io/providers/terraform-provider-graylog/graylog/latest/docs)
- [API Changes Documentation](API_MAPPING.md)

---

**Last Updated:** 2025-11-12
**Provider Version:** 3.0.0+
**Target Graylog Version:** 7.0+
