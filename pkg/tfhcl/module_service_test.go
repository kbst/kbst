package tfhcl

import (
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestBlockModuleService(t *testing.T) {
	fa := hclwrite.NewEmptyFile()

	cfgs := []Configuration{
		{
			EnvironmentKey: "apps",
			Attributes:     make(map[string]cty.Value),
		},
		{
			EnvironmentKey: "ops",
			Attributes:     make(map[string]cty.Value),
		},
	}

	ModuleService(fa, "test_service", "test_cluster", "kbst.xyz/catalog/test/kustomization", "0.0.0-test.0", cfgs)

	fn := "fixtures/test_block_module_service.tf"
	d, _ := os.ReadFile(fn)
	ft, _ := hclwrite.ParseConfig(d, fn, hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}
