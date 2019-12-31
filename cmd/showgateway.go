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

var showGatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Show a gateway",
	Long: AddAppName(`Show a gateway
    Example:
    $AppName show gateway name --organization org --environment env`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Gateway Name needed")
		}
		name := args[0]

		logging.Info("Showing gateway %v in environment %v in organization %v\n", name, environment, organization)

		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}
		gateway, showGatewayError := core.GetGateway(namespacedEnvironment, name)
		if showGatewayError != nil {
			return showGatewayError
		}
		fmt.Printf("%v", gateway)
		return nil
	},
}

func init() {
	showCmd.AddCommand(showGatewayCmd)

	showGatewayCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the gateway")
	showGatewayCmd.MarkFlagRequired("organization")

	showGatewayCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the gateway")
	showGatewayCmd.MarkFlagRequired("environment")
}
