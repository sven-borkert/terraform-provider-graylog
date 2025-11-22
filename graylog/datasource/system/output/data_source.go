package output

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// DataSource exposes Graylog outputs for lookup by id or title.
func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,

		Schema: map[string]*schema.Schema{
			"output_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"output_id", "title"},
				ConflictsWith: []string{"title"},
			},
			"title": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"output_id", "title"},
				ConflictsWith: []string{"output_id"},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"configuration": {
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
			"content_pack": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
