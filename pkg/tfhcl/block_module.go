package tfhcl

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func BlockModule(f *hclwrite.File, module_name string, providers map[string]string, source string, version string, attributes map[string]hclwrite.Tokens, configurations []Configuration) {
	fb := f.Body()

	bl := fb.AppendNewBlock("module", []string{module_name})
	b := bl.Body()

	// providers
	if len(providers) > 0 {
		b.SetAttributeRaw("providers", AttributeProviders(providers))
		b.AppendNewline()
	}

	// source and version
	if source != "" {
		b.SetAttributeValue("source", cty.StringVal(source))
	}
	if version != "" {
		b.SetAttributeValue("version", cty.StringVal(version))
	}
	if source != "" || version != "" {
		b.AppendNewline()
	}

	// attributes
	for k, v := range attributes {
		b.SetAttributeRaw(k, v)
	}
	if len(attributes) > 0 {
		b.AppendNewline()
	}

	// configuration
	if len(configurations) > 0 {
		cbk := configurations[0].EnvironmentKey
		if cbk != "" && cbk != "apps" {
			b.SetAttributeValue(
				"configuration_base_key",
				cty.StringVal(cbk),
			)
		}

		b.SetAttributeRaw("configuration", AttributeConfiguration(configurations))
	}
}
