# Graylog 7.0.4 Pipeline Rule Caching Bug

## Environment

- Graylog version: 7.0.4+ac94792
- Deployment: Docker container
- Search backend: OpenSearch
- Pipeline management: Terraform provider (API-driven)
- Message processor order: Message Filter Chain → Stream Rule Processor → Pipeline Processor (correct order)

## Summary

Pipeline rules updated or created via the Graylog REST API are stored correctly in MongoDB but the Pipeline Interpreter continues executing stale, cached compiled versions of the rules. This behavior persists across multiple Graylog container restarts.

## Observations

### 1. Initial pipeline rule creation works

A rule created during initial Terraform setup and loaded on first Graylog startup executes correctly:

```
rule "parse syslog message"
when true
then
  set_field("syslog_parsed", "yes");
end
```

Result: `syslog_parsed: "yes"` appears on all syslog messages. Confirmed working.

### 2. Rule source code updates via API are ignored

The same rule was updated via `PUT /api/system/pipelines/rule/{id}` to add additional `set_field` calls:

```
rule "parse syslog message"
when true
then
  set_field("syslog_parsed", "yes");
  set_field("enrichment_test", "via_parse");
  set_field("source_type", "syslog_test");
end
```

The API confirms the update (`modified_at` timestamp changes, `errors: null`). However, after a Graylog container restart:
- `syslog_parsed: "yes"` still appears (from the original compiled version)
- `enrichment_test` and `source_type` changes do NOT appear
- The Pipeline Interpreter uses the original compiled rule, not the updated database version

### 3. New rules in existing pipelines don't execute

A brand new rule (`syslog enrichment`, new ID, new name) was created and referenced by the pipeline in stage 0. After restart:
- The new rule does NOT execute (its fields don't appear on messages)
- The existing `parse syslog message` rule in stage 1 continues to work

This was tested with:
- `when true` condition (unconditional)
- Simple `set_field("enrichment_test", "works")` action
- Both single-stage and multi-stage pipeline configurations
- Rule swap between stages (rules that don't execute in stage 0 also don't execute in stage 1)

### 4. Deleted rules continue to execute

The most alarming observation: after deleting a rule that set `source_type: "syslog_v2"` and creating a new pipeline with a completely different rule that sets `source_type: "syslog_v3"`:
- Messages continue to show `source_type: "syslog_v2"` (from the deleted rule)
- The new rule's `source_type: "syslog_v3"` never appears
- Multiple container restarts do not resolve the issue

### 5. Same pipeline works on different streams

To rule out pipeline logic issues, the same test rule was connected to both the syslog stream and the B2Bi TE stream:
- On TE stream: `pipeline_test: "working"` - the rule executes correctly
- On syslog stream: `pipeline_test: "-"` - the rule does NOT execute

This proves the rule itself is valid and the Pipeline Interpreter can execute it, but fails to do so on certain streams.

### 6. Pipeline processor order is correct

Verified via `GET /api/system/messageprocessors/config`:
1. AWS Instance Name Lookup
2. GeoIP Resolver
3. Message Filter Chain
4. Stream Rule Processor
5. Pipeline Processor

The Pipeline Processor runs last, after stream routing. This is the recommended configuration.

## Reproduction Steps

1. Create a Graylog 7.0.4 instance with a Syslog TCP input
2. Create a stream with rules matching the syslog input
3. Create a pipeline rule via REST API:
   ```json
   {
     "title": "test rule",
     "source": "rule \"test rule\"\nwhen true\nthen\n  set_field(\"test\", \"v1\");\nend\n"
   }
   ```
4. Create a pipeline referencing this rule and connect it to the stream
5. Restart Graylog - verify `test: "v1"` appears on messages
6. Update the rule via `PUT` to set `test: "v2"`
7. Restart Graylog
8. Observe that `test: "v1"` still appears (stale cache)

## Impact

- Pipeline rules managed via the REST API (Terraform, automation, etc.) cannot be updated after initial creation
- This makes API-driven pipeline management unreliable in Graylog 7.0.4
- The only workaround found so far is unknown - neither restarts, rule deletion/recreation, nor pipeline deletion/recreation resolve the caching issue for affected streams

## Workaround Attempts (all failed)

- Updating rule source code via PUT → changes stored in DB but not picked up
- Deleting and recreating rules with new IDs → new rules not executed
- Deleting and recreating pipelines with new IDs → stale rules still execute
- Deleting and recreating pipeline connections → no effect
- Multiple Graylog container restarts → cache persists
- Creating brand new rules with different names → not executed on affected stream
- Saving/updating pipeline via PUT (same content) to trigger re-resolution → no effect

## Possibly Related

- Early in the session, the PipelineResolver logged: `WARN: Cannot resolve rule <parse syslog message> referenced by stage #1 within pipeline <698f0203644d3f6fca889122>` - this occurred because the pipeline was created before the rule existed (race condition in Terraform apply). This initial resolution failure may have permanently corrupted the Pipeline Interpreter's cache for this stream.
- The TE stream pipeline (which was created later, without any resolution failures) works correctly

## Graylog Version

```
Graylog Server 7.0.4+ac94792 (Noir)
Timezone: Etc/UTC
OS: Linux 4.18.0-553.97.1.el8_10.x86_64
```
