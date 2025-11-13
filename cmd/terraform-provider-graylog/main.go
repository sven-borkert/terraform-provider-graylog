package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/sven-borkert/terraform-provider-graylog/graylog"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: graylog.Provider,
	})
}
