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

var upsertReleasePlanCmd = &cobra.Command{
	Use:   "releaseplan",
	Short: "Upsert a release plan",
	Long: AddAppName(`Upsert a release plan
    Example:
    $AppName upsert releaseplan name --application-id 10 --service-id 5 --file ./releaseplandefinition.json`),
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

		releasePlanBytes, readErr := util.UseSourceUrl(configPath)
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
	upsertCmd.AddCommand(upsertReleasePlanCmd)

	upsertReleasePlanCmd.Flags().Uint64VarP(&applicationID, "application-id", "", "", "ID of the application")
	upsertReleasePlanCmd.MarkFlagRequired("application-id")

	upsertReleasePlanCmd.Flags().Uint64VarP(&serviceID, "service-id", "", "", "ID of the service")
	upsertReleasePlanCmd.MarkFlagRequired("service-id")

	upsertReleasePlanCmd.Flags().StringVarP(&configPath, "file", "", "", "Release plan configuration file path")
	upsertReleasePlanCmd.MarkFlagRequired("file")
	upsertReleasePlanCmd.Flags().StringVarP(&configFileType, "input", "i", "json", "Release plan configuration file type yaml or json")
}
