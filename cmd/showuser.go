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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/logging"
	"github.com/spf13/cobra"
)

// showenvironmentCmd represents the showenvironment command
var showUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Show an user",
	Long: AddAppName(`Show an user
    Example:
    $AppName show user name --organization org`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, User Name needed.")
		}
		name := args[0]

		logging.Info("Showing user %v in organization %v\n", name, organization)

		namespaced := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}
		user, showEnvironmentError := core.GetUser(namespaced, name)
		if showEnvironmentError != nil {
			return showEnvironmentError
		}
		jsonUser, jsonMarshallError := json.Marshal(user)
		if jsonMarshallError != nil {
			return jsonMarshallError
		}
		fmt.Printf("%v", string(jsonUser))
		return nil
	},
}

func init() {
	showCmd.AddCommand(showUserCmd)

	showUserCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	showUserCmd.MarkFlagRequired("organization")
}
