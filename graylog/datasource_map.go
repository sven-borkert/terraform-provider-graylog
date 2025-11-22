package graylog

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/dashboard"
	dashboardwidget "github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/dashboard/widget"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/search/saved"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/sidecar"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/stream"
	streamrule "github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/stream/rule"
	dgrok "github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/grok"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/indices/indexset"
	indextemplate "github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/indices/template"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/input"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/output"
	ppipeline "github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/pipeline/pipeline"
	ppipelinerule "github.com/sven-borkert/terraform-provider-graylog/graylog/datasource/system/pipeline/rule"
)

var dataSourcesMap = map[string]*schema.Resource{
	"graylog_dashboard":           dashboard.DataSource(),
	"graylog_dashboard_widget":    dashboardwidget.DataSource(),
	"graylog_index_set":           indexset.DataSource(),
	"graylog_input":               input.DataSource(),
	"graylog_sidecar":             sidecar.DataSource(),
	"graylog_stream":              stream.DataSource(),
	"graylog_stream_rule":         streamrule.DataSource(),
	"graylog_pipeline":            ppipeline.DataSource(),
	"graylog_pipeline_rule":       ppipelinerule.DataSource(),
	"graylog_saved_search":        saved.DataSource(),
	"graylog_grok_pattern":        dgrok.DataSource(),
	"graylog_grok_patterns":       dgrok.DataSourceList(),
	"graylog_output":              output.DataSource(),
	"graylog_index_set_template":  indextemplate.DataSourceBuiltIn(),
	"graylog_index_set_templates": indextemplate.DataSourceList(),
}
