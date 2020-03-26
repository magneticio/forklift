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
	"github.com/spf13/cobra"
)

var deletePolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Delete a policy",
	Long: AddAppName(`Delete existing policy
    Example:
    $AppName delete policy name --organization org --environment env`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Policy name needed.")
		}
		name := args[0]

		logging.Info("Deleteing new policy to environment %v in organization %v\n", organization, environment)
		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		createpolicyError := core.DeleteReleasePolicy(namespacedOrganization, namespacedEnvironment, name)
		if createpolicyError != nil {
			return createpolicyError
		}

		fmt.Printf("Policy has been deleted\n")

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deletePolicyCmd)

	deletePolicyCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the workflow")
	deletePolicyCmd.MarkFlagRequired("organization")

	deletePolicyCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the workflow")
	deletePolicyCmd.MarkFlagRequired("environment")

}
