package indexset

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

const (
	// deflectorTimeout is how long to wait for the deflector to be ready
	deflectorTimeout = 30 * time.Second
	// deflectorPollInterval is how often to check the deflector status
	deflectorPollInterval = 500 * time.Millisecond
)

func create(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	data, err := getDataFromResourceData(d)
	if err != nil {
		return err
	}
	delete(data, keyDefault)

	is, _, err := cl.IndexSet.Create(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to create a index set: %w", err)
	}

	id := is[keyID].(string)
	d.SetId(id)

	// Wait for the deflector to be ready before returning
	// This ensures the index and alias are properly created before any
	// dependent resources (like streams) start routing data to this index set
	if err := cl.IndexSet.WaitForDeflector(ctx, id, deflectorTimeout, deflectorPollInterval); err != nil {
		return fmt.Errorf("index set created but deflector not ready: %w", err)
	}

	return util.ReadAfterCreate(d, m, id, read)
}
