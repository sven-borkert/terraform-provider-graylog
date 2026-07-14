# Code Analysis: terraform-provider-graylog

*Analysis date: 2026-07-14 — at commit `aadace6` (Add graylog_index_set_field_type resource)*

## Overall health

The project is in good shape for a codebase mid-refactor: `go build`, `go vet`, and the entire test suite (48 test files) pass cleanly, and `golangci-lint` reports only 3 issues — all unchecked `d.Set` calls in the newest resource, `graylog/resource/system/indices/fieldtype/resource.go:107-109`. The codebase is ~22k lines of Go across 312 files, structured as a clean three-layer design: 31 thin HTTP sub-clients under `graylog/client/`, per-resource CRUD packages under `graylog/resource/` and `graylog/datasource/`, and shared helpers in `graylog/util/util.go` and `graylog/convert/`.

The Graylog 7 refactor is roughly 60% documented (14 of 23 resource docs moved out of `todo/`), but the code is further along than the docs suggest: event_definition, event_notification, and the sidecar resources already use the Graylog 7 patterns (`WrapEntityForCreation`, `RemoveComputedFields`) while their docs still sit in `docs/resources/todo/`.

## Highest-priority findings

### 1. The provider credential is not marked sensitive

`auth_password` in `graylog/provider/provider.go:37-43` lacks `Sensitive: true`, so it can appear in plan output and logs. Only the user-resource password and LDAP password are marked. This is a one-line fix and the most user-visible issue found.

### 2. Panic risk from unchecked type assertions on API responses

There are ~248 unchecked `.(T)` assertions outside tests. Many are schema-safe `d.Get(...)` casts, but a meaningful subset assert on raw decoded JSON from the API and will panic (crashing the whole provider plugin) if Graylog returns an unexpected shape. The hotspot is `graylog/resource/dashboard/util.go` (lines 50, 54, 105, 107, 160, 175, 262, 264), which is exactly the hand-built `map[string]interface{}` state-graph code for the Views API — the newest and most complex logic. Others: `graylog/resource/stream/util.go:33`, `graylog/resource/system/input/extractor/create.go:23`, and the alert condition state upgrader (`graylog/resource/stream/alert/condition/state_upgrader.go`).

### 3. The newest code is the least tested

Test coverage overall is thin (~2% per package; the entire client layer has 0%) because the flute-mock tests mostly exercise schemas. But four packages have *no* tests at all, and they are precisely the recent, complex Graylog 7 work:

- `graylog/resource/saved_search/` — the trickiest package; it splits one resource across two API objects
- `graylog/resource/view/`
- `graylog/resource/system/indices/template/`
- `graylog/resource/system/indices/fieldtype/`

The risk concentration is inverted from where you'd want it: the older stream/system/user/role resources are both simpler and better covered.

### 4. Orphan-on-partial-failure in saved_search

In `graylog/resource/saved_search/create.go:181-184`, if the `viewResp["id"]` assertion fails after `View.Create` succeeded, the function errors out without deleting the just-created search — unlike the `err != nil` branch right above it, which does roll back.

## Architectural observations

- **Two CRUD paradigms coexist.** 27 resources use the deprecated non-context SDKv2 signatures (`Create:`, `Read:` returning `error`), fabricating `ctx := context.Background()` in ~118 places, so Terraform's cancellation and timeouts never reach HTTP calls. Only the newest resource (`graylog/resource/system/indices/fieldtype/resource.go`) uses `CreateContext` + `diag.Diagnostics`. The provider itself uses the deprecated `ConfigureFunc` (`graylog/provider.go:14`). If new resources keep being added in the new style, the split will keep widening — worth deciding whether to migrate the rest mechanically.
- **Dashboard split-brain.** Dashboards read single items via the Graylog 7 `/views` API but list via the legacy `/dashboards` client, with a triple-fallback for `elements`/`dashboards`/`views` response keys in `graylog/datasource/dashboard/read.go:48-57`. Works, but it's the fragile spot if Graylog drops the legacy endpoint.
- **Asymmetric transform layers.** `WrapEntityForCreation` is applied inside six *clients* (stream, view, event/notification, event/definition, sidecar/configuration, sidecar/collector), while the mirror-image `RemoveComputedFields` is applied in 23 *resources* — the create and update payload transforms live at different layers, which makes the data flow harder to follow.
- **Heavy duplication in the client layer.** The 31 sub-clients implement near-identical `Get/Gets/Create/Update/Delete` differing only in URL path and whether they wrap the entity — a strong candidate for one generic client (Go generics or a path-parameterized base) that would delete well over a thousand lines.
- **Inert plumbing:** `config.LoadAndValidate()` is a no-op stub (`graylog/config/config.go:12`), and the `api_version` setting defaults to `"v3"` but nothing branches on it — Graylog 7 behavior is hard-coded. `graylog_view` is fully implemented but commented out of the resource map (`graylog/resource_map.go:64`, `// TODO support view`).

## Minor cleanups

- Leftover debug logging in dashboard and saved_search create/update paths (e.g. `graylog/resource/dashboard/create.go:68`, and an unconditional `log.Printf("dashboard flatten input state type %T", ...)` at `graylog/resource/dashboard/util.go:159`) — no other resource logs.
- ~20 `_ = d.Set(...)` swallowed errors in datasources.
- The lint config (`.golangci.yml`) enables only errcheck, govet, ineffassign, and unused — adding `staticcheck` would catch the deprecated-SDK-field usage and more.
- Naming drift: `graylog/client/role/role.go` and `graylog/client/user/user.go` break the `client.go` convention.

## Suggested priority order

1. Mark `auth_password` as `Sensitive` and fix the 3 lint errors (minutes of work).
2. Add ok-checks to the type assertions in `dashboard/util.go` and `saved_search` (the panic-prone API-response paths), fixing the saved_search orphan case while there.
3. Add flute-mock tests for saved_search and fieldtype — the untested, highest-complexity packages.
4. Longer term: migrate resources to `CreateContext`-style signatures (mechanical, per-resource), collapse the client layer with a generic implementation, and enable staticcheck.
