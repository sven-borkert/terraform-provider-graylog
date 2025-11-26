# graylog_grok_pattern Resource

Use this resource to manage Graylog grok patterns.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/system/grok)

## Example Usage

### Basic Grok Pattern

```tf
resource "graylog_grok_pattern" "custom_date" {
  name    = "CUSTOM_DATE"
  pattern = "\\d{4}-\\d{2}-\\d{2}"
}
```

### Pattern Using Other Patterns

```tf
resource "graylog_grok_pattern" "custom_timestamp" {
  name    = "CUSTOM_TIMESTAMP"
  pattern = "%%{DATE}[T ]%%{TIME}"
}
```

## Argument Reference

* `name` - (Required) The name of the grok pattern.
* `pattern` - (Required) The grok pattern definition.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the grok pattern.

## Import

Grok patterns can be imported using the pattern ID:

```
$ terraform import graylog_grok_pattern.example 5c4acaefc9e77bbbbbbbbbbb
```

## Notes

- Content packs are not currently supported for grok patterns.
- When using references to other patterns with `%{`, you must escape them as `%%{` in HCL:
  ```tf
  pattern = "%%{DATE}[- ]%%{TIME}"  # References DATE and TIME patterns
  ```
