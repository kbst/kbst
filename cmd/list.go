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
	"fmt"
	"log"

	"github.com/kbst/kbst/pkg/stack"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters, node pools and services",
	Run: func(cmd *cobra.Command, args []string) {
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot()
		s := stack.NewStack(r, cj)
		err = s.FromPath(path)
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range s.Clusters {
			fmt.Println(c.Name())

			for _, np := range s.NodePools {
				if np.ClusterName == c.Name() {
					fmt.Println(np.Name())
				}
			}

			for _, svc := range s.Services {
				if svc.ClusterName == c.Name() {
					fmt.Println(svc.Name())
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
