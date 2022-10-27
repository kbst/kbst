package tfhcl

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/maps"
)

func BlockTerraform(f *hclwrite.File, providers map[string]map[string]string) {
	rootBody := f.Body()

	tfBlock := rootBody.AppendNewBlock("terraform", nil)
	tfBody := tfBlock.Body()

	rpBlock := tfBody.AppendNewBlock("required_providers", nil)
	rpBody := rpBlock.Body()

	keys := maps.Keys(providers)
	sort.Strings(keys)

	for _, k := range keys {
		pVal := map[string]cty.Value{}
		for ik, iv := range providers[k] {
			pVal[ik] = cty.StringVal(iv)
		}
		rpBody.SetAttributeValue(k, cty.MapVal(pVal))
	}
}
