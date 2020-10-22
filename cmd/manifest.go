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

	"github.com/kbst/kbst/cli"
	"github.com/kbst/kbst/util"
	"github.com/spf13/cobra"
)

var manifestRelease string
var manifestOverlay string
var manifestGitRef string
var manifestSkipEditKustomization bool

// devCmd represents the dev command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Add, update and remove services from the catalog",
}

var manifestInstallCmd = &cobra.Command{
	Use:   "install <entry> <variant>",
	Short: "Install and vendor a manifest from the catalog",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		entry := args[0]
		variant := args[1]
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}
		m := cli.Manifest{
			Catalog:    cj.Catalog,
			Downloader: util.CachedDownloader{},
		}
		err = m.Install(entry, variant, manifestOverlay, manifestRelease, manifestGitRef, path, manifestSkipEditKustomization)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var manifestRemoveCmd = &cobra.Command{
	Use:   "remove <entry>",
	Short: "Remove a vendored manifest from all environments",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entry := args[0]
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}
		m := cli.Manifest{
			Catalog:    cj.Catalog,
			Downloader: util.CachedDownloader{},
		}
		err = m.Remove(entry, path, manifestSkipEditKustomization)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var manifestUpdateCmd = &cobra.Command{
	Use:   "update <entry>",
	Short: "Update vendored manifests from the catalog",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entry := args[0]
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}
		m := cli.Manifest{
			Catalog:    cj.Catalog,
			Downloader: util.CachedDownloader{},
		}
		err = m.Update(entry, manifestOverlay, manifestRelease, manifestGitRef, path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(manifestCmd)

	// InstallCmd
	manifestCmd.AddCommand(manifestInstallCmd)
	manifestInstallCmd.Flags().StringVarP(&manifestRelease, "release", "r", "latest", "desired release version")
	manifestInstallCmd.Flags().StringVarP(&manifestOverlay, "overlay", "o", "apps", "overlay to add resources to")
	manifestInstallCmd.Flags().BoolVar(&manifestSkipEditKustomization, "skip-edit-kustomization", false, "skip editing kustomization resources")
	manifestInstallCmd.Flags().StringVar(&manifestGitRef, "gitref", "", "git ref to download a dev artifact")
	manifestInstallCmd.Flags().MarkHidden("gitref")

	// RemoveCmd
	manifestCmd.AddCommand(manifestRemoveCmd)
	manifestRemoveCmd.Flags().BoolVar(&manifestSkipEditKustomization, "skip-edit-kustomization", false, "skip editing kustomization resources")

	// UpdateCmd
	manifestCmd.AddCommand(manifestUpdateCmd)
	manifestUpdateCmd.Flags().StringVarP(&manifestRelease, "release", "r", "latest", "desired release version")
	manifestUpdateCmd.Flags().BoolVar(&manifestSkipEditKustomization, "skip-edit-kustomization", false, "skip editing kustomization resources")
	manifestUpdateCmd.Flags().StringVar(&manifestGitRef, "gitref", "", "git ref to download a dev artifact")
	manifestUpdateCmd.Flags().MarkHidden("gitref")
}
