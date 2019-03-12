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
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// environmentsCmd represents the environments command
var listArtifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "list environments in an organization",
	Long: AddAppName(`List artifacts
    Example:
    $AppName list artifacts --kind k --organization org --environment env`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		logging.Info.Log("Listing %v in environment %v in organization %v\n", kind, environment, organization)

		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}

		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}
		environmentList, listOrganizationsError := core.ListArtifacts(namespacedOrganization, namespacedEnvironment, kind)
		if listOrganizationsError != nil {
			return listOrganizationsError
		}
		output, marshallError := yaml.Marshal(environmentList)
		if marshallError != nil {
			return marshallError
		}
		fmt.Print(string(output))
		return nil
	},
}

func init() {
	listCmd.AddCommand(listArtifactsCmd)

	listArtifactsCmd.Flags().StringVarP(&kind, "kind", "", "", "Kind of the artifacts")
	listArtifactsCmd.MarkFlagRequired("kind")

	listArtifactsCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the artifacts")
	listArtifactsCmd.MarkFlagRequired("organization")

	listArtifactsCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the artifacts")
	listArtifactsCmd.MarkFlagRequired("environment")
}