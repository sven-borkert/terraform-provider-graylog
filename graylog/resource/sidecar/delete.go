package sidecar

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

func destroy(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	// Clear assignments only for the sidecars this resource manages, not every
	// sidecar registered on the server.
	managed := d.Get(keySidecars).(*schema.Set).List()
	nodes := make([]interface{}, 0, len(managed))
	for _, s := range managed {
		sm, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		nodes = append(nodes, map[string]interface{}{
			keyNodeID:      sm[keyNodeID],
			keyAssignments: []interface{}{},
		})
	}
	if len(nodes) == 0 {
		return nil
	}

	if _, err := cl.SidecarConfiguration.Assign(ctx, map[string]interface{}{
		keyNodes: nodes,
	}); err != nil {
		return fmt.Errorf("failed to remove configuration assignments from managed sidecars: %w", err)
	}
	return nil
}
