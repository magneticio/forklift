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
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

// adminCmd represents the admin command
var addAdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Add a new admin",
	Long: AddAppName(`Add a new admin
    Example:
    $AppName add admin --organization org --configuration ./somepath.json`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		namespaced := Config.Namespace + "-" + organization
		fmt.Printf("name: %v , configPath: %v , configFileType %v\n", namespaced, configPath, configFileType)

		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}

		adminBytes, adminErr := util.UseSourceUrl(configPath) // just pass the file name
		if adminErr != nil {
			return adminErr
		}

		adminText := string(adminBytes)

		createAdminError := core.AddAdmin(namespaced, adminText)
		if createAdminError != nil {
			return createAdminError
		}
		fmt.Printf("Admin is added\n")

		return nil
	},
}

func init() {
	addCmd.AddCommand(addAdminCmd)

	addAdminCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the environment")
	addAdminCmd.MarkFlagRequired("organization")

	addAdminCmd.Flags().StringVarP(&configPath, "configuration", "", "", "Admin configuration file path")
	addAdminCmd.MarkFlagRequired("configuration")
	addAdminCmd.Flags().StringVarP(&configFileType, "input", "i", "yaml", "Admin configuration file type yaml or json")

}
