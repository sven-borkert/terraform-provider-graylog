package fieldtype

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	clientPkg "github.com/sven-borkert/terraform-provider-graylog/graylog/client"
	ftClient "github.com/sven-borkert/terraform-provider-graylog/graylog/client/system/indices/fieldtype"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"index_set_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"field": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"string", "string_fts", "long", "double", "date",
					"boolean", "binary", "ip", "geo-point",
				}, false),
			},
			"rotate": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl, err := clientPkg.New(m)
	if err != nil {
		return diag.FromErr(err)
	}

	indexSetID := d.Get("index_set_id").(string)
	field := d.Get("field").(string)
	fieldType := d.Get("type").(string)
	rotate := d.Get("rotate").(bool)

	resp, err := cl.FieldType.ChangeFieldType(ctx, ftClient.FieldTypeChangeRequest{
		Field:     field,
		Type:      fieldType,
		IndexSets: []string{indexSetID},
		Rotate:    rotate,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to set field type: %w (status: %v)", err, resp))
	}

	d.SetId(indexSetID + "/" + field)
	return resourceRead(ctx, d, m)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl, err := clientPkg.New(m)
	if err != nil {
		return diag.FromErr(err)
	}

	indexSetID := d.Get("index_set_id").(string)
	field := d.Get("field").(string)

	// Handle import: parse ID
	if indexSetID == "" || field == "" {
		parts := splitID(d.Id())
		if parts == nil {
			return diag.Errorf("invalid resource ID: %s", d.Id())
		}
		indexSetID = parts[0]
		field = parts[1]
	}

	ft, _, err := cl.FieldType.GetFieldType(ctx, indexSetID, field)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read field type: %w", err))
	}

	if ft == nil {
		d.SetId("")
		return nil
	}

	d.Set("index_set_id", indexSetID)
	d.Set("field", field)
	d.Set("type", ft.Type)

	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceCreate(ctx, d, m)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	cl, err := clientPkg.New(m)
	if err != nil {
		return diag.FromErr(err)
	}

	indexSetID := d.Get("index_set_id").(string)
	field := d.Get("field").(string)
	rotate := d.Get("rotate").(bool)

	_, err = cl.FieldType.RemoveCustomMapping(ctx, ftClient.CustomFieldMappingRemovalRequest{
		Fields:    []string{field},
		IndexSets: []string{indexSetID},
		Rotate:    rotate,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to remove field type mapping: %w", err))
	}

	d.SetId("")
	return nil
}

func splitID(id string) []string {
	for i, c := range id {
		if c == '/' {
			return []string{id[:i], id[i+1:]}
		}
	}
	return nil
}
