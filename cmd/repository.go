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
)

var initRelease string
var initGitRef string
var initEnvNames string

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo"},
	Short:   "Create and change Kubestack repositories",
	Hidden:  true,
}

var initCmd = &cobra.Command{
	Use:   "init <starter>",
	Short: "Scaffold a new Kubestack repository",
}

var initAKSCmd = &cobra.Command{
	Use:   "aks <base_domain>",
	Short: "Scaffold a repository with one AKS cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initStarter("aks", args[0])
	},
}

var initEKSCmd = &cobra.Command{
	Use:   "eks <base_domain>",
	Short: "Scaffold a repository with one EKS cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initStarter("eks", args[0])
	},
}

var initGKECmd = &cobra.Command{
	Use:   "gke <base_domain>",
	Short: "Scaffold a repository with one GKE cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initStarter("gke", args[0])
	},
}

var initMultiCloudCmd = &cobra.Command{
	Use:   "multi-cloud <base_domain>",
	Short: "Scaffold a repository with one AKS, one EKS and one GKE cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initStarter("multi-cloud", args[0])
	},
}

func initStarter(starter, baseDomain string) {
	cj := util.CliJSON{}
	err := cj.Load(util.CachedDownloader{})
	if err != nil {
		log.Fatal(err)
	}
	r := cli.Repo{
		Framework:  cj.Framework,
		Downloader: util.CachedDownloader{},
	}
	err = r.Init(starter, baseDomain, strings.Split(initEnvNames, ","), initRelease, initGitRef, path)
	if err != nil {
		log.Fatal(err)
	}
}

var repositoryGenerateCmd = &cobra.Command{
	Use:    "generate <json_path>",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		json_path := args[0]
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}
		r := cli.Repo{
			Catalog:    cj.Catalog,
			Framework:  cj.Framework,
			Downloader: util.CachedDownloader{},
		}
		err = r.Generate(json_path, path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.PersistentFlags().StringVar(&initEnvNames, "environment-names", "apps,ops", "list of environment names, mission critical first")
	initCmd.PersistentFlags().StringVarP(&initRelease, "release", "r", "latest", "desired release version")
	initCmd.PersistentFlags().StringVar(&initGitRef, "gitref", "", "git ref to download a dev artifact")
	initCmd.PersistentFlags().MarkHidden("gitref")

	initCmd.AddCommand(initAKSCmd)
	initAKSCmd.Flags().AddFlagSet(nodePoolAddAKSCmd.Flags())

	initCmd.AddCommand(initEKSCmd)
	initEKSCmd.Flags().AddFlagSet(nodePoolAddEKSCmd.Flags())

	initCmd.AddCommand(initGKECmd)
	initGKECmd.Flags().AddFlagSet(nodePoolAddGKECmd.Flags())

	initCmd.AddCommand(initMultiCloudCmd)
	initMultiCloudCmd.Flags().AddFlagSet(nodePoolAddAKSCmd.Flags())
	initMultiCloudCmd.Flags().AddFlagSet(nodePoolAddEKSCmd.Flags())
	initMultiCloudCmd.Flags().AddFlagSet(nodePoolAddGKECmd.Flags())

	rootCmd.AddCommand(repositoryCmd)
	repositoryCmd.AddCommand(repositoryGenerateCmd)
}
