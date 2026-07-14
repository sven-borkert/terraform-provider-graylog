# Graylog 7 API Verification — Provider Client Implementation

*Verification date: 2026-07-14, at provider commit `aadace6`.*

## Scope and method

This report verifies that the provider's API client implementation (`graylog/client/`, plus the request/response handling in `graylog/resource/` and `graylog/datasource/`) matches the actual Graylog 7 REST API.

- **Target server:** Graylog **7.0.8+38f8cc8** (the running test instance, verified via its unauthenticated version endpoint).
- **Ground truth:** the `Graylog2/graylog2-server` source at git tag `7.0.8` (commit `38f8cc8...` — exactly the build the server reports). The live swagger cache (`docs/api-docs/`, gitignored) could not be refetched because `graylog.auto.tfvars` is absent on this machine and Vault was unavailable, so the JAX-RS resource classes were used directly.
- **Coverage:** all ~120 endpoint calls across all 31 client packages were checked for: path/method existence, required query params, request body shape (especially `{"entity": ..., "share_request": ...}` wrapping via `util.WrapEntityForCreation`), and response shape vs. what the provider parses.

Note: this repo contains no MCP server; the "implementation" verified here is the provider's Graylog API client, which is the only code that talks to the Graylog API.

## Summary

| Domain | Client packages | Verdict |
|---|---|---|
| Streams, stream rules, stream outputs | `client/stream`, `stream/rule`, `stream/output` | ✅ All match |
| Legacy alerting | `stream/alarmcallback`, `stream/alert/condition` | ❌ API removed — resources dead |
| Views, saved searches | `client/view`, `view/search`, `search/saved` | ✅ Match, except nonexistent DELETE used for rollback |
| Legacy dashboards | `client/dashboard`, `dashboard/widget`, `dashboard/position` | ❌ Mostly removed — widget/position resources dead |
| Index sets, field types, templates, deflector | `system/indices/*` | ✅ All match (several state-handling issues) |
| Inputs, static fields, outputs, grok | `system/input`, `system/output`, `system/grok` | ✅ Match |
| Extractors | `system/input/extractor` | ❌ Create/update send pre-5.0 wire format |
| Pipelines, rules, connections | `system/pipeline/*` | ✅ Match, except data source list endpoint |
| Events (definitions, notifications) | `client/event/*` | ✅ All match |
| Sidecar | `client/sidecar/*` | ❌ Both create endpoints wrongly entity-wrapped |
| Users | `client/user` | ❌ Update targets wrong path semantics; password change is a no-op |
| Roles | `client/role` | ✅ All match |
| LDAP settings | `system/ldap/setting` | ❌ API removed in Graylog 4 — resource dead |

## Critical: resources targeting removed APIs (every call 404s)

These are registered in `graylog/resource_map.go` / `graylog/datasource_map.go` but their entire API surface no longer exists in 7.0.8. They should be removed or hard-deprecated:

1. **`graylog_alarm_callback`** — `/streams/{id}/alarmcallbacks` (`StreamAlarmCallbackResource`) is gone; legacy alerting was removed. Successor: `graylog_event_notification` (already implemented).
2. **`graylog_alert_condition`** — `/streams/{id}/alerts/conditions` (`StreamAlertConditionResource`) is gone. Successor: `graylog_event_definition` (already implemented).
3. **`graylog_dashboard_widget`** (resource *and* data source) — no class serves `/dashboards/{id}/widgets*` in 7.0.8; widgets live inside the view `state` now (which the current `graylog_dashboard` resource already manages).
4. **`graylog_dashboard_widget_positions`** — `PUT /dashboards/{id}/positions` and `GET /dashboards/{id}` (used by its read, `resource/dashboard/position/read.go:19`) are both gone.
5. **`graylog_ldap_setting`** — `/system/ldap/settings` was removed in Graylog 4.0. Successor API: authentication services under `/system/authentication/services/backends` (backend CRUD, LDAP/AD as config types) plus `/system/authentication/services/configuration` (activation). Needs a full rewrite or removal.

This matches the `docs/resources/todo/` list almost exactly — the un-refactored components are precisely the broken ones.

## Critical: broken operations on live resources

