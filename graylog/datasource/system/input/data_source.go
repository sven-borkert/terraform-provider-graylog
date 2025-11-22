package input

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,

		Schema: map[string]*schema.Schema{
			"input_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"input_id", "title"},
				ConflictsWith: []string{"title"},
			},
			"title": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"input_id", "title"},
				ConflictsWith: []string{"input_id"},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
