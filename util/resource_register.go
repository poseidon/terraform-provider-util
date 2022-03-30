package util

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRegister() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   registerRead,
		UpdateContext: registerUpdate,
		DeleteContext: registerDelete,
		Schema: map[string]*schema.Schema{
			"set": &schema.Schema{
				Type:                  schema.TypeString,
				Optional:              true,
				ForceNew:              false,
				DiffSuppressOnRefresh: true,
				DiffSuppressFunc:      registerDiffSuppress,
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last set register value",
			},
		},
	}
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	set := d.Get("set").(string)
	d.Set("value", set)
	d.SetId(strconv.Itoa(hashcode.String(set)))
	return diags
}

func registerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func registerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	set := d.Get("set").(string)
	if set != "" {
		d.Set("value", set)
	}
	return diags
}

func registerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func registerDiffSuppress(k, oldV, newV string, d *schema.ResourceData) bool {
	return newV == ""
}
