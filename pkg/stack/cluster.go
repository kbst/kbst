package stack

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
)

type Cluster struct {
	NamePrefix     string
	Provider       string
	Region         string
	Version        string
	Configurations []Configuration
}

func (c *Cluster) ToHCL() map[string]*hclwrite.File {
	files := make(map[string]*hclwrite.File)

	// _cluster.tf
	fc := hclwrite.NewEmptyFile()

	tfhcl.ModuleCluster(fc, c.Name(), c.Provider, c.Name(), c.Version, convertToTfhclConfiguration(c.Configurations))
	files[fmt.Sprintf("%s_cluster.tf", c.Name())] = fc

	// _providers.tf
	fp := hclwrite.NewEmptyFile()
	tfhcl.BlockProvider(fp, c.Provider, c.Name(), c.Region)
	files[fmt.Sprintf("%s_providers.tf", c.Name())] = fp

	return files
}

func (c *Cluster) cloudK8sPrefix() string {
	k8sServiceName := map[string]string{
		"aws":     "eks",
		"azurerm": "aks",
		"google":  "gke",
	}

	return k8sServiceName[c.Provider]
}

func (c *Cluster) Name() string {
	return fmt.Sprintf("%s_%s_%s", c.cloudK8sPrefix(), c.NamePrefix, c.Region)
}
