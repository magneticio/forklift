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
	"fmt"

	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

var addArtifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Add a new artifact",
	Long: AddAppName(`Add a new artifact
    Example:
    $AppName add artifact --organization org --environment env --file ./somepath.yaml`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.Log("Adding new artifact to environment %v in organization %v\n", organization, environment)

		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}

		artifactBytes, readErr := util.UseSourceUrl(configPath) // just pass the file name
		if readErr != nil {
			return readErr
		}

		artifactJson, jsonError := util.Convert(configFileType, "json", artifactBytes)
		if jsonError != nil {
			return jsonError
		}

		artifactText := string(artifactJson)

		createArtifactError := core.AddArtifact(namespacedOrganization, namespacedEnvironment, artifactText)
		if createArtifactError != nil {
			return createArtifactError
		}
		fmt.Printf("Artifact is added\n")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addArtifactCmd)

	addArtifactCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the workflow")
	addArtifactCmd.MarkFlagRequired("organization")

	addArtifactCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the workflow")
	addArtifactCmd.MarkFlagRequired("environment")

	addArtifactCmd.Flags().StringVarP(&configPath, "file", "", "", "Artifact configuration file path")
	addArtifactCmd.MarkFlagRequired("file")
	addArtifactCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Artifact configuration file type yaml or json")

}
