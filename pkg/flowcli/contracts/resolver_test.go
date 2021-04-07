package contracts

import (
	"regexp"
	"testing"

	"github.com/onflow/flow-cli/pkg/flowcli/project"
	"github.com/onflow/flow-go-sdk"

	"github.com/stretchr/testify/assert"
)

func cleanCode(code []byte) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(string(code), " ")
}

func TestResolver(t *testing.T) {

	contracts := []project.Contract{{
		Name:   "Kibble",
		Source: "./tests/Kibble.cdc",
		Target: flow.HexToAddress("0x1"),
	}, {
		Name:   "FT",
		Source: "./tests/FT.cdc",
		Target: flow.HexToAddress("0x2"),
	}, {
		Name:   "NFT",
		Source: "./tests/NFT.cdc",
		Target: flow.HexToAddress("0x3"),
	}}

	scripts := [][]byte{
		[]byte(`
			import Kibble from "./Kibble.cdc"
      import FT from "./FT.cdc"
			pub fun main() {}
    `),
	}

	t.Run("Import exists", func(t *testing.T) {
		resolver, err := NewResolver([]byte(`
      import Kibble from "./Kibble.cdc"
      pub fun main() {}
    `))
		assert.NoError(t, err)
		assert.True(t, resolver.ImportExists())
	})

	t.Run("Import doesn't exists", func(t *testing.T) {
		resolver, err := NewResolver([]byte(`
      pub fun main() {}
    `))
		assert.NoError(t, err)
		assert.False(t, resolver.ImportExists())
	})

	t.Run("Parse imports", func(t *testing.T) {
		resolver, err := NewResolver(scripts[0])
		assert.NoError(t, err)
		assert.Equal(t, resolver.parseImports(), []string{
			"./Kibble.cdc", "./FT.cdc",
		})
	})

	t.Run("Resolve imports", func(t *testing.T) {
		resolver, err := NewResolver(scripts[0])
		assert.NoError(t, err)

		code, err := resolver.ResolveImports("./tests/foo.cdc", contracts, make(map[string]string))

		assert.NoError(t, err)
		assert.Equal(t, cleanCode(code), cleanCode([]byte(`
			import Kibble from 0x0000000000000001 
			import FT from 0x0000000000000002 
			pub fun main() {}
		`)))
	})

}