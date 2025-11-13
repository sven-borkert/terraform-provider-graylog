# Terraform Provider Graylog - Component Inventory

**Generated:** 2025-11-12
**Purpose:** Complete inventory of all provider components for Graylog 7.0 migration

---

## SUMMARY STATISTICS

| Category | Count |
|----------|-------|
| **Total Resources** | 25 |
| **Total Data Sources** | 4 |
| **API Clients** | 24 |
| **Total Go Files** | 252 |

---

## 1. TERRAFORM RESOURCES

### 1.1 Stream Resources (5)
| Resource Name | File Path | Current Endpoints | Status |
|---------------|-----------|-------------------|--------|
| `graylog_stream` | `graylog/resource/stream/` | `POST /streams`<br>`GET /streams/{id}`<br>`PUT /streams/{id}`<br>`DELETE /streams/{id}`<br>`POST /streams/{id}/pause`<br>`POST /streams/{id}/resume` | ⚠️ **REQUIRES UPDATE** - Create endpoint needs CreateEntityRequest wrapper |
| `graylog_stream_rule` | `graylog/resource/stream/rule/` | `POST /streams/{streamId}/rules`<br>`GET /streams/{streamId}/rules/{ruleId}`<br>`PUT /streams/{streamId}/rules/{ruleId}`<br>`DELETE /streams/{streamId}/rules/{ruleId}` | ✅ **LIKELY COMPATIBLE** - No entity creation wrapper needed |
| `graylog_stream_output` | `graylog/resource/stream/output/` | `POST /streams/{streamId}/outputs`<br>`GET /streams/{streamId}/outputs`<br>`DELETE /streams/{streamId}/outputs/{outputId}` | ✅ **LIKELY COMPATIBLE** - Association endpoint |
| `graylog_alarm_callback` | `graylog/resource/stream/alarmcallback/` | `POST /streams/{streamId}/alarmcallbacks`<br>`GET /streams/{streamId}/alarmcallbacks/{id}`<br>`PUT /streams/{streamId}/alarmcallbacks/{id}`<br>`DELETE /streams/{streamId}/alarmcallbacks/{id}` | ⚠️ **DEPRECATED** - Replaced by Event System in Graylog 3.0+ |
| `graylog_alert_condition` | `graylog/resource/stream/alert/condition/` | `POST /streams/{streamId}/alerts/conditions`<br>`GET /streams/{streamId}/alerts/conditions/{id}`<br>`PUT /streams/{streamId}/alerts/conditions/{id}`<br>`DELETE /streams/{streamId}/alerts/conditions/{id}` | ⚠️ **DEPRECATED** - Replaced by Event System in Graylog 3.0+ |

### 1.2 Dashboard Resources (3)
| Resource Name | File Path | Current Endpoints | Status |
|---------------|-----------|-------------------|--------|
| `graylog_dashboard` | `graylog/resource/dashboard/` | `POST /dashboards`<br>`GET /dashboards/{id}`<br>`PUT /dashboards/{id}`<br>`DELETE /dashboards/{id}` | ⚠️ **REQUIRES UPDATE** - Create endpoint needs CreateEntityRequest wrapper |
| `graylog_dashboard_widget` | `graylog/resource/dashboard/widget/` | `POST /dashboards/{dashboardId}/widgets`<br>`GET /dashboards/{dashboardId}/widgets/{widgetId}`<br>`PUT /dashboards/{dashboardId}/widgets/{widgetId}`<br>`DELETE /dashboards/{dashboardId}/widgets/{widgetId}` | ⚠️ **LIKELY COMPATIBLE** - Nested resource, verify if wrapper needed |
| `graylog_dashboard_widget_positions` | `graylog/resource/dashboard/position/` | `PUT /dashboards/{dashboardId}/positions` | ✅ **LIKELY COMPATIBLE** - Update-only endpoint |