6. **`graylog_extractor` create/update fail on 7.0.8** — `resource/system/input/extractor/util.go:26` renames `cursor_strategy` → `cut_or_copy` (the pre-Graylog-5 field name), and `util.go:36-43` sends `converters` as a map `{type: config}` where 7.0.8 requires a **list** of `{"type", "config"}` objects. Graylog's ObjectMapper rejects unknown properties, so the request fails multiple ways. The read path already handles the new shapes; only the request builder needs fixing (drop the rename, send the converters list).
7. **`graylog_sidecar_collector` and `graylog_sidecar_configuration` creates fail** — `client/sidecar/collector/client.go:44` and `client/sidecar/configuration/client.go:44` apply `util.WrapEntityForCreation`, but 7.0.8's `createCollector`/`createConfiguration` take the **bare DTO** (no `CreateEntityRequest` support) and unknown properties are rejected → HTTP 400. The wrapping must be removed for these two clients. (Consistent with sidecar docs still being in `todo/` and no sidecar resources in the e2e config — these were never tested against Graylog 7.)
8. **`graylog_user` update is broken** — 7.0.8's update is `PUT /users/{userId}` resolved strictly by Mongo ObjectId (or the literal `local:admin`), but `client/user/user.go` + `resource/user/update.go` send the **username**. Compounding it, the provider never captures the real ID: the API returns `"id"` but the schema key is `user_id`, and `convert.SetResourceData` only copies exact schema-key matches, so `user_id` stays empty in state. Additionally, **password changes are silent no-ops**: `ChangeUserRequest` has no `password` field and ignores unknowns; 7.0.8 changes passwords only via `PUT /users/{userId}/password`, which the provider never calls.
9. **`graylog_pipeline` data source by-title lookup fails** — `client/system/pipeline/pipeline/client.go:32` (`Gets`) calls `GET /system/pipelines/pipeline`, which returns a **bare JSON array**; decoding into `map[string]interface{}` errors in go-httpclient, and `datasource/system/pipeline/pipeline/read.go:31` then expects a `"pipelines"` key that only `GET /system/pipelines/pipeline/paginated` returns. Fix: decode the array, or switch to `/paginated` (mind `per_page=50` default). The by-ID path is fine.
10. **Search rollback calls a nonexistent endpoint** — `DELETE /views/search/{id}` does not exist in 7.0.8 (the only DELETE on `SearchResource` is `cancel/{jobId}`). The three best-effort rollback calls (`resource/dashboard/create.go:77`, `resource/saved_search/create.go:177`, `resource/saved_search/update.go:178`) silently get 405 and clean up nothing; orphaned Search documents are left to the server's periodic cleanup job. The calls should be dropped or replaced.

## Data-correctness and drift issues

11. **`graylog_index_set` update wipes field type profiles** — `resource/system/indices/indexset/util.go:61` deletes `field_type_profile` from the update body, and the server's `IndexSetUpdateRequest.toIndexSetConfig` writes `null` when the field is absent — so any profile assigned via the UI is removed on every apply that updates the index set.
12. **`writable` is always sent, defaulting to `false`** — `convert.GetFromResourceData` emits every schema key, so an unset `writable` creates a read-only index set, and updating the default index set without `writable = true` yields HTTP 409 ("Default index set must be writable").
13. **Event definitions are force-enabled on every apply** — the provider never sends `?schedule` (server default `true`) and doesn't model `state`, so a definition disabled in the UI is silently re-enabled by any `terraform apply` touching it. Similarly, `storage` is absent from the schema, so non-default storage handler configs are reset on update.
14. **Sidecar configuration `tags` are wiped on update** — the server supports `tags` (nullable → empty set) but the provider schema omits it, so every update clears tags set in the UI.
15. **`graylog_sidecar` manages all sidecars on the server** — its read imports every sidecar from `GET /sidecars/all` (including tag-based assignments, which it re-submits as manual ones), and its delete clears assignments for **every** sidecar, not just managed ones. Not an API mismatch, but behaviorally hazardous.
16. **Index set data source mishandles Graylog 7 fields** — `datasource/system/indices/indexset/util.go` JSON-encodes only `rotation_strategy`/`retention_strategy`; `data_tiering` and `field_restrictions` come back as objects and get set into `TypeString` attributes (fails/mangles for data-tiering index sets). It also reuses the resource schema, so the datasource-only `index_template_type` is never populated. The resource-side read handles all this correctly.
17. **`graylog_index_set_field_type` read ignores `origin`** — `client/system/indices/fieldtype/client.go` matches any `field_name` regardless of origin, so a field typed by the index itself (`INDEX`) or a profile (`PROFILE`) is mistaken for an existing custom mapping; state is never cleared after out-of-band removal. 7.0.8 exposes `origin` (`OVERRIDDEN_INDEX`/`OVERRIDDEN_PROFILE` = custom) precisely for this distinction. Also, the DTO field is `is_reserved` (there is no `is_custom`).
18. **Pagination caps on lookup paths (50 items)** — three by-title/by-name scans read only the first default page: dashboards (`datasource/dashboard/read.go`, `GET /dashboards`), saved searches (`datasource/search/saved/read.go`, `GET /search/saved`), and field types (`GetFieldType`, `per_page=50`, page 1 only). Items beyond the first 50 are not found. Fixes: pass `query=`/`per_page`, or for field types use the unpaginated `GET /system/indices/index_sets/types/{id}/all?fieldNameQuery=`.
19. **`graylog_index_set_templates` lists built-ins only** — the data source calls only `/system/indices/index_sets/templates/built-in`; custom templates never appear although `GET .../templates/paginated` exists.
20. **Input reads mask secrets** — 7.0.8 masks password-type attribute values in input GET responses, which can cause perpetual Terraform diffs on inputs with secret config values. Server behavior; would need diff suppression to hide.
21. **Rule updates null out builder metadata** — pipeline rule PUT omits `rule_builder`/`simulator_message`, which the server overwrites with null; harmless for source-managed rules, but wipes metadata on imported UI-built rules.

