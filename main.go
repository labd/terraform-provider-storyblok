package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/labd/terraform-provider-storyblok/internal"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name storyblok

func main() {
	providerserver.Serve(context.Background(), internal.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/labd/storyblok",
	})
}
