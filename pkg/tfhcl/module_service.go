package tfhcl

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ModuleService(f *hclwrite.File, module_name string, cluster_name string, source string, version string, configurations []Configuration) {
	providers := make(map[string]string)
	providers["kustomization"] = cluster_name

	BlockModule(f, module_name, providers, source, version, map[string]hclwrite.Tokens{}, configurations)
}
