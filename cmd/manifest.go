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
	"github.com/spf13/cobra"
)

var manifestRelease string
var manifestOverlay string
var manifestGitRef string

// devCmd represents the dev command
var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Add, update and remove services from the catalog",
}

var manifestInstallCmd = &cobra.Command{
	Use:   "install <entry> <variant>",
	Short: "Install manifests from the catalog",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		entry := args[0]
		variant := args[1]
		err := cli.ManifestInstall(entry, variant, manifestOverlay, manifestRelease, manifestGitRef, path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.AddCommand(manifestInstallCmd)

	manifestInstallCmd.Flags().StringVarP(&manifestRelease, "release", "r", "latest", "desired release version")
	manifestInstallCmd.Flags().StringVarP(&manifestOverlay, "overlay", "o", "apps", "overlay to add resources to")
	manifestInstallCmd.Flags().StringVar(&manifestGitRef, "gitref", "", "git ref to download a dev artifact")
	manifestInstallCmd.Flags().MarkHidden("gitref")
}