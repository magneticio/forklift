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
	yaml "gopkg.in/yaml.v2"
)

// showorganizationCmd represents the showorganization command
var showorganizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Show an organization configuration",
	Long: AddAppName(`Show an organization
    Example:
    $AppName show organization organization1`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Organization Name needed.")
		}
		name := args[0]

		if !util.ValidateName(name) {
			return errors.New("Organization name should be only lowercase alphanumerics")
		}

		namespaced := Config.Namespace + "-" + name
		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}
		configuration, showOrganizationError := core.ShowOrganization(namespaced)
		if showOrganizationError != nil {
			return showOrganizationError
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
	showCmd.AddCommand(showorganizationCmd)
}
