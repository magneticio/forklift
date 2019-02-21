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

	hocon "github.com/go-akka/configuration"
	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

var configPath string
var configFileType string

// organizationCmd represents the organization command
var organizationCmd = &cobra.Command{
	Use:   "organization",
	Short: "Create a new organization",
	Long: AddAppName(`Create a new organization
    Example:
    $AppName create organization org1 --configuration ./somepath.json`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Organization Name needed.")
		}
		name := args[0]
		fmt.Printf("name: %v , configPath: %v , configFileType %v\n", name, configPath, configFileType)

		configBtye, readErr := util.UseSourceUrl(configPath) // just pass the file name
		if readErr != nil {
			return readErr
		}
		configText := string(configBtye)
		conf := hocon.ParseString(configText)

		basePath := "vamp.persistence.database.sql."
		sqlConfiguration := core.SqlConfiguration{
			Database:          conf.GetString(basePath + "database"),
			Table:             conf.GetString(basePath + "table"),
			User:              conf.GetString(basePath + "user"),
			Password:          conf.GetString(basePath + "password"),
			Url:               conf.GetString(basePath + "url"),
			DatabaseServerUrl: conf.GetString(basePath + "database-server-url"),
		}
		config := core.Configuration{
			Sql: sqlConfiguration,
		}
		core, coreError := core.NewCore(config)
		if coreError != nil {
			return coreError
		}
		createOrganizationError := core.CreateOrganization(name)
		if createOrganizationError != nil {
			return createOrganizationError
		}
		fmt.Printf("Organization %v is created\n", name)
		return nil
	},
}

func init() {
	createCmd.AddCommand(organizationCmd)

	organizationCmd.Flags().StringVarP(&configPath, "configuration", "", "", "Organization configuration file path")
	organizationCmd.MarkFlagRequired("configuration")
	organizationCmd.Flags().StringVarP(&configFileType, "input", "i", "hocon", "Configuration file type yaml or json, hocon")
}
