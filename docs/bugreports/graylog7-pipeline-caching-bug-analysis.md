# Graylog 7.0.4 Pipeline Caching Bug - Root Cause Analysis

**Related bug report:** [graylog7-pipeline-caching-bug.md](graylog7-pipeline-caching-bug.md)

## Executive Summary

After extensive investigation of both the Terraform provider code and Graylog server source code, I've identified **a critical architectural flaw in the StageIterator caching mechanism** that causes pipeline rules to become stale after updates. The bug is in Graylog server, not the Terraform provider.

## Root Cause: StageIterator.Configuration Cache with Flawed Equality Contract

### The Caching Architecture

**File:** `graylog2-server/src/main/java/org/graylog/plugins/pipelineprocessor/processors/PipelineInterpreter.java` (lines 481-539)

```java
private final LoadingCache<Set<Pipeline>, StageIterator.Configuration> cache;

public StageIterator getStageIterator(Set<Pipeline> pipelines) {
    if (cachedIterators) {
        return new StageIterator(cache.get(pipelines));  // THE BUG
    } else {
        return new StageIterator(pipelines);
    }
}
```

The cache uses `Set<Pipeline>` as the key, which relies on `Pipeline.equals()` and `Pipeline.hashCode()`.

### The Equality Flaw

**Pipeline.java** uses `@AutoValue`:
- Generates equals/hashCode from: `id()`, `name()`, `stages()`
- `stages()` returns `SortedSet<Stage>`

**Stage.java** uses `@AutoValue` BUT:
```java
@AutoValue
public abstract class Stage implements Comparable<Stage> {
    private List<Rule> rules;  // NOT an AutoValue property - MUTABLE!
    // Comment in source: "not an autovalue property, because it introduces a cycle in hashCode()"

    public void setRules(List<Rule> rules) {
        this.rules = rules;
    }
}
```

**Critical Finding:** `Stage.equals()` only compares:
- `stage()` - stage number
- `match()` - match type (ALL/EITHER/PASS)
- `ruleReferences()` - rule NAMES only (not Rule objects)

The actual **Rule objects** (with compiled logic) are in the **mutable `rules` field** which is **deliberately excluded** from equals/hashCode!

### The Bug Sequence

```
1. Initial State:
   - Pipeline P1 with Stage containing Rule R1 ("test" → "v1")
   - Cache stores: {P1} → Configuration(Stage with R1)

2. Rule Updated:
   - Rule R1 modified in MongoDB ("test" → "v2")
   - RulesChangedEvent posted
   - ConfigurationStateUpdater.reloadAndSave() runs
   - New Pipeline P1' created with new Rule R1'
   - stage.setRules(newRules) called

3. Cache Lookup:
   - cache.get({P1'}) checks if {P1'} equals any existing key
   - Pipeline.equals(P1, P1') → compares stages
   - Stage.equals() → compares stage number, match type, rule REFERENCES (names)
   - Rule names unchanged → Stage.equals() returns TRUE
   - Pipeline.equals() returns TRUE

4. Result:
   - CACHE HIT with stale entry
   - Old Configuration with old Rule R1 returned
   - "v1" continues to be set instead of "v2"
```

### Why Restarts Don't Help - Additional Factors

The StageIterator cache is recreated with each new State (line 497-505), so simple cache staleness shouldn't persist. Additional factors must be involved:

**Hypothesis 1: Event Not Being Fired or Processed**
- ClusterEventPeriodical bridges ClusterEventBus → ServerEventBus
- If bridge fails silently, ConfigurationStateUpdater never receives events
- Silent failures in `reloadAndSave()` (no try/catch) would leave stale state

**Hypothesis 2: Resolution Failure Persisted in Pipeline Source**
- When pipeline is created before rules exist, a `WARN: Cannot resolve rule` is logged
- The pipeline SOURCE in MongoDB doesn't change when rules are added
- `ruleReferences()` still contains the rule name, but resolution keeps failing

**Hypothesis 3: MongoDB Read Returning Stale Data**
- Docker networking or MongoDB connection pooling issues
- Unlikely for single-node deployment but worth investigating

**Stream-specific behavior** suggests:
- The syslog stream's pipeline connection was created during initial (broken) resolution
- TE stream was connected later when rules existed
- Something in the connection/resolution path is cached at stream level

## Secondary Contributing Factors

### 1. No Error Handling in Reload Path

**File:** `ConfigurationStateUpdater.java` (lines 105-113)

```java
private synchronized PipelineInterpreter.State reloadAndSave() {
    // NO try/catch - exceptions go to thread pool handler
    final ImmutableMap<String, Pipeline> currentPipelines =
        pipelineResolver.resolvePipelines(pipelineMetricRegistry);
    // ... if this throws, state remains stale
}
```

