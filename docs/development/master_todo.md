# Terraform Provider Graylog 7.0 Modernization - Master TODO

**Generated:** 2025-11-12
**Project:** Complete Graylog 7.0 compatibility refactoring
**Status:** Planning Complete, Ready for Implementation

---

## PROGRESS TRACKING

| Phase | Status | Progress |
|-------|--------|----------|
| **Phase 0: Planning & Analysis** | ✅ COMPLETE | 100% |
| **Phase 1: Critical Updates** | ⏳ PENDING | 0% |
| **Phase 2: Verification** | ⏳ PENDING | 0% |
| **Phase 3: Testing** | ⏳ PENDING | 0% |
| **Phase 4: Documentation** | ⏳ PENDING | 0% |
| **Phase 5: Validation** | ⏳ PENDING | 0% |

---

## PHASE 0: PLANNING & ANALYSIS ✅ COMPLETE

- [x] Research Graylog 7.0 API documentation
- [x] Identify breaking changes
- [x] Create component inventory (INVENTORY.md)
- [x] Create API mapping analysis (API_MAPPING.md)
- [x] Generate master TODO list (this file)

---

## PHASE 1: CRITICAL UPDATES (High Priority)

### 1.1 Core Utilities

#### [ ] Task 1.1.1: Create Entity Wrapper Helper Functions
**File:** `graylog/util/util.go`
**Priority:** 🔴 CRITICAL
**Complexity:** LOW
**Time Estimate:** 30 minutes

**Requirements:**
- Add `WrapEntityForCreation` function
- Add `UnwrapEntityResponse` function
- Add comprehensive documentation
- Add unit tests

**Implementation:**
```go
// WrapEntityForCreation wraps entity data in CreateEntityRequest structure
func WrapEntityForCreation(entityData map[string]interface{}) map[string]interface{} {
    return map[string]interface{}{
        "entity": entityData,
        "share_request": map[string]interface{}{
            "selected_grantee_capabilities": map[string]interface{}{},
        },
    }
}
```

**Validation:**
- [ ] Function correctly wraps entity data
- [ ] Handles nil entity data gracefully
- [ ] Unit tests pass
- [ ] Documentation complete

---

#### [ ] Task 1.1.2: Add Field Cleanup Helper
**File:** `graylog/util/util.go`
**Priority:** 🔴 CRITICAL
**Complexity:** LOW
**Time Estimate:** 20 minutes

**Requirements:**
- Add function to remove computed/read-only fields before updates
- Handle common fields: `id`, `created_at`, `creator_user_id`, `last_modified`

**Implementation:**
```go
// RemoveComputedFields removes read-only fields that cause issues in Graylog 7.0
func RemoveComputedFields(data map[string]interface{}) {
    delete(data, "id")
    delete(data, "created_at")
    delete(data, "creator_user_id")
    delete(data, "last_modified")
}
```

**Validation:**
- [ ] Function correctly removes fields
- [ ] Doesn't fail on missing fields
- [ ] Unit tests pass

---

### 1.2 Stream Resources (CONFIRMED - Requires Wrapper)

#### [ ] Task 1.2.1: Update Stream Client Create Method
**File:** `graylog/client/stream/client.go`
**Priority:** 🔴 CRITICAL
**Complexity:** LOW
**Time Estimate:** 15 minutes
**Dependencies:** Task 1.1.1

