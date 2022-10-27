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
	"os"

	"github.com/kbst/kbst/pkg/util"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

var Version string

//var cfgFile string
var path string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kbst",
	Short: "Kubestack Framework CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// check if a newer CLI version is available
		cj := util.CliJSON{}
		err := cj.Load(util.CachedDownloader{})
		if err != nil {
			return
		}

		if len(cj.Cli.Versions) > 1 {
			current := Version
			latest := cj.Cli.Versions[0].Name

			if semver.Compare(current, latest) == -1 {
				fmt.Fprintf(cmd.OutOrStderr(), "The latest version %s of `kbst` is newer than your current version %s\n", latest, current)
				fmt.Fprintf(cmd.OutOrStderr(), "To update visit: https://github.com/kbst/kbst/releases/tag/%v\n", latest)
				fmt.Fprint(cmd.OutOrStderr(), "\n")
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Version = Version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kbst.yaml)")

	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "path to the working directory")
}

// initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	if cfgFile != "" {
// 		// Use config file from the flag.
// 		viper.SetConfigFile(cfgFile)
// 	} else {
// 		// Find home directory.
// 		home, err := homedir.Dir()
// 		if err != nil {
// 			fmt.Println(err)
// 			os.Exit(1)
// 		}
//
// 		// Search config in home directory with name ".kbst" (without extension).
// 		viper.AddConfigPath(home)
// 		viper.SetConfigName(".kbst")
// 	}
//
// 	viper.AutomaticEnv() // read in environment variables that match
//
// 	// If a config file is found, read it in.
// 	if err := viper.ReadInConfig(); err == nil {
// 		fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	}
// }
//
