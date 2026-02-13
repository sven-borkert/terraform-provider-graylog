# Terraform Provider Bug Report: sven-borkert/graylog

## Bug 1: Event Notification update fails with "Notification IDs don't match"

### Summary

`graylog_event_notification` resource fails on any update operation with a 400 error. The Graylog API rejects the request because the notification ID is not included in the PUT request body.

### Error

```
Error: failed to update a event notification 698f0f5e0792fb97ec207206: status code: 400,
{"type":"ApiError","message":"Notification IDs don't match"}: status code >= 300

  with graylog_event_notification.email,
  on alerts.tf line 1, in resource "graylog_event_notification" "email":
   1: resource "graylog_event_notification" "email" {
```

### Root Cause

The Graylog REST API `PUT /events/notifications/{notificationId}` requires the `id` field in the JSON request body to match the URL path parameter. The provider's update function does not include the resource ID in the request body when sending the update, causing the API to reject it with "Notification IDs don't match".

### How to Fix

In the resource update function for `graylog_event_notification`, ensure the resource ID is included in the JSON body sent to the API. The typical pattern is:

```go
// Before sending the PUT request, set the ID in the body:
data["id"] = d.Id()
```

Look for the update function in the resource implementation (likely `resource_event_notification.go` or similar). The create function probably works because Graylog generates the ID on creation. The update function needs to echo the ID back.

### Affected API Endpoint

- `PUT /events/notifications/{notificationId}`
- Graylog 7.0+

### Workaround

Taint the resource to force destroy+recreate instead of update:

```bash
terraform taint graylog_event_notification.email
terraform apply
```

---

## Bug 2: Event Definition update fails with "Event definition IDs don't match"

### Summary

`graylog_event_definition` resource fails on any update operation with a 400 error. Same root cause as Bug 1 — the resource ID is not included in the PUT request body.

### Error

```
Error: failed to update a event definition 698f0d5a6e7b08cafd71caec: status code: 400,
{"errors":{"id":["Event definition IDs don't match"]},"error_context":{},"failed":true}: status code >= 300

  with graylog_event_definition.high_error_rate,
  on alerts.tf line 19, in resource "graylog_event_definition" "high_error_rate":
  19: resource "graylog_event_definition" "high_error_rate" {
```

### Root Cause

Same pattern as Bug 1. The Graylog REST API `PUT /events/definitions/{definitionId}` requires the `id` field in the JSON request body. The provider's update function omits it.

### How to Fix

In the resource update function for `graylog_event_definition`, include the resource ID in the body:

```go
data["id"] = d.Id()
```

Look for the update function in the resource implementation (likely `resource_event_definition.go` or similar).

### Affected API Endpoint

- `PUT /events/definitions/{definitionId}`
- Graylog 7.0+

### Workaround

```bash
terraform taint graylog_event_definition.high_error_rate
terraform apply
```

---

## Bug 3: Event Definition series uses wrong field name in docs/config

### Summary

The `series` block inside `graylog_event_definition` config uses `function` as the key name for the aggregation type, but Graylog 7.0+ expects `type`. When `function` is used, the API stores `"type": null` which causes a runtime error: "No series handler registered for: null".

### Error in Graylog UI

```
Event definition High Error Rate failed: No series handler registered for: null.
```

And on the event definition page:

```
Condition is not valid: Function must be set
```

### What Graylog stored via API

```json
{
  "series": [
    {
      "type": null,
      "id": "count-",
      "function": "count"
    }
  ]
}
```

The API accepted `function` but stored `type` as `null`. The event processor looks up the handler by `type`, finds `null`, and fails.

### Correct Configuration

```hcl
series = [{
  id   = "count-"
  type = "count"
}]
```

For series with a field (e.g., avg, min, max):

```hcl
series = [{
  id    = "avg-response_time"
  type  = "avg"
  field = "response_time"
}]
```

### How to Fix

This may be a documentation issue, a schema issue, or both. Options:

1. **Update docs/examples** to use `type` instead of `function` in series blocks
2. **Add provider-level mapping** that translates `function` to `type` before sending to the API if backwards compatibility is desired
3. **Add validation** that warns if `function` is used instead of `type`

### Affected Resource

- `graylog_event_definition` — the `config` JSON's `series` array
- Graylog 7.0+

---

## Environment

- Graylog version: 7.0+
- Terraform provider: sven-borkert/graylog (local dev build)
- Terraform version: (using dev overrides)
- OS: Linux (WSL2)
