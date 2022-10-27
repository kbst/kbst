package tfhcl

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/maps"
)

type Configuration struct {
	EnvironmentKey string
	Attributes     map[string]cty.Value
}

func AttributeConfiguration(configurations []Configuration) hclwrite.Tokens {
	body := hclwrite.NewEmptyFile().Body()
	cfgs := []hclwrite.ObjectAttrTokens{}

	for _, c := range configurations {
		cfg := hclwrite.ObjectAttrTokens{
			Name:  hclwrite.TokensForIdentifier(c.EnvironmentKey),
			Value: hclwrite.TokensForObject(configTokens(c.Attributes)),
		}

		cfgs = append(cfgs, cfg)
	}

	tokens := body.BuildTokens(hclwrite.TokensForObject(cfgs))

	return tokens
}

func configTokens(in map[string]cty.Value) (out []hclwrite.ObjectAttrTokens) {
	keys := maps.Keys(in)
	sort.Strings(keys)

	for _, k := range keys {
		var vt hclwrite.Tokens

		kt := hclwrite.TokensForIdentifier(k)
		switch k {
		case "base_domain":
			vt = hclwrite.TokensForTraversal(hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: "base_domain"},
			})
		case "tfref_project_id":
			kt = hclwrite.TokensForIdentifier("project_id")
			vt = hclwrite.TokensForTraversal(hcl.Traversal{
				hcl.TraverseRoot{
					Name: "module",
				},
				hcl.TraverseAttr{
					Name: in[k].AsString(),
				},
				hcl.TraverseAttr{
					Name: "current_config[\"project_id\"]",
				},
			})
		default:
			ctyv := in[k]
			if ctyv.IsNull() || !ctyv.IsWhollyKnown() {
				continue
			}
			vt = hclwrite.TokensForValue(ctyv)

		}

		out = append(out, hclwrite.ObjectAttrTokens{
			Name:  kt,
			Value: vt,
		})
	}

	return out
}
