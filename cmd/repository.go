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

var repoRelease string
var repoGitRef string

// repositoryCmd represents the repository command
var repositoryCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repo"},
	Short:   "Create and change Kubestack repositories",
}

// repositoryInitCmd represents the repository init command
var repositoryInitCmd = &cobra.Command{
	Use:   "init <starter>",
	Short: "Scaffold a new repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		starter := args[0]
		err := cli.RepoInit(starter, repoRelease, repoGitRef, path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(repositoryCmd)

	repositoryCmd.AddCommand(repositoryInitCmd)

	repositoryInitCmd.Flags().StringVarP(&repoRelease, "release", "r", "latest", "desired release version")
	repositoryInitCmd.Flags().StringVar(&repoGitRef, "gitref", "", "git ref to download a dev artifact")
	repositoryInitCmd.Flags().MarkHidden("gitref")
}
