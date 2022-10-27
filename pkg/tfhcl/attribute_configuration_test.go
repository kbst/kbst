package tfhcl

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestConfigTokensString(t *testing.T) {
	cfg := map[string]cty.Value{
		"testattr": cty.StringVal("testvalue"),
	}

	expAttrs := []hclwrite.ObjectAttrTokens{
		{
			Name:  hclwrite.TokensForIdentifier("testattr"),
			Value: hclwrite.TokensForValue(cty.StringVal("testvalue")),
		},
	}

	attrs := configTokens(cfg)

	assert.Equal(t, expAttrs, attrs)
}

func TestConfigTokensNumber(t *testing.T) {
	cfg := map[string]cty.Value{
		"testattr": cty.NumberIntVal(5),
	}

	expAttrs := []hclwrite.ObjectAttrTokens{
		{
			Name:  hclwrite.TokensForIdentifier("testattr"),
			Value: hclwrite.TokensForValue(cty.NumberIntVal(5)),
		},
	}

	attrs := configTokens(cfg)

	assert.Equal(t, expAttrs, attrs)
}

func TestConfigTokensBaseDomain(t *testing.T) {
	cfg := map[string]cty.Value{
		"base_domain": cty.StringVal("should_be_replaced"),
	}

	expAttrs := []hclwrite.ObjectAttrTokens{
		{
			Name: hclwrite.TokensForIdentifier("base_domain"),
			Value: hclwrite.TokensForTraversal(hcl.Traversal{
				hcl.TraverseRoot{Name: "var"},
				hcl.TraverseAttr{Name: "base_domain"},
			}),
		},
	}

	attrs := configTokens(cfg)

	assert.Equal(t, expAttrs, attrs)
}