### 1.3 Event Resources (2)
| Resource Name | File Path | Current Endpoints | Status |
|---------------|-----------|-------------------|--------|
| `graylog_event_definition` | `graylog/resource/event/definition/` | `POST /events/definitions`<br>`GET /events/definitions/{id}`<br>`PUT /events/definitions/{id}`<br>`DELETE /events/definitions/{id}` | ⚠️ **REQUIRES UPDATE** - Create endpoint needs CreateEntityRequest wrapper |
| `graylog_event_notification` | `graylog/resource/event/notification/` | `POST /events/notifications`<br>`GET /events/notifications/{id}`<br>`PUT /events/notifications/{id}`<br>`DELETE /events/notifications/{id}` | ⚠️ **REQUIRES UPDATE** - Create endpoint needs CreateEntityRequest wrapper |

### 1.4 System Resources (11)
| Resource Name | File Path | Current Endpoints | Status |
|---------------|-----------|-------------------|--------|
| `graylog_index_set` | `graylog/resource/system/indices/indexset/` | `POST /system/indices/index_sets`<br>`GET /system/indices/index_sets/{id}`<br>`PUT /system/indices/index_sets/{id}`<br>`DELETE /system/indices/index_sets/{id}` | ⚠️ **VERIFY** - May need entity wrapper, check Graylog 7 docs |
| `graylog_input` | `graylog/resource/system/input/` | `POST /system/inputs`<br>`GET /system/inputs/{id}`<br>`PUT /system/inputs/{id}`<br>`DELETE /system/inputs/{id}` | ⚠️ **VERIFY** - May need entity wrapper, check Graylog 7 docs |
| `graylog_input_static_fields` | `graylog/resource/system/input/staticfield/` | `POST /system/inputs/{inputId}/staticfields`<br>`GET /system/inputs/{inputId}/staticfields`<br>`DELETE /system/inputs/{inputId}/staticfields/{key}` | ✅ **LIKELY COMPATIBLE** - Association endpoint |
| `graylog_extractor` | `graylog/resource/system/input/extractor/` | `POST /system/inputs/{inputId}/extractors`<br>`GET /system/inputs/{inputId}/extractors/{extractorId}`<br>`PUT /system/inputs/{inputId}/extractors/{extractorId}`<br>`DELETE /system/inputs/{inputId}/extractors/{extractorId}` | ⚠️ **VERIFY** - Nested resource, check if wrapper needed |
| `graylog_grok_pattern` | `graylog/resource/system/grok/` | `POST /system/grok`<br>`GET /system/grok/{id}`<br>`PUT /system/grok/{id}`<br>`DELETE /system/grok/{id}` | ⚠️ **VERIFY** - System resource, check if wrapper needed |
| `graylog_output` | `graylog/resource/system/output/` | `POST /system/outputs`<br>`GET /system/outputs/{id}`<br>`PUT /system/outputs/{id}`<br>`DELETE /system/outputs/{id}` | ⚠️ **VERIFY** - May need entity wrapper |
| `graylog_pipeline` | `graylog/resource/system/pipeline/pipeline/` | `POST /system/pipelines/pipeline`<br>`GET /system/pipelines/pipeline/{id}`<br>`PUT /system/pipelines/pipeline/{id}`<br>`DELETE /system/pipelines/pipeline/{id}` | ⚠️ **VERIFY** - May need entity wrapper |
| `graylog_pipeline_connection` | `graylog/resource/system/pipeline/connection/` | `POST /system/pipelines/connections`<br>`GET /system/pipelines/connections`<br>`PUT /system/pipelines/connections` | ⚠️ **VERIFY** - Connection management endpoint |
| `graylog_pipeline_rule` | `graylog/resource/system/pipeline/rule/` | `POST /system/pipelines/rule`<br>`GET /system/pipelines/rule/{id}`<br>`PUT /system/pipelines/rule/{id}`<br>`DELETE /system/pipelines/rule/{id}` | ⚠️ **VERIFY** - May need entity wrapper |
| `graylog_ldap_setting` | `graylog/resource/system/ldap/setting/` | `PUT /system/ldap/settings/{id}`<br>`GET /system/ldap/settings/{id}`<br>`DELETE /system/ldap/settings/{id}` | ✅ **LIKELY COMPATIBLE** - Update-focused endpoint |
| `graylog_grok_pattern` | `graylog/resource/system/grok/` | `POST /system/grok`<br>`GET /system/grok/{id}`<br>`PUT /system/grok/{id}`<br>`DELETE /system/grok/{id}` | ⚠️ **VERIFY** - System resource |

