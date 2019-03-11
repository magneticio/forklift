// Copyright © 2019 Developer <developer@vamp.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

// updateenvironmentCmd represents the environment command
var updateenvironmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Update a new environment",
	Long: AddAppName(`Update a new environment
    Example:
    $AppName update environment env1 --organization org --file ./somepath.json`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Organization Name needed.")
		}
		name := args[0]
		namespaced := Config.Namespace + "-" + organization + "-" + name
		namespacedOrganization := Config.Namespace + "-" + organization

		configBtye, readErr := util.UseSourceUrl(configPath) // just pass the file name
		if readErr != nil {
			return readErr
		}
		configText := string(configBtye)
		configJson, convertErr := util.Convert(configFileType, "json", configText)
		if convertErr != nil {
			return convertErr
		}
		var vampConfig core.Configuration
		unmarshallError := json.Unmarshal([]byte(configJson), &vampConfig)
		if unmarshallError != nil {
			return unmarshallError
		}
		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}
		params := []string{}
		if artficatsPath != "" {
			artifacts, artifactsReadError := util.ReadFilesIndirectory(artficatsPath)
			if artifactsReadError != nil && artficatsPath != "" {
				return artifactsReadError
			}
			for file, content := range artifacts {
				fmt.Printf("Using file as artifact %v\n", file)
				contentJson, jsonError := util.Convert("yaml", "json", content)
				if jsonError != nil {
					return jsonError
				}
				params = append(params, contentJson)
			}
		}
		createOrganizationError := core.UpdateEnvironment(namespaced, namespacedOrganization, params, vampConfig)
		if createOrganizationError != nil {
			return createOrganizationError
		}
		fmt.Printf("Environment %v is created\n", name)
		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateenvironmentCmd)

	updateenvironmentCmd.Flags().StringVarP(&configPath, "file", "", "", "Environment configuration file path")
	updateenvironmentCmd.MarkFlagRequired("file")
	updateenvironmentCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Configuration file type yaml or json")
	updateenvironmentCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	updateenvironmentCmd.MarkFlagRequired("organization")
	updateenvironmentCmd.Flags().StringVarP(&artficatsPath, "artifacts", "", "", "Path to the folder containing artifacts")

}
