package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/saymedia/terraform-beanstalk/beanstalk"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: beanstalk.Provider,
	})
}
