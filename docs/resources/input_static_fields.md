# graylog_input_static_fields Resource

Use this resource to manage static fields on a Graylog input. Static fields are added to every message that comes through the input.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/system/input/staticfield)

## Example Usage

### Basic Static Fields

```tf
resource "graylog_input" "gelf_udp" {
  title  = "GELF UDP"
  type   = "org.graylog2.inputs.gelf.udp.GELFUDPInput"
  global = true

  attributes = jsonencode({
    bind_address = "0.0.0.0"
    port         = 12201
  })
}

resource "graylog_input_static_fields" "gelf_udp" {
  input_id = graylog_input.gelf_udp.id
  fields = {
    environment = "production"
    source_type = "application"
  }
}
```

### Multiple Static Fields

```tf
resource "graylog_input_static_fields" "example" {
  input_id = graylog_input.example.id
  fields = {
    datacenter = "us-east-1"
    team       = "platform"
    service    = "api-gateway"
    version    = "2.0"
  }
}
```

## Argument Reference

* `input_id` - (Required, Forces new resource) The ID of the input to add static fields to.
* `fields` - (Optional) A map of static field names to values. All values must be strings.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the input (same as `input_id`).

## Import

Input static fields can be imported using the input ID:

```
$ terraform import graylog_input_static_fields.example 5c4acaefc9e77bbbbbbbbbbb
```

## Notes

- Static fields are added to every message received by the input.
- This resource manages all static fields for an input. If you need to manage static fields separately, use multiple inputs.
- Changing the `input_id` will force recreation of the resource.
