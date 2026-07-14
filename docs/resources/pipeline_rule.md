# Resource: graylog_pipeline_rule

* [Example](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/examples/v0.12/pipeline.tf)
* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/system/pipeline/rule/resource.go)

## Argument Reference

* `source` - (Required) The source of the Pipeline Rule. The data type is `string`.
* `description` - (Optional) description of the Pipeline Rule. The data type is `string`.

## Attributes Reference

Nothing.

## Known limitations

This resource manages the rule `source` only. Graylog 7 stores an additional
`rule_builder` representation for rules created with the UI's visual rule builder.
Because the provider does not model `rule_builder`, updating a rule that was
originally created in the rule builder clears that visual representation on the
server (the rule continues to work from its `source`). Manage such rules entirely
in Terraform, or edit them only in the Graylog UI.

## Import

`graylog_pipeline_rule` can be imported using the Pipeline Rule id, e.g.

```console
$ terraform import graylog_pipeline_rule.test 5c4acaefc9e77bbbbbbbbbbb
```
