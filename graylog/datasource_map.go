package graylog

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/dashboard"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/sidecar"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/stream"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/indices/indexset"
	indextemplate "github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/indices/template"
)

var dataSourcesMap = map[string]*schema.Resource{
	"graylog_dashboard": dashboard.DataSource(),
	"graylog_index_set": indexset.DataSource(),
	"graylog_sidecar":   sidecar.DataSource(),
	"graylog_stream":    stream.DataSource(),
	"graylog_index_set_template": indextemplate.DataSourceBuiltIn(),
	"graylog_index_set_templates": indextemplate.DataSourceList(),
}
