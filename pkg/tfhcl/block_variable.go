package tfhcl

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func BlockVariable(f *hclwrite.File, name string, tftype string, description string) {
	r := f.Body()

	bl := r.AppendNewBlock("variable", []string{name})
	b := bl.Body()

	b.SetAttributeTraversal("type", hcl.Traversal{
		hcl.TraverseRoot{
			Name: tftype,
		}})
	b.SetAttributeValue("description", cty.StringVal(description))
}
