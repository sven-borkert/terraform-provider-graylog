package dashboard

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,

		Schema: map[string]*schema.Schema{
			"dashboard_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ExactlyOneOf:  []string{"dashboard_id", "title"},
				ConflictsWith: []string{"title"},
				Description:   "The ID of the dashboard. Either dashboard_id or title must be set.",
			},
			"title": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ExactlyOneOf:  []string{"dashboard_id", "title"},
				ConflictsWith: []string{"dashboard_id"},
				Description:   "The title of the dashboard. Either dashboard_id or title must be set.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the dashboard.",
			},
			"summary": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The summary of the dashboard.",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the dashboard.",
			},
			"search_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The search ID associated with the dashboard.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation timestamp of the dashboard.",
			},
		},
	}
}
