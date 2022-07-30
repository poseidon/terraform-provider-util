package util

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Based on Terraform's builtin replace()
// https://www.terraform.io/docs/language/functions/replace.html
// https://github.com/hashicorp/terraform/blob/main/internal/lang/funcs/string.go
func datasourceReplace() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceReplaceRead,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replacements": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				Required: true,
			},
			"replaced": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "content with replacements performed",
			},
		},
	}
}

func datasourceReplaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	replaced, err := replaceContent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("replaced", replaced); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(Hashcode(replaced))
	return diags
}

// Replace content by performing a set of replace operations.
func replaceContent(d *schema.ResourceData) (string, error) {
	// unchecked assertions seem to be the norm in Terraform :S
	content := d.Get("content").(string)
	replacements := map[string]string{}
	for k, v := range d.Get("replacements").(map[string]interface{}) {
		replacements[k] = v.(string)
	}

	for substr, replacement := range replacements {
		// search/replace using regexp if key surrounded by forward slashes
		if len(substr) > 1 && substr[0] == '/' && substr[len(substr)-1] == '/' {
			re, err := regexp.Compile(substr[1 : len(substr)-1])
			if err != nil {
				return content, err
			}
			content = re.ReplaceAllString(content, replacement)
		} else {
			content = strings.Replace(content, substr, replacement, -1)
		}
	}
	return content, nil
}
