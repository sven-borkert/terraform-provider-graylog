package dashboard

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func schemaMap() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Type:     schema.TypeString,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"summary": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"search_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"state": schemaState,
		"owner": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

var schemaState = &schema.Schema{
	Type:     schema.TypeList,
	Required: true,
	MinItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id":           {Type: schema.TypeString, Optional: true, Computed: true},
			"tab_title":    {Type: schema.TypeString, Optional: true},
			"query_string": {Type: schema.TypeString, Optional: true, Default: ""},
			"widgets":      schemaWidgets,
			"widget_mapping": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: util.SchemaDiffSuppressJSONString,
				ValidateFunc:     util.ValidateIsJSON,
			},
			"positions": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: util.SchemaDiffSuppressJSONString,
				ValidateFunc:     util.ValidateIsJSON,
			},
			"titles": schemaTitles,
		},
	},
}

var schemaTitles = &schema.Schema{
	Type:             schema.TypeString,
	Optional:         true,
	DiffSuppressFunc: util.SchemaDiffSuppressJSONString,
	ValidateFunc:     util.ValidateIsJSON,
}

var schemaQuery = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	MinItems: 1,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type":         {Type: schema.TypeString, Required: true},
			"query_string": {Type: schema.TypeString, Required: true},
		},
	},
}

var schemaWidgets = &schema.Schema{
	Type:     schema.TypeList,
	Required: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"widget_id": {Type: schema.TypeString, Optional: true},
			"type":      {Type: schema.TypeString, Required: true},
			"config": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: util.SchemaDiffSuppressJSONString,
				ValidateFunc:     util.ValidateIsJSON,
			},
			"timerange": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: util.SchemaDiffSuppressJSONString,
				ValidateFunc:     util.ValidateIsJSON,
			},
			"query": schemaQuery,
			"streams": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}
