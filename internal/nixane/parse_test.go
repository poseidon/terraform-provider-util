package nixane

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseModule(t *testing.T) {
	cases := []*NixOSModule{
		moduleA,
		moduleB,
		moduleC,
	}

	for _, c := range cases {
		module, err := ParseContent(c.Content)
		assert.Nil(t, err)
		assert.Equal(t, c.Imports, module.Imports)
		assert.Equal(t, c.Content, module.Content)
	}
}
