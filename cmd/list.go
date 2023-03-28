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
	"os"
	"strings"
	"text/tabwriter"

	"github.com/kbst/kbst/pkg/stack"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var showAll bool

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

		w := tabwriter.NewWriter(os.Stdout, 4, 8, 2, '\t', 0)
		line := "%s\t%s\t%s\n"

		filterCustomModules := []string{}
		for _, c := range s.Clusters {
			fmt.Fprintf(w, line, "TYPE", "NAME", "VERSION")
			fmt.Fprintf(w, line, "cluster", c.Name(), c.Version)

			for _, np := range s.NodePools {
				if np.ClusterName == c.Name() {
					fmt.Fprintf(w, line, "node-pool", np.Name(), np.Version)
				}
			}

			for _, svc := range s.Services {
				if svc.ClusterName == c.Name() {
					fmt.Fprintf(w, line, "service", svc.Name(), svc.Version)
				}
			}

			if showAll {
				for _, m := range s.Modules {
					if strings.HasPrefix(m.Name(), fmt.Sprintf("%s_", c.Name())) {
						// if custom modules are prefixed
						// with the name of a cluster show them here
						filterCustomModules = append(filterCustomModules, m.Name())
						fmt.Fprintf(w, line, "custom", m.Name(), "")
					}
				}
			}

			fmt.Fprintf(w, line, "", "", "")
		}

		// show custom modules that were not
		// prefixed with a cluster name
		if showAll && len(s.Modules) > 0 {
			fmt.Fprintf(w, line, "TYPE", "NAME", "VERSION")
			for _, m := range s.Modules {
				if !slices.Contains(filterCustomModules, m.Name()) {
					fmt.Fprintf(w, line, "custom", m.Name(), "")
				}
			}
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&showAll, "all", "a", false, "List all modules")
}
