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
	"encoding/json"
	"log"

	"github.com/kbst/kbst/cli"
	"github.com/kbst/kbst/pkg/export"
	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:    "import JSON",
	Short:  "Create new repository imported from JSON string",
	Args:   cobra.ExactArgs(1),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}
		r := cli.Repo{
			Framework:  cj.Framework,
			Downloader: util.CachedDownloader{},
		}

		iS := export.Stack{}
		err = json.Unmarshal([]byte(args[0]), &iS)
		if err != nil {
			log.Fatal(err)
		}

		err = r.Import(iS, path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
}
