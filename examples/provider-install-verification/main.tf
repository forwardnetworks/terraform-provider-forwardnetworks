terraform {
  required_providers {
    fwdnet = {
      source = "registry.terraform.io/fracticated/fwdnet"
    }
  }
}

provider "fwdnet" {}

data "fwdnet_version" "version" {}

data "fwdnet_external_id" "example" {
  network_id = "159780"
}

output "external_id" {
  value = data.fwdnet_external_id.example.id
}

output "fwdnet_version" {
  value = data.fwdnet_version.version.id
  description = "The Forward Networks version."
}
