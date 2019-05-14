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

var addReleasingpolicyCmd = &cobra.Command{
	Use:   "releasingpolicy",
	Short: "Add a new releasingpolicy",
	Long: AddAppName(`Add a new releasingpolicy
    Example:
    $AppName add releasingpolicy --organization org --environment env --file ./releasingpolicydefinition.yson -i json`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.Info("Adding new releasingpolicy to environment %v in organization %v\n", organization, environment)
		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		releasingpolicyBytes, readErr := util.UseSourceUrl(configPath) // just pass the file name
		if readErr != nil {
			return readErr
		}

		releasingpolicyJSON, jsonError := util.Convert(configFileType, "json", releasingpolicyBytes)
		if jsonError != nil {
			return jsonError
		}

		releasingpolicyText := string(releasingpolicyJSON)

		createreleasingpolicyError := core.AddReleasingPolicy(namespacedOrganization, namespacedEnvironment, releasingpolicyText)
		if createreleasingpolicyError != nil {
			return createreleasingpolicyError
		}

		fmt.Printf("Releasing Policy is added\n")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addReleasingpolicyCmd)

	addReleasingpolicyCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the workflow")
	addReleasingpolicyCmd.MarkFlagRequired("organization")

	addReleasingpolicyCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the workflow")
	addReleasingpolicyCmd.MarkFlagRequired("environment")

	addReleasingpolicyCmd.Flags().StringVarP(&configPath, "file", "", "", "Releasing policy configuration file path")
	addReleasingpolicyCmd.MarkFlagRequired("file")
	addReleasingpolicyCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Releasing policy configuration file type yaml or json")

}
