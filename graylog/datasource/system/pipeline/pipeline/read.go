package pipeline

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	if id, ok := d.GetOk("pipeline_id"); ok {
		data, resp, err := cl.Pipeline.Get(ctx, id.(string))
		if err != nil {
			return err
		}
		return setDataToResourceData(d, data, resp)
	}

	if title, ok := d.GetOk("title"); ok {
		list, _, err := cl.Pipeline.Gets(ctx)
		if err != nil {
			return err
		}
		var hit map[string]interface{}
		matches := 0
		for _, p := range list["pipelines"].([]interface{}) {
			pm := p.(map[string]interface{})
			if pm["title"] == title {
				hit = pm
				matches++
				if matches > 1 {
					return errors.New("matched multiple pipelines; title not unique")
				}
			}
		}
		if matches == 0 {
			return errors.New("matched pipeline is not found")
		}
		return setDataToResourceData(d, hit, nil)
	}

	return errors.New("one of pipeline_id or title must be set")
}
