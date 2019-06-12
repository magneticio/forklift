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

var addReleasepolicyCmd = &cobra.Command{
	Use:   "releasepolicy",
	Short: "Add a new releasepolicy",
	Long: AddAppName(`Add a new releasepolicy
    Example:
    $AppName add releasepolicy name --organization org --environment env --file ./releasepolicydefinition.json -i json`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Release Policy name needed.")
		}
		name := args[0]

		logging.Info("Adding new releasepolicy to environment %v in organization %v\n", organization, environment)
		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		releasepolicyBytes, readErr := util.UseSourceUrl(configPath) // just pass the file name
		if readErr != nil {
			return readErr
		}

		releasepolicyJSON, jsonError := util.Convert(configFileType, "json", releasepolicyBytes)
		if jsonError != nil {
			return jsonError
		}

		releasepolicyText := string(releasepolicyJSON)

		createreleasepolicyError := core.AddReleasePolicy(namespacedOrganization, namespacedEnvironment, name, releasepolicyText)
		if createreleasepolicyError != nil {
			return createreleasepolicyError
		}

		fmt.Printf("Release Policy is added\n")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addReleasepolicyCmd)

	addReleasepolicyCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the workflow")
	addReleasepolicyCmd.MarkFlagRequired("organization")

	addReleasepolicyCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the workflow")
	addReleasepolicyCmd.MarkFlagRequired("environment")

	addReleasepolicyCmd.Flags().StringVarP(&configPath, "file", "", "", "Release policy configuration file path")
	addReleasepolicyCmd.MarkFlagRequired("file")
	addReleasepolicyCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Release policy configuration file type yaml or json")

}
