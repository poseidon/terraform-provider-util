# terraform-provider-util
[![GoDoc](https://pkg.go.dev/badge/github.com/poseidon/terraform-provider-util.svg)](https://pkg.go.dev/github.com/poseidon/terraform-provider-util) [![Workflow](https://github.com/poseidon/terraform-provider-util/actions/workflows/test.yaml/badge.svg)](https://github.com/poseidon/terraform-provider-util/actions/workflows/test.yaml?query=branch%3Amain) ![Downloads](https://img.shields.io/github/downloads/poseidon/terraform-provider-util/total) [![Sponsors](https://img.shields.io/github/sponsors/poseidon?logo=github)](https://github.com/sponsors/poseidon) [![Twitter](https://img.shields.io/badge/follow-news-1da1f2?logo=twitter)](https://twitter.com/poseidonlabs)

`terraform-provider-util` provides some low-level utility functions.

## Usage

Configure the `util` provider (e.g. `providers.tf`).

```hcl
provider "util" {}

terraform {
  required_providers {
    ct = {
      source  = "poseidon/util"
      version = "0.2.2"
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

Store a content value that persists in state until changed to a non-empty value.

```tf
resource "util_register" "example" {
  content = "a1b2c3"
}
```

Later, the register's content may be updated, but empty values (`null` or `""`) are ignored.

```tf
resource "util_register" "example" {
  content = null
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

To develop the provider plugin locally, build an executable with Go v1.18+.

```
make
```