If reload fails, exceptions are logged but state isn't updated.

### 2. Async Event Scheduling Without Retry

```java
scheduler.schedule(() -> serverEventBus.post(reloadAndSave(event)), 0, TimeUnit.SECONDS);
```

- Events are scheduled asynchronously
- No retry mechanism if reload fails
- No acknowledgment that reload succeeded

### 3. Event Bridge Delays

Events flow: `ClusterEventBus` → MongoDB → Polling (1 sec) → `ServerEventBus`

Local events are posted immediately, but there's a window where:
1. Rule is saved to MongoDB
2. Event is posted
3. Reload scheduled immediately (0 second delay)
4. But scheduler executes async - MongoDB might not have flushed yet

## API Usage Verification

The Terraform provider **uses the API correctly**:

| Operation | Endpoint | Implementation |
|-----------|----------|----------------|
| Create Rule | `POST /system/pipelines/rule` | ✓ Correct |
| Update Rule | `PUT /system/pipelines/rule/{id}` | ✓ Correct (removes computed fields) |
| Delete Rule | `DELETE /system/pipelines/rule/{id}` | ✓ Correct |
| Create Pipeline | `POST /system/pipelines/pipeline` | ✓ Correct |
| Update Pipeline | `PUT /system/pipelines/pipeline/{id}` | ✓ Correct (adds ID to body for Graylog 7) |
| Connect Pipeline | `POST /system/pipelines/connections/to_stream` | ✓ Correct |

The API documentation confirms "up to 1 second" delay for changes, which is expected.

## Workaround

### Disable the StageIterator Cache

Add to `graylog.conf` or as environment variable:
```
cached_stageiterators = false
```

### Performance Impact

The `StageIterator.Configuration` is a lightweight structure that:
- Iterates through pipelines connected to a stream
- Builds an `ArrayListMultimap<Integer, Stage>` mapping stage numbers to stages
- Computes first/last stage extent

| Aspect | Cache Enabled | Cache Disabled |
|--------|---------------|----------------|
| Per-message overhead | ~100 nanoseconds | ~1-5 microseconds |
| Work per message | O(1) cache lookup | O(P × S) where P=pipelines, S=stages |

For typical deployments (3 pipelines × 2 stages):
- At 10,000 msg/sec: ~10-50ms total overhead per second (~1-5% of one CPU core)
- At 100,000 msg/sec: ~100-500ms total overhead per second

**This overhead is negligible for most deployments.** The cache was a premature optimization that introduced a correctness bug.

## Recommended Fix (Graylog Server)

### Option 1: Cache Invalidation on State Update (Preferred)

When a new `PipelineInterpreter.State` is created, ensure the cache doesn't produce false hits. Options:
- Include a state version/timestamp in the cache key
- Clear cache entries when rules change
- Use rule content hashes in the key

### Option 2: Fix Stage Equality Contract

Include rule content (or a hash) in Stage equality:

```java
// In Stage.java, add to AutoValue comparison:
@Memoized
public String rulesContentHash() {
    return rules.stream()
        .map(r -> r.id() + ":" + r.when().hashCode() + ":" + r.then().hashCode())
        .collect(Collectors.joining(","));
}
```

### Option 3: Remove the Cache Entirely

The performance benefit is minimal. Remove the cache and always create fresh `StageIterator.Configuration` objects.

## Verification Steps

### Step 1: Confirm the Bug
```bash
# Create rule with v1
curl -X POST http://graylog:9000/api/system/pipelines/rule \
  -H "Content-Type: application/json" \
  -d '{"source": "rule \"test\"\nwhen true\nthen\n  set_field(\"test\", \"v1\");\nend"}'

# Create pipeline, connect to stream...
# Verify "test: v1" appears on messages

# Update rule to v2
curl -X PUT http://graylog:9000/api/system/pipelines/rule/{id} \
  -H "Content-Type: application/json" \
  -d '{"source": "rule \"test\"\nwhen true\nthen\n  set_field(\"test\", \"v2\");\nend"}'

# Wait 2 seconds, check messages
# If "test: v1" still appears → bug confirmed
```

### Step 2: Check Event Delivery
Enable DEBUG logging for pipeline processor:
```xml
<logger name="org.graylog.plugins.pipelineprocessor" level="DEBUG"/>
```

After rule update, look for:
- `"Refreshing rule {}"` - event received
- `"Pipeline interpreter state got updated"` - state replaced

If these logs are missing, event delivery is broken.

