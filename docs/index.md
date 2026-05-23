# terraform-provider-util

`terraform-provider-util` provides some low-level utility functions.

## Usage

Configure the `util` provider (e.g. `providers.tf`).

```hcl
provider "util" {}

terraform {
  required_providers {
    util = {
      source  = "poseidon/util"
      version = "0.4.0"
    }
  }
}
```

Run `terraform init` to ensure plugin version requirements are met.

```
$ terraform init
```
