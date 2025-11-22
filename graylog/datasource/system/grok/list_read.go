package grok

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func readList(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	data, _, err := cl.Grok.Gets(ctx)
	if err != nil {
		return err
	}

	patterns, ok := data["patterns"]
	if !ok {
		patterns = []interface{}{}
	}
	j, err := json.Marshal(patterns)
	if err != nil {
		return err
	}
	_ = d.Set("patterns_json", string(j))
	d.SetId("all")
	return nil
}
