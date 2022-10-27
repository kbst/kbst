package tfhcl

import (
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

func TestBlockTerraformSingle(t *testing.T) {
	p := make(map[string]map[string]string)
	p["testprovider"] = make(map[string]string)
	p["testprovider"]["source"] = "test/testprovider"

	fa := hclwrite.NewEmptyFile()
	BlockTerraform(fa, p)

	fn := "fixtures/test_block_terraform_single.tf"
	d, _ := os.ReadFile(fn)
	ft, _ := hclwrite.ParseConfig(d, fn, hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}

func TestBlockTerraformDouble(t *testing.T) {
	p := make(map[string]map[string]string)
	p["testprovider"] = make(map[string]string)
	p["testprovider"]["source"] = "test/testprovider"
	p["testprovider2"] = make(map[string]string)
	p["testprovider2"]["source"] = "test/testprovider"

	fa := hclwrite.NewEmptyFile()
	BlockTerraform(fa, p)

	fn := "fixtures/test_block_terraform_double.tf"
	d, _ := os.ReadFile(fn)
	ft, _ := hclwrite.ParseConfig(d, fn, hcl.InitialPos)

	assert.Equal(t, string(ft.Bytes()), string(fa.Bytes()))
}
