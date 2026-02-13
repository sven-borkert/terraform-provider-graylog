package saved_search

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/util"
)

func read(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	cl, err := client.New(m)
	if err != nil {
		return err
	}

	data, resp, err := cl.View.Get(ctx, d.Id())
	if err != nil {
		return util.HandleGetResourceError(
			d, resp, fmt.Errorf("failed to get saved search %s: %w", d.Id(), err))
	}

	// Set basic fields
	if v, ok := data["title"].(string); ok {
		if err := d.Set("title", v); err != nil {
			return err
		}
	}
	if v, ok := data["description"].(string); ok {
		if err := d.Set("description", v); err != nil {
			return err
		}
	}
	if v, ok := data["summary"].(string); ok {
		if err := d.Set("summary", v); err != nil {
			return err
		}
	}
	if v, ok := data["owner"].(string); ok {
		if err := d.Set("owner", v); err != nil {
			return err
		}
	}
	if v, ok := data["created_at"].(string); ok {
		if err := d.Set("created_at", v); err != nil {
			return err
		}
	}
	if v, ok := data["search_id"].(string); ok {
		if err := d.Set("search_id", v); err != nil {
			return err
		}
	}

	// Extract state information
	if stateMap, ok := data["state"].(map[string]interface{}); ok {
		for stateID, stateData := range stateMap {
			if err := d.Set("state_id", stateID); err != nil {
				return err
			}

			state, ok := stateData.(map[string]interface{})
			if !ok {
				continue
			}

			// Get selected fields
			if fields, ok := state["selected_fields"].([]interface{}); ok {
				selectedFields := make([]string, 0, len(fields))
				for _, f := range fields {
					if s, ok := f.(string); ok {
						selectedFields = append(selectedFields, s)
					}
				}
				if err := d.Set("selected_fields", selectedFields); err != nil {
					return err
				}
			}

			// Get query from first widget
			if widgets, ok := state["widgets"].([]interface{}); ok && len(widgets) > 0 {
				if widget, ok := widgets[0].(map[string]interface{}); ok {
					// Get timerange
					if tr, ok := widget["timerange"].(map[string]interface{}); ok {
						if t, ok := tr["type"].(string); ok {
							if err := d.Set("timerange_type", t); err != nil {
								return err
							}
						}
						if r, ok := tr["range"].(float64); ok {
							if err := d.Set("timerange_range", int(r)); err != nil {
								return err
							}
						}
					}

					// Get streams
					if streams, ok := widget["streams"].([]interface{}); ok {
						streamIDs := make([]string, 0, len(streams))
						for _, s := range streams {
							if id, ok := s.(string); ok {
								streamIDs = append(streamIDs, id)
							}
						}
						if err := d.Set("streams", streamIDs); err != nil {
							return err
						}
					}

					// Get query
					if q, ok := widget["query"].(map[string]interface{}); ok {
						if qs, ok := q["query_string"].(string); ok {
							if err := d.Set("query", qs); err != nil {
								return err
							}
						}
					}

					// Get sort from config
					if config, ok := widget["config"].(map[string]interface{}); ok {
						if sorts, ok := config["sort"].([]interface{}); ok && len(sorts) > 0 {
							if sort, ok := sorts[0].(map[string]interface{}); ok {
								if f, ok := sort["field"].(string); ok {
									if err := d.Set("sort_field", f); err != nil {
										return err
									}
								}
								if o, ok := sort["order"].(string); ok {
									if err := d.Set("sort_order", o); err != nil {
										return err
									}
								}
							}
						}
					}
				}
			}
			break // Only process first state
		}
	}

	return nil
}
