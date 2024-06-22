package nixane

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectModules(t *testing.T) {
	testFS := setupTestFilesystem()
	cases := []struct {
		path    string
		modules []*NixOSModule
	}{
		{
			path:    moduleA.Path,
			modules: []*NixOSModule{moduleA, moduleB, moduleC},
		},
		{
			path:    moduleB.Path,
			modules: []*NixOSModule{moduleB, moduleC},
		},
		{
			path:    moduleC.Path,
			modules: []*NixOSModule{moduleC},
		},
	}
	for _, c := range cases {
		modules, err := CollectModules(testFS, c.path)
		assert.Nil(t, err)
		assert.Equal(t, c.modules, modules)
	}
}
