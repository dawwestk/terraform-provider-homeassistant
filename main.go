package main

import (
	"github.com/dawwestk/terraform-provider-homeassistant/homeassistant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: homeassistant.Provider,
	})
}
