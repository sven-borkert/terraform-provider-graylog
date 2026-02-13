package input

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

	if id, ok := d.GetOk("input_id"); ok {
		return readFromID(ctx, d, cl, id.(string))
	}

	if t, ok := d.GetOk("title"); ok {
		return readFromTitle(ctx, d, cl, t.(string))
	}

	return errors.New("one of input_id or title must be set")
}

func readFromID(ctx context.Context, d *schema.ResourceData, cl client.Client, id string) error {
	if _, ok := d.GetOk("title"); ok {
		return errors.New("only one of input_id or title must be set")
	}
	data, resp, err := cl.Input.Get(ctx, id)
	if err != nil {
		return err
	}
	normalizeConfiguration(data)
	return setDataToResourceData(d, data, resp)
}

func readFromTitle(ctx context.Context, d *schema.ResourceData, cl client.Client, title string) error {
	inputs, _, err := cl.Input.Gets(ctx)
	if err != nil {
		return err
	}

	raw, ok := inputs["inputs"]
	if !ok {
		return errors.New("unexpected API response: 'inputs' field missing")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return errors.New("unexpected API response: 'inputs' is not a list")
	}

	cnt := 0
	var data map[string]interface{}
	filterType, hasType := d.GetOk("type")

	for _, in := range list {
		a, ok := in.(map[string]interface{})
		if !ok {
			continue
		}
		if name, _ := a["title"].(string); name != title {
			continue
		}
		if hasType {
			if t, _ := a["type"].(string); t != filterType.(string) {
				continue
			}
		}
		data = a
		cnt++
		if cnt > 1 {
			return errors.New("matched multiple inputs; narrow by type")
		}
	}

	if cnt == 0 {
		return errors.New("matched input is not found")
	}

	normalizeConfiguration(data)
	return setDataToResourceData(d, data, nil)
}
