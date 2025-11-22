# Refactoring Status (Graylog 7 Provider)

Date: 2025-11-21

## Completed
- Upgraded tooling: Go 1.22, terraform-plugin-sdk v2.38.1; Makefile with build/dev-install/fmt/test/acc-test/venv targets.
- MCP setup documented and working via local proxy; credentials kept local/ignored.
- Index Set resource/data source aligned to Graylog 7:
  - Supports data_tiering, field_restrictions, index_set_template_id, default removal of computed/unsupported fields.
  - Create/Update fixed (includes id on update; uses template ID).
- Added template client, data sources (`graylog_index_set_template`, `graylog_index_set_templates`), and resource for custom templates.
- Example `examples/graylog7-e2e` provisions an index set using built-in “30 Days Hot” template via data source; dev overrides workflow documented.
- Cleanup: removed legacy `example-local-usage/`, scripts/, build-local.sh, local-registry; `.venv` and `docs/api-docs` gitignored; API fetch script added.
- OpenSearch health fixed (replicas set to 0 for top_queries-*; cluster green).

## In Progress / To Do
- Refactor remaining resources for Graylog 7: inputs, streams, pipelines, outputs, events/notifications, sidecars, users/roles, LDAP.
- Add acceptance/E2E coverage per resource (create/read/update/delete) using the Graylog 7 test instance; expand `examples/graylog7-e2e`.
- Consider migration notes per resource (breaking fields removed/renamed).
- Update docs/index to link new guides (MCP, templates, API docs).
- Add tests for template data sources/resource; add index set update tests with template IDs.
- Review provider default schemas for any deprecated fields (alerts/alarm callbacks) in Graylog 7.

## How to Test Locally (current workflow)
- Dev overrides in `~/.terraformrc` point to `bin/`.
- Build: `make build`
- Example: `cd examples/graylog7-e2e && TF_CLI_CONFIG_FILE=$HOME/.terraformrc terraform apply -auto-approve`
- Template discovery: outputs `builtin_templates` from `graylog_index_set_templates`.

## Notes
- Graylog built-in template IDs (current cluster): 7d `692067a9371785137c9b2b73`, 14d `692067a9371785137c9b2b74`, 30d `692067a9371785137c9b2b75`.
- Index template “type” string is no longer used; use template IDs.
