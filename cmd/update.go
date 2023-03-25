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

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/kbst/kbst/pkg/tfhcl"
	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update cluster, node pool or service module versions",
	Run: func(cmd *cobra.Command, args []string) {
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			log.Fatal(err)
		}

		r := tfhcl.NewRoot()
		err = r.Read(path)
		if err != nil {
			log.Fatalln(err)
		}

		for n, f := range r.Parser.Files() {
			wf, _ := hclwrite.ParseConfig(f.Bytes, n, hcl.InitialPos)

			for i, b := range wf.Body().Blocks() {
				if b.Type() != "module" {
					continue
				}

				source := r.Modules[n][i].Source
				var latestVersion string
				switch ref := strings.Contains(source, "?ref="); ref {
				case true:
					// module has version in source
					// like framework modules currently do
					latestVersion = cj.Framework.Versions[0].Name
				case false:
					// module has source and version
					// like cluster service modules currently do
					name := strings.Split(source, "/")[2]
					latestVersion = cj.Catalog[name].Versions[0].Name
				}

				if b.Body().GetAttribute("version") != nil {
					v := strings.TrimPrefix(latestVersion, "v")
					b.Body().SetAttributeValue("version", cty.StringVal(v))
				} else {
					currentRef := strings.Split(source, "?ref=")[1]
					s := strings.Replace(source, currentRef, latestVersion, 1)
					b.Body().SetAttributeValue("source", cty.StringVal(s))
				}
			}

			f.Bytes = wf.Bytes()
		}

		r.Write()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
