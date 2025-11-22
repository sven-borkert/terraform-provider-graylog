# Graylog MCP Connection (Local Dev)

This repo includes a local MCP config to talk to the Graylog 7 server.

## Files
- `.mcp/graylog.json` — MCP server definition pointing to `https://graylog.internal.borkert.net/api/mcp` with Basic auth header derived from the provided API token (useful for tools that support HTTP MCP directly).
- `.gitignore` includes `.mcp/` so secrets stay local.

## How to Use (examples)
- Claude Code CLI:
  ```bash
  claude mcp add --config .mcp/graylog.json
  claude mcp list
  ```
- LM Studio (`~/.lmstudio/mcp.json`):
  - Merge the contents of `.mcp/graylog.json` into your `mcp.json` under `"mcpServers"`.
  - Restart or reload settings.

## Codex CLI (needs STDIO proxy)
Codex only speaks STDIO MCP, so use `mcp-proxy` to bridge HTTP→STDIO. We vendor a local venv with `mcp-proxy` installed.

`~/.codex/config.toml` entry:
```toml
[mcp_servers.graylog]
command = "/Users/borkert/Repos/terraform-provider-graylog/.venv/bin/mcp-proxy"
args = [
  "--transport", "streamablehttp",
  "--headers", "Authorization", "Basic <base64_token>",
  "https://graylog.internal.borkert.net/api/mcp"
]
startup_timeout_sec = 20
tool_timeout_sec = 60
```
Then restart Codex and run `codex mcp list` to confirm the server is configured. No `login` step is needed (Basic auth only).

## Security Notes
- Token is scoped to a read-only MCP user. Regenerate if shared or committed by mistake.
- For other environments, create a new API token and update `.mcp/graylog.json` locally.
