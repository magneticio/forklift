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
	"github.com/spf13/cobra"
)

var deleteReleasePlanCmd = &cobra.Command{
	Use:   "releaseplan",
	Short: "Delete existing release plan",
	Long: AddAppName(`Delete existing release plan
    Example:
    $AppName delete releaseplan name`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not enough arguments - release plan name needed")
		}
		serviceVersion := args[0]

		logging.Info("Deleteing release plan for service version '%s'\n", serviceVersion)

		core, err := core.NewCore(Config)
		if err != nil {
			return err
		}

		err = core.DeleteReleasePlan(serviceID, serviceVersion)
		if err != nil {
			return err
		}

		fmt.Printf("Release plan for service version '%s' has been deleted\n", serviceVersion)

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteReleasePlanCmd)

	deleteReleasePlanCmd.Flags().Uint64VarP(&serviceID, "service-id", "", 0, "ID of the service")
	deleteReleasePlanCmd.MarkFlagRequired("service-id")
}