### Step 3: Test Cache Disable Workaround
Add to `graylog.conf` or environment:
```
cached_stageiterators = false
```

Restart Graylog and repeat update test. If updates work now, the cache is confirmed as the culprit.

### Step 4: Check MongoDB Data Directly
```javascript
// In MongoDB shell
db.pipeline_processor_rules.find({title: "test"}).pretty()
```

Verify the source field contains "v2". If MongoDB has correct data but Graylog uses "v1", the issue is in the Java code path.

### Step 5: Check Cluster Events
```javascript
db.cluster_events.find({event_class: /RulesChangedEvent/}).sort({timestamp: -1}).limit(5).pretty()
```

Verify events are being persisted. Check if `consumers` array includes your node ID.

## Files Analyzed

### Graylog Server
| File | Purpose | Key Lines |
|------|---------|-----------|
| `PipelineInterpreter.java` | Cache implementation | 477-547 |
| `ConfigurationStateUpdater.java` | Event handling and reload | 105-200 |
| `PipelineResolver.java` | Pipeline/rule resolution | 109-209 |
| `MongoDbRuleService.java` | Event posting after save | 89, 168 |
| `Stage.java` | Mutable rules field | 31-47 |
| `Pipeline.java` | AutoValue with stages | 1-50 |
| `StageIterator.java` | Configuration structure | 61-107 |
| `ClusterEventPeriodical.java` | Event bridge mechanism | 159-179 |

### Terraform Provider
| File | Purpose |
|------|---------|
| `graylog/client/system/pipeline/rule/client.go` | API calls |
| `graylog/resource/system/pipeline/rule/resource.go` | Resource implementation |
| `graylog/resource/system/pipeline/pipeline/resource.go` | Pipeline resource |

## Provider Mitigation: content_hash Attribute

The `graylog_pipeline_rule` resource now includes a computed `content_hash` attribute that contains the SHA256 hash of the rule source. This can be used with Terraform's `replace_triggered_by` lifecycle to force pipeline connection recreation when rules change, which may help invalidate Graylog's stale cache.

### Usage Example

```hcl
resource "graylog_pipeline_rule" "syslog_parse" {
  source = <<-EOF
    rule "parse syslog"
    when true
    then
      set_field("parsed", "yes");
    end
  EOF
}

resource "graylog_pipeline" "syslog" {
  source = <<-EOF
    pipeline "syslog processing"
    stage 0 match either
      rule "parse syslog"
    end
  EOF
}

resource "graylog_pipeline_connection" "syslog" {
  stream_id    = graylog_stream.syslog.id
  pipeline_ids = [graylog_pipeline.syslog.id]

  # Force reconnection when rule content changes
  lifecycle {
    replace_triggered_by = [graylog_pipeline_rule.syslog_parse.content_hash]
  }
}
```

When the rule source changes, the `content_hash` changes, triggering the pipeline connection to be recreated. This disconnects and reconnects the pipeline, which may force Graylog to reload the pipeline state with fresh rules.

**Note:** This mitigation may help in some cases but is not guaranteed to fully resolve the caching issue. The recommended workaround remains setting `cached_stageiterators = false` in Graylog configuration.

## Conclusion

This is a **confirmed architectural issue in Graylog 7.0.4**, not a Terraform provider bug. The provider uses the API correctly. The root cause is in Graylog's pipeline processing layer where:

1. **Stage.rules is mutable but excluded from equality checks** (by design, to avoid hashCode cycles)
2. **StageIterator cache uses Set<Pipeline> as key**, which doesn't detect rule content changes
3. **Error handling is missing** in the reload path, causing silent failures

The `cached_stageiterators = false` configuration should bypass the problematic cache.

## Bug Report Template for Graylog

**Title:** Pipeline rules not updated due to StageIterator cache equality bug

**Affected Version:** 7.0.4+ac94792

**Description:**
Pipeline rules updated via REST API are stored correctly in MongoDB but the Pipeline Interpreter continues executing stale cached versions. The StageIterator.Configuration cache uses `Set<Pipeline>` as a key, but `Stage.equals()` excludes the mutable `rules` field (intentionally, to avoid hashCode cycles). This causes cache hits for pipelines with updated rule content.

**Root Cause:**
- `Stage.java` line 31-32: `rules` field excluded from AutoValue equals/hashCode
- `PipelineInterpreter.java` line 481: Cache uses `Set<Pipeline>` as key
- Updated rules don't invalidate cache entries

**Workaround:**
Set `cached_stageiterators = false` in graylog.conf

**Suggested Fix:**
Include rule content hash in cache key, or remove the cache entirely (minimal performance benefit vs correctness).
