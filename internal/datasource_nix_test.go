package internal

import (
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const nixExample1 = `
data "util_nix" "example" {
	name = "configuration.nix"
	path = "./testdata/foo.nix"
}
`

const nixExpected1 = `#### configuration.nix
{ modulesPath, lib, pkgs, ...}:
{
  system.stateVersion = "24.05";
  imports = [ "${toString modulesPath}/nixpkgs" ./modules/bar.nix ];
}
#### modules/bar.nix
{
  imports = [ ./baz.nix ];

  options = {};
  config = {};
}
#### modules/baz.nix
{ pkgs }:
let
  x = 1;
{
  inherit x;
}
`

const nixExample2 = `
data "util_nix" "example" {
	name = "configuration.nix"
	path = "./testdata/foo.nix"
	overlay = {
		"testdata/foo.nix" = <<-EOT
{ modulesPath, lib, pkgs, ...}:
{
  system.stateVersion = "24.05";
  imports = [
    "$${toString modulesPath}/nixpkgs"
    ./bar.nix
  ];
}
EOT
	}
}
`

func TestNix(t *testing.T) {
	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: nixExample1,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.util_nix.example", "rendered", nixExpected1),
				),
			},
			{
				Config: nixExample2,
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr("data.util_nix.example", "rendered", nixExpected1),
				),
			},
		},
	})
}
