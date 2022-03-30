package util

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a util Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"util_register": resourceRegister(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"util_replace": datasourceReplace(),
		},
	}
}
