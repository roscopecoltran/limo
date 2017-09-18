// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"fmt" // go-core
	//"github.com/roscopecoltran/sniperkit-limo/config" 										// app-config
	"github.com/sirupsen/logrus" // logs-logrus
	"github.com/spf13/cobra"     // cli-cmd
)

var PluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "manages plugins",
	Long:  "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plugin called")
	},
}

func init() {
	log.WithFields(
		logrus.Fields{
			"prefix":      "app-action",
			"src.file":    "action/plugin.go",
			"cmd.name":    "PluginCmd",
			"method.name": "init()",
			"var.options": options,
		}).Info("registering command...")
	RootCmd.AddCommand(PluginCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// PluginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// PluginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
