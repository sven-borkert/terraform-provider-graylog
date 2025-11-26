# graylog_user Resource

Use this resource to manage Graylog users.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/user)

## Example Usage

### Basic User

```tf
resource "graylog_user" "example" {
  username   = "example-user"
  email      = "example@example.com"
  first_name = "Example"
  last_name  = "User"
  password   = "SecurePassword123"
  roles      = ["Reader"]
}
```

### User with Custom Settings

```tf
resource "graylog_user" "operator" {
  username           = "operator"
  email              = "operator@example.com"
  first_name         = "System"
  last_name          = "Operator"
  password           = "SecurePassword456"
  roles              = ["Reader", "Dashboard Creator"]
  timezone           = "Europe/Berlin"
  session_timeout_ms = 7200000  # 2 hours
}
```

### Service Account

```tf
resource "graylog_user" "service" {
  username        = "api-service"
  email           = "api-service@example.com"
  first_name      = "API"
  last_name       = "Service"
  password        = "ServicePassword789"
  roles           = ["Reader"]
  service_account = true
}
```

## Argument Reference

* `username` - (Required, Forces new resource) The username. Cannot be changed after creation.
* `email` - (Required) The email address of the user.
* `first_name` - (Required) The first name of the user.
* `last_name` - (Required) The last name of the user.
* `password` - (Optional, Sensitive) The password. Required when creating a new user.
* `roles` - (Optional) Set of roles assigned to the user.
* `permissions` - (Optional) Set of permissions assigned to the user.
* `timezone` - (Optional) The timezone of the user. Defaults to server timezone.
* `session_timeout_ms` - (Optional) Session timeout in milliseconds. Defaults to `3600000` (1 hour).
* `service_account` - (Optional) Whether this is a service account.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `user_id` - The ID of the user.
* `full_name` - The full name (computed from first_name and last_name).
* `external` - Whether the user is external (e.g., LDAP).
* `read_only` - Whether the user is read-only.
* `session_active` - Whether the user has an active session.
* `last_activity` - The timestamp of the user's last activity.
* `client_address` - The client address of the user's last session.
* `account_status` - The account status (enabled, disabled, deleted).

## Import

Users can be imported using the username:

```
$ terraform import graylog_user.example example-user
```

## Notes

- The `password` attribute is write-only. It cannot be read back from the API.
- The `full_name` attribute is computed from `first_name` and `last_name`.
- Users with `read_only = true` cannot be modified or deleted (e.g., built-in users).
