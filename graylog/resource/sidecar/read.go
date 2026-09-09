package sidecar

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func read(d *schema.ResourceData, m interface{}) error {
	if err := validateAssignmentScope(d); err != nil {
		return err
	}
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	data, resp, err := cl.Sidecar.GetAll(ctx)
	if err != nil {
		return util.HandleGetResourceError(
			d, resp, fmt.Errorf("failed to get all sidecars: %w", err))
	}
	// The endpoint returns every node. Only refresh assignments for nodes
	// already managed by this resource, so destroy cannot affect other nodes.
	managed := make(map[string]bool)
	for _, value := range d.Get(keySidecars).(*schema.Set).List() {
		node := value.(map[string]interface{})
		managed[node[keyNodeID].(string)] = true
	}
	nodes, ok := data[keySidecars].([]interface{})
	if !ok {
		return fmt.Errorf("unexpected sidecars response: %T", data[keySidecars])
	}
	filtered := make([]interface{}, 0, len(managed))
	for _, value := range nodes {
		node, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unexpected sidecar: %T", value)
		}
		id, _ := node[keyNodeID].(string)
		if managed[id] {
			filtered = append(filtered, node)
		}
	}
	data[keySidecars] = filtered
	return setDataToResourceData(d, data)
}
