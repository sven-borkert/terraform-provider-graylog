# Graylog 7 Index Mapping Technical Report

## Executive Summary

This report details how Graylog 7 creates and manages Elasticsearch/OpenSearch index mappings, specifically addressing why timestamp fields may incorrectly be stored as `text` instead of `date` type. Understanding this mechanism is critical for building reliable REST API configurations.

---

## 1. How Graylog Creates Index Mappings

### 1.1 The Two-Layer Mapping System

Graylog uses a **two-layer mapping system**:

1. **Graylog Index Templates**: Graylog creates its own index templates (e.g., `graylog-internal`) that define base mappings for known fields
2. **Elasticsearch Dynamic Mapping**: For fields not explicitly defined, Elasticsearch's dynamic mapping infers types from the first document

### 1.2 When Mappings Are Created

Mappings are established at **two critical moments**:

#### A. Index Creation Time (Template Application)
- When Graylog creates a new index (during rotation), it applies its internal index template
- The template defines mappings for Graylog's reserved fields (`timestamp`, `message`, `source`, `gl2_*`, etc.)
- Custom field type profiles and mappings configured in Graylog are merged into this template

```
Index Rotation → New Index Created → Template Applied → Base Mappings Set
```

#### B. First Document Indexing (Dynamic Mapping)
- For any field **NOT** covered by the template, Elasticsearch uses dynamic mapping
- **The first document containing a new field determines its type permanently for that index**
- This is the root cause of the timestamp-as-text problem

```
New Field Encountered → Elasticsearch Analyzes First Value → Type Locked for Index
```

### 1.3 The Timestamp Problem Explained

**Scenario: Why a timestamp field gets stored as `text`:**

1. A new index is created via rotation
2. Graylog's template defines the `timestamp` field (Graylog's internal field) as `date`
3. Your custom field (e.g., `my_timestamp`, `event_time`, `log_time`) is **NOT** in the template
4. The first log message containing this field arrives
5. If the value doesn't match Elasticsearch's date detection patterns, it's mapped as `text`

**Example:**
```json
// First message arrives with:
{"my_timestamp": "28/Nov/2025 14:30:00"}

// Elasticsearch cannot parse this format → maps as "text"
// All subsequent messages with proper ISO dates STILL get mapped as "text"
```

---

## 2. Elasticsearch Dynamic Mapping Rules

### 2.1 Automatic Type Detection

Elasticsearch uses these rules for dynamic field detection:

| JSON Type | Elasticsearch Type |
|-----------|-------------------|
| `null` | No field added |
| `true`/`false` | `boolean` |
| Floating point number | `float` |
| Integer | `long` |
| Object | `object` |
| Array | Depends on first non-null element |
| String (passes date detection) | `date` |
| String (passes numeric detection) | `float` or `long` |
| String (other) | `text` with `keyword` subfield |

### 2.2 Date Detection

By default, Elasticsearch has date detection **enabled**. It checks strings against these patterns:

```
strict_date_optional_time || epoch_millis
```

This means it recognizes:
- ISO 8601: `2025-11-28T14:30:00.000Z`
- Epoch milliseconds: `1732802400000`

**It does NOT recognize by default:**
- `28/Nov/2025 14:30:00` (Apache/HTTPD format)
- `Nov 28, 2025 2:30 PM`
- `28-11-2025 14:30:00`
- Unix timestamps in seconds: `1732802400`

### 2.3 The "First Document Wins" Problem

```
Index: graylog_42 (new)
├── Message 1: {"custom_time": "invalid-date-format"} → custom_time: text
├── Message 2: {"custom_time": "2025-11-28T14:30:00Z"} → STILL text (mapping locked)
└── Message 3: {"custom_time": "2025-11-28T15:00:00Z"} → STILL text (mapping locked)
```

**Mappings cannot be changed for existing indices** - only new indices get new mappings.

---

## 3. Graylog 7 Field Type Management

### 3.1 Graylog's Field Type System

Graylog 7 introduced native field type management with these components:

1. **Index Set Field Type Configuration**: Per-index-set field mappings
2. **Field Type Profiles**: Reusable sets of field mappings
3. **Custom Field Mappings**: Override dynamic mapping for specific fields

### 3.2 How Graylog Updates Mappings

When you change a field type in Graylog:

1. Graylog stores the mapping in MongoDB
2. The index template is updated
3. **Changes only apply to NEW indices** (after rotation)
4. Existing indices keep their original mappings

```
Field Type Change → Template Updated → Index Rotation Required → New Index Gets New Mapping
```

### 3.3 The Field Type Refresh Interval

Graylog periodically refreshes field type information from the active write index. This is configured per index set:

- **Setting**: `Field type refresh interval`
- **Default**: Updates every few seconds/minutes
- **Purpose**: Synchronizes Graylog's field type cache with Elasticsearch

