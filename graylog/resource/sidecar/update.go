package sidecar

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func update(d *schema.ResourceData, m interface{}) error {
	if err := validateAssignmentScope(d); err != nil {
		return err
	}
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	data, err := getDataFromResourceData(d)
	if err != nil {
		return err
	}

	// Removing a node from configuration relinquishes its assignments too.
	// Compare IDs rather than whole set elements, since assignments may change.
	old, newValue := d.GetChange(keySidecars)
	retained := make(map[string]bool)
	for _, value := range newValue.(*schema.Set).List() {
		retained[value.(map[string]interface{})[keyNodeID].(string)] = true
	}
	nodes := data[keyNodes].([]interface{})
	for _, value := range old.(*schema.Set).List() {
		id := value.(map[string]interface{})[keyNodeID].(string)
		if !retained[id] {
			nodes = append(nodes, map[string]interface{}{
				keyNodeID: id, keyAssignments: []interface{}{},
			})
		}
	}
	data[keyNodes] = nodes

	// Remove computed fields for Graylog 7.0 compatibility
	util.RemoveComputedFields(data)

	if _, err := cl.SidecarConfiguration.Assign(ctx, data); err != nil {
		return fmt.Errorf("failed to update sidecars's assignments: %w", err)
	}
	return nil
}
