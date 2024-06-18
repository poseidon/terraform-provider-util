package internal

import (
	"regexp"
	"strings"
)

// Regular expression to match NixOS module's imports blocks.
//
//	  imports = [
//		   capture group, [^\]] matches any character except ]
//	  ]
var moduleImportsRegexp = regexp.MustCompile(`imports[\s]*=[\s]*\[[\s]*(?<imports>[^\]]*)[\s]*\];`)

// ParseContent parses NixOS Module imports.
func ParseContent(content string) (*NixOSModule, error) {
	// all captured imports
	imports := []*NixOSModuleImport{}

	// find earliest match, only one imports block discovered for now
	matches := moduleImportsRegexp.FindStringSubmatch(content)
	if len(matches) > 1 {
		// capture group captures raw contents (e.g. imports = [ contents ];)
		capture := matches[1]
		// TODO: It would be nice to be able to use fields, but some imports have
		// spaces in them
		//fields := strings.Fields(capture)
		fields := strings.Split(capture, "\n")

		for _, field := range fields {
			line := strings.TrimSpace(field)
			// most imports are "in-tree" file references to be traversed. Exceptions:
			// - interpolations are assumed to resolve at runtime
			if line != "" && !strings.HasPrefix(line, "#") && !strings.Contains(line, "$") {
				imports = append(imports, &NixOSModuleImport{
					path:   line,
					inTree: true,
				})
			} else if line != "" && !strings.HasPrefix(line, "#") {
				// preserve non-comment imports as-is. Comments are removed.
				imports = append(imports, &NixOSModuleImport{
					path: line,
				})
			}
		}
	}

	return &NixOSModule{
		Content: content,
		Imports: imports,
	}, nil
}
