package stack

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
)

type NodePool struct {
	NameSuffix     string
	ClusterName    string
	Provider       string
	Region         string
	Version        string
	Configurations []Configuration
}

func (np *NodePool) ToHCL() map[string]*hclwrite.File {
	files := make(map[string]*hclwrite.File)
	f := hclwrite.NewEmptyFile()

	tfhcl.ModuleNodePool(f, np.Name(), np.Provider, np.ClusterName, np.Version, convertToTfhclConfiguration(np.Configurations))
	files[fmt.Sprintf("%s.tf", np.Name())] = f

	return files
}

func (np *NodePool) Name() string {
	return fmt.Sprintf("%s_%s_%s", np.ClusterName, "node_pool", np.NameSuffix)
}
