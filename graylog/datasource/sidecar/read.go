package sidecar

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

	if id, ok := d.GetOk(keyNodeID); ok {
		if _, ok := d.GetOk(keyNodeName); ok {
			return errors.New("both node_id and node_name must not be set")
		}
		data, _, err := cl.Sidecar.Get(ctx, id.(string))
		if err != nil {
			return err
		}
		return setDataToResourceData(d, data)
	}

	if t, ok := d.GetOk(keyNodeName); ok {
		nodeName := t.(string)
		sidecars, _, err := cl.Sidecar.GetAll(ctx)
		if err != nil {
			return err
		}
		raw, ok := sidecars[keySidecars]
		if !ok {
			return errors.New("unexpected API response: 'sidecars' field missing")
		}
		list, ok := raw.([]interface{})
		if !ok {
			return errors.New("unexpected API response: 'sidecars' is not a list")
		}
		for _, a := range list {
			sidecar, ok := a.(map[string]interface{})
			if !ok {
				continue
			}
			if name, _ := sidecar[keyNodeName].(string); name == nodeName {
				return setDataToResourceData(d, sidecar)
			}
		}
		return errors.New("matched sidecar is not found")
	}
	return errors.New("one of node_id or node_name must be set")
}
