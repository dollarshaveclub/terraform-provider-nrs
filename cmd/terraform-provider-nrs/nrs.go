package main

import (
	"github.com/dollarshaveclub/terraform-provider-nrs/pkg/provider"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
