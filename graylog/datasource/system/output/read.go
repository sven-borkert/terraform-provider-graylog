package output

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

	if id, ok := d.GetOk("output_id"); ok {
		return readFromID(ctx, d, cl, id.(string))
	}

	if title, ok := d.GetOk("title"); ok {
		return readFromTitle(ctx, d, cl, title.(string))
	}

	return errors.New("one of output_id or title must be set")
}

func readFromID(ctx context.Context, d *schema.ResourceData, cl client.Client, id string) error {
	if _, ok := d.GetOk("title"); ok {
		return errors.New("only one of output_id or title must be set")
	}

	data, resp, err := cl.Output.Get(ctx, id)
	if err != nil {
		return err
	}
	return setDataToResourceData(d, data, resp)
}

func readFromTitle(ctx context.Context, d *schema.ResourceData, cl client.Client, title string) error {
	outputs, _, err := cl.Output.Gets(ctx, nil)
	if err != nil {
		return err
	}

	filterType, hasType := d.GetOk("type")
	var data map[string]interface{}
	count := 0

	for _, raw := range outputs["outputs"].([]interface{}) {
		o := raw.(map[string]interface{})
		if o["title"].(string) != title {
			continue
		}
		if hasType && o["type"].(string) != filterType.(string) {
			continue
		}
		data = o
		count++
		if count > 1 {
			return errors.New("matched multiple outputs; narrow by type")
		}
	}

	if count == 0 {
		return errors.New("matched output is not found")
	}

	return setDataToResourceData(d, data, nil)
}
