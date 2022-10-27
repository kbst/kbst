package tfhcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func ModuleNodePool(f *hclwrite.File, module_name string, cluster_provider string, cluster_name string, version string, configurations []Configuration) {
	providers := make(map[string]string)
	if cluster_provider == "aws" {
		providers[cluster_provider] = cluster_name
	}

	source := fmt.Sprintf("github.com/kbst/terraform-kubestack//%s/cluster/node-pool?ref=%s", cluster_provider, version)

	attributes := make(map[string]hclwrite.Tokens)

	switch cluster_provider {
	case "aws":
		attributes["cluster_name"] = hclwrite.TokensForTraversal(hcl.Traversal{
			hcl.TraverseRoot{
				Name: "module",
			},
			hcl.TraverseAttr{
				Name: cluster_name,
			},
			hcl.TraverseAttr{
				Name: "current_metadata[\"name\"]",
			},
		})
	case "azurerm":
		attributes["cluster_name"] = hclwrite.TokensForTraversal(hcl.Traversal{
			hcl.TraverseRoot{
				Name: "module",
			},
			hcl.TraverseAttr{
				Name: cluster_name,
			},
			hcl.TraverseAttr{
				Name: "current_metadata[\"name\"]",
			},
		})
		attributes["resource_group"] = hclwrite.TokensForTraversal(hcl.Traversal{
			hcl.TraverseRoot{
				Name: "module",
			},
			hcl.TraverseAttr{
				Name: cluster_name,
			},
			hcl.TraverseAttr{
				Name: "current_config[\"resource_group\"]",
			},
		})
	case "google":
		attributes["cluster_metadata"] = hclwrite.TokensForTraversal(hcl.Traversal{
			hcl.TraverseRoot{
				Name: "module",
			},
			hcl.TraverseAttr{
				Name: cluster_name,
			},
			hcl.TraverseAttr{
				Name: "current_metadata",
			},
		})

		// hack: handle project_id traversal special case
		configurations[0].Attributes["tfref_project_id"] = cty.StringVal(cluster_name)
	}

	BlockModule(f, module_name, providers, source, "", attributes, configurations)
}
