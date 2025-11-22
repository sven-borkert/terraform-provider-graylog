package template

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func DataSourceBuiltIn() *schema.Resource {
	return &schema.Resource{
		Read: readBuiltIn,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
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

func readBuiltIn(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	title := d.Get("title").(string)
	templates, _, err := cl.IndexSetTemplate.BuiltIns(ctx, nil)
	if err != nil {
		return err
	}
	for _, t := range templates {
		if t["title"] == title {
			id := t["id"].(string)
			d.SetId(id)
			_ = d.Set("id", id)
			_ = d.Set("description", t["description"])
			return nil
		}
	}
	return errors.New("built-in template with given title not found")
}
