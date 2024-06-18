package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var awkballABC = `#### configuration.nix
{ modulesPath, lib, pkgs, ...}:
{
  system.stateVersion = "24.05";
  imports = [ "${toString modulesPath}/nixpkgs" ./modules/server.nix ];
}
#### modules/server.nix
{ lib, pkgs, ...}:
{
  imports = [ ./feature.nix ];
  services.openssh = {
    enable = true;
  };
}
#### modules/feature.nix
{ pkgs }:
let
  x = 1;
{
  inherit x;
  options = {};
  config = {};
}
`

var awkballBC = `#### configuration.nix
{ lib, pkgs, ...}:
{
  imports = [ ./modules/feature.nix ];
  services.openssh = {
    enable = true;
  };
}
#### modules/feature.nix
{ pkgs }:
let
  x = 1;
{
  inherit x;
  options = {};
  config = {};
}
`

var awkballC = `#### configuration.nix
{ pkgs }:
let
  x = 1;
{
  inherit x;
  options = {};
  config = {};
}
`

func TestEncode(t *testing.T) {
	cases := []struct {
		modules  []*NixOSModule
		expected string
	}{
		{
			[]*NixOSModule{moduleA, moduleB, moduleC},
			awkballABC,
		},
		{
			[]*NixOSModule{moduleB, moduleC},
			awkballBC,
		},
		{
			[]*NixOSModule{moduleC},
			awkballC,
		},
	}

	for _, c := range cases {
		awkball := EncodeToAwkball("configuration.nix", c.modules)
		assert.Equal(t, c.expected, awkball)
	}
}
