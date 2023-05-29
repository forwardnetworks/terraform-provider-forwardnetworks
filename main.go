package main

import (
    "context"
    "terraform-provider-fwdnet/fwdnet"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name fwdnet

func main() {
    providerserver.Serve(context.Background(), fwdnet.New, providerserver.ServeOpts{
        Address: "registry.terraform.io/fracticated/fwdnet",
    })
}
