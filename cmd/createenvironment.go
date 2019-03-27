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
	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

// environmentCmd represents the environment command
var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Create a new environment",
	Long: AddAppName(`Create a new environment
    Example:
    $AppName create environment env1 --organization org --file ./environmentdefinition.yaml --artifacts ./artifactsfolder`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Organization Name needed.")
		}
		name := args[0]

		if !util.ValidateName(name) {
			return errors.New("Environment name should be only lowercase alphanumerics")
		}

		logging.Info("Creating environment %v in organization %v\n", name, organization)

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
		var vampConfig models.VampConfiguration
		unmarshallError := json.Unmarshal([]byte(configJson), &vampConfig)
		if unmarshallError != nil {
			return unmarshallError
		}

		core, coreError := core.NewCore(Config)
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
				logging.Info("Using file as artifact %v\n", file)
				contentJson, jsonError := util.Convert("yaml", "json", content)
				if jsonError != nil {
					return jsonError
				}
				params = append(params, contentJson)
			}
		}
		createOrganizationError := core.CreateEnvironment(namespaced, namespacedOrganization, params, vampConfig)
		if createOrganizationError != nil {
			return createOrganizationError
		}
		fmt.Printf("Environment %v is created\n", name)
		return nil
	},
}

func init() {
	createCmd.AddCommand(environmentCmd)

	environmentCmd.Flags().StringVarP(&configPath, "file", "", "", "Environment configuration file path")
	environmentCmd.MarkFlagRequired("file")
	environmentCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Configuration file type yaml or json")
	environmentCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	environmentCmd.MarkFlagRequired("organization")
	environmentCmd.Flags().StringVarP(&artficatsPath, "artifacts", "", "", "Path to the folder containing artifacts")

}
