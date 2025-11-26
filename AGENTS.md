# Repository Guidelines

## Project Structure & Module Organization
- Root Go module (`go.mod`); provider code under `graylog/`:
  - `client/`, `config/`, `convert/`, `util/` – shared client/helpers.
  - `provider/`, `resource/`, `datasource/` – Terraform schemas & CRUD; maps in `resource_map.go`, `datasource_map.go`.
  - `testutil/` – helper fixtures.
- CLI entrypoint in `cmd/terraform-provider-graylog/`.
- Documentation in `docs/` (resources, data sources). Graylog REST API docs can be cached in `docs/api-docs`; if missing or stale, run `docs/fetch_api_docs.py` to download fresh copies.
- Examples: `examples/graylog7-e2e/` (current Graylog 7 testing), `examples/v0.12/` (legacy reference for future refactoring).

## Build, Test, and Development Commands
- `make build` – build provider binary into `bin/`.
- `make test` – run unit tests with race detection.
- `make acc-test` – run acceptance tests (requires `TF_ACC=1` and running Graylog instance).
- `make lint` – run golangci-lint.
- `make fmt` – format Go code with gofumpt/gofmt.
- `bin/terraform-dev` – wrapper script to run Terraform with local provider override.

## Local Development Workflow
1. Build the provider: `make build`
2. Use `bin/terraform-dev` instead of `terraform` to test with local binary:
   ```bash
   cd examples/graylog7-e2e
   ../../bin/terraform-dev plan
   ../../bin/terraform-dev apply
   ```
3. To use the registry version instead, use regular `terraform` commands.

## Coding Style & Naming Conventions
- Go 1.22+ idioms; standard Go formatting.
- Run `make fmt` before committing; prefer explicit struct field names in JSON payloads.
- Resource/data source names match Graylog REST nouns; keep consistent prefixes (`graylog_*`).
- Errors: wrap with context via `fmt.Errorf("context: %w", err)`.
- Dashboards (Graylog 7): implemented via Views API with `type=DASHBOARD`; widget mapping/positions/timerange must be JSON strings; saved search datasource returns `search_id`/`state_id` for wiring.

## Testing Guidelines
- Prefer table-driven tests for converters and client logic.
- Add coverage for CRUD read/update paths and 404 handling.
- Acceptance tests should provision minimal Graylog objects then clean up; keep identifiers unique (use `acctest.RandStringFromCharSet` pattern).

## Commit & Pull Request Guidelines
- Commit messages: present tense, concise scope (e.g., `update pipeline resource schema`).
- PRs: include summary of resource/API changes, test evidence (`go test`, acceptance runs), and note any breaking changes to existing schemas.
- Reference related issues/tasks when available.

## Security & Configuration Tips
- Never commit tokens; `.mcp/` and `.tfvars` are ignored by default. Use `graylog.auto.tfvars` locally for credentials.
- TLS: example curl uses `-k` only for self-signed labs—avoid in production.
- Use least-privilege Graylog users for testing; rotate tokens after test runs.
