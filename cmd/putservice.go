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
	"github.com/magneticio/forklift/util"
	"github.com/spf13/cobra"
)

var putServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Put a service",
	Long: AddAppName(`Put a service
    Usage:
    $AppName put service --cluster <cluster_id> --file <service_config_file_path>`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.Info("Putting service\n")
		core, err := core.NewCore(Config)
		if err != nil {
			return err
		}

		serviceConfigBytes, err := util.UseSourceUrl(configPath)
		if err != nil {
			return err
		}

		serviceConfigJSON, err := util.Convert(configFileType, "json", serviceConfigBytes)
		if err != nil {
			return err
		}

		serviceConfigText := string(serviceConfigJSON)

		err = core.PutServiceConfig(serviceConfigText)
		if err != nil {
			return err
		}

		fmt.Printf("Service has been put\n")

		return nil
	},
}

func init() {
	putCmd.AddCommand(putServiceCmd)

	putServiceCmd.Flags().StringVarP(&configPath, "file", "", "", "Service configuration file path")
	putServiceCmd.MarkFlagRequired("file")
	putServiceCmd.Flags().StringVarP(&configFileType, "input", "i", "json", "Service configuration file type yaml or json")
}