---

## 4. Solutions for the Timestamp Problem

### 4.1 Solution 1: Pre-define Field Types in Graylog (Recommended)

**Before any data arrives**, configure the field type in Graylog:

#### Via UI:
1. Navigate to `System → Indices`
2. Select your index set
3. Click `Configuration → Configure index field types`
4. Add mapping: `my_timestamp` → `date`
5. Rotate the index

#### Via REST API:
```bash
# Get index set ID
curl -u admin:password \
  -H 'Accept: application/json' \
  'http://graylog:9000/api/system/indices/index_sets'

# The field type API endpoints (check API browser for exact paths):
# GET /api/system/indices/index_sets/types/{indexSetId}
# PUT /api/system/indices/index_sets/types/{indexSetId}
```

### 4.2 Solution 2: Create Field Type Profiles

For consistent mappings across multiple index sets:

```bash
# Create a profile via UI or API
# Profile contains: {"my_timestamp": "date", "event_time": "date", ...}
# Assign profile to index sets
# Rotate all affected indices
```

### 4.3 Solution 3: Custom Elasticsearch Index Template

Create a custom template that Elasticsearch applies alongside Graylog's template:

```bash
# Create custom template file
cat > graylog-custom-mapping.json << 'EOF'
{
  "index_patterns": ["graylog_*", "your_prefix_*"],
  "priority": 1,
  "template": {
    "mappings": {
      "properties": {
        "my_timestamp": {
          "type": "date",
          "format": "strict_date_optional_time||epoch_millis||dd/MMM/yyyy:HH:mm:ss"
        },
        "event_time": {
          "type": "date"
        },
        "log_time": {
          "type": "date",
          "format": "yyyy-MM-dd HH:mm:ss.SSS||strict_date_optional_time"
        }
      }
    }
  }
}
EOF

# Apply template to Elasticsearch/OpenSearch
curl -X PUT \
  -H 'Content-Type: application/json' \
  -d @graylog-custom-mapping.json \
  'http://elasticsearch:9200/_index_template/graylog-custom-mapping'

# Rotate index to apply
curl -X POST \
  -u admin:password \
  -H 'X-Requested-By: api' \
  'http://graylog:9000/api/system/deflector/cycle'
```

### 4.4 Solution 4: Pipeline Processing (Normalize Before Indexing)

Use Graylog pipelines to normalize date formats BEFORE indexing:

```
rule "normalize timestamp"
when
  has_field("my_timestamp")
then
  let parsed = parse_date(to_string($message.my_timestamp), "dd/MMM/yyyy:HH:mm:ss");
  set_field("my_timestamp", parsed);
end
```

---

## 5. REST API Reference for Field Type Management

### 5.1 Authentication

```bash
# Using access token (recommended)
curl -u YOUR_TOKEN:token -H 'Accept: application/json' ...

# Using session
curl -u admin:password -H 'Accept: application/json' ...
```

### 5.2 Index Set Management

```bash
# List all index sets
GET /api/system/indices/index_sets

# Get specific index set
GET /api/system/indices/index_sets/{indexSetId}

# Update index set (including profile assignment)
PUT /api/system/indices/index_sets/{indexSetId}
```

### 5.3 Index Rotation

```bash
# Rotate default index set
POST /api/system/deflector/cycle
Header: X-Requested-By: api

# Rotate specific index set
POST /api/system/deflector/{indexSetId}/cycle
Header: X-Requested-By: api
```

### 5.4 Field Type Configuration

```bash
# Get field types for index set
GET /api/system/indices/index_sets/types/{indexSetId}

# Note: The exact API endpoints for field type manipulation should be 
# verified via the API browser at:
# http://your-graylog:9000/api/api-browser
```

### 5.5 Template Management (Direct Elasticsearch)

```bash
# List all templates
GET http://elasticsearch:9200/_index_template

# Get Graylog's internal template
GET http://elasticsearch:9200/_index_template/graylog-internal*

# Create/update custom template
PUT http://elasticsearch:9200/_index_template/custom-template-name

# Delete custom template
DELETE http://elasticsearch:9200/_index_template/custom-template-name
```

---

## 6. Best Practices for Coding Agents

### 6.1 Pre-flight Checks

Before creating configurations:

1. **Check existing mappings**: Query the current index mapping to see what types are already set
2. **Check existing templates**: List templates to avoid conflicts
3. **Get index set IDs**: You'll need these for most operations

```bash
# Check current mapping for active index
curl 'http://elasticsearch:9200/graylog_deflector/_mapping?pretty'

# List all templates
curl 'http://elasticsearch:9200/_index_template?pretty'
```

### 6.2 Creating Reliable Configurations