### 1.5 User Management Resources (2)
| Resource Name | File Path | Current Endpoints | Status |
|---------------|-----------|-------------------|--------|
| `graylog_user` | `graylog/resource/user/` | `POST /users`<br>`GET /users/{username}`<br>`PUT /users/{username}`<br>`DELETE /users/{username}` | ⚠️ **VERIFY** - User management, check if wrapper needed |
| `graylog_role` | `graylog/resource/role/` | `POST /roles`<br>`GET /roles/{roleId}`<br>`PUT /roles/{roleId}`<br>`DELETE /roles/{roleId}` | ⚠️ **VERIFY** - Role management, check if wrapper needed |

### 1.6 Sidecar Resources (3)
| Resource Name | File Path | Current Endpoints | Status |
|---------------|-----------|-------------------|--------|
| `graylog_sidecars` | `graylog/resource/sidecar/` | `GET /sidecars/{id}`<br>`PUT /sidecars/{id}` | ✅ **LIKELY COMPATIBLE** - Read/Update only |
| `graylog_sidecar_collector` | `graylog/resource/sidecar/collector/` | `POST /sidecar/collectors`<br>`GET /sidecar/collectors/{id}`<br>`PUT /sidecar/collectors/{id}`<br>`DELETE /sidecar/collectors/{id}` | ⚠️ **VERIFY** - Check if wrapper needed |
| `graylog_sidecar_configuration` | `graylog/resource/sidecar/configuration/` | `POST /sidecar/configurations`<br>`GET /sidecar/configurations/{id}`<br>`PUT /sidecar/configurations/{id}`<br>`DELETE /sidecar/configurations/{id}` | ⚠️ **VERIFY** - Check if wrapper needed |

---

## 2. TERRAFORM DATA SOURCES

| Data Source Name | File Path | Current Endpoints | Status |
|------------------|-----------|-------------------|--------|
| `graylog_dashboard` | `graylog/datasource/dashboard/` | `GET /dashboards/{id}` | ✅ **COMPATIBLE** - Read-only |
| `graylog_index_set` | `graylog/datasource/system/indices/indexset/` | `GET /system/indices/index_sets/{id}` | ✅ **COMPATIBLE** - Read-only |
| `graylog_sidecar` | `graylog/datasource/sidecar/` | `GET /sidecars/{id}` | ✅ **COMPATIBLE** - Read-only |
| `graylog_stream` | `graylog/datasource/stream/` | `GET /streams/{id}` | ✅ **COMPATIBLE** - Read-only |

---

## 3. API CLIENT IMPLEMENTATIONS

### 3.1 Client Structure
- **Main Client:** `graylog/client/client.go`
- **Architecture:** Single aggregated client with 24 specialized sub-clients
- **HTTP Library:** `github.com/suzuki-shunsuke/go-httpclient`
- **Authentication:** HTTP Basic Auth via `SetRequest` function

