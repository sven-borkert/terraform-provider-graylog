# Repository Guidelines

## Project Structure & Module Organization
- Root Go module (`go.mod`); provider code under `graylog/`:
  - `client/`, `config/`, `convert/`, `util/` – shared client/helpers.
  - `provider/`, `resource/`, `datasource/` – Terraform schemas & CRUD; maps in `resource_map.go`, `datasource_map.go`.
  - `testutil/` – helper fixtures.
- CLI entrypoints in `cmd/`.
- Docs in `docs/` (see `refactoring-plan-graylog7.md`, `mcp-usage.md`). Graylog REST API docs are cached in `docs/api-docs`; if they are missing or stale, run `docs/fetch_api_docs.py` to download fresh copies.
- Examples: `examples/` and `example-local-usage/`.
- Scripts: `build-local.sh`, `scripts/` (lint/test helpers).

## Build, Test, and Development Commands
- `go test ./...` – run unit tests.
- `./build-local.sh` – build provider binary into `./terraform-provider-graylog` (used for local Terraform testing).
- `go test ./graylog/provider -run TestAcc` (once acceptance tests are added) – targeted acceptance/API tests.
- `terraform init/plan/apply` inside example projects with local provider override (see docs).

## Coding Style & Naming Conventions
- Go 1.22+ idioms; 2-space indentation (Go fmt defaults).
- Run `gofmt`/`goimports` on change; prefer explicit struct field names in JSON payloads.
- Resource/data source names match Graylog REST nouns; keep consistent prefixes (`graylog_*`).
- Errors: wrap with context via `fmt.Errorf("context: %w", err)`.

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
- MCP usage policy: MCP is read-only; use it only to inspect Graylog state. All mutations must go through Terraform resources/providers, not MCP calls.
