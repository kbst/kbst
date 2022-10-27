package tfhcl

import (
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"golang.org/x/exp/maps"
)

func AttributeProviders(providers map[string]string) hclwrite.Tokens {
	body := hclwrite.NewEmptyFile().Body()

	prvdrs := []hclwrite.ObjectAttrTokens{}
	keys := maps.Keys(providers)
	sort.Strings(keys)
	for _, k := range keys {
		prvdrs = append(prvdrs, hclwrite.ObjectAttrTokens{
			Name: hclwrite.TokensForIdentifier(k),
			Value: hclwrite.TokensForTraversal(
				hcl.Traversal{
					hcl.TraverseRoot{
						Name: k,
					},
					hcl.TraverseAttr{
						Name: providers[k],
					},
				},
			),
		})
	}

	tokens := body.BuildTokens(hclwrite.TokensForObject(prvdrs))

	return tokens
}