**Current Code (lines 41-56):**
```go
func (cl Client) Create(
    ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
    if data == nil {
        return nil, nil, errors.New("request body is nil")
    }

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

**New Code:**
```go
func (cl Client) Create(
    ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
    if data == nil {
        return nil, nil, errors.New("request body is nil")
    }

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

**Validation:**
- [ ] Create operation works with Graylog 7.0
- [ ] Response correctly unwrapped
- [ ] Error handling preserved
- [ ] Resource ID correctly extracted

---

#### [ ] Task 1.2.2: Update Stream Resource Create
**File:** `graylog/resource/stream/create.go`
**Priority:** 🔴 CRITICAL
**Complexity:** LOW
**Time Estimate:** 10 minutes
**Dependencies:** Task 1.2.1

**Changes:**
- Verify resource create still works with updated client
- Ensure ID extraction from response works
- Test pause/resume logic after creation

**Validation:**
- [ ] Stream creation succeeds
- [ ] ID correctly set in resource data
- [ ] Pause/resume still functional
- [ ] No regression in existing behavior

---

#### [ ] Task 1.2.3: Update Stream Tests
**File:** `graylog/resource/stream/resource_test.go`
**Priority:** 🔴 CRITICAL
**Complexity:** MEDIUM
**Time Estimate:** 30 minutes
**Dependencies:** Task 1.2.2

**Changes:**
- Update mock HTTP responses to expect wrapped request
- Update expected request bodies
- Add test for wrapper structure
- Verify all CRUD tests pass

**Validation:**
- [ ] All unit tests pass
- [ ] Mock expects correct request structure
- [ ] Create, Read, Update, Delete all tested
- [ ] Edge cases covered

---

### 1.3 Dashboard Resources (CONFIRMED - Requires Wrapper)

#### [ ] Task 1.3.1: Update Dashboard Client Create Method
**File:** `graylog/client/dashboard/client.go`
**Priority:** 🔴 CRITICAL
**Complexity:** LOW
**Time Estimate:** 15 minutes
**Dependencies:** Task 1.1.1

**Changes:** Same pattern as Task 1.2.1 for streams

**Validation:**
- [ ] Create operation works with Graylog 7.0
- [ ] Response correctly handled
- [ ] Error handling preserved

---

#### [ ] Task 1.3.2: Update Dashboard Resource Create
**File:** `graylog/resource/dashboard/create.go`
**Priority:** 🔴 CRITICAL
**Complexity:** LOW
**Time Estimate:** 10 minutes
**Dependencies:** Task 1.3.1

**Validation:**
- [ ] Dashboard creation succeeds
- [ ] ID correctly set
- [ ] No regressions

---

#### [ ] Task 1.3.3: Update Dashboard Tests
**File:** `graylog/resource/dashboard/resource_test.go`
**Priority:** 🔴 CRITICAL
**Complexity:** MEDIUM
**Time Estimate:** 30 minutes
**Dependencies:** Task 1.3.2

**Validation:**
- [ ] All tests pass
- [ ] Wrapper structure verified

---

### 1.4 Event Definition Resources (CONFIRMED - Requires Wrapper)

#### [ ] Task 1.4.1: Update Event Definition Client Create
**File:** `graylog/client/event/definition/client.go`
**Priority:** 🔴 CRITICAL
**Complexity:** MEDIUM (complex entity structure)
**Time Estimate:** 20 minutes
**Dependencies:** Task 1.1.1

**Notes:** Event definitions have complex nested structures - verify all fields properly wrapped

**Validation:**
- [ ] Create operation works
- [ ] Complex nested structure preserved
- [ ] No data loss in wrapping

---

#### [ ] Task 1.4.2: Update Event Definition Resource Create
**File:** `graylog/resource/event/definition/create.go`
**Priority:** 🔴 CRITICAL
**Complexity:** MEDIUM
**Time Estimate:** 15 minutes
**Dependencies:** Task 1.4.1

**Validation:**
- [ ] Event definition creation succeeds
- [ ] All configuration fields preserved
- [ ] ID correctly extracted

---

#### [ ] Task 1.4.3: Update Event Definition Tests
**File:** `graylog/resource/event/definition/resource_test.go`
**Priority:** 🔴 CRITICAL
**Complexity:** HIGH (complex test cases)
**Time Estimate:** 45 minutes
**Dependencies:** Task 1.4.2

**Validation:**
- [ ] All test cases pass
- [ ] Complex configurations tested
- [ ] Edge cases covered

---

### 1.5 Event Notification Resources (CONFIRMED - Requires Wrapper)

#### [ ] Task 1.5.1: Update Event Notification Client Create
**File:** `graylog/client/event/notification/client.go`
**Priority:** 🔴 CRITICAL
**Complexity:** MEDIUM
**Time Estimate:** 15 minutes
**Dependencies:** Task 1.1.1

**Validation:**
- [ ] Create operation works
- [ ] Response correctly handled

---

#### [ ] Task 1.5.2: Update Event Notification Resource Create
**File:** `graylog/resource/event/notification/create.go`
**Priority:** 🔴 CRITICAL
**Complexity:** MEDIUM
**Time Estimate:** 15 minutes
**Dependencies:** Task 1.5.1

**Validation:**
- [ ] Notification creation succeeds
- [ ] Configuration preserved

---

#### [ ] Task 1.5.3: Update Event Notification Tests
**File:** `graylog/resource/event/notification/resource_test.go`
**Priority:** 🔴 CRITICAL
**Complexity:** MEDIUM
**Time Estimate:** 30 minutes
**Dependencies:** Task 1.5.2

**Validation:**
- [ ] All tests pass
- [ ] Different notification types tested

---

## PHASE 2: VERIFICATION (Medium Priority)

### 2.1 Index Set Resource

#### [ ] Task 2.1.1: Verify Index Set Create Endpoint
**File:** `graylog/client/system/indices/indexset/client.go`
**Priority:** 🟡 HIGH
**Complexity:** LOW
**Time Estimate:** 15 minutes

**Action:**
1. Test current implementation against Graylog 7.0
2. Determine if CreateEntityRequest wrapper needed
3. Document findings

**Possible Outcomes:**
- **If wrapper needed:** Follow Task 1.2.1 pattern
- **If not needed:** Mark as compatible, verify unknown fields cleaned

---

#### [ ] Task 2.1.2: Update or Verify Index Set Create
**File:** `graylog/resource/system/indices/indexset/create.go`
**Priority:** 🟡 HIGH
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.1.1

**Validation:**
- [ ] Create operation tested
- [ ] Status documented in validation report

---

#### [ ] Task 2.1.3: Update Index Set Tests
**File:** `graylog/resource/system/indices/indexset/resource_test.go`
**Priority:** 🟡 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 20 minutes
**Dependencies:** Task 2.1.2

**Validation:**
- [ ] Tests updated if needed
- [ ] All tests pass

---

### 2.2 Input Resource

#### [ ] Task 2.2.1: Verify Input Create Endpoint
**File:** `graylog/client/system/input/client.go`
**Priority:** 🟡 HIGH
**Complexity:** LOW
**Time Estimate:** 15 minutes

**Action:** Same as Task 2.1.1 for inputs

---

#### [ ] Task 2.2.2: Update or Verify Input Create
**File:** `graylog/resource/system/input/create.go`
**Priority:** 🟡 HIGH
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.2.1

**Validation:**
- [ ] Create operation tested
- [ ] Complex input attributes handled

---

#### [ ] Task 2.2.3: Update Input Tests
**File:** `graylog/resource/system/input/resource_test.go`
**Priority:** 🟡 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 25 minutes
**Dependencies:** Task 2.2.2

---

### 2.3 Output Resource

#### [ ] Task 2.3.1: Verify Output Create Endpoint
**File:** `graylog/client/system/output/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.3.2: Update or Verify Output Create
**File:** `graylog/resource/system/output/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.3.1

---

#### [ ] Task 2.3.3: Update Output Tests
**File:** `graylog/resource/system/output/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 25 minutes
**Dependencies:** Task 2.3.2

---

### 2.4 Pipeline Resource

#### [ ] Task 2.4.1: Verify Pipeline Create Endpoint
**File:** `graylog/client/system/pipeline/pipeline/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.4.2: Update or Verify Pipeline Create
**File:** `graylog/resource/system/pipeline/pipeline/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.4.1

---

#### [ ] Task 2.4.3: Update Pipeline Tests
**File:** `graylog/resource/system/pipeline/pipeline/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 25 minutes
**Dependencies:** Task 2.4.2

---

### 2.5 Pipeline Rule Resource

#### [ ] Task 2.5.1: Verify Pipeline Rule Create Endpoint
**File:** `graylog/client/system/pipeline/rule/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.5.2: Update or Verify Pipeline Rule Create
**File:** `graylog/resource/system/pipeline/rule/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.5.1

---

#### [ ] Task 2.5.3: Update Pipeline Rule Tests
**File:** `graylog/resource/system/pipeline/rule/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 25 minutes
**Dependencies:** Task 2.5.2

---

### 2.6 Grok Pattern Resource

#### [ ] Task 2.6.1: Verify Grok Pattern Create Endpoint
**File:** `graylog/client/system/grok/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.6.2: Update or Verify Grok Pattern Create
**File:** `graylog/resource/system/grok/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.6.1

---

#### [ ] Task 2.6.3: Update Grok Pattern Tests
**File:** `graylog/resource/system/grok/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 20 minutes
**Dependencies:** Task 2.6.2

---

### 2.7 Role Resource

#### [ ] Task 2.7.1: Verify Role Create Endpoint
**File:** `graylog/client/role/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.7.2: Update or Verify Role Create
**File:** `graylog/resource/role/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.7.1

---

#### [ ] Task 2.7.3: Update Role Tests
**File:** `graylog/resource/role/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 20 minutes
**Dependencies:** Task 2.7.2

---

### 2.8 User Resource

#### [ ] Task 2.8.1: Verify User Create Endpoint
**File:** `graylog/client/user/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.8.2: Update or Verify User Create
**File:** `graylog/resource/user/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.8.1

---

#### [ ] Task 2.8.3: Update User Tests
**File:** `graylog/resource/user/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 25 minutes
**Dependencies:** Task 2.8.2

---

### 2.9 Sidecar Resources

#### [ ] Task 2.9.1: Verify Sidecar Collector Create Endpoint
**File:** `graylog/client/sidecar/collector/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.9.2: Update or Verify Sidecar Collector Create
**File:** `graylog/resource/sidecar/collector/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.9.1

---

#### [ ] Task 2.9.3: Update Sidecar Collector Tests
**File:** `graylog/resource/sidecar/collector/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 20 minutes
**Dependencies:** Task 2.9.2

---

#### [ ] Task 2.9.4: Verify Sidecar Configuration Create Endpoint
**File:** `graylog/client/sidecar/configuration/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

---

#### [ ] Task 2.9.5: Update or Verify Sidecar Configuration Create
**File:** `graylog/resource/sidecar/configuration/create.go`
**Priority:** 🟡 MEDIUM
**Complexity:** VARIES
**Time Estimate:** 10-20 minutes
**Dependencies:** Task 2.9.4

---

#### [ ] Task 2.9.6: Update Sidecar Configuration Tests
**File:** `graylog/resource/sidecar/configuration/resource_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 20 minutes
**Dependencies:** Task 2.9.5

---

## PHASE 3: TESTING & VALIDATION

### 3.1 Unit Tests

#### [ ] Task 3.1.1: Create Wrapper Helper Tests
**File:** `graylog/util/util_test.go`
**Priority:** 🔴 HIGH
**Complexity:** LOW
**Time Estimate:** 20 minutes
**Dependencies:** Task 1.1.1, 1.1.2

**Test Cases:**
- [ ] WrapEntityForCreation with valid data
- [ ] WrapEntityForCreation with nil data
- [ ] WrapEntityForCreation with empty map
- [ ] RemoveComputedFields with all fields present
- [ ] RemoveComputedFields with some fields missing
- [ ] RemoveComputedFields with no fields present

---

#### [ ] Task 3.1.2: Run All Unit Tests
**Priority:** 🔴 HIGH
**Complexity:** LOW
**Time Estimate:** 10 minutes
**Dependencies:** All Phase 1 & 2 tasks

**Command:**
```bash
go test ./graylog/... -v
```

**Validation:**
- [ ] All tests pass
- [ ] No test failures
- [ ] No skipped tests (unless intentional)
- [ ] Coverage report generated

---

### 3.2 Integration Tests

#### [ ] Task 3.2.1: Set Up Graylog 7.0 Test Instance
**Priority:** 🟡 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 1 hour

**Requirements:**
- Docker container with Graylog 7.0
- MongoDB and Elasticsearch/OpenSearch dependencies
- Admin credentials configured
- API accessible on localhost

**Validation:**
- [ ] Graylog 7.0 running
- [ ] API accessible via curl
- [ ] Can authenticate with test credentials

---

#### [ ] Task 3.2.2: Manual Testing - Stream Resource
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 30 minutes
**Dependencies:** Task 3.2.1, Task 1.2.3

**Test Cases:**
- [ ] Create stream with minimal configuration
- [ ] Create stream with full configuration
- [ ] Read stream
- [ ] Update stream
- [ ] Delete stream
- [ ] Pause/resume stream
- [ ] Import stream

---

#### [ ] Task 3.2.3: Manual Testing - Dashboard Resource
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 30 minutes
**Dependencies:** Task 3.2.1, Task 1.3.3

**Test Cases:**
- [ ] Create dashboard
- [ ] Read dashboard
- [ ] Update dashboard
- [ ] Delete dashboard
- [ ] Import dashboard

---

#### [ ] Task 3.2.4: Manual Testing - Event Definition Resource
**Priority:** 🔴 HIGH
**Complexity:** HIGH
**Time Estimate:** 45 minutes
**Dependencies:** Task 3.2.1, Task 1.4.3

**Test Cases:**
- [ ] Create event definition with aggregation
- [ ] Create event definition with correlation
- [ ] Read event definition
- [ ] Update event definition
- [ ] Delete event definition

---

#### [ ] Task 3.2.5: Manual Testing - Event Notification Resource
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 30 minutes
**Dependencies:** Task 3.2.1, Task 1.5.3

**Test Cases:**
- [ ] Create email notification
- [ ] Create HTTP notification
- [ ] Read notification
- [ ] Update notification
- [ ] Delete notification

---

#### [ ] Task 3.2.6: Manual Testing - All Other Resources
**Priority:** 🟡 HIGH
**Complexity:** HIGH
**Time Estimate:** 3-4 hours
**Dependencies:** Task 3.2.1, All Phase 2 tasks

**Resources to Test:**
- [ ] Index Set
- [ ] Input
- [ ] Input Static Fields
- [ ] Extractor
- [ ] Output
- [ ] Pipeline
- [ ] Pipeline Rule
- [ ] Pipeline Connection
- [ ] Grok Pattern
- [ ] LDAP Setting
- [ ] Role
- [ ] User
- [ ] Sidecar
- [ ] Sidecar Collector
- [ ] Sidecar Configuration
- [ ] Stream Rule
- [ ] Stream Output
- [ ] Alarm Callback (deprecated)
- [ ] Alert Condition (deprecated)
- [ ] Dashboard Widget
- [ ] Dashboard Widget Positions

---

### 3.3 Acceptance Tests

#### [ ] Task 3.3.1: Update Provider Test Configuration
**File:** `graylog/provider/provider_test.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

**Changes:**
- Update test provider configuration for Graylog 7.0
- Ensure test credentials work

---

#### [ ] Task 3.3.2: Run Acceptance Tests
**Priority:** 🟡 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 1 hour
**Dependencies:** Task 3.2.1, All code updates

**Command:**
```bash
TF_ACC=1 go test ./graylog/... -v -timeout 120m
```

**Validation:**
- [ ] All acceptance tests pass
- [ ] No failures or errors
- [ ] Test coverage adequate

---

## PHASE 4: DOCUMENTATION

### 4.1 Code Documentation

#### [ ] Task 4.1.1: Update Provider Configuration Docs
**File:** `graylog/provider/provider.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 15 minutes

**Changes:**
- Add comments about Graylog 7.0 compatibility
- Document any new configuration options
- Update version compatibility notes

---

#### [ ] Task 4.1.2: Update Client Documentation
**Files:** All `graylog/client/*/client.go`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 30 minutes

**Changes:**
- Document CreateEntityRequest wrapper usage
- Add Graylog 7.0 compatibility notes
- Update function comments

---

### 4.2 User Documentation

#### [ ] Task 4.2.1: Create Migration Guide
**File:** `docs/MIGRATION_GUIDE_v7.md`
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 1-2 hours

**Contents:**
- Overview of Graylog 7.0 changes
- Breaking changes affecting users
- Required provider version update
- Step-by-step migration instructions
- Common issues and solutions
- Rollback procedures

**Validation:**
- [ ] Guide is comprehensive
- [ ] All breaking changes documented
- [ ] Examples provided
- [ ] Reviewed by second person

---

#### [ ] Task 4.2.2: Update README.md
**File:** `README.md`
**Priority:** 🟡 HIGH
**Complexity:** LOW
**Time Estimate:** 20 minutes

**Changes:**
- Update compatibility matrix
- Add Graylog 7.0 support note
- Link to migration guide
- Update installation instructions

---

#### [ ] Task 4.2.3: Update Resource Documentation
**Files:** `docs/resources/*.md`
**Priority:** 🟡 MEDIUM
**Complexity:** MEDIUM
**Time Estimate:** 2-3 hours

**Changes:**
- Update each resource documentation
- Add Graylog 7.0 compatibility notes
- Update examples if needed
- Mark deprecated resources

**Resources:**
- [ ] graylog_stream.md
- [ ] graylog_stream_rule.md
- [ ] graylog_stream_output.md
- [ ] graylog_dashboard.md
- [ ] graylog_dashboard_widget.md
- [ ] graylog_dashboard_widget_positions.md
- [ ] graylog_event_definition.md
- [ ] graylog_event_notification.md
- [ ] graylog_index_set.md
- [ ] graylog_input.md
- [ ] graylog_input_static_fields.md
- [ ] graylog_extractor.md
- [ ] graylog_output.md
- [ ] graylog_pipeline.md
- [ ] graylog_pipeline_rule.md
- [ ] graylog_pipeline_connection.md
- [ ] graylog_grok_pattern.md
- [ ] graylog_ldap_setting.md
- [ ] graylog_role.md
- [ ] graylog_user.md
- [ ] graylog_sidecars.md
- [ ] graylog_sidecar_collector.md
- [ ] graylog_sidecar_configuration.md
- [ ] graylog_alarm_callback.md (mark deprecated)
- [ ] graylog_alert_condition.md (mark deprecated)

---

#### [ ] Task 4.2.4: Update Data Source Documentation
**Files:** `docs/data-sources/*.md`
**Priority:** 🟡 LOW
**Complexity:** LOW
**Time Estimate:** 30 minutes

**Data Sources:**
- [ ] graylog_dashboard.md
- [ ] graylog_index_set.md
- [ ] graylog_sidecar.md
- [ ] graylog_stream.md

---

#### [ ] Task 4.2.5: Create Examples
**Directory:** `examples/graylog7/`
**Priority:** 🟡 MEDIUM
**Complexity:** HIGH
**Time Estimate:** 3-4 hours

**Examples to Create:**
- [ ] `basic_stream.tf` - Simple stream creation
- [ ] `complete_stream.tf` - Stream with rules, outputs, alerts
- [ ] `dashboard.tf` - Dashboard with widgets
- [ ] `event_system.tf` - Event definitions and notifications
- [ ] `inputs.tf` - Various input types
- [ ] `pipelines.tf` - Pipeline processing
- [ ] `complete_setup.tf` - Full Graylog configuration
- [ ] `migration_example.tf` - Before/after migration
- [ ] `README.md` - Examples documentation

---

### 4.3 API Documentation

#### [ ] Task 4.3.1: Create API Changes Document
**File:** `docs/API_CHANGES_v7.md`
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 30 minutes

**Contents:**
- Copy relevant sections from API_MAPPING.md
- Format for end users
- Focus on what changed and why
- Provide examples

---

#### [ ] Task 4.3.2: Update CHANGELOG
**File:** `CHANGELOG.md`
**Priority:** 🔴 HIGH
**Complexity:** LOW
**Time Estimate:** 30 minutes

**Format:**
```markdown
## [Unreleased]

### Added
- Graylog 7.0 API compatibility
- CreateEntityRequest wrapper support
- [List all new features]

### Changed
- Updated all create endpoints for Graylog 7.0
- [List all changes]

### Deprecated
- graylog_alarm_callback (use graylog_event_notification)
- graylog_alert_condition (use graylog_event_definition)

### Fixed
- Unknown properties validation errors
- [List all fixes]

### Breaking Changes
- Requires Graylog 7.0 or later
- [List all breaking changes]
```

---

## PHASE 5: FINAL VALIDATION & RELEASE

### 5.1 Code Quality

#### [ ] Task 5.1.1: Run Linter
**Priority:** 🟡 HIGH
**Complexity:** LOW
**Time Estimate:** 15 minutes

**Command:**
```bash
golangci-lint run
```

**Validation:**
- [ ] No linter errors
- [ ] All warnings addressed or suppressed with comments

---

#### [ ] Task 5.1.2: Run Go Vet
**Priority:** 🟡 HIGH
**Complexity:** LOW
**Time Estimate:** 10 minutes

**Command:**
```bash
go vet ./...
```

**Validation:**
- [ ] No vet errors
- [ ] All issues resolved

---

#### [ ] Task 5.1.3: Format Code
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 5 minutes

**Command:**
```bash
go fmt ./...
```

**Validation:**
- [ ] All code formatted
- [ ] No formatting inconsistencies

---

#### [ ] Task 5.1.4: Update Dependencies
**Priority:** 🟡 MEDIUM
**Complexity:** LOW
**Time Estimate:** 20 minutes

**Commands:**
```bash
go mod tidy
go mod verify
```

**Validation:**
- [ ] go.mod clean
- [ ] go.sum accurate
- [ ] No unused dependencies

---

### 5.2 Testing

#### [ ] Task 5.2.1: Full Test Suite
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 1 hour

**Commands:**
```bash
go test ./... -v -race -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

**Validation:**
- [ ] All tests pass
- [ ] No race conditions
- [ ] Coverage > 70%
- [ ] Coverage report reviewed

---

#### [ ] Task 5.2.2: Full Acceptance Tests
**Priority:** 🔴 HIGH
**Complexity:** HIGH
**Time Estimate:** 2 hours

**Command:**
```bash
TF_ACC=1 go test ./graylog/... -v -timeout 180m
```

**Validation:**
- [ ] All acceptance tests pass
- [ ] Tested against Graylog 7.0
- [ ] No flaky tests

---

### 5.3 Audit & Documentation

#### [ ] Task 5.3.1: Create VALIDATION_REPORT.md
**File:** `VALIDATION_REPORT.md`
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 1 hour

**Contents:**
```markdown
# Graylog 7.0 Validation Report

## Test Results
- Unit Tests: [PASS/FAIL] - [X/Y tests]
- Integration Tests: [PASS/FAIL] - [X/Y resources]
- Acceptance Tests: [PASS/FAIL] - [X/Y tests]

## Resource Status
- [Resource Name]: ✅ VALIDATED | ❌ FAILED | ⚠️ PARTIAL

## API Compatibility
- Stream: ✅ COMPATIBLE
- Dashboard: ✅ COMPATIBLE
- [etc...]

## Known Issues
- [List any outstanding issues]

## Recommendations
- [List any recommendations]
```

---

#### [ ] Task 5.3.2: Create AUDIT_REPORT.md
**File:** `AUDIT_REPORT.md`
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 1 hour

**Template from instructions:**
```markdown
# Terraform Provider Graylog 7 Audit Report

## Coverage Report
- Total Resources: 25
- Resources Updated: [X]
- Resources Tested: [X]
- Data Sources Updated: 4
- Examples Created: [X]

## API Compatibility
- All endpoints verified: [YES/NO]
- All schemas updated: [YES/NO]
- Authentication updated: [YES/NO]

## Testing Status
- Unit Tests: [X/Y passing]
- Acceptance Tests: [X/Y passing]
- Manual Tests: [Complete/Incomplete]

## Documentation Status
- Resource Docs: [Complete/Incomplete]
- Data Source Docs: [Complete/Incomplete]
- Examples: [Complete/Incomplete]
- Migration Guide: [Complete/Incomplete]

## Outstanding Issues
[List any remaining issues]
```

---

#### [ ] Task 5.3.3: Final README Update
**File:** `README.md`
**Priority:** 🔴 HIGH
**Complexity:** LOW
**Time Estimate:** 15 minutes

**Changes:**
- Remove EOL warning if fully maintained
- Update compatibility matrix
- Add badges for Graylog 7.0 support
- Update quick start guide

---

### 5.4 Release Preparation

#### [ ] Task 5.4.1: Version Bump
**Files:** Various
**Priority:** 🔴 HIGH
**Complexity:** LOW
**Time Estimate:** 10 minutes

**Changes:**
- Update version in relevant files
- Tag as major version if breaking changes
- Update version in documentation

---

#### [ ] Task 5.4.2: Create Release Notes
**File:** `RELEASE_NOTES.md`
**Priority:** 🔴 HIGH
**Complexity:** MEDIUM
**Time Estimate:** 30 minutes

**Contents:**
- Summary of changes
- Upgrade instructions
- Breaking changes
- New features
- Bug fixes
- Known issues
- Contributors

---

#### [ ] Task 5.4.3: Git Commit and Push
**Priority:** 🔴 HIGH
**Complexity:** LOW
**Time Estimate:** 10 minutes

**Commands:**
```bash
git add .
git commit -m "feat: Add Graylog 7.0 compatibility

- Update all create endpoints with CreateEntityRequest wrapper
- Add support for entity sharing in Graylog 7.0
- Fix unknown property validation errors
- Update all tests for Graylog 7.0 API
- Add comprehensive migration guide
- Mark deprecated resources (alarm_callback, alert_condition)

BREAKING CHANGE: Requires Graylog 7.0 or later"

git push -u origin claude/checkout-and-read-011CV4i1kHTeeTJd3hcivfmb
```

**Validation:**
- [ ] All files committed
- [ ] Commit message follows conventions
- [ ] Pushed to correct branch

---

#### [ ] Task 5.4.4: Create Pull Request
**Priority:** 🔴 HIGH
**Complexity:** LOW
**Time Estimate:** 20 minutes

**PR Title:** `feat: Add Graylog 7.0 API compatibility`

**PR Description Template:**
```markdown
## Summary
This PR adds full compatibility with Graylog 7.0 API by implementing the required CreateEntityRequest wrapper for entity creation endpoints and addressing all breaking changes.

## Changes
- ✅ Updated all resource create operations for Graylog 7.0
- ✅ Added CreateEntityRequest wrapper helper functions
- ✅ Fixed unknown property validation issues
- ✅ Updated all client implementations
- ✅ Updated all tests for new API structure
- ✅ Added comprehensive documentation and migration guide
- ✅ Marked deprecated resources

## Testing
- ✅ All unit tests passing ([X/Y])
- ✅ All acceptance tests passing ([X/Y])
- ✅ Manually tested against Graylog 7.0 instance
- ✅ Full integration test suite completed

## Documentation
- ✅ Migration guide created
- ✅ API changes documented
- ✅ Examples updated
- ✅ README updated

## Breaking Changes
- Requires Graylog 7.0 or later
- [List any other breaking changes]

## Checklist
- [ ] Code follows project style guidelines
- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG updated
- [ ] Migration guide complete
- [ ] Examples working
```

---

## SUMMARY

### Total Tasks: 100+

### Time Estimates by Phase:
| Phase | Estimated Time |
|-------|---------------|
| Phase 0 | ✅ Complete |
| Phase 1 | 4-5 hours |
| Phase 2 | 5-7 hours |
| Phase 3 | 8-10 hours |
| Phase 4 | 6-8 hours |
| Phase 5 | 3-4 hours |
| **TOTAL** | **26-34 hours** |

### Critical Path:
1. Phase 1 Tasks 1.1.1 → 1.2.1 → 1.2.2 → 1.2.3 (Streams)
2. Parallel: Tasks 1.3.x (Dashboards), 1.4.x (Events)
3. Phase 2: All verification tasks
4. Phase 3: Testing
5. Phase 4: Documentation
6. Phase 5: Release

### Success Criteria:
- ✅ All resources compatible with Graylog 7.0
- ✅ All tests passing
- ✅ Complete documentation
- ✅ Migration guide available
- ✅ No breaking changes for end users (if possible)

---

**END OF MASTER TODO**
