package dashboard

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func read(d *schema.ResourceData, m interface{}) error {
	// log.Printf("dashboard read id=%s", d.Id())
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}
	data, resp, err := cl.View.Get(ctx, d.Id())
	if err != nil {
		return util.HandleGetResourceError(
			d, resp, fmt.Errorf("failed to get a dashboard %s: %w", d.Id(), err))
	}

	// Fetch the search to extract per-widget streams and per-tab query strings
	if searchID, ok := data["search_id"].(string); ok && searchID != "" {
		if searchData, _, err := cl.ViewSearch.Get(ctx, searchID); err == nil {
			injectStreamsFromSearch(data, searchData)
			injectQueryFromSearch(data, searchData)
		}
	}

	return setDataToResourceData(d, data)
}
