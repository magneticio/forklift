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

var addGatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Add a new gateway",
	Long: AddAppName(`Add a new gateway
    Example:
    $AppName add gateway name --organization org --environment env --file ./gatewaydefinition.yaml`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Gateway name needed")
		}
		name := args[0]

		logging.Info("Adding new gateway to environment %v in organization %v\n", organization, environment)

		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		gatewayBytes, readErr := util.UseSourceUrl(configPath)
		if readErr != nil {
			return readErr
		}

		gatewayJSON, jsonError := util.Convert(configFileType, "json", gatewayBytes)
		if jsonError != nil {
			return jsonError
		}

		gatewayText := string(gatewayJSON)

		createGatewayError := core.AddGateway(namespacedEnvironment, name, gatewayText)
		if createGatewayError != nil {
			return createGatewayError
		}
		fmt.Printf("Gateway is added\n")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addGatewayCmd)

	addGatewayCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the workflow")
	addGatewayCmd.MarkFlagRequired("organization")

	addGatewayCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the workflow")
	addGatewayCmd.MarkFlagRequired("environment")

	addGatewayCmd.Flags().StringVarP(&configPath, "file", "", "", "Gateway configuration file path")
	addGatewayCmd.MarkFlagRequired("file")
	addGatewayCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Gateway configuration file type yaml or json")

}
