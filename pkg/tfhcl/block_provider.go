package tfhcl

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"golang.org/x/exp/slices"
)

func BlockProvider(f *hclwrite.File, cluster_provider string, alias string, region string) {
	rootBody := f.Body()

	// cluster provider
	switch cluster_provider {
	case "aws":
		aws := rootBody.AppendNewBlock("provider", []string{cluster_provider})
		awsb := aws.Body()

		awsb.SetAttributeValue("alias", cty.StringVal(alias))
		awsb.AppendNewline()

		awsb.SetAttributeValue("region", cty.StringVal(region))
		rootBody.AppendNewline()
	case "azurerm":
		azurerm := rootBody.AppendNewBlock("provider", []string{cluster_provider})
		azrmb := azurerm.Body()

		azrmb.AppendNewBlock("features", []string{})
		rootBody.AppendNewline()
	}

	// kubeconfig traversal
	kct := hcl.Traversal{
		hcl.TraverseRoot{
			Name: "module",
		},
		hcl.TraverseAttr{
			Name: alias,
		},
		hcl.TraverseAttr{
			Name: "kubeconfig",
		},
	}

	// kustomization
	kustomization := rootBody.AppendNewBlock("provider", []string{"kustomization"})
	kstmzb := kustomization.Body()

	kstmzb.SetAttributeValue("alias", cty.StringVal(alias))
	kstmzb.AppendNewline()

	kstmzb.SetAttributeTraversal("kubeconfig_raw", kct)
	rootBody.AppendNewline()

	// kubernetes
	if slices.Contains([]string{"aws", "google"}, cluster_provider) {
		locals := rootBody.AppendNewBlock("locals", []string{})
		lb := locals.Body()
		lb.SetAttributeRaw(fmt.Sprintf("%s_kubeconfig", alias), hclwrite.TokensForFunctionCall("yamldecode", hclwrite.TokensForTraversal(kct)))
		rootBody.AppendNewline()

		kubernetes := rootBody.AppendNewBlock("provider", []string{"kubernetes"})
		k8sb := kubernetes.Body()

		k8sb.SetAttributeValue("alias", cty.StringVal(alias))
		k8sb.AppendNewline()

		k8sb.SetAttributeTraversal("host", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "local",
			},
			hcl.TraverseAttr{
				Name: fmt.Sprintf("%s_kubeconfig[\"clusters\"][0][\"cluster\"][\"server\"]", alias),
			},
		},
		)
		k8sb.SetAttributeRaw("cluster_ca_certificate", hclwrite.TokensForFunctionCall("base64decode", hclwrite.TokensForTraversal(
			hcl.Traversal{
				hcl.TraverseRoot{
					Name: "local",
				},
				hcl.TraverseAttr{
					Name: fmt.Sprintf("%s_kubeconfig[\"clusters\"][0][\"cluster\"][\"certificate-authority-data\"]", alias),
				},
			},
		)))
		k8sb.SetAttributeTraversal("token", hcl.Traversal{
			hcl.TraverseRoot{
				Name: "local",
			},
			hcl.TraverseAttr{
				Name: fmt.Sprintf("%s_kubeconfig[\"users\"][0][\"user\"][\"token\"]", alias),
			},
		},
		)
	}
}
