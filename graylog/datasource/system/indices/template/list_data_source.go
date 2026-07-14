package template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func DataSourceList() *schema.Resource {
	return &schema.Resource{
		Read: readList,
		Schema: map[string]*schema.Schema{
			// Optional: include built-in or custom filtering later
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"built_in": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readList(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	builtIns, _, err := cl.IndexSetTemplate.BuiltIns(ctx, nil)
	if err != nil {
		return err
	}
	custom, _, err := cl.IndexSetTemplate.Paginated(ctx)
	if err != nil {
		return err
	}

	seen := map[string]bool{}
	var result []map[string]interface{}
	appendTemplate := func(t map[string]interface{}) {
		if id, ok := t["id"].(string); ok && id != "" {
			if seen[id] {
				return
			}
			seen[id] = true
		}
		result = append(result, map[string]interface{}{
			"id":          t["id"],
			"title":       t["title"],
			"description": t["description"],
			"built_in":    t["built_in"],
			"default":     t["default"],
		})
	}
	for _, t := range builtIns {
		appendTemplate(t)
	}
	for _, t := range custom {
		appendTemplate(t)
	}

	if err := d.Set("templates", result); err != nil {
		return err
	}
	d.SetId("templates")
	return nil
}
