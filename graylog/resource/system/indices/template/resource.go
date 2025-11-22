package template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Create: create,
		Read:   read,
		Update: update,
		Delete: destroy,
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"index_set_config": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     util.ValidateIsJSON,
				DiffSuppressFunc: util.SchemaDiffSuppressJSONString,
			},
		},
	}
}

func create(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	req := map[string]interface{}{
		"title":            d.Get("title").(string),
		"description":      d.Get("description"),
	}
	cfg, err := convert.StringJSONToData(d.Get("index_set_config").(string))
	if err != nil {
		return err
	}
	req["index_set_config"] = cfg

	body, _, err := cl.IndexSetTemplate.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create index set template: %w", err)
	}
	if id, ok := body["id"].(string); ok {
		d.SetId(id)
	}
	return read(d, m)
}

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	body, resp, err := cl.IndexSetTemplate.Get(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}
	_ = d.Set("title", body["title"])
	_ = d.Set("description", body["description"])
	if cfg, ok := body["index_set_config"]; ok {
		if b, err := json.Marshal(cfg); err == nil {
			_ = d.Set("index_set_config", string(b))
		}
	}
	return nil
}

func update(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	req := map[string]interface{}{
		"title":            d.Get("title").(string),
		"description":      d.Get("description"),
	}
	cfg, err := convert.StringJSONToData(d.Get("index_set_config").(string))
	if err != nil {
		return err
	}
	req["index_set_config"] = cfg

	if _, err := cl.IndexSetTemplate.Update(ctx, d.Id(), req); err != nil {
		return fmt.Errorf("failed to update index set template %s: %w", d.Id(), err)
	}
	return read(d, m)
}

func destroy(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	if _, err := cl.IndexSetTemplate.Delete(ctx, d.Id()); err != nil {
		return fmt.Errorf("failed to delete index set template %s: %w", d.Id(), err)
	}
	return nil
}
