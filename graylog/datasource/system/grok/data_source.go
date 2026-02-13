package grok

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"name", "pattern_id"},
				ConflictsWith: []string{"pattern_id"},
			},
			"pattern_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"name", "pattern_id"},
				ConflictsWith: []string{"name"},
			},
			"pattern": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
