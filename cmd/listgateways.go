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
var listGatewayCmd = &cobra.Command{
	Use:   "gateways",
	Short: "list gateways in an environment",
	Long: AddAppName(`List gateways
    Example:
    $AppName list gateways --kind k --organization org --environment env`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		logging.Info("Listing %v in environment %v in organization %v\n", kind, environment, organization)

		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}
		gatewaysList, listGatewaysError := core.ListGateways(namespacedEnvironment)
		if listGatewaysError != nil {
			return listGatewaysError
		}
		output, marshallError := yaml.Marshal(gatewaysList)
		if marshallError != nil {
			return marshallError
		}
		fmt.Print(string(output))
		return nil
	},
}

func init() {
	listCmd.AddCommand(listGatewayCmd)

	listGatewayCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the gateways")
	listGatewayCmd.MarkFlagRequired("organization")

	listGatewayCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the gateways")
	listGatewayCmd.MarkFlagRequired("environment")
}
