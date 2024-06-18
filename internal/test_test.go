package internal

import (
	"io/fs"
	"testing/fstest"
)

var (
	moduleA = &NixOSModule{
		Path: "infra/nix/machines/configuration.nix",
		Content: `{ modulesPath, lib, pkgs, ...}:
{
  system.stateVersion = "24.05";
  imports = [
    "${toString modulesPath}/nixpkgs"
    ../modules/server.nix
		# ../modules/other.nix
  ];
}
`,
		Imports: []*NixOSModuleImport{
			{
				path:   "\"${toString modulesPath}/nixpkgs\"",
				inTree: false,
			},
			{
				path:   "../modules/server.nix",
				inTree: true,
			},
		},
	}

	moduleB = &NixOSModule{
		Path: "infra/nix/modules/server.nix",
		Content: `{ lib, pkgs, ...}:
{
  imports = [
    ./feature.nix
  ];
  services.openssh = {
    enable = true;
  };
}
`,
		Imports: []*NixOSModuleImport{
			{
				path:   "./feature.nix",
				inTree: true,
			},
		},
	}

	moduleC = &NixOSModule{
		Path: "infra/nix/modules/feature.nix",
		Content: `{ pkgs }:
let
  x = 1;
{
  inherit x;
  options = {};
  config = {};
}
`,
		Imports: []*NixOSModuleImport{},
	}
)

func setupTestFilesystem() fs.FS {
	return fstest.MapFS{
		moduleA.Path: &fstest.MapFile{
			Data: []byte(moduleA.Content),
		},
		moduleB.Path: &fstest.MapFile{
			Data: []byte(moduleB.Content),
		},
		moduleC.Path: &fstest.MapFile{
			Data: []byte(moduleC.Content),
		},
	}
}
