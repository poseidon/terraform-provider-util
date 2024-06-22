package nixane

import (
	"io/fs"
	"path/filepath"
)

// CollectModules parses NixOS module import blocks recursively.
func CollectModules(filesystem fs.FS, path string) ([]*NixOSModule, error) {
	modules := []*NixOSModule{}
	stack := []string{path}

	for len(stack) > 0 {
		path = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// fs.FS does not have a direct os implementation that allows rooted
		// paths (e.g. /home/user/...). afero was better at this.
		// https://github.com/golang/go/issues/47803
		data, err := fs.ReadFile(filesystem, path)
		if err != nil {
			return nil, err
		}

		// parse NixOS module imports
		module, err := ParseContent(string(data))
		if err != nil {
			return nil, err
		}
		module.Path = path

		modules = append(modules, module)

		// traverse in-tree imported NixOS modules
		for _, nmi := range module.Imports {
			if !nmi.inTree {
				continue
			}
			resolvedPath := filepath.Join(filepath.Dir(path), nmi.path)
			stack = append(stack, resolvedPath)
		}
	}
	return modules, nil
}