## Confirmed correct (notable)

- **Entity wrapping is exactly right where it's used for streams, views, and events**: `POST /streams`, `POST /views`, `PUT /views/{id}` (update is also wrapped — verified required), `POST /events/definitions`, `POST /events/notifications` all take `CreateEntityRequest{entity, share_request}` in 7.0.8, matching `util.WrapEntityForCreation` field-for-field. `POST /views/search` is correctly *not* wrapped.
- Response-key parsing all verified: `{"streams", "total"}`, `{"stream_id"}`, `{"streamrule_id"}`, `{"outputs"}`, `{"inputs"}`, `{"patterns"}`, `{"index_sets"}`, `{"sidecars"}`, `{"users"}`, `{"roles"}`, `elements` for `/dashboards`, `/search/saved`, and field types.
- The `elements` fallback in `datasource/dashboard/read.go` is the correct branch on 7.0.8; the `dashboards`/`views` fallbacks are dead but harmless.
- Index set create/update field sets (including `data_tiering`, `use_legacy_rotation`), field-type mapping DTOs and the allowed type list, index set templates, deflector (`GET /system/deflector/{indexSetId}` → `{is_up, current_target}`), grok (incl. single-pattern POST semantics and `PUT /system/grok/{id}`), static fields (DELETE by field *key* — correct despite the parameter name), stream pause/resume, pipeline/rule/connection CRUD, sidecar assignment (`PUT /sidecars/configurations` with `nodes`/`assignments`), and role CRUD all match 7.0.8 exactly.
- User create correctly strips `full_name` and sends `first_name`/`last_name` (7.0.8 rejects unknown properties — `full_name` would 400).

## Deprecated-but-working endpoints (future risk)

Still functional in 7.0.8 but `@Deprecated` upstream; paginated successors exist:

| Used endpoint | Successor |
|---|---|
| `GET /streams` (stream datasource) | `GET /streams/paginated` |
| `GET /users`, `GET /users/{username}` | `GET /users/paginated` |
| `GET /system/grok` | `GET /system/grok/paginated` |

## Minor cleanups

- Misleading "Graylog 7 Update requires id in body" comments: for streams (`resource/stream/update.go:29`) the claim is false and the field is ignored; for outputs (`resource/system/output/update.go:26-27`) it's false **and** the extra `id` key is written verbatim into the Mongo document via `$set` — remove it. (For index sets the body `id` *is* genuinely required; for pipelines/grok it's accepted/ignored.)
- Dead code: `dashboard.Client.Create/Update/Delete` (legacy `/dashboards` writes, no call sites), `user.Client.Gets`, `role.Client.Gets`, `output` client's `GetAllParams` (skip/limit/stats — params don't exist on the 7.0.8 endpoint and are never sent).
- Sidecar collector/configuration DELETE returns 304 when nothing was deleted; the provider treats it as success.

## Suggested priority order

1. **Remove/deprecate the five dead components** (#1–5): `graylog_alarm_callback`, `graylog_alert_condition`, `graylog_dashboard_widget` (+data source), `graylog_dashboard_widget_positions`, `graylog_ldap_setting`.
2. **Fix the four broken write paths** (#6–9): extractor request format, sidecar create wrapping (2 clients), user update by ID (store `user_id`, use it in PUT, add the password endpoint), pipeline data source list.
3. **Drop the dead rollback DELETE** (#10) and fix the state-correctness issues most likely to bite: index set `field_type_profile` wipe (#11), `writable` default (#12), index set datasource `data_tiering` (#16).
4. Address drift/pagination issues (#13–15, #17–19) as the affected resources get their Graylog 7 refactor pass.
