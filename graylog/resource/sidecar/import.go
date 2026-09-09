package sidecar

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
)

// Import explicitly takes ownership of all current sidecar assignments.
// Ordinary refresh must never adopt additional nodes.
func importState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	cl, err := client.New(m)
	if err != nil {
		return nil, err
	}
	data, _, err := cl.Sidecar.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to import sidecars: %w", err)
	}
	if err := setDataToResourceData(d, data); err != nil {
		return nil, err
	}
	if err := d.Set(keyAssignmentScopeVersion, 1); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
