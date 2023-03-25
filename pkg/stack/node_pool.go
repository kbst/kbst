package stack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"golang.org/x/exp/slices"
)

type NodePool struct {
	PoolName       string
	ClusterName    string
	Provider       string
	Region         string
	Version        string
	Configurations []Configuration
}

func (np *NodePool) Validate(cj util.CliJSON) error {
	var instanceType string
	var zones []string

	//
	//
	// Reject empty configuration
	if len(np.Configurations) == 0 {
		return fmt.Errorf("invalid empty configuration %+v", np.Configurations)
	}

	//
	//
	// Validate provider, region, instance type, zone combinations
	baseCfg := np.Configurations[0].Attributes

	switch np.Provider {
	case "aws":
		its := baseCfg["instance_types"].AsString()
		instanceType = strings.Split(its, ",")[0]

		azs, found := baseCfg["availability_zones"]
		if found {
			zones = strings.Split(azs.AsString(), ",")
		}
	case "azurerm":
		vms, found := baseCfg["default_node_pool_vm_size"]
		if found {
			instanceType = vms.AsString()
		}

		azs, found := baseCfg["availability_zones"]
		if found {
			zones = strings.Split(azs.AsString(), ",")
		}
	case "google":
		instanceType = baseCfg["machine_type"].AsString()

		nl, found := baseCfg["node_locations"]
		if found {
			for _, z := range nl.AsValueSlice() {
				zones = append(zones, z.AsString())
			}
		}
	}

	regionOptions := cj.CloudInfo.Regions(np.Provider)
	if !slices.Contains(regionOptions, np.Region) {
		return fmt.Errorf("invalid region %q for provider %q: choose one of %q", np.Region, np.Provider, regionOptions)
	}

	instanceOptions := cj.CloudInfo.Instances(np.Provider, np.Region)
	if !slices.Contains(instanceOptions, instanceType) {
		return fmt.Errorf("invalid instance type %q for region %q and provider %q: choose one of %q", instanceType, np.Region, np.Provider, instanceOptions)
	}

	zoneOptions := cj.CloudInfo.Zones(np.Provider, np.Region, instanceType)
	for _, z := range zones {
		if !slices.Contains(zoneOptions, z) {
			return fmt.Errorf("invalid zone %q for instance type %q, region %q and provider %q: choose one of %q", z, instanceType, np.Region, np.Provider, zoneOptions)
		}
	}

	return nil
}

func (np *NodePool) ToHCL() map[string]*hclwrite.File {
	files := make(map[string]*hclwrite.File)
	f := hclwrite.NewEmptyFile()

	tfhcl.ModuleNodePool(f, np.Name(), np.Provider, np.ClusterName, np.Version, convertToTfhclConfiguration(np.Configurations))
	files[fmt.Sprintf("%s.tf", np.Name())] = f

	return files
}

func (np *NodePool) Name() string {
	return fmt.Sprintf("%s_%s_%s", np.ClusterName, "node_pool", np.PoolName)
}
