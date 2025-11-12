package staticfield

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-graylog/terraform-provider-graylog/graylog/util"
)

func update(d *schema.ResourceData, m interface{}) error {
	return create(d, m)
}
