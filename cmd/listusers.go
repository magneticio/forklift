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
var listUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "list environments in an organization",
	Long: AddAppName(`List users
    Example:
    $AppName list users --organization organizaionname`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		logging.Info.Log("Listing users in organization %v\n", organization)

		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}

		namespaced := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}
		environmentList, listOrganizationsError := core.ListUsers(namespaced)
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
	listCmd.AddCommand(listUsersCmd)

	listUsersCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	listUsersCmd.MarkFlagRequired("organization")
}
