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
	"github.com/spf13/cobra"
)

var deletePolicyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Delete existing policy",
	Long: AddAppName(`Delete existing policy
    Example:
    $AppName delete policy id`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Not enough arguments - policy id needed")
		}
		policyIDString := args[0]

		policyID, err := strconv.ParseUint(policyIDString, 10, 64)
		if err != nil {
			return fmt.Errorf("Policy id '%s' must be a natural number", policyIDString)
		}

		logging.Info("Deleting policy '%d'\n", policyID)

		core, err := core.NewCore(Config)
		if err != nil {
			return err
		}

		err = core.DeletePolicy(policyID)
		if err != nil {
			return err
		}

		fmt.Printf("Policy '%d' has been deleted\n", policyID)

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deletePolicyCmd)
}
