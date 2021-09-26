package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/poseidon/terraform-provider-util/util"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: util.Provider,
	})
}
