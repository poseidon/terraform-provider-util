package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

var awkballPrefix = "modules"

// Encode encodes the NixOS modules in awkball format, with each file
// separated by a `#### path` marker. An awkball can be decoded to
// files using an awk one-liner.
//
//	#### configuration.nix
//	content
//	#### modules/a.nix
//	content
//
// The first module is considered the entrypoint and assigned the given
// entrypoint filename. Others are organized under modules. Imports are
// modified to match.
func EncodeToAwkball(entrypoint string, modules []*NixOSModule) string {
	var sb strings.Builder
	for i, module := range modules {
		var name string
		var content string
		if i == 0 {
			name = entrypoint
			content = module.RewriteImports(awkballPrefix)
		} else {
			name = filepath.Join(awkballPrefix, filepath.Base(module.Path))
			content = module.RewriteImports("")
		}
		sb.WriteString(fmt.Sprintf("#### %s\n%s", name, content))
	}
	return sb.String()
}

// NixOSModule represents info from a NixOS module.
type NixOSModule struct {
	Path    string
	Content string
	Imports []*NixOSModuleImport
}

type NixOSModuleImport struct {
	path   string
	inTree bool
}

// Rewrite imports returns the NixOS module content with imports rewritten
// for a remote system.
func (m *NixOSModule) RewriteImports(prefix string) string {
	imports := []string{}
	// evaluate all nixos module imports
	for _, nmi := range m.Imports {
		// source/intree NixOS module
		if nmi.inTree {
			rewrite := fmt.Sprintf("./%s", filepath.Join(prefix, filepath.Base(nmi.path)))
			imports = append(imports, rewrite)
		} else {
			imports = append(imports, nmi.path)
		}
	}

	return moduleImportsRegexp.ReplaceAllString(
		m.Content,
		fmt.Sprintf("imports = [ %s ];", strings.Trim(strings.Join(imports, " "), " ")),
	)
}
