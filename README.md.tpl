# Forward Networks Provider

{{ .Provider.Description }}

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x

## Provider Versioning

This provider follows [Semantic Versioning](http://semver.org/).

## Usage

```hcl
provider "forwardnetworks" {
  // Provider configuration here
}

** Provider Configuration
{{ .Provider.ConfigurationBlock | markdownTable }}

## Resources

{{ range .Resources }}
- [`{{ .Type }}`]({{ .Type }}/README.md){{ end }}

## Data Sources

{{ range .DataSources }}
- [`{{ .Type }}`]({{ .Type }}/README.md){{ end }}
