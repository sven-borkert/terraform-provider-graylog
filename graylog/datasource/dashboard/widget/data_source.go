package widget

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,
		Schema: map[string]*schema.Schema{
			"dashboard_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"widget_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cache_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
