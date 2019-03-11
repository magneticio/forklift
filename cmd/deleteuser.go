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
	"github.com/spf13/cobra"
)

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Delete a user",
	Long: AddAppName(`Delete a user
    Example:
    $AppName delete user name --organization org`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Name needed.")
		}
		name := args[0]

		namespaced := Config.Namespace + "-" + organization
		fmt.Printf("name: %v , configPath: %v\n", name, namespaced)

		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}

		deleteUserError := core.DeleteUser(namespaced, name)
		if deleteUserError != nil {
			return deleteUserError
		}
		fmt.Printf("User is deleted\n")

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteUserCmd)

	deleteUserCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	deleteUserCmd.MarkFlagRequired("organization")

}
