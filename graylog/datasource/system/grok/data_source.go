package grok

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pattern_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pattern": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
