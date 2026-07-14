# Graylog 7.0 API Mapping Analysis

**Generated:** 2025-11-12
**Purpose:** Map current provider API usage to Graylog 7.0 requirements
**Source:** Graylog 7.0 upgrade documentation and API browser

---

## TABLE OF CONTENTS
1. [Breaking Changes Summary](#1-breaking-changes-summary)
2. [CreateEntityRequest Wrapper](#2-createentityrequest-wrapper)
3. [Endpoint Mapping](#3-endpoint-mapping)
4. [Authentication Changes](#4-authentication-changes)
5. [Request/Response Changes](#5-requestresponse-changes)
6. [Deprecated Endpoints](#6-deprecated-endpoints)
7. [Required Updates by Resource](#7-required-updates-by-resource)

---

## 1. BREAKING CHANGES SUMMARY

### 1.1 Critical Breaking Changes

| Change | Impact | Affected Resources |
|--------|--------|-------------------|
| **CreateEntityRequest wrapper required** | All entity creation must wrap data in `{entity: {...}, share_request: {...}}` | 🔴 HIGH - 11+ resources |
| **Unknown properties rejected** | APIs reject unknown JSON fields | 🔴 HIGH - All resources |
| **URL allowlist rename** | `/system/urlwhitelist` → `/system/urlallowlist` | 🟡 MEDIUM - Not used by provider |
| **Data lake endpoint rename** | `/data_warehouse/` → `/data_lake/` | 🟡 MEDIUM - Not used by provider |
| **API Browser permission** | Requires `api_browser:read` permission | 🟢 LOW - Documentation only |

### 1.2 Affected Entity Types
Per Graylog 7.0 documentation, the following entity types **REQUIRE** CreateEntityRequest wrapper:
- ✅ Searches
- ✅ Dashboards
- ✅ Filters
- ✅ Reports
- ✅ Event Definitions
- ✅ Streams
- ✅ Notifications
- ✅ Sigma Rules
- ✅ Event Procedures
- ✅ Steps
- ✅ Content Packs
- ✅ Teams
- ✅ Illuminate Packs

---

## 2. CREATEENTITYREQUEST WRAPPER

### 2.1 Structure Change

#### OLD FORMAT (Pre-Graylog 7.0)
```json
{
  "title": "My Stream",
  "description": "An example stream",
  "index_set_id": "65b7ba138cdb8c534a953fef",
  "remove_matches_from_default_stream": false,
  "matching_type": "AND"
}
```

#### NEW FORMAT (Graylog 7.0+)
```json
{
  "entity": {
    "title": "My Stream",
    "description": "An example stream",
    "index_set_id": "65b7ba138cdb8c534a953fef",
    "remove_matches_from_default_stream": false,
    "matching_type": "AND"
  },
  "share_request": {
    "selected_grantee_capabilities": {
      "grn::::user:admin": "own"
    }
  }
}
```

### 2.2 Share Request Structure

The `share_request` object controls entity permissions:

```json
{
  "share_request": {
    "selected_grantee_capabilities": {
      "<GRN>": "<capability>"
    }
  }
}
```

**Capabilities:**
- `"own"` - Full ownership
- `"manage"` - Management permissions
- `"view"` - Read-only access

**GRN Format:** `grn::::type:identifier`
- Example: `grn::::user:admin`
- Example: `grn::::team:developers`

### 2.3 Implementation Requirements

**For Create Operations:**
```go
// OLD CODE (will fail in Graylog 7.0)
data := map[string]interface{}{
    "title": "My Stream",
    "index_set_id": "...",
}
resp, err := cl.Stream.Create(ctx, data)

// NEW CODE (Graylog 7.0 compatible)
requestData := map[string]interface{}{
    "entity": map[string]interface{}{
        "title": "My Stream",
        "index_set_id": "...",
    },
    "share_request": map[string]interface{}{
        "selected_grantee_capabilities": map[string]interface{}{},
    },
}
resp, err := cl.Stream.Create(ctx, requestData)
```

**For Update Operations:**
- ✅ NO CHANGE - Update endpoints continue to accept entity data directly
- ✅ NO CHANGE - Delete endpoints unchanged
- ✅ NO CHANGE - Read endpoints unchanged

---

## 3. ENDPOINT MAPPING

### 3.1 Stream Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /streams` | `POST /streams` | ⚠️ REQUIRES WRAPPER | Must use CreateEntityRequest |
| Read | `GET /streams/{id}` | `GET /streams/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /streams/{id}` | `PUT /streams/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /streams/{id}` | `DELETE /streams/{id}` | ✅ COMPATIBLE | No changes |
| List | `GET /streams` | `GET /streams` | ✅ COMPATIBLE | No changes |
| Pause | `POST /streams/{id}/pause` | `POST /streams/{id}/pause` | ✅ COMPATIBLE | No changes |
| Resume | `POST /streams/{id}/resume` | `POST /streams/{id}/resume` | ✅ COMPATIBLE | No changes |

**Required Changes:**
- ✅ `resource/stream/create.go` - Wrap entity in CreateEntityRequest
- ✅ `client/stream/client.go` - Update Create method to accept wrapper

### 3.2 Dashboard Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /dashboards` | `POST /dashboards` | ⚠️ REQUIRES WRAPPER | Must use CreateEntityRequest |
| Read | `GET /dashboards/{id}` | `GET /dashboards/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /dashboards/{id}` | `PUT /dashboards/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /dashboards/{id}` | `DELETE /dashboards/{id}` | ✅ COMPATIBLE | No changes |

**Required Changes:**
- ✅ `resource/dashboard/create.go` - Wrap entity in CreateEntityRequest
- ✅ `client/dashboard/client.go` - Update Create method

### 3.3 Event Definition Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /events/definitions` | `POST /events/definitions` | ⚠️ REQUIRES WRAPPER | Must use CreateEntityRequest |
| Read | `GET /events/definitions/{id}` | `GET /events/definitions/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /events/definitions/{id}` | `PUT /events/definitions/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /events/definitions/{id}` | `DELETE /events/definitions/{id}` | ✅ COMPATIBLE | No changes |

**Required Changes:**
- ✅ `resource/event/definition/create.go` - Wrap entity in CreateEntityRequest
- ✅ `client/event/definition/client.go` - Update Create method

### 3.4 Event Notification Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /events/notifications` | `POST /events/notifications` | ⚠️ REQUIRES WRAPPER | Must use CreateEntityRequest |
| Read | `GET /events/notifications/{id}` | `GET /events/notifications/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /events/notifications/{id}` | `PUT /events/notifications/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /events/notifications/{id}` | `DELETE /events/notifications/{id}` | ✅ COMPATIBLE | No changes |

**Required Changes:**
- ✅ `resource/event/notification/create.go` - Wrap entity in CreateEntityRequest
- ✅ `client/event/notification/client.go` - Update Create method

### 3.5 Index Set Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /system/indices/index_sets` | `POST /system/indices/index_sets` | ❓ VERIFY | Check if wrapper needed |
| Read | `GET /system/indices/index_sets/{id}` | `GET /system/indices/index_sets/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /system/indices/index_sets/{id}` | `PUT /system/indices/index_sets/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /system/indices/index_sets/{id}` | `DELETE /system/indices/index_sets/{id}` | ✅ COMPATIBLE | No changes |

**Action Required:**
- 🔍 Test against Graylog 7 API to verify if CreateEntityRequest wrapper needed
- 🔍 Check official API browser documentation

### 3.6 Input Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /system/inputs` | `POST /system/inputs` | ❓ VERIFY | Check if wrapper needed |
| Read | `GET /system/inputs/{id}` | `GET /system/inputs/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /system/inputs/{id}` | `PUT /system/inputs/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /system/inputs/{id}` | `DELETE /system/inputs/{id}` | ✅ COMPATIBLE | No changes |

**Action Required:**
- 🔍 Test against Graylog 7 API to verify if CreateEntityRequest wrapper needed

### 3.7 Pipeline Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /system/pipelines/pipeline` | `POST /system/pipelines/pipeline` | ❓ VERIFY | Check if wrapper needed |
| Read | `GET /system/pipelines/pipeline/{id}` | `GET /system/pipelines/pipeline/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /system/pipelines/pipeline/{id}` | `PUT /system/pipelines/pipeline/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /system/pipelines/pipeline/{id}` | `DELETE /system/pipelines/pipeline/{id}` | ✅ COMPATIBLE | No changes |

**Action Required:**
- 🔍 Test against Graylog 7 API to verify if CreateEntityRequest wrapper needed

### 3.8 Output Endpoints

| Operation | Current Endpoint | Graylog 7.0 Endpoint | Status | Notes |
|-----------|------------------|----------------------|--------|-------|
| Create | `POST /system/outputs` | `POST /system/outputs` | ❓ VERIFY | Check if wrapper needed |
| Read | `GET /system/outputs/{id}` | `GET /system/outputs/{id}` | ✅ COMPATIBLE | No changes |
| Update | `PUT /system/outputs/{id}` | `PUT /system/outputs/{id}` | ✅ COMPATIBLE | No wrapper needed |
| Delete | `DELETE /system/outputs/{id}` | `DELETE /system/outputs/{id}` | ✅ COMPATIBLE | No changes |

**Action Required:**
- 🔍 Test against Graylog 7 API to verify if CreateEntityRequest wrapper needed

### 3.9 User & Role Endpoints

| Resource | Operation | Current Endpoint | Graylog 7.0 Endpoint | Status |
|----------|-----------|------------------|----------------------|--------|
| User | Create | `POST /users` | `POST /users` | ❓ VERIFY |
| User | Read | `GET /users/{username}` | `GET /users/{username}` | ✅ COMPATIBLE |
| User | Update | `PUT /users/{username}` | `PUT /users/{username}` | ✅ COMPATIBLE |
| User | Delete | `DELETE /users/{username}` | `DELETE /users/{username}` | ✅ COMPATIBLE |
| Role | Create | `POST /roles` | `POST /roles` | ❓ VERIFY |
| Role | Read | `GET /roles/{id}` | `GET /roles/{id}` | ✅ COMPATIBLE |
| Role | Update | `PUT /roles/{id}` | `PUT /roles/{id}` | ✅ COMPATIBLE |
| Role | Delete | `DELETE /roles/{id}` | `DELETE /roles/{id}` | ✅ COMPATIBLE |

**Action Required:**
- 🔍 Test user/role creation to verify if wrapper needed

### 3.10 Sidecar Endpoints

| Resource | Operation | Current Endpoint | Graylog 7.0 Endpoint | Status |
|----------|-----------|------------------|----------------------|--------|
| Sidecar | Read | `GET /sidecars/{id}` | `GET /sidecars/{id}` | ✅ COMPATIBLE |
| Sidecar | Update | `PUT /sidecars/{id}` | `PUT /sidecars/{id}` | ✅ COMPATIBLE |
| Collector | Create | `POST /sidecar/collectors` | `POST /sidecar/collectors` | ❓ VERIFY |
| Collector | Read | `GET /sidecar/collectors/{id}` | `GET /sidecar/collectors/{id}` | ✅ COMPATIBLE |
| Collector | Update | `PUT /sidecar/collectors/{id}` | `PUT /sidecar/collectors/{id}` | ✅ COMPATIBLE |
| Collector | Delete | `DELETE /sidecar/collectors/{id}` | `DELETE /sidecar/collectors/{id}` | ✅ COMPATIBLE |
| Configuration | Create | `POST /sidecar/configurations` | `POST /sidecar/configurations` | ❓ VERIFY |
| Configuration | Read | `GET /sidecar/configurations/{id}` | `GET /sidecar/configurations/{id}` | ✅ COMPATIBLE |
| Configuration | Update | `PUT /sidecar/configurations/{id}` | `PUT /sidecar/configurations/{id}` | ✅ COMPATIBLE |
| Configuration | Delete | `DELETE /sidecar/configurations/{id}` | `DELETE /sidecar/configurations/{id}` | ✅ COMPATIBLE |

---

## 4. AUTHENTICATION CHANGES

### 4.1 Current Authentication
```go
// graylog/client/client.go:74-78
httpClient.SetRequest = func(req *http.Request) error {
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Requested-By", xRequestedBy)
    req.SetBasicAuth(cfg.AuthName, cfg.AuthPassword)
    return nil
}
```

### 4.2 Graylog 7.0 Authentication
- ✅ **COMPATIBLE** - HTTP Basic Auth still supported
- ✅ **COMPATIBLE** - API tokens still supported
- ✅ **COMPATIBLE** - Session authentication still supported

### 4.3 New Permission: API Browser Access
- **Requirement:** Users need `api_browser:read` permission to access API browser
- **Impact:** Documentation only - provider doesn't use API browser
- **Action:** None required for provider code

---

## 5. REQUEST/RESPONSE CHANGES

### 5.1 Unknown Properties Handling

#### OLD BEHAVIOR (Pre-Graylog 7.0)
- Unknown JSON properties were **silently ignored**
- Allowed sending extra/deprecated fields without errors

#### NEW BEHAVIOR (Graylog 7.0+)
- Unknown JSON properties are **rejected with error**
- Must remove all deprecated/unsupported fields before sending

### 5.2 Fields to Remove Before Requests

**Computed/Read-Only Fields** (must delete before PUT requests):
```go
// These fields are returned by GET but rejected by PUT
delete(data, "id")
delete(data, "created_at")
delete(data, "creator_user_id")
delete(data, "last_modified")
```

**Current Implementation:**
- ✅ Stream resource deletes `creator_user_id` and `created_at` (update.go:24-25)
- ❓ Other resources may need similar cleanup

### 5.3 Response Structure

**No changes to response structure** - responses continue to return entity data directly:

```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "My Stream",
  "description": "...",
  "created_at": "2024-01-01T00:00:00.000Z",
  ...
}
```

---

## 6. DEPRECATED ENDPOINTS

### 6.1 Renamed Endpoints

| Old Endpoint | New Endpoint | Provider Usage |
|--------------|--------------|----------------|
| `GET /system/urlwhitelist` | `GET /system/urlallowlist` | ❌ Not used |
| `PUT /system/urlwhitelist` | `PUT /system/urlallowlist` | ❌ Not used |
| `POST /system/urlwhitelist/check` | `POST /system/urlallowlist/check` | ❌ Not used |
| `POST /system/urlwhitelist/generate_regex` | `POST /system/urlallowlist/generate_regex` | ❌ Not used |
| `/api/plugins/.../datawarehouse/data_warehouse/...` | `/api/plugins/.../datalake/data_lake/...` | ❌ Not used |

**Impact:** ✅ NONE - Provider doesn't use these endpoints

### 6.2 Removed Endpoints

| Endpoint | Reason | Provider Usage |
|----------|--------|----------------|
| `/api/plugins/.../assets/history/...` | Migrated to Asset History index set/stream | ❌ Not used |

**Impact:** ✅ NONE - Provider doesn't use these endpoints

### 6.3 Removed Resources (No Longer Supported in Graylog 7)

The following legacy resources were removed from the provider because their underlying Graylog 7 REST APIs no longer exist:

| Resource | Replacement |
|----------|-------------|
| `graylog_alarm_callback` | `graylog_event_notification` |
| `graylog_alert_condition` | `graylog_event_definition` |
| `graylog_dashboard_widget` | Managed inline via the `graylog_dashboard` `state` block |
| `graylog_dashboard_widget_positions` | Managed inline via the `graylog_dashboard` `state` block |
| `graylog_ldap_setting` | Authentication-services API |

**Note:** These resources (and the `graylog_dashboard_widget` data source) were removed in Graylog 7. Use the Events System for alerting and manage dashboard widgets inline through `graylog_dashboard`.

---

## 7. REQUIRED UPDATES BY RESOURCE

### 7.1 CONFIRMED - Requires CreateEntityRequest Wrapper

| Resource | File | Priority | Complexity |
|----------|------|----------|------------|
| `graylog_stream` | `resource/stream/create.go` | 🔴 HIGH | MEDIUM |
| `graylog_dashboard` | `resource/dashboard/create.go` | 🔴 HIGH | MEDIUM |
| `graylog_event_definition` | `resource/event/definition/create.go` | 🔴 HIGH | HIGH |
| `graylog_event_notification` | `resource/event/notification/create.go` | 🔴 HIGH | MEDIUM |

### 7.2 VERIFY - May Require Wrapper

| Resource | File | Priority | Action |
|----------|------|----------|--------|
| `graylog_index_set` | `resource/system/indices/indexset/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_input` | `resource/system/input/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_output` | `resource/system/output/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_pipeline` | `resource/system/pipeline/pipeline/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_pipeline_rule` | `resource/system/pipeline/rule/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_grok_pattern` | `resource/system/grok/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_role` | `resource/role/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_user` | `resource/user/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_sidecar_collector` | `resource/sidecar/collector/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |
| `graylog_sidecar_configuration` | `resource/sidecar/configuration/create.go` | 🟡 MEDIUM | Test against Graylog 7 API |

### 7.3 COMPATIBLE - No Changes Required

| Resource | Reason |
|----------|--------|
| `graylog_stream_rule` | No entity creation, just rules |
| `graylog_stream_output` | Association endpoint, not entity creation |
| `graylog_extractor` | Nested resource under input |
| `graylog_input_static_fields` | Association endpoint |
| `graylog_pipeline_connection` | Connection management |
| `graylog_sidecars` | Read/Update only, no create |

### 7.4 ALL Data Sources

| Data Source | Status | Reason |
|-------------|--------|--------|
| `graylog_dashboard` | ✅ COMPATIBLE | Read-only |
| `graylog_index_set` | ✅ COMPATIBLE | Read-only |
| `graylog_sidecar` | ✅ COMPATIBLE | Read-only |
| `graylog_stream` | ✅ COMPATIBLE | Read-only |

---

## 8. IMPLEMENTATION STRATEGY

### 8.1 Wrapper Helper Function

Create a helper function in `graylog/util/util.go`:

```go
// WrapEntityForCreation wraps entity data in CreateEntityRequest structure
// required by Graylog 7.0+ for entity creation endpoints
func WrapEntityForCreation(entityData map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "entity": entityData,
        "share_request": map[string]interface{}{
            "selected_grantee_capabilities": map[string]interface{}{},
        },
    }
}

// UnwrapEntityResponse extracts entity from CreateEntityRequest response
func UnwrapEntityResponse(response map[string]interface{}) map[string]interface{}{
    // Graylog 7.0 returns the created entity directly, not wrapped
    // This function exists for consistency and future-proofing
    return response
}
```

### 8.2 Client Update Pattern

**Before (Current):**
```go
func (cl Client) Create(ctx context.Context, data map[string]interface{}) (map[string]interface{}, *http.Response, error) {
    body := map[string]interface{}{}
    resp, err := cl.Client.Call(ctx, httpclient.CallParams{
        Method:       "POST",
        Path:         "/streams",
        RequestBody:  data,
        ResponseBody: &body,
    })
    return body, resp, err
}
```

**After (Graylog 7.0):**
```go
func (cl Client) Create(ctx context.Context, data map[string]interface{}) (map[string]interface{}, *http.Response, error) {
    // Wrap entity for Graylog 7.0 CreateEntityRequest
    requestData := map[string]interface{}{
        "entity": data,
        "share_request": map[string]interface{}{
            "selected_grantee_capabilities": map[string]interface{}{},
        },
    }

    body := map[string]interface{}{}
    resp, err := cl.Client.Call(ctx, httpclient.CallParams{
        Method:       "POST",
        Path:         "/streams",
        RequestBody:  requestData,
        ResponseBody: &body,
    })
    return body, resp, err
}
```

### 8.3 Resource Update Pattern

**Resource create.go files require no changes** if client is updated properly:

```go
// resource/stream/create.go - NO CHANGES NEEDED
func create(d *schema.ResourceData, m interface{}) error {
    ctx := context.Background()
    cl, err := client.New(m)
    if err != nil {
        return err
    }
    data, err := getDataFromResourceData(d)
    if err != nil {
        return err
    }

    // Client.Create now handles wrapping internally
    stream, _, err := cl.Stream.Create(ctx, data)
    if err != nil {
        return fmt.Errorf("failed to create a stream: %w", err)
    }
    d.SetId(stream["id"].(string))
    return nil
}
```

---

**END OF API MAPPING**
