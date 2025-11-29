package saved_search

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: create,
		Read:   read,
		Update: update,
		Delete: destroy,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the saved search",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the saved search",
			},
			"summary": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Summary of the saved search",
			},
			"query": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
				Description: "The search query (Lucene syntax)",
			},
			"streams": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Stream IDs to search in",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"timerange_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "relative",
				Description: "Type of time range: relative, absolute, or keyword",
			},
			"timerange_range": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "For relative time range, the number of seconds to look back",
			},
			"selected_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of fields to display in the search results",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sort_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "timestamp",
				Description: "Field to sort results by",
			},
			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Descending",
				Description: "Sort order: Ascending or Descending",
			},

			// Computed
			"search_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the underlying search object",
			},
			"state_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state ID within the view",
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Owner of the saved search",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp",
			},
		},
	}
}
