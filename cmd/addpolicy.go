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
	"errors"
	"fmt"

	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

var addPolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Add a new policy",
	Long: AddAppName(`Add a new policy
    Example:
    $AppName add policy name --organization org --environment env --file ./policydefinition.json -i json`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Policy name needed.")
		}
		name := args[0]

		logging.Info("Adding new policy to environment %v in organization %v\n", organization, environment)
		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		policyBytes, readErr := util.UseSourceUrl(configPath) // just pass the file name
		if readErr != nil {
			return readErr
		}

		policyJSON, jsonError := util.Convert(configFileType, "json", policyBytes)
		if jsonError != nil {
			return jsonError
		}

		policyText := string(policyJSON)

		createpolicyError := core.AddReleasePolicy(namespacedOrganization, namespacedEnvironment, name, policyText)
		if createpolicyError != nil {
			return createpolicyError
		}

		fmt.Printf("Policy has been added\n")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addPolicyCmd)

	addPolicyCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the workflow")
	addPolicyCmd.MarkFlagRequired("organization")

	addPolicyCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the workflow")
	addPolicyCmd.MarkFlagRequired("environment")

	addPolicyCmd.Flags().StringVarP(&configPath, "file", "", "", "Policy configuration file path")
	addPolicyCmd.MarkFlagRequired("file")
	addPolicyCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Policy configuration file type yaml or json")

}
