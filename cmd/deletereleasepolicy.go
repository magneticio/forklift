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

var deleteReleasepolicyCmd = &cobra.Command{
	Use:   "releasepolicy",
	Short: "Delete a releasepolicy",
	Long: AddAppName(`Delete a new releasepolicy
    Example:
    $AppName delete releasepolicy name --organization org --environment env`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Release Policy name needed.")
		}
		name := args[0]

		if !util.ValidateName(name) {
			return errors.New("Release Policy name should be only lowercase alphanumerics")
		}

		logging.Info("Deleteing new releasepolicy to environment %v in organization %v\n", organization, environment)
		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		createreleasepolicyError := core.DeleteReleasePolicy(namespacedOrganization, namespacedEnvironment, name)
		if createreleasepolicyError != nil {
			return createreleasepolicyError
		}

		fmt.Printf("Release Policy is deleted\n")

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteReleasepolicyCmd)

	deleteReleasepolicyCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the workflow")
	deleteReleasepolicyCmd.MarkFlagRequired("organization")

	deleteReleasepolicyCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the workflow")
	deleteReleasepolicyCmd.MarkFlagRequired("environment")

}
