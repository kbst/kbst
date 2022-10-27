package tfhcl

import (
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestBlockModuleClusterAKS(t *testing.T) {
	fa := hclwrite.NewEmptyFile()

	cfgs := []Configuration{
		{
			EnvironmentKey: "apps",
			Attributes: map[string]cty.Value{
				"base_domain": cty.StringVal("replaced"),
				"string":      cty.StringVal("testvalue"),
				"number":      cty.NumberIntVal(5),
			},
		},
		{
			EnvironmentKey: "ops",
			Attributes:     map[string]cty.Value{},
		},
	}

	ModuleCluster(fa, "aks_kbst_westeurope", "azurerm", "test_cluster", "0.0.0-test.0", cfgs)

	fn := "fixtures/test_block_module_cluster_aks.tf"
	d, _ := os.ReadFile(fn)
	ft, _ := hclwrite.ParseConfig(d, fn, hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}
