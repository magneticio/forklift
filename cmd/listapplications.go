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

	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/logging"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var listApplicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "List existing applications",
	Long: AddAppName(`List existing applications
    Usage:
    $AppName list applications --cluster <cluster_id>`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.Info("Listing applications")
		core, err := core.NewCore(Config)
		if err != nil {
			return err
		}

		applications, err := core.ListApplications()
		if err != nil {
			return err
		}

		output, err := yaml.Marshal(applications)
		if err != nil {
			return err
		}

		fmt.Print(string(output))

		return nil
	},
}

func init() {
	listCmd.AddCommand(listApplicationsCmd)
}
