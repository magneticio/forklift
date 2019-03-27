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

var addUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Add a new user",
	Long: AddAppName(`Add a new user
    Example:
    $AppName add user --organization org --file ./userdefinition.yaml`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		logging.Info("Adding new user to organization %v\n", organization)

		namespaced := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		userBytes, userErr := util.UseSourceUrl(configPath) // just pass the file name
		if userErr != nil {
			return userErr
		}

		userJson, jsonError := util.Convert(configFileType, "json", userBytes)
		if jsonError != nil {
			return jsonError
		}

		userText := string(userJson)

		createUserError := core.AddUser(namespaced, userText)
		if createUserError != nil {
			return createUserError
		}
		fmt.Printf("User is added\n")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addUserCmd)

	addUserCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	addUserCmd.MarkFlagRequired("organization")

	addUserCmd.Flags().StringVarP(&configPath, "file", "", "", "User configuration file path")
	addUserCmd.MarkFlagRequired("file")
	addUserCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "User configuration file type yaml or json")

}