### 3.2 Client Listing
| Client Name | File Path | Endpoints Served |
|-------------|-----------|------------------|
| `AlarmCallback` | `graylog/client/stream/alarmcallback/` | Stream alarm callbacks |
| `AlertCondition` | `graylog/client/stream/alert/condition/` | Stream alert conditions |
| `Collector` | `graylog/client/sidecar/collector/` | Sidecar collectors |
| `Dashboard` | `graylog/client/dashboard/` | Dashboards |
| `DashboardWidget` | `graylog/client/dashboard/widget/` | Dashboard widgets |
| `DashboardWidgetPosition` | `graylog/client/dashboard/position/` | Widget positions |
| `EventDefinition` | `graylog/client/event/definition/` | Event definitions |
| `EventNotification` | `graylog/client/event/notification/` | Event notifications |
| `Extractor` | `graylog/client/system/input/extractor/` | Input extractors |
| `Grok` | `graylog/client/system/grok/` | Grok patterns |
| `IndexSet` | `graylog/client/system/indices/indexset/` | Index sets |
| `Input` | `graylog/client/system/input/` | Inputs |
| `InputStaticField` | `graylog/client/system/input/staticfield/` | Input static fields |
| `LDAPSetting` | `graylog/client/system/ldap/setting/` | LDAP settings |
| `Output` | `graylog/client/system/output/` | Outputs |
| `Pipeline` | `graylog/client/system/pipeline/pipeline/` | Pipelines |
| `PipelineConnection` | `graylog/client/system/pipeline/connection/` | Pipeline connections |
| `PipelineRule` | `graylog/client/system/pipeline/rule/` | Pipeline rules |
| `Role` | `graylog/client/role/` | Roles |
| `Sidecar` | `graylog/client/sidecar/` | Sidecars |
| `SidecarConfiguration` | `graylog/client/sidecar/configuration/` | Sidecar configurations |
| `Stream` | `graylog/client/stream/` | Streams |
| `StreamOutput` | `graylog/client/stream/output/` | Stream outputs |
| `StreamRule` | `graylog/client/stream/rule/` | Stream rules |
| `User` | `graylog/client/user/` | Users |
| `View` | `graylog/client/view/` | Views (not implemented) |

---

## 4. HELPER PACKAGES

### 4.1 Configuration
- **Path:** `graylog/config/config.go`
- **Purpose:** Provider configuration structure
- **Fields:**
  - `Endpoint` - Graylog API endpoint URL
  - `AuthName` - Username for authentication
  - `AuthPassword` - Password for authentication
  - `XRequestedBy` - Custom header value (default: "terraform-provider-graylog")
  - `APIVersion` - API version (default: "v3", **unused**)
- **Status:** ⚠️ **NEEDS UPDATE** - LoadAndValidate() is no-op

### 4.2 Conversion Utilities
- **Path:** `graylog/convert/`
- **Files:**
  - `resource_data.go` - Convert between ResourceData and maps
  - `schema.go` - Schema conversion utilities
  - `json.go` - JSON handling utilities
- **Status:** ✅ **FUNCTIONAL** - May need updates for new wrapper structure

### 4.3 Utilities
- **Path:** `graylog/util/util.go`
- **Functions:**
  - `HandleGetResourceError` - Handle 404 gracefully
  - `SchemaDiffSuppressJSONString` - Suppress JSON formatting diffs
  - Import helpers
- **Status:** ✅ **FUNCTIONAL** - Core utilities

### 4.4 Test Utilities
- **Path:** `graylog/testutil/util.go`
- **Purpose:** HTTP mocking with flute library
- **Status:** ⚠️ **NEEDS UPDATE** - Tests need updating for new API structure

---

## 5. TEST FILES

### 5.1 Test Coverage
- **Total Test Files:** 26 resource tests + 4 datasource tests = **30 test files**
- **Framework:** Terraform Plugin SDK Testing + Flute (HTTP mocking)
- **Pattern:** Each resource has `resource_test.go` or `data_source_test.go`

### 5.2 Test File Listing
All test files follow pattern: `<resource_path>/resource_test.go` or `data_source_test.go`

**Status:** ⚠️ **ALL TESTS NEED UPDATE** - Must accommodate CreateEntityRequest structure

---

## 6. EXAMPLE CONFIGURATIONS

### 6.1 Examples Directory
- **Path:** `examples/`
- **Status:** ⚠️ **NEEDS REVIEW** - Check compatibility with Graylog 7

