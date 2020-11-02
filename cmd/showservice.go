// Copyright © 2020 Developer <developer@vamp.io>
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
	"strconv"

	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

var showServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Show existing service",
	Long: AddAppName(`Show existing service
    Usage:
    $AppName show service <service_id> --cluster <cluster_id> --application <application_id>`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Not enough arguments - service id needed")
		}
		serviceIDString := args[0]

		serviceID, err := strconv.ParseUint(serviceIDString, 10, 64)
		if err != nil {
			return fmt.Errorf("Service id '%s' must be a natural number", serviceIDString)
		}

		logging.Info("Showing service '%d'\n", serviceID)
		core, err := core.NewCore(Config)
		if err != nil {
			return err
		}

		serviceConfigText, err := core.GetServiceConfigText(serviceID, applicationID)
		if err != nil {
			return err
		}

		prettyServiceConfigText, err := util.Convert("json", "json", serviceConfigText)
		if err != nil {
			return err
		}

		fmt.Print(prettyServiceConfigText)

		return nil
	},
}

func init() {
	showCmd.AddCommand(showServiceCmd)

	showServiceCmd.Flags().Uint64VarP(&applicationID, "application", "a", 0, "ID of the application")
	showServiceCmd.MarkFlagRequired("application")
}