```python
# Pseudocode for reliable field type configuration

def ensure_field_type(graylog_api, es_api, field_name, field_type, index_set_id):
    """
    Ensures a field has the correct type in both Graylog and Elasticsearch.
    
    Steps:
    1. Check if field already exists in current mapping
    2. If wrong type exists, log warning (cannot change existing index)
    3. Update Graylog field type configuration
    4. Create/update Elasticsearch template
    5. Trigger index rotation
    6. Verify new index has correct mapping
    """
    
    # Step 1: Check current state
    current_mapping = es_api.get(f"/{index_set_id}_deflector/_mapping")
    
    # Step 2: Warn if mismatch in existing index
    if field_exists_with_wrong_type(current_mapping, field_name, field_type):
        logger.warning(f"Field {field_name} already mapped as wrong type in current index")
    
    # Step 3: Update Graylog configuration
    graylog_api.put(
        f"/api/system/indices/index_sets/types/{index_set_id}",
        json={"field_name": field_name, "type": field_type}
    )
    
    # Step 4: Update Elasticsearch template (belt and suspenders)
    template = build_template(field_name, field_type)
    es_api.put(f"/_index_template/custom-{field_name}", json=template)
    
    # Step 5: Rotate index
    graylog_api.post(
        f"/api/system/deflector/{index_set_id}/cycle",
        headers={"X-Requested-By": "api"}
    )
    
    # Step 6: Wait and verify
    time.sleep(5)  # Wait for rotation
    new_mapping = es_api.get(f"/{index_set_id}_deflector/_mapping")
    return verify_field_type(new_mapping, field_name, field_type)
```

### 6.3 Template Priority

When multiple templates match an index pattern, priority determines which wins:

- Higher number = higher priority
- Graylog's internal templates typically have default priority
- Set your custom templates with `"priority": 1` or higher if you need to override

### 6.4 Common Date Formats

For maximum compatibility, use multiple formats in date field definitions:

```json
{
  "type": "date",
  "format": "strict_date_optional_time||epoch_millis||epoch_second||yyyy-MM-dd HH:mm:ss||dd/MMM/yyyy:HH:mm:ss Z"
}
```

### 6.5 Error Handling

Watch for these common errors:

| Error | Cause | Solution |
|-------|-------|----------|
| `mapper_parsing_exception` | Data doesn't match mapping | Fix data format or change mapping type |
| `illegal_argument_exception` | Invalid date format | Add format to date field definition |
| `strict_dynamic_mapping_exception` | Field not in mapping (strict mode) | Add field to template |
| `index_not_found_exception` | Index doesn't exist yet | Wait for first document or pre-create |

---

## 7. Verification Commands

### 7.1 Check Current Mapping

```bash
# Active index mapping
curl 'http://elasticsearch:9200/graylog_deflector/_mapping?pretty' | jq '.[] .mappings.properties.YOUR_FIELD'

# Specific index
curl 'http://elasticsearch:9200/graylog_42/_mapping?pretty'
```

### 7.2 Check Templates

```bash
# All templates
curl 'http://elasticsearch:9200/_index_template?pretty'

# Specific template
curl 'http://elasticsearch:9200/_index_template/graylog-internal?pretty'
```

### 7.3 Check Graylog Field Types

```bash
# Via Graylog API
curl -u admin:password 'http://graylog:9000/api/system/indices/index_sets/types/{indexSetId}?pretty'
```

### 7.4 Test Date Parsing

```bash
# Test if a date string can be parsed
curl -X POST 'http://elasticsearch:9200/_analyze' \
  -H 'Content-Type: application/json' \
  -d '{"analyzer": "standard", "text": "2025-11-28T14:30:00Z"}'
```

---

## 8. Summary: Key Points for Implementation

1. **Field types are determined by the FIRST document** if not pre-defined
2. **Mappings cannot be changed** on existing indices - only new ones
3. **Always pre-define date fields** before data arrives
4. **Use both Graylog field types AND Elasticsearch templates** for redundancy
5. **Index rotation is required** after any mapping change
6. **Use pipeline rules** to normalize date formats before indexing
7. **Verify mappings** after rotation to confirm changes applied

---

## 9. Quick Reference: API Endpoints

| Operation | Method | Endpoint |
|-----------|--------|----------|
| List index sets | GET | `/api/system/indices/index_sets` |
| Get index set | GET | `/api/system/indices/index_sets/{id}` |
| Update index set | PUT | `/api/system/indices/index_sets/{id}` |
| Rotate index | POST | `/api/system/deflector/{id}/cycle` |
| Get field types | GET | `/api/system/indices/index_sets/types/{id}` |
| ES: Get mapping | GET | `http://es:9200/{index}/_mapping` |
| ES: Create template | PUT | `http://es:9200/_index_template/{name}` |
| ES: Delete template | DELETE | `http://es:9200/_index_template/{name}` |

---

*Report generated for Graylog 7.x with Elasticsearch/OpenSearch backend*
*Last updated: November 2025*
