# Resource: graylog_sidecars

* [Example](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/examples/v0.12/sidecar.tf)
* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/sidecar/resource.go)

Manages configuration assignments for the node IDs listed in `sidecars`.
Refresh does not adopt other nodes. Removing a managed node or destroying this
resource clears its assignments. Use a single resource for the assignments you
want Terraform to manage.

Good

```hcl
resource "graylog_sidecars" "all" {
  sidecars {
    # ...
  }
}
```

NG

```hcl
resource "graylog_sidecars" "foo" {
  sidecars {
    # ...
  }
}

resource "graylog_sidecars" "bar" {
  sidecars {
    # ...
  }
}
```

## Argument Reference

* `sidecars` - (Required) The data type is `[]object (set)`
* `sidecars[].node_id` - (Required) The data type is `string`
* `sidecars[].assignments` - (Required) The data type is `[]object (set)`
* `sidecars[].assignments[].collector_id` - (Required) The data type is `string`
* `sidecars[].assignments[].configuration_id` - (Required) The data type is `string`

## Attributes Reference

* `assignment_scope_version` - Internal ownership tracking, computed by the provider.

## Upgrading existing state

Older provider versions could copy unrelated nodes into this resource's state.
The provider rejects legacy state without ownership tracking to avoid clearing
assignments on those nodes. Re-establish ownership from configuration:

1. Back up state with `terraform state pull > sidecar-state-backup.json` and keep the backup private.
2. Check that the resource configuration lists only the nodes and assignments you intend to manage.
3. Run `terraform state rm graylog_sidecars.all`, replacing the address with your resource address. This removes Terraform's tracking entry; it does not change Graylog.
4. Run `terraform plan`, review it, then `terraform apply`. The create operation assigns configuration only to the listed nodes and records their ownership.

Use this procedure before changing membership or destroying a legacy resource.

## Import

Unlike other resources, the given ID is ignored so please specify any string as ID.
Import explicitly adopts all sidecars currently returned by Graylog, including
their assignments. Use the upgrade procedure above when managing only a subset.

e.g.

```console
$ terraform import graylog_sidecars.all all
```
