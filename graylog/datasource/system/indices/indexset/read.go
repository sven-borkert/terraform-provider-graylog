package indexset

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

	if id, ok := d.GetOk("index_set_id"); ok {
		return readFromID(ctx, d, cl, id.(string))
	}

	if t, ok := d.GetOk("title"); ok {
		return readFromTitle(ctx, d, cl, t.(string))
	}

	if p, ok := d.GetOk("index_prefix"); ok {
		return readFromPrefix(ctx, d, cl, p.(string))
	}
	return errors.New("one of index_set_id or title or prefix must be set")
}

func readFromID(ctx context.Context, d *schema.ResourceData, cl client.Client, id string) error {
	if _, ok := d.GetOk("title"); ok {
		return errors.New("only one of index_set_id or title or index_prefix must be set")
	}
	if _, ok := d.GetOk("index_prefix"); ok {
		return errors.New("only one of index_set_id or title or index_prefix must be set")
	}
	data, _, err := cl.IndexSet.Get(ctx, id)
	if err != nil {
		return err
	}
	return setDataToResourceData(d, data)
}

func readFromTitle(ctx context.Context, d *schema.ResourceData, cl client.Client, title string) error {
	if _, ok := d.GetOk("index_prefix"); ok {
		return errors.New("only one of index_set_id or title or index_prefix must be set")
	}
	indexSets, _, err := cl.IndexSet.Gets(ctx, nil)
	if err != nil {
		return err
	}
	raw, ok := indexSets["index_sets"]
	if !ok {
		return errors.New("unexpected API response: 'index_sets' field missing")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return errors.New("unexpected API response: 'index_sets' is not a list")
	}
	cnt := 0
	var data map[string]interface{}
	for _, is := range list {
		a, ok := is.(map[string]interface{})
		if !ok {
			continue
		}
		if name, _ := a["title"].(string); name == title {
			data = a
			cnt++
			if cnt > 1 {
				return errors.New("title isn't unique")
			}
		}
	}
	if cnt == 0 {
		return errors.New("matched index set is not found")
	}
	return setDataToResourceData(d, data)
}

func readFromPrefix(ctx context.Context, d *schema.ResourceData, cl client.Client, prefix string) error {
	indexSets, _, err := cl.IndexSet.Gets(ctx, nil)
	if err != nil {
		return err
	}
	raw, ok := indexSets["index_sets"]
	if !ok {
		return errors.New("unexpected API response: 'index_sets' field missing")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return errors.New("unexpected API response: 'index_sets' is not a list")
	}
	for _, is := range list {
		data, ok := is.(map[string]interface{})
		if !ok {
			continue
		}
		if p, _ := data["index_prefix"].(string); p == prefix {
			return setDataToResourceData(d, data)
		}
	}
	return errors.New("matched index prefix is not found")
}
