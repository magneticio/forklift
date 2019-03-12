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
	yaml "gopkg.in/yaml.v2"
)

// showenvironmentCmd represents the showenvironment command
var showenvironmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Show an environment configuration",
	Long: AddAppName(`Show an environment
    Example:
    $AppName show environment environment1`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Environment Name needed.")
		}
		name := args[0]

		logging.Info("Showing environment %v in organization %v\n", name, organization)

		if !util.ValidateName(name) {
			return errors.New("Environment name should be only lowercase alphanumerics")
		}

		namespaced := Config.Namespace + "-" + organization + "-" + name
		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}
		configuration, showEnvironmentError := core.ShowEnvironment(namespaced)
		if showEnvironmentError != nil {
			return showEnvironmentError
		}
		configYaml, marshalError := yaml.Marshal(configuration)
		if marshalError != nil {
			return marshalError
		}
		fmt.Printf("%v", string(configYaml))
		return nil
	},
}

func init() {
	showCmd.AddCommand(showenvironmentCmd)

	showenvironmentCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	showenvironmentCmd.MarkFlagRequired("organization")
}
