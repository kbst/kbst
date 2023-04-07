/*
Copyright © 2020 Kubestack <hello@kubestack.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"strings"

	"github.com/kbst/kbst/pkg/stack"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zclconf/go-cty/cty"
)

var sharedFlags pflag.FlagSet

var clusterNamePrefix string
var clusterRegion string

var clusterAKSInstanceType string
var clusterAKSMinNodes int64
var clusterAKSMaxNodes int64
var clusterAKSZones string

var clusterEKSInstanceType string
var clusterEKSMinNodes int64
var clusterEKSMaxNodes int64
var clusterEKSZones string

var clusterGKEInstanceType string
var clusterGKEMinNodes int64
var clusterGKEMaxNodes int64
var clusterGKEZones string

var nodePoolAKSInstanceType string
var nodePoolAKSMinNodes int64
var nodePoolAKSMaxNodes int64
var nodePoolAKSDiskSize int64
var nodePoolAKSZones string

var nodePoolEKSInstanceType string
var nodePoolEKSMinNodes int64
var nodePoolEKSMaxNodes int64
var nodePoolEKSDiskSize int64
var nodePoolEKSZones string
var nodePoolEKSAMIType string

var nodePoolGKEInstanceType string
var nodePoolGKEImageType string
var nodePoolGKEMinNodes int64
var nodePoolGKEMaxNodes int64
var nodePoolGKEDiskType string
var nodePoolGKEDiskSize int64
var nodePoolGKEZones string

var serviceRelease string
var serviceClusterName string

var addCmd = &cobra.Command{
	Use:   "add command [flags]",
	Short: "Add clusters, node pools or services",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ErrMissingCommand
	},
}

var clusterAddCmd = &cobra.Command{
	Use:     "cluster command [flags]",
	Aliases: []string{"c"},
	Short:   "Add an AKS, EKS or GKE cluster module",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ErrMissingCommand
	},
}

var clusterAddAKSCmd = &cobra.Command{
	Use:   "aks <name-prefix> <region> <resource-group>",
	Short: "Add an AKS cluster module",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		namePrefix := args[0]
		region := args[1]
		resourceGroup := args[2]

		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot(path)
		s := stack.NewStack(r, cj)
		err = s.FromPath()
		if err != nil {
			log.Fatal(err)
		}

		zones := strings.Split(clusterAKSZones, ",")
		if clusterAKSZones == "" {
			zones = cj.CloudInfo.Zones("azurerm", region, clusterAKSInstanceType)
		}

		baseCfg := map[string]cty.Value{
			"name_prefix":                  cty.StringVal(namePrefix),
			"resource_group":               cty.StringVal(resourceGroup),
			"default_node_pool_vm_size":    cty.StringVal(clusterAKSInstanceType),
			"default_node_pool_min_count":  cty.NumberIntVal(clusterAKSMinNodes),
			"default_node_pool_node_count": cty.NumberIntVal(clusterAKSMinNodes),
			"default_node_pool_max_count":  cty.NumberIntVal(clusterAKSMaxNodes),
			"availability_zones":           cty.StringVal(strings.Join(zones, ",")),
		}

		_, err = s.AddCluster(namePrefix, "azurerm", region, "", stack.GenerateConfigurations(s.Environments, baseCfg))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var clusterAddEKSCmd = &cobra.Command{
	Use:   "eks <name-prefix> <region>",
	Short: "Add an EKS cluster module",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		namePrefix := args[0]
		region := args[1]

		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot(path)
		s := stack.NewStack(r, cj)
		err = s.FromPath()
		if err != nil {
			log.Fatal(err)
		}

		zones := strings.Split(clusterEKSZones, ",")
		if clusterEKSZones == "" {
			zones = cj.CloudInfo.Zones("aws", region, clusterEKSInstanceType)
		}

		baseCfg := map[string]cty.Value{
			"name_prefix":                cty.StringVal(namePrefix),
			"cluster_availability_zones": cty.StringVal(strings.Join(zones, ",")),
			"cluster_instance_type":      cty.StringVal(clusterEKSInstanceType),
			"cluster_min_size":           cty.NumberIntVal(clusterEKSMinNodes),
			"cluster_desired_capacity":   cty.NumberIntVal(clusterEKSMinNodes),
			"cluster_max_size":           cty.NumberIntVal(clusterEKSMaxNodes),
		}

		_, err = s.AddCluster(namePrefix, "aws", region, "", stack.GenerateConfigurations(s.Environments, baseCfg))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var clusterAddGKECmd = &cobra.Command{
	Use:   "gke <name-prefix> <region> <project-id>",
	Short: "Add a GKE cluster module",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		namePrefix := args[0]
		region := args[1]
		projectID := args[2]

		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot(path)
		s := stack.NewStack(r, cj)
		err = s.FromPath()
		if err != nil {
			log.Fatal(err)
		}

		zones := strings.Split(clusterGKEZones, ",")
		if clusterGKEZones == "" {
			zones = cj.CloudInfo.Zones("google", region, clusterGKEInstanceType)
		}

		baseCfg := map[string]cty.Value{
			"name_prefix":                cty.StringVal(namePrefix),
			"project_id":                 cty.StringVal(projectID),
			"region":                     cty.StringVal(region),
			"cluster_min_node_count":     cty.NumberIntVal(clusterGKEMinNodes),
			"cluster_initial_node_count": cty.NumberIntVal(clusterGKEMinNodes),
			"cluster_max_node_count":     cty.NumberIntVal(clusterGKEMaxNodes),
			"cluster_node_locations":     cty.StringVal(strings.Join(zones, ",")),
			"cluster_machine_type":       cty.StringVal(clusterGKEInstanceType),
			"cluster_min_master_version": cty.StringVal("1.25"),
		}

		_, err = s.AddCluster(namePrefix, "google", region, "", stack.GenerateConfigurations(s.Environments, baseCfg))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var nodePoolAddCmd = &cobra.Command{
	Use:     "node-pool command [flags]",
	Aliases: []string{"np"},
	Short:   "Add an AKS, EKS or GKE node pool module",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ErrMissingCommand
	},
}

var nodePoolAddAKSCmd = &cobra.Command{
	Use:   "aks <cluster-name> <pool-name>",
	Short: "Add a AKS node pool module",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		poolName := args[1]

		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot(path)
		s := stack.NewStack(r, cj)
		err = s.FromPath()
		if err != nil {
			log.Fatal(err)
		}

		baseCfg := map[string]cty.Value{
			"node_pool_name": cty.StringVal(poolName),
			"vm_size":        cty.StringVal(nodePoolAKSInstanceType),
			"node_count ":    cty.NumberIntVal(nodePoolAKSMinNodes),
			"min_count":      cty.NumberIntVal(nodePoolAKSMinNodes),
			"max_count":      cty.NumberIntVal(nodePoolAKSMaxNodes),
		}

		if len(nodePoolAKSZones) > 0 {
			baseCfg["availability_zones "] = cty.StringVal(nodePoolAKSZones)
		}

		if nodePoolAKSDiskSize != 0 {
			baseCfg["os_disk_size_gb"] = cty.NumberIntVal(nodePoolAKSDiskSize)
		}

		_, err = s.AddNodePool(clusterName, poolName, stack.GenerateConfigurations(s.Environments, baseCfg))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var nodePoolAddEKSCmd = &cobra.Command{
	Use:   "eks <cluster-name> <pool-name>",
	Short: "Add a EKS node pool module",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		poolName := args[1]

		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot(path)
		s := stack.NewStack(r, cj)
		err = s.FromPath()
		if err != nil {
			log.Fatal(err)
		}

		baseCfg := map[string]cty.Value{
			"name":              cty.StringVal(poolName),
			"instance_types":    cty.StringVal(nodePoolEKSInstanceType),
			"desired_capacity ": cty.NumberIntVal(nodePoolEKSMinNodes),
			"min_size":          cty.NumberIntVal(nodePoolEKSMinNodes),
			"max_size":          cty.NumberIntVal(nodePoolEKSMaxNodes),
		}

		if len(nodePoolEKSZones) > 0 {
			baseCfg["availability_zones "] = cty.StringVal(nodePoolEKSZones)
		}

		if nodePoolEKSAMIType != "" {
			baseCfg["ami_type"] = cty.StringVal(nodePoolEKSAMIType)
		}

		if nodePoolEKSDiskSize != 0 {
			baseCfg["disk_size"] = cty.NumberIntVal(nodePoolEKSDiskSize)
		}

		_, err = s.AddNodePool(clusterName, poolName, stack.GenerateConfigurations(s.Environments, baseCfg))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var nodePoolAddGKECmd = &cobra.Command{
	Use:   "gke <cluster-name> <pool-name>",
	Short: "Add a GKE node pool module",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clusterName := args[0]
		poolName := args[1]

		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot(path)
		s := stack.NewStack(r, cj)
		err = s.FromPath()
		if err != nil {
			log.Fatal(err)
		}

		baseCfg := map[string]cty.Value{
			"name":               cty.StringVal(poolName),
			"min_node_count":     cty.NumberIntVal(nodePoolGKEMinNodes),
			"initial_node_count": cty.NumberIntVal(nodePoolGKEMinNodes),
			"max_node_count":     cty.NumberIntVal(nodePoolGKEMaxNodes),
			"machine_type":       cty.StringVal(nodePoolGKEInstanceType),
		}

		var zones []cty.Value
		if len(nodePoolGKEZones) > 0 {
			for _, z := range strings.Split(nodePoolGKEZones, ",") {
				zones = append(zones, cty.StringVal(z))
			}
			baseCfg["node_locations"] = cty.ListVal(zones)
		}

		if nodePoolGKEImageType != "" {
			baseCfg["image_type"] = cty.StringVal(nodePoolGKEImageType)
		}

		if nodePoolGKEDiskType != "" {
			baseCfg["disk_type"] = cty.StringVal(nodePoolGKEDiskType)
		}

		if nodePoolGKEDiskSize != 0 {
			baseCfg["disk_size_gb"] = cty.NumberIntVal(nodePoolGKEDiskSize)
		}

		_, err = s.AddNodePool(clusterName, poolName, stack.GenerateConfigurations(s.Environments, baseCfg))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var serviceAddCmd = &cobra.Command{
	Use:     "service <name>",
	Aliases: []string{"svc"},
	Short:   "Add a service module",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entryName := args[0]

		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot(path)
		s := stack.NewStack(r, cj)
		err = s.FromPath()
		if err != nil {
			log.Fatal(err)
		}

		clusters := s.Clusters()

		for _, c := range clusters {
			currentClusterName := c.Name()

			if serviceClusterName != "" && serviceClusterName != currentClusterName {
				continue
			}

			_, err = s.AddService(currentClusterName, entryName, serviceRelease)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	sharedFlags.StringVarP(&clusterNamePrefix, "name-prefix", "n", "", "cluster name prefix")
	sharedFlags.StringVarP(&clusterRegion, "region", "r", "", "cluster region")

	// Clusters
	addCmd.AddCommand(clusterAddCmd)
	clusterAddCmd.PersistentFlags().AddFlagSet(&sharedFlags)

	clusterAddCmd.AddCommand(clusterAddAKSCmd)
	clusterAddAKSCmd.Flags().StringVar(&clusterAKSInstanceType, "aks-vm-size", "Standard_D2_v4", "vm size of nodes")
	clusterAddAKSCmd.Flags().Int64Var(&clusterAKSMinNodes, "aks-min", 3, "min number of nodes")
	clusterAddAKSCmd.Flags().Int64Var(&clusterAKSMaxNodes, "aks-max", 9, "max number of nodes")
	clusterAddAKSCmd.Flags().StringVar(&clusterAKSZones, "aks-availability-zones", "", "zones to use for nodes (default inherit cluster zones)")

	clusterAddCmd.AddCommand(clusterAddEKSCmd)
	clusterAddEKSCmd.Flags().StringVar(&clusterEKSInstanceType, "eks-instance-type", "t3a.xlarge", "instance type of nodes")
	clusterAddEKSCmd.Flags().Int64Var(&clusterEKSMinNodes, "eks-min", 3, "min number of nodes")
	clusterAddEKSCmd.Flags().Int64Var(&clusterEKSMaxNodes, "eks-max", 9, "max number of nodes")
	clusterAddEKSCmd.Flags().StringVar(&clusterEKSZones, "eks-availability-zones", "", "zones to use for nodes (default inherit cluster zones)")

	clusterAddCmd.AddCommand(clusterAddGKECmd)
	clusterAddGKECmd.Flags().StringVar(&clusterGKEInstanceType, "gke-machine-type", "e2-standard-8", "machine type of nodes")
	clusterAddGKECmd.Flags().Int64Var(&clusterGKEMinNodes, "gke-min", 1, "min number of nodes per zone")
	clusterAddGKECmd.Flags().Int64Var(&clusterGKEMaxNodes, "gke-max", 3, "max number of nodes per zone")
	clusterAddGKECmd.Flags().StringVar(&clusterGKEZones, "gke-node-locations", "", "zones to use for nodes (default inherit cluster zones)")

	// Node Pools
	addCmd.AddCommand(nodePoolAddCmd)

	nodePoolAddCmd.AddCommand(nodePoolAddAKSCmd)
	nodePoolAddAKSCmd.Flags().StringVar(&nodePoolAKSInstanceType, "aks-vm-size", "Standard_D2_v4", "vm size of nodes")
	nodePoolAddAKSCmd.Flags().Int64Var(&nodePoolAKSMinNodes, "aks-min", 3, "min number of nodes")
	nodePoolAddAKSCmd.Flags().Int64Var(&nodePoolAKSMaxNodes, "aks-max", 9, "max number of nodes")
	nodePoolAddAKSCmd.Flags().Int64Var(&nodePoolAKSDiskSize, "aks-disk-size", 0, "disk size of nodes in GB")
	nodePoolAddAKSCmd.Flags().StringVar(&nodePoolAKSZones, "aks-availability-zones", "", "zones to use for nodes (default 3 zones from the cluster's region)")

	nodePoolAddCmd.AddCommand(nodePoolAddEKSCmd)
	nodePoolAddEKSCmd.Flags().StringVar(&nodePoolEKSInstanceType, "eks-instance-type", "t3a.xlarge", "instance type of nodes")
	nodePoolAddEKSCmd.Flags().Int64Var(&nodePoolEKSMinNodes, "eks-min", 3, "min number of nodes")
	nodePoolAddEKSCmd.Flags().Int64Var(&nodePoolEKSMaxNodes, "eks-max", 9, "max number of nodes")
	nodePoolAddEKSCmd.Flags().Int64Var(&nodePoolEKSDiskSize, "eks-disk-size", 0, "disk size of nodes in GB")
	nodePoolAddEKSCmd.Flags().StringVar(&nodePoolEKSZones, "eks-availability-zones", "", "zones to use for nodes (default 3 zones from the cluster's region)")
	nodePoolAddEKSCmd.Flags().StringVar(&nodePoolEKSAMIType, "eks-ami-type", "", "AMI type of nodes (default EKS)")

	nodePoolAddCmd.AddCommand(nodePoolAddGKECmd)
	nodePoolAddGKECmd.Flags().StringVar(&nodePoolGKEInstanceType, "gke-machine-type", "e2-standard-8", "machine type of nodes")
	nodePoolAddGKECmd.Flags().StringVar(&nodePoolGKEImageType, "gke-image-type", "", "image type of nodes")
	nodePoolAddGKECmd.Flags().Int64Var(&nodePoolGKEMinNodes, "gke-min", 1, "min number of nodes per zone")
	nodePoolAddGKECmd.Flags().Int64Var(&nodePoolGKEMaxNodes, "gke-max", 3, "max number of nodes per zone")
	nodePoolAddGKECmd.Flags().StringVar(&nodePoolGKEDiskType, "gke-disk-type", "", "disk type of nodes")
	nodePoolAddGKECmd.Flags().Int64Var(&nodePoolGKEDiskSize, "gke-disk-size", 0, "disk size of nodes in GB")
	nodePoolAddGKECmd.Flags().StringVar(&nodePoolGKEZones, "gke-node-locations", "", "zones to use for nodes (default cluster's zones)")

	// Services
	addCmd.AddCommand(serviceAddCmd)
	serviceAddCmd.Flags().StringVarP(&serviceRelease, "release", "r", "latest", "desired release version")
	serviceAddCmd.Flags().StringVarP(&serviceClusterName, "cluster-name", "c", "", "add service to single cluster (default add to all clusters)")
}
