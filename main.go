package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return Provider()
		},
	}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "registry.terraform.io/yurii-kysil/azure-github-runners"
	}

	plugin.Serve(opts)
}
