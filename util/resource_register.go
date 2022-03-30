package util

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
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
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Computed register value",
			},
		},
		// Changes to set (that are non-empty) mark value as computed
		CustomizeDiff: customdiff.ComputedIf("value", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("set") && d.Get("set").(string) != ""
		}),
	}
}

// registerCreate stores content as the register value.
func resourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	set := d.Get("set").(string)
	d.Set("value", set)
	d.SetId(strconv.Itoa(hashcode.String(set)))
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

	set := d.Get("set").(string)
	if set != "" {
		d.Set("value", set)
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
