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
      version = "0.2.0"
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

Store a value in state that persists until changed to a non-empty value.

```tf
resource "util_register" "example" {
  set = "a1b2c3"
}
```

Later, the register's value may be updated, but setting it to `null` or `""` is ignored.

```tf
resource "util_register" "example" {
  set = null
}

output "sha" {
  value = util_register.example.value  # "a1b2c3"
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
