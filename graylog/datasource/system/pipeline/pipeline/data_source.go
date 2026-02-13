package pipeline

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,
		Schema: map[string]*schema.Schema{
			"pipeline_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"pipeline_id", "title"},
				ConflictsWith: []string{"title"},
			},
			"title": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"pipeline_id", "title"},
				ConflictsWith: []string{"pipeline_id"},
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
