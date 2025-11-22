package saved

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	title := d.Get("title").(string)
	list, _, err := cl.SavedSearch.Gets(ctx)
	if err != nil {
		return err
	}

	var match map[string]interface{}
	count := 0
	elems, ok := list["elements"].([]interface{})
	if !ok {
		return errors.New("unexpected response: elements missing")
	}
	for _, e := range elems {
		m := e.(map[string]interface{})
		if m["title"] == title {
			match = m
			count++
			if count > 1 {
				return errors.New("matched multiple saved searches; title not unique")
			}
		}
	}
	if count == 0 {
		return errors.New("matched saved search not found")
	}

	if err := setDataToResourceData(d, match); err != nil {
		return err
	}

	view, _, err := cl.View.Get(ctx, match["id"].(string))
	if err != nil {
		return fmt.Errorf("failed to fetch saved search view %s: %w", match["id"], err)
	}
	if state, ok := view["state"].(map[string]interface{}); ok {
		for id := range state {
			_ = d.Set("state_id", id)
			break
		}
	}
	return nil
}
