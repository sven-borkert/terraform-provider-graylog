package role

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Read: read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the role.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the role.",
			},
			"permissions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The permissions assigned to the role.",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the role is read-only (built-in).",
			},
		},
	}
}
