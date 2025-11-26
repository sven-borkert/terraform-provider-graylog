# graylog_user Data Source

Use this data source to retrieve information about an existing Graylog user.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/datasource/user)

## Example Usage

### Lookup by Username

```tf
data "graylog_user" "admin" {
  username = "admin"
}

output "admin_roles" {
  value = data.graylog_user.admin.roles
}
```

### Lookup by User ID

```tf
data "graylog_user" "by_id" {
  user_id = "6926e4342562186bc3ea1a26"
}

output "user_email" {
  value = data.graylog_user.by_id.email
}
```

### Reference from a User Resource

```tf
resource "graylog_user" "example" {
  username   = "example-user"
  email      = "example@example.com"
  first_name = "Example"
  last_name  = "User"
  password   = "SecurePassword123"
  roles      = ["Reader"]
}

data "graylog_user" "example" {
  username = graylog_user.example.username
}

output "user_full_name" {
  value = data.graylog_user.example.full_name
}
```

## Argument Reference

One of `user_id` or `username` must be set. Both cannot be set at the same time.

* `user_id` - (Optional) The ID of the user.
* `username` - (Optional) The username.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the user (same as `user_id`).
* `username` - The username.
* `email` - The email address of the user.
* `first_name` - The first name of the user.
* `last_name` - The last name of the user.
* `full_name` - The full name of the user.
* `timezone` - The timezone of the user.
* `session_timeout_ms` - The session timeout in milliseconds.
* `roles` - The roles assigned to the user.
* `permissions` - The permissions assigned to the user.
* `external` - Whether the user is external (e.g., LDAP).
* `read_only` - Whether the user is read-only.
* `session_active` - Whether the user has an active session.
* `last_activity` - The timestamp of the user's last activity.
* `client_address` - The client address of the user's last session.
* `account_status` - The account status (enabled, disabled, deleted).
* `service_account` - Whether this is a service account.
