package tfhcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ModuleCluster(f *hclwrite.File, module_name string, cluster_provider string, cluster_name string, version string, configurations []Configuration) {
	providers := make(map[string]string)
	source := fmt.Sprintf("github.com/kbst/terraform-kubestack//%s/cluster?ref=%s", cluster_provider, version)

	// AWS
	if cluster_provider == "aws" {
		providers[cluster_provider] = cluster_name
		providers["kubernetes"] = cluster_name
	}

	// AzureRM
	//if cluster_provider == "azurerm" {}

	// Google
	if cluster_provider == "google" {
		providers["kubernetes"] = cluster_name
	}

	BlockModule(f, module_name, providers, source, "", map[string]hclwrite.Tokens{}, configurations)
}
