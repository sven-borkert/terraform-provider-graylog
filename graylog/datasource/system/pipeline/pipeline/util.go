package pipeline

import (
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sven-borkert/terraform-provider-graylog/graylog/convert"
	rpipeline "github.com/sven-borkert/terraform-provider-graylog/graylog/resource/system/pipeline/pipeline"
)

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}, _ *http.Response) error {
	if err := convert.SetResourceData(d, rpipeline.Resource(), data); err != nil {
		return err
	}
	if id, ok := data["id"]; ok {
		d.SetId(id.(string))
		_ = d.Set("pipeline_id", id.(string))
	}
	if t, ok := data["title"]; ok {
		_ = d.Set("title", t)
	}
	return nil
}
