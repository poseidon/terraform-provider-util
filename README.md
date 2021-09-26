# terraform-provider-util [![Build Status](https://github.com/poseidon/terraform-provider-util/workflows/test/badge.svg)](https://github.com/poseidon/terraform-provider-util/actions?query=workflow%3Atest+branch%3Amain)

`terraform-provider-util` provides some low-level utility functions.

## Usage

Configure the `util` provider (e.g. `providers.tf`).

```hcl
provider "util" {}

terraform {
  required_providers {
    ct = {
      source  = "poseidon/util"
      version = "0.1.0"
    }
  }
}
```

Perform a set of replacements on content with `replace`.

```hcl
data "util_replace" "example" {
  content      = "hello world"
  replacements = {
		"/(h|H)ello/": "Hallo",
		"world": "Welt",
  }
}

# Hallo Welt
output "example" {
  value = data.ct_replace.example.replaced
}
```

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```

## Requirements

* Terraform v0.13+ [installed](https://www.terraform.io/downloads.html)

## Development

### Binary

To develop the provider plugin locally, build an executable with Go v1.16+.

```
make
```

