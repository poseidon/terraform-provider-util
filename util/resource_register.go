package util

import (
	"context"

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
			"content": {
				Type: schema.TypeString,
				// Allow content to be null
				Optional: true,
				// ForceNew would create a new resource on any change,
				// but this resource should ignore empty string or null
				// changes
				ForceNew: false,
				// Suppress plan diffs setting content to an empty string
				// or null. Empty changes are ignored / NoOps
				DiffSuppressOnRefresh: true,
				DiffSuppressFunc:      registerDiffSuppress,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Computed register value",
			},
		},
		// Changes to content (that are non-empty) mark value as computed
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			if content := d.Get("content").(string); d.HasChange("content") && content != "" {
				if err := d.SetNew("value", content); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// registerCreate stores content as the register value.
func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	content := d.Get("content").(string)
	if err := d.Set("value", content); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(Hashcode(content))
	return diags
}

// registerRead sets register attributes
func registerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// register has no remote API equivalent to check or attributes to set
	return nil
}

// registerUpdate applies non-empty content changes to the value attribute.
func registerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	content := d.Get("content").(string)
	if content != "" {
		if err := d.Set("value", content); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

// registerDelete removes the resource from state.
func registerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

// registerDiffSuppress supresses plan diffs setting content to an empty
// string or null (converts to empty string in ResourceData). Empty content
// does not alter the register value.
func registerDiffSuppress(k, oldV, newV string, d *schema.ResourceData) bool {
	return newV == ""
}
