# Graylog 7 Terraform Provider Refactor Plan

## Goals
1. Update the provider to target Graylog 7 only; drop backward compatibility concerns.
2. Align every resource and data source with the Graylog 7 REST API (including request/response bodies, IDs, and error handling).
3. Provide a local test harness (Terraform example project + scripts) to exercise all resources against a live Graylog 7 instance without publishing to a registry.
4. Document the refactor and testing flow for future contributors.

## Current Surface (from `resource_map.go` / `datasource_map.go`)
- Resources: alarm_callback, alert_condition, dashboard, dashboard_widget, dashboard_widget_positions, event_definition, event_notification, extractor, grok_pattern, index_set, input, input_static_fields, ldap_setting, output, pipeline, pipeline_connection, pipeline_rule, role, sidecars, sidecar_collector, sidecar_configuration, stream, stream_output, stream_rule, user.
- Data sources: dashboard, index_set, sidecar, stream.
- Potentially disabled: view (TODO in code).

## Current Environment
- Graylog host/version: `https://graylog.internal.borkert.net` running Graylog 7.0.0 (via MCP `get_system_status`).
- MCP access: read-only through `mcp-proxy` bridge (HTTP → STDIO). Use MCP only for inspection; all mutations must go through Terraform.
- Known index sets (MCP `list_index_sets`): Default (`graylog`, id `691e31ee945ebfd39f28bb94`), Events (`gl-events`, id `691e31ef945ebfd39f28bca4`), System Events (`gl-system-events`, id `691e31ef945ebfd39f28bca7`).
- SSH available (if needed) at `root@graylog.internal.borkert.net` for server-side API inspection/logs.

## Work Phases
### Phase 0 — Baseline & Tooling
- Confirm Go version, module deps, terraform-plugin-sdk/v2 level; upgrade if needed for TF 1.5+ compatibility.
- Map build entry points (`cmd/`, `build-local.sh`) and ensure reproducible local build (no registry publish).
- Keep local provider override flow documented (dev overrides in `~/.terraformrc` / `terraform.rc`).
- Maintain `.venv` with `mcp-proxy` for read-only checks.

### Phase 1 — API Recon for Graylog 7
- For each resource, map Graylog 7 REST endpoints, required/optional fields, and response shapes. Highlight changes vs prior versions (alerts → events, sidecar path changes, etc.).
- Use MCP read-only calls to capture live shapes/IDs; supplement with Graylog 7 docs and, if needed, direct API/SSH inspection.
- Record required permissions per endpoint; flag immutable fields and server-side defaults.
- Produce a compatibility matrix per resource (old field → new field / removed / renamed).

### Phase 2 — Client/Foundation Updates
- Update client layer (`graylog/client`, `config`, `convert`, `util`) to match Graylog 7 endpoints, media types, paging, and error semantics.
- Add consistent request builders, context/timeouts, and structured errors.
- Remove compatibility branches for older API versions.

### Phase 3 — Resource/Data Source Refactors
- Per resource: rewrite CRUD to match Graylog 7 payloads; verify ID handling; ensure Read detects not-found and handles immutable fields; align schema types and validations.
- Streamline nested schemas (e.g., dashboard widgets/positions, pipeline connections, sidecar collectors/configs).
- Decide on view support (implement or remove TODO) based on Graylog 7 capabilities.

### Phase 4 — Acceptance & Example Suite
- Create `tests/e2e/` Terraform project (or `examples/graylog7-e2e/`) that provisions: index sets → inputs (+extractors/static fields) → pipelines/rules/connections → streams/rules/outputs → dashboards/widgets/positions → events/notifications → sidecars/collectors/configs → users/roles/LDAP setting.
- Provide `graylog.auto.tfvars.example` (URL, admin creds, optional tokens) and ignore real `graylog.auto.tfvars`.
- Add thin shell helper to build provider locally and wire Terraform dev overrides.
- During e2e runs, use MCP read-only tools (`get_system_status`, `list_index_sets`, `list_streams`, etc.) to verify resulting state without direct mutations.

### Phase 5 — Local Build & Provider Install Flow
- Standardize: `make build` or `scripts/build-local.sh` to `./bin/terraform-provider-graylog_v0.0.0`.
- Terraform CLI config snippet (`provider_installation { dev_overrides { "local/graylog" = "<repo>/bin" } direct {} }`) documented in README and example project so provider is used in-place.

### Phase 6 — Documentation & Migration Notes
- Update README/docs per resource with Graylog 7 field coverage and examples.
- Add changelog entry summarizing Graylog 7-only target and breaking changes.

### Phase 7 — MCP-Assisted Verification (once available)
- Use MCP connection to interrogate live Graylog 7 while running the example project, iterating resource by resource.

## Deliverables
- Updated codebase aligned with Graylog 7.
- Example Terraform project + tfvars template + build/test scripts.
- Documentation update and checklist of verified resources.
