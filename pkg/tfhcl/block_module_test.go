package tfhcl

import (
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestBlockModuleParseValid(t *testing.T) {
	fa := hclwrite.NewEmptyFile()
	name := "test_mod"
	src := "test_src"
	v := "v0"
	BlockModule(fa, name, map[string]string{}, src, v, map[string]hclwrite.Tokens{}, []Configuration{})

	p := hclparse.NewParser()
	b, _ := p.ParseHCL(fa.Bytes(), "test.tf")

	var fb Blocks
	gohcl.DecodeBody(b.Body, &hcl.EvalContext{}, &fb)

	assert.Equal(t, 1, len(fb.Modules))

	assert.Equal(t, name, fb.Modules[0].Name)
	assert.Equal(t, src, fb.Modules[0].Source)
	assert.Equal(t, v, fb.Modules[0].Version)
}

func TestBlockModuleNoProviders(t *testing.T) {
	fa := hclwrite.NewEmptyFile()
	BlockModule(fa, "test_service", map[string]string{}, "", "", map[string]hclwrite.Tokens{}, []Configuration{})

	exp := "module \"test_service\" {\n}\n"
	ft, _ := hclwrite.ParseConfig([]byte(exp), "inline", hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}

func TestBlockModuleProviders(t *testing.T) {
	fa := hclwrite.NewEmptyFile()
	BlockModule(fa, "test_service", map[string]string{"testprovider": "alias"}, "", "", map[string]hclwrite.Tokens{}, []Configuration{})

	exp := "module \"test_service\" {\n  providers = {\n    testprovider = testprovider.alias\n  }\n\n}\n"
	ft, _ := hclwrite.ParseConfig([]byte(exp), "inline", hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}

func TestBlockModuleAttributes(t *testing.T) {
	fa := hclwrite.NewEmptyFile()
	attrs := make(map[string]hclwrite.Tokens)
	attrs["testattr"] = hclwrite.TokensForValue(cty.StringVal("test"))
	BlockModule(fa, "test_service", map[string]string{}, "", "", attrs, []Configuration{})

	exp := "module \"test_service\" {\n  testattr = \"test\"\n\n}\n"
	ft, _ := hclwrite.ParseConfig([]byte(exp), "inline", hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}

func TestBlockModuleAppsOps(t *testing.T) {
	fa := hclwrite.NewEmptyFile()

	cfgs := []Configuration{
		{EnvironmentKey: "apps", Attributes: make(map[string]cty.Value)},
		{EnvironmentKey: "ops", Attributes: make(map[string]cty.Value)},
	}

	BlockModule(fa, "test_service", map[string]string{}, "test_source", "0.0.0-test.0", map[string]hclwrite.Tokens{}, cfgs)

	fn := "fixtures/test_block_module_apps_ops.tf"
	d, _ := os.ReadFile(fn)
	ft, _ := hclwrite.ParseConfig(d, fn, hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}

func TestBlockModuleAppsProdAppsOps(t *testing.T) {
	fa := hclwrite.NewEmptyFile()

	cfgs := []Configuration{
		{EnvironmentKey: "apps-prod", Attributes: make(map[string]cty.Value)},
		{EnvironmentKey: "apps", Attributes: make(map[string]cty.Value)},
		{EnvironmentKey: "ops", Attributes: make(map[string]cty.Value)},
	}

	BlockModule(fa, "test_service", map[string]string{}, "test_source", "0.0.0-test.0", map[string]hclwrite.Tokens{}, cfgs)

	fn := "fixtures/test_block_module_apps-prod_apps_ops.tf"
	d, _ := os.ReadFile(fn)
	ft, _ := hclwrite.ParseConfig(d, fn, hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}
