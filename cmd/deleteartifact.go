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
	"github.com/spf13/cobra"
)

var deleteArtifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Delete an artifact",
	Long: AddAppName(`Delete an artifact
    Example:
    $AppName delete artifact name --kind k --organization org --environment env`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not Enough Arguments, Name needed.")
		}
		name := args[0]

		namespacedEnvironment := Config.Namespace + "-" + organization + "-" + environment
		namespacedOrganization := Config.Namespace + "-" + organization

		coreConfig := core.Configuration{
			VampConfiguration: Config.VampConfiguration,
		}
		core, coreError := core.NewCore(coreConfig)
		if coreError != nil {
			return coreError
		}

		deleteArtifactError := core.DeleteArtifact(namespacedOrganization, namespacedEnvironment, name, kind)
		if deleteArtifactError != nil {
			return deleteArtifactError
		}
		fmt.Printf("Artifact is deleted\n")

		return nil
	},
}

func init() {
	deleteCmd.AddCommand(deleteArtifactCmd)

	deleteArtifactCmd.Flags().StringVarP(&kind, "kind", "", "", "Kind of the artifact")
	deleteArtifactCmd.MarkFlagRequired("kind")

	deleteArtifactCmd.Flags().StringVarP(&organization, "organization", "", "", "Organization of the artifact")
	deleteArtifactCmd.MarkFlagRequired("organization")

	deleteArtifactCmd.Flags().StringVarP(&environment, "environment", "", "", "Environment of the artifact")
	deleteArtifactCmd.MarkFlagRequired("environment")

}
