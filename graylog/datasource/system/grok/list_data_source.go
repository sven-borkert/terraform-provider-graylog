package grok

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSourceList() *schema.Resource {
	return &schema.Resource{
		Read: readList,
		Schema: map[string]*schema.Schema{
			"patterns_json": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
