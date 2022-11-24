package stack

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"golang.org/x/exp/slices"
)

type Cluster struct {
	NamePrefix     string
	Provider       string
	Region         string
	Version        string
	Configurations []Configuration
}

func (c *Cluster) Validate(cj util.CliJSON) error {
	var instanceType string
	var zones []string

	//
	//
	// Validate framework version
	var versionOptions []string
	for _, fv := range cj.Framework.Versions {
		versionOptions = append(versionOptions, fv.Name)
	}

	if !slices.Contains(versionOptions, c.Version) {
		return fmt.Errorf("invalid version %q, choose one of %q", c.Version, versionOptions)
	}

	//
	//
	// Validate provider, region, instance type, zone combinations
	baseCfg := c.Configurations[0].Attributes

	switch c.Provider {
	case "aws":
		instanceType = baseCfg["cluster_instance_type"].AsString()
		zones = strings.Split(baseCfg["cluster_availability_zones"].AsString(), ",")
	case "azurerm":
		instanceType = baseCfg["default_node_pool_vm_size"].AsString()
		zones = strings.Split(baseCfg["availability_zones"].AsString(), ",")
	case "google":
		instanceType = baseCfg["cluster_machine_type"].AsString()
		zones = strings.Split(baseCfg["cluster_node_locations"].AsString(), ",")
	}

	regionOptions := cj.CloudInfo.Regions(c.Provider)
	if !slices.Contains(regionOptions, c.Region) {
		return fmt.Errorf("invalid region %q for provider %q: choose one of %q", c.Region, c.Provider, regionOptions)
	}

	instanceOptions := cj.CloudInfo.Instances(c.Provider, c.Region)
	if !slices.Contains(instanceOptions, instanceType) {
		return fmt.Errorf("invalid instance type %q for region %q and provider %q: choose one of %q", instanceType, c.Region, c.Provider, instanceOptions)
	}

	zoneOptions := cj.CloudInfo.Zones(c.Provider, c.Region, instanceType)
	for _, z := range zones {
		if !slices.Contains(zoneOptions, z) {
			return fmt.Errorf("invalid zone %q for instance type %q, region %q and provider %q: choose one of %q", z, instanceType, c.Region, c.Provider, zoneOptions)
		}
	}

	return nil
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
