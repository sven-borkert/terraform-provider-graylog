# graylog_grok_pattern Data Source

Use this data source to retrieve information about an existing Graylog grok pattern.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/datasource/system/grok)

## Example Usage

### Lookup by Name

```tf
data "graylog_grok_pattern" "year" {
  name = "YEAR"
}

output "year_pattern" {
  value = data.graylog_grok_pattern.year.pattern
}
```

### Lookup by Pattern ID

```tf
data "graylog_grok_pattern" "by_id" {
  pattern_id = "691e31ef945ebfd39f28bbbb"
}

output "pattern_name" {
  value = data.graylog_grok_pattern.by_id.name
}
```

### Reference from a Grok Pattern Resource

```tf
resource "graylog_grok_pattern" "custom" {
  name    = "CUSTOM_PATTERN"
  pattern = "\\d{4}-\\d{2}-\\d{2}"
}

data "graylog_grok_pattern" "custom" {
  name = graylog_grok_pattern.custom.name
}

output "custom_pattern_id" {
  value = data.graylog_grok_pattern.custom.pattern_id
}
```

## Argument Reference

One of `pattern_id` or `name` must be set. Both cannot be set at the same time.

* `pattern_id` - (Optional) The ID of the grok pattern.
* `name` - (Optional) The name of the grok pattern.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the grok pattern (same as `pattern_id`).
* `pattern_id` - The ID of the grok pattern.
* `name` - The name of the grok pattern.
* `pattern` - The grok pattern definition.
