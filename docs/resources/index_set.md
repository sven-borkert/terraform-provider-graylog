# Resource: graylog_index_set

Manages Graylog index sets for storing log data.

* [Source Code](https://github.com/sven-borkert/terraform-provider-graylog/blob/master/graylog/resource/system/indices/indexset/resource.go)

## Example Usage

```hcl
# Get a built-in index set template
data "graylog_index_set_template" "hot7" {
  title = "7 Days Hot"
}

# Create an index set with 7 days retention
resource "graylog_index_set" "application_logs" {
  title                               = "Application Logs"
  description                         = "Index set for application log data"
  index_prefix                        = "applogs"
  rotation_strategy_class             = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
  rotation_strategy                   = jsonencode({
    type               = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  retention_strategy_class            = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
  retention_strategy                  = jsonencode({
    type                  = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig"
    max_number_of_indices = 5
  })
  data_tiering                        = jsonencode({
    type               = "hot_only"
    index_lifetime_min = "P7D"
    index_lifetime_max = "P8D"
  })
  index_analyzer                      = "standard"
  index_set_template_id               = data.graylog_index_set_template.hot7.id
  shards                              = 1
  replicas                            = 0
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
  index_optimization_disabled         = false
  writable                            = true
  use_legacy_rotation                 = false
}
```

## Argument Reference

* `title` - (Required) The title of the Index Set.
* `index_prefix` - (Required, Forces new resource) The index prefix. Must be unique across all index sets.
* `rotation_strategy_class` - (Required) The rotation strategy class. Common values:
  - `org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy` - Time-based with size optimization (recommended)
  - `org.graylog2.indexer.rotation.strategies.TimeBasedRotationStrategy` - Time-based rotation
  - `org.graylog2.indexer.rotation.strategies.SizeBasedRotationStrategy` - Size-based rotation
  - `org.graylog2.indexer.rotation.strategies.MessageCountRotationStrategy` - Message count based
* `rotation_strategy` - (Required) JSON string with rotation strategy configuration.
* `retention_strategy_class` - (Required) The retention strategy class. Common values:
  - `org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy` - Delete old indices
  - `org.graylog2.indexer.retention.strategies.ClosingRetentionStrategy` - Close old indices
  - `org.graylog2.indexer.retention.strategies.NoopRetentionStrategy` - Keep all indices
* `retention_strategy` - (Required) JSON string with retention strategy configuration.
* `index_analyzer` - (Required) The Elasticsearch/OpenSearch analyzer. Usually `"standard"`.
* `shards` - (Required) Number of shards per index.
* `description` - (Optional) Description of the Index Set.
* `replicas` - (Optional) Number of replicas. Default: 0.
* `index_optimization_disabled` - (Optional) Whether to disable index optimization. Default: false.
* `index_optimization_max_num_segments` - (Required) Maximum number of segments after optimization.
* `default` - (Optional) Whether this is the default index set. Default: false.
* `field_type_refresh_interval` - (Optional) Field type refresh interval in milliseconds.
* `writable` - (Optional) Whether the index set is writable. Default: true.
* `use_legacy_rotation` - (Optional) Use legacy rotation. Default: false for Graylog 7.
* `data_tiering` - (Optional) JSON string with data tiering configuration for Graylog 7.
* `index_set_template_id` - (Optional) ID of an index set template to use.
* `field_restrictions` - (Optional) JSON string with field restrictions.

## Attributes Reference

* `id` - The Index Set ID.
* `creation_date` - The date time when the Index Set was created.

## Deflector Initialization

When an index set is created, Graylog initializes a deflector alias that points to the active write index. The provider waits for the deflector to be ready before completing the create operation, ensuring dependent resources (like streams) can immediately route data to the index set.

## Import

`graylog_index_set` can be imported using the Index Set ID:

```console
$ terraform import graylog_index_set.example 5c4acaefc9e77bbbbbbbbbbb
```
