# terraform-provider-forwardnetworks

This is a [Terraform](https://www.terraform.io)©
[provider](https://developer.hashicorp.com/terraform/language/providers?page=providers)
for Forward Networks. It relies on a Go client library at [https://github.com/forwardnetworks/terraform-provider-forwardnetworks/forwardnetworks](https://github.com/forwardnetworks/terraform-provider-forwardnetworks/tree/main/forwardnetworks)

## Getting Started

### Install Terraform

Instructions for popular operating systems can be found [here](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli).

### Create a Terraform configuration

The terraform configuration must:
- be named with a `.tf` file extension.
- reference this provider by its global address.
  *registry.terraform.io/forwardnetworks/forwardnetworks* or just: *forwardnetworks/forwardnetworks*.
- include a provider configuration block which tells the provider where to
find the forwardnetworks service.

```hcl
terraform {
  required_providers {
    forwardnetworks = {
      source = "forwardnetworks/forwardnetworks"
    }
  }
}

provider "forwardnetworks" {
  apphost = "<forwardnetworks-server-ip>"
}
```

### Terraform Init

Run the following at a command prompt while in the same directory as the
configuration file to fetch the forwardnetworks provider plugin.
```shell
terraform init
```

### Supply forwardnetworks credentials
forwardnetworks credentials can be supplied through environment variables:
```shell
export forwardnetworks_USER=<username>
export forwardnetworks_PASS=<password>
```

### Start configuring resources

Full documentation for provider, resources and data sources can be found
[here](https://registry.terraform.io/providers/forwardnetworks/forwardnetworks/latest/docs).
