package rule

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
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// content_hash is a SHA256 hash of the rule source.
			// Use this with lifecycle { replace_triggered_by } on pipeline connections
			// to force cache invalidation when rules change.
			// This works around a caching bug in Graylog 7.0.4 where updated rules
			// may not take effect due to stale StageIterator cache entries.
			"content_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SHA256 hash of the rule source. Use with replace_triggered_by to refresh pipeline connections on rule changes.",
			},

			// We don't define the attribute "title",
			// because the request parameter "title" is ignored in create and update API.
		},
	}
}
