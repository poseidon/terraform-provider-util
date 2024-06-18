package util

import (
	"context"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/poseidon/terraform-provider-util/internal"
)

func datasourceNix() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceNixRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"overlay": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
			"rendered": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "render merged modules",
			},
		},
	}
}

func datasourceNixRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	rendered, err := renderContent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rendered", rendered); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(Hashcode(rendered))
	return diags
}

func renderContent(d *schema.ResourceData) (string, error) {
	// unchecked assertions seem to be the norm in Terraform :S
	name := d.Get("name").(string)
	path := d.Get("path").(string)
	overlays := map[string]string{}
	for k, v := range d.Get("overlay").(map[string]interface{}) {
		overlays[k] = v.(string)
	}
	path = filepath.Clean(path)

	fsys := internal.NewOverlayFS(overlays)
	modules, err := internal.CollectModules(fsys, path)
	if err != nil {
		return "", err
	}
	return internal.EncodeToAwkball(name, modules), nil
}
