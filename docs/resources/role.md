# graylog_role Resource

Use this resource to manage Graylog roles.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/role)

## Example Usage

### Basic Role

```tf
resource "graylog_role" "example" {
  name        = "custom-reader"
  description = "Custom read-only role"
  permissions = ["streams:read"]
}
```

### Role with Multiple Permissions

```tf
resource "graylog_role" "operator" {
  name        = "stream-operator"
  description = "Role with stream read and write permissions"
  permissions = [
    "streams:read",
    "streams:edit",
    "streams:changestate"
  ]
}
```

## Argument Reference

* `name` - (Required) The name of the role. This is used as the identifier.
* `permissions` - (Required) Set of permissions assigned to the role.
* `description` - (Optional) A description of the role.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `read_only` - Whether the role is read-only (built-in).

## Import

Roles can be imported using the role name:

```
$ terraform import graylog_role.example custom-reader
```

## Notes

- Built-in roles (e.g., `Admin`, `Reader`) have `read_only = true` and cannot be modified or deleted.
- Role names must be unique within a Graylog instance.
