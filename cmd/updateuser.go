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
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

var updateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Update a user",
	Long: AddAppName(`Update a user
    Example:
    $AppName update user name --role r --organization org`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, User Name needed.")
		}
		name := args[0]

		namespaced := Config.Namespace + "-" + organization

		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}

		passwd, passwdError := util.GetParameterFromTerminalAsSecret(
			"Enter your password (password will not be visible):",
			"Enter your password again (password will not be visible):",
			"Passwords do not match.")
		if passwdError != nil {
			return passwdError
		}

		if len(passwd) < 6 {
			return errors.New("Password should be 6 or more characters")
		}

		updateUserError := core.UpdateUser(namespaced, name, role, passwd)
		if updateUserError != nil {
			return updateUserError
		}
		fmt.Printf("User is updated\n")

		return nil
	},
}

func init() {
	updateCmd.AddCommand(updateUserCmd)

	updateUserCmd.Flags().StringVarP(&role, "role", "", "", "Role of the user")
	updateUserCmd.MarkFlagRequired("role")
	updateUserCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	updateUserCmd.MarkFlagRequired("organization")

}