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

var addReleasePlanCmd = &cobra.Command{
	Use:   "releaseplan",
	Short: "Add a new release plan",
	Long: AddAppName(`Add a new release plan
    Example:
    $AppName add releaseplan name --file ./releaseplandefinition.json -i json`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Release Plan name needed")
		}
		name := args[0]

		logging.Info("Adding new release plan: '%s'\n", name)

		core, coreError := core.NewCore(Config)
		if coreError != nil {
			return coreError
		}

		releasePlanBytes, readErr := util.UseSourceUrl(configPath) // just pass the file name
		if readErr != nil {
			return readErr
		}

		releasePlanJSON, jsonError := util.Convert(configFileType, "json", releasePlanBytes)
		if jsonError != nil {
			return jsonError
		}

		releasePlanText := string(releasePlanJSON)

		createReleasePlanError := core.AddReleasePlan(name, releasePlanText)
		if createReleasePlanError != nil {
			return createReleasePlanError
		}

		fmt.Printf("Release Plan '%s' is added\n", name)

		return nil
	},
}

func init() {
	addCmd.AddCommand(addReleasePlanCmd)

	addReleasePlanCmd.Flags().StringVarP(&configPath, "file", "", "", "Release policy configuration file path")
	addReleasePlanCmd.MarkFlagRequired("file")
	addReleasePlanCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Release policy configuration file type yaml or json")

}
