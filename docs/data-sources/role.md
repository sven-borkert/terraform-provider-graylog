# graylog_role Data Source

Use this data source to retrieve information about an existing Graylog role.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/datasource/role)

## Example Usage

### Lookup Built-in Role

```tf
data "graylog_role" "reader" {
  name = "Reader"
}

output "reader_permissions" {
  value = data.graylog_role.reader.permissions
}
```

### Reference from a Role Resource

```tf
resource "graylog_role" "custom" {
  name        = "custom-role"
  description = "A custom role"
  permissions = ["streams:read"]
}

data "graylog_role" "custom" {
  name = graylog_role.custom.name
}

output "custom_role_read_only" {
  value = data.graylog_role.custom.read_only
}
```

## Argument Reference

* `name` - (Required) The name of the role.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The name of the role (same as `name`).
* `description` - The description of the role.
* `permissions` - The permissions assigned to the role.
* `read_only` - Whether the role is read-only (built-in).
