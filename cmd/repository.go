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

	"github.com/kbst/kbst/cli"
	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

var initRelease string
var initGitRef string
var initEnvNames string

var initCmd = &cobra.Command{
	Use:   "init command [flags]",
	Short: "Scaffold a new Kubestack repository",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ErrMissingCommand
	},
}

var initAKSCmd = &cobra.Command{
	Use:   "aks <base-domain> <name-prefix> <region> <resource-group>",
	Short: "Scaffold a repository with one AKS cluster",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		initStarter("aks", args)
	},
}

var initEKSCmd = &cobra.Command{
	Use:   "eks <base-domain> <name-prefix> <region>",
	Short: "Scaffold a repository with one EKS cluster",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		initStarter("eks", args)
	},
}

var initGKECmd = &cobra.Command{
	Use:   "gke <base-domain> <name-prefix> <region> <project-id>",
	Short: "Scaffold a repository with one GKE cluster",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		initStarter("gke", args)
	},
}

func initStarter(starter string, args []string) {
	cj := util.CliJSON{}
	err := cj.Load(util.CachedDownloader{})
	if err != nil {
		log.Fatal(err)
	}
	r := cli.Repo{
		Framework:  cj.Framework,
		Downloader: util.CachedDownloader{},
	}

	baseDomain := args[0]
	namePrefix := args[1]
	region := args[2]

	var baseCfg map[string]cty.Value

	switch starter {
	case "aks":
		zones := strings.Split(clusterAKSZones, ",")
		if clusterAKSZones == "" {
			zones = cj.CloudInfo.Zones("azurerm", region, clusterAKSInstanceType)
		}

		baseCfg = map[string]cty.Value{
			"name_prefix":                  cty.StringVal(namePrefix),
			"resource_group":               cty.StringVal(args[3]),
			"default_node_pool_vm_size":    cty.StringVal(clusterAKSInstanceType),
			"default_node_pool_min_count":  cty.NumberIntVal(clusterAKSMinNodes),
			"default_node_pool_node_count": cty.NumberIntVal(clusterAKSMinNodes),
			"default_node_pool_max_count":  cty.NumberIntVal(clusterAKSMaxNodes),
			"availability_zones":           cty.StringVal(strings.Join(zones, ",")),
		}
	case "eks":
		zones := strings.Split(clusterEKSZones, ",")
		if clusterEKSZones == "" {
			zones = cj.CloudInfo.Zones("aws", region, clusterEKSInstanceType)
		}

		baseCfg = map[string]cty.Value{
			"name_prefix":                cty.StringVal(namePrefix),
			"cluster_availability_zones": cty.StringVal(strings.Join(zones, ",")),
			"cluster_instance_type":      cty.StringVal(clusterEKSInstanceType),
			"cluster_min_size":           cty.NumberIntVal(clusterEKSMinNodes),
			"cluster_desired_capacity":   cty.NumberIntVal(clusterEKSMinNodes),
			"cluster_max_size":           cty.NumberIntVal(clusterEKSMaxNodes),
		}
	case "gke":
		zones := strings.Split(clusterGKEZones, ",")
		if clusterGKEZones == "" {
			zones = cj.CloudInfo.Zones("google", region, clusterGKEInstanceType)
		}

		baseCfg = map[string]cty.Value{
			"name_prefix":                cty.StringVal(namePrefix),
			"project_id":                 cty.StringVal(args[3]),
			"region":                     cty.StringVal(region),
			"cluster_min_node_count":     cty.NumberIntVal(clusterGKEMinNodes),
			"cluster_initial_node_count": cty.NumberIntVal(clusterGKEMinNodes),
			"cluster_max_node_count":     cty.NumberIntVal(clusterGKEMaxNodes),
			"cluster_node_locations":     cty.StringVal(strings.Join(zones, ",")),
			"cluster_machine_type":       cty.StringVal(clusterGKEInstanceType),
			"cluster_min_master_version": cty.StringVal("1.20"),
		}
	default:
		log.Fatalf("unexpected error: starter: '%s' exists as archive, but is not implemented in CLI", starter)
	}

	err = r.Init(starter, baseDomain, namePrefix, region, strings.Split(initEnvNames, ","), baseCfg, initRelease, initGitRef, path)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().AddFlagSet(&sharedFlags)
	initCmd.PersistentFlags().StringVar(&initEnvNames, "environment-names", "apps,ops", "list of environment names, mission critical first")

	initCmd.PersistentFlags().StringVar(&initRelease, "release", "latest", "desired release version")
	initCmd.PersistentFlags().StringVar(&initGitRef, "gitref", "", "git ref to download a dev artifact")
	initCmd.PersistentFlags().MarkHidden("gitref")

	initCmd.AddCommand(initAKSCmd)
	initAKSCmd.Flags().AddFlagSet(clusterAddAKSCmd.Flags())

	initCmd.AddCommand(initEKSCmd)
	initEKSCmd.Flags().AddFlagSet(clusterAddEKSCmd.Flags())

	initCmd.AddCommand(initGKECmd)
	initGKECmd.Flags().AddFlagSet(clusterAddGKECmd.Flags())
}
