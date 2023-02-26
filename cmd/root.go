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
	"errors"
	"fmt"

	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

var ErrMissingCommand = errors.New("missing command")

var path string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kbst command [flags]",
	Short: "Kubestack Framework CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Version == "" {
			return
		}

		// check if a newer CLI version is available
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			return
		}

		if len(cj.Cli.Versions) > 1 {
			latest := cj.Cli.Versions[0].Name

			if semver.Compare(cmd.Version, latest) == -1 {
				fmt.Fprintf(cmd.OutOrStderr(), "The latest version %s of `kbst` is newer than your current version %s\n", latest, cmd.Version)
				fmt.Fprintf(cmd.OutOrStderr(), "To update visit: https://github.com/kbst/kbst/releases/tag/%v\n", latest)
				fmt.Fprint(cmd.OutOrStderr(), "\n")
			}
		}
	},
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return ErrMissingCommand
	},
}

func Execute(version, commit string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "path to the working directory")
}
