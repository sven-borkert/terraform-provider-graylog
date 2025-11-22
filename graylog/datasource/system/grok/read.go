package grok

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

	if id, ok := d.GetOk("pattern_id"); ok {
		data, resp, err := cl.Grok.Get(ctx, id.(string))
		if err != nil {
			return err
		}
		return setDataToResourceData(d, data, resp)
	}

	if name, ok := d.GetOk("name"); ok {
		list, _, err := cl.Grok.Gets(ctx)
		if err != nil {
			return err
		}
		var match map[string]interface{}
		cnt := 0
		for _, p := range list["patterns"].([]interface{}) {
			pm := p.(map[string]interface{})
			if pm["name"] == name {
				match = pm
				cnt++
				if cnt > 1 {
					return errors.New("matched multiple grok patterns; name not unique")
				}
			}
		}
		if cnt == 0 {
			return errors.New("matched grok pattern not found")
		}
		return setDataToResourceData(d, match, nil)
	}

	return errors.New("one of pattern_id or name must be set")
}
