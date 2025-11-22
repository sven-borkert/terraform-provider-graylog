# Graylog MCP with Codex CLI (Current Setup)

## Context & Limitation
- Graylog MCP is an HTTP JSON-RPC endpoint (`/api/mcp`); Codex CLI only speaks STDIO MCP.
- We bridge HTTP → STDIO using a local proxy (`mcp-proxy`).

## What’s Configured Now
- Graylog host: `https://graylog.internal.borkert.net/api/mcp`
- Auth: Basic header built from token `1d0g4tgpg5h3v3jkcdr9mr1rn9ner9ilc7iuosm8hef054m875df:token` → `MWQwZzR0Z3BnNWgzdjNqa2NkcjltcjFybjluZXI5aWxjN2l1b3NtOGhlZjA1NG04NzVkZjp0b2tlbg==`
- Local venv with proxy: `.venv/` contains `mcp-proxy` (installed via `python3 -m venv .venv && .venv/bin/pip install mcp-proxy`).
- Codex config (`~/.codex/config.toml`):
  ```toml
  [features]
  rmcp_client = true  # replaces deprecated experimental_use_rmcp_client

  [mcp_servers.graylog]
  command = "/Users/borkert/Repos/terraform-provider-graylog/.venv/bin/mcp-proxy"
  args = [
    "--transport", "streamablehttp",
    "--headers", "Authorization", "Basic MWQwZzR0Z3BnNWgzdjNqa2NkcjltcjFybjluZXI5aWxjN2l1b3NtOGhlZjA1NG04NzVkZjp0b2tlbg==",
    "https://graylog.internal.borkert.net/api/mcp"
  ]
  startup_timeout_sec = 20
  tool_timeout_sec = 60
  ```
- Secrets ignored: `.mcp/` and `.venv/` are in `.gitignore`.

## How to Enable on Your Machine
1) Ensure MCP is enabled in Graylog (`System → Configurations → MCP`) and create a read-only token; build the Basic header: `echo -n "<token>:token" | base64`.
2) Create venv and install proxy:
   ```bash
   python3 -m venv .venv
   .venv/bin/pip install mcp-proxy
   ```
3) Add the `mcp_servers.graylog` block above to `~/.codex/config.toml` (adjust paths/host/token as needed).
4) Restart Codex CLI. Check with `codex mcp list` (it just ensures config is loaded; no OAuth step needed).

## Verifying Connectivity
Use Codex (or curl via proxy bypass) to call a tool:
```bash
curl -k -H "Content-Type: application/json" \
  -H "Authorization: Basic <base64_token>" \
  --data '{"id":"1","jsonrpc":"2.0","method":"tools/call","params":{"name":"get_system_status","arguments":{}}}' \
  https://graylog.internal.borkert.net/api/mcp
```
Expected: JSON with version/hostname. Through Codex, ask to call MCP tool `get_system_status` or `list_index_sets`.

## Tips & Troubleshooting
- If Codex fails: ensure the `rmcp_client` feature flag is enabled (set `[features].rmcp_client = true` in `config.toml` or start Codex with `--enable rmcp_client`); the old `experimental_use_rmcp_client` flag is deprecated.
- If you get “Method not allowed” or metadata errors, confirm you’re using `tools/call` and POST, not GET.
- Permissions: token must have read access to index sets/streams; otherwise tools return empty sets.

## Useful MCP tools for provider debugging
- `get_system_status` — verify Graylog version/leader/processing state.
- `list_index_sets` — confirm current index sets and IDs (used in acceptance/E2E checks).
- `list_indices` / `list_fields` — inspect index/field metadata without mutating state.
- `aggregate_messages` / `search_messages` — spot-check data flowing through inputs/streams during tests.