---

## 7. DOCUMENTATION

### 7.1 Documentation Files
- **Path:** `docs/`
- **Status:** ⚠️ **NEEDS UPDATE** - Update for Graylog 7 changes

---

## 8. STATE UPGRADERS

### 8.1 Resources with State Upgraders
Several resources implement state upgraders for backward compatibility:

| Resource | Upgrader Path | Purpose |
|----------|---------------|---------|
| `graylog_stream_rule` | `resource/stream/rule/state_upgrader.go` | Change ID format |
| `graylog_alarm_callback` | `resource/stream/alarmcallback/state_upgrader.go` | Change ID format |
| `graylog_alert_condition` | `resource/stream/alert/condition/state_upgrader.go` | Change ID format |
| `graylog_dashboard_widget` | `resource/dashboard/widget/state_upgrader.go` | Change ID format |
| `graylog_dashboard_widget_positions` | `resource/dashboard/position/state_upgrader.go` | Schema changes |
| `graylog_extractor` | `resource/system/input/extractor/state_upgrader.go` | Change ID format |
| `graylog_input` | `resource/system/input/state_upgrader.go` | Schema changes |
| `graylog_index_set` | `resource/system/indices/indexset/state_upgrader.go` | Schema changes |

**Status:** ✅ **MAINTAIN** - May need new upgraders for Graylog 7 changes

---

## 9. IDENTIFIED ISSUES

### 9.1 Critical Issues
1. ⛔ **CreateEntityRequest Wrapper Required**
   - Affects: streams, dashboards, event definitions, event notifications
   - Impact: **BREAKING CHANGE** - All create operations will fail
   - Priority: **HIGH**

2. ⛔ **Unknown Properties Validation**
   - Graylog 7 now rejects unknown JSON properties
   - Impact: Provider may send deprecated fields causing errors
   - Priority: **HIGH**

3. ⚠️ **Deprecated Resources**
   - `graylog_alarm_callback` - Replaced by Events System
   - `graylog_alert_condition` - Replaced by Events System
   - Impact: Should be marked as deprecated in docs
   - Priority: **MEDIUM**

### 9.2 Medium Priority Issues
1. ⚠️ **API Version Parameter Unused**
   - Config accepts `api_version` but never uses it
   - Impact: Cannot route to version-specific endpoints
   - Priority: **MEDIUM**

2. ⚠️ **No Connection Validation**
   - `LoadAndValidate()` is empty
   - Impact: Errors discovered late at first API call
   - Priority: **MEDIUM**

3. ⚠️ **Context Cancellation**
   - All operations use `context.Background()`
   - Impact: Cannot cancel long operations
   - Priority: **LOW**

---

## 10. MIGRATION PRIORITIES

### 10.1 Phase 1 - Critical (Immediate)
1. ✅ Research Graylog 7 API changes (COMPLETED)
2. ⏳ Update create operations with CreateEntityRequest wrapper
3. ⏳ Remove deprecated fields from all requests
4. ⏳ Update API client to handle new structure

### 10.2 Phase 2 - High Priority
1. ⏳ Update all resource schemas
2. ⏳ Update all tests with new structure
3. ⏳ Add proper validation in LoadAndValidate()

### 10.3 Phase 3 - Documentation
1. ⏳ Create migration guide
2. ⏳ Update all resource documentation
3. ⏳ Update examples
4. ⏳ Mark deprecated resources

### 10.4 Phase 4 - Enhancements
1. ⏳ Implement API version routing
2. ⏳ Add context cancellation support
3. ⏳ Implement retry logic

---

## 11. NEXT STEPS

1. ✅ Complete API mapping analysis (API_MAPPING.md)
2. ✅ Generate detailed TODO list (MASTER_TODO.md)
3. ⏳ Begin systematic refactoring starting with critical resources
4. ⏳ Update tests incrementally
5. ⏳ Create comprehensive documentation

---

**END OF INVENTORY**
