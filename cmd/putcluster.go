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

var natsChannelName string
var optimiserNatsChannelName string
var natsToken string

var putClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Put a cluster",
	Long: AddAppName(`Put a cluster
    Usage:
    $AppName put cluster <cluster_id> --nats-channel-name <nats_channel_name> --optimiser-nats-channel-name <optimiser_nats_channel_name> --nats-token <nats_token>`),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("Not enough arguments - cluster id needed")
		}
		clusterIDString := args[0]

		clusterID, err := strconv.ParseUint(clusterIDString, 10, 64)
		if err != nil {
			return fmt.Errorf("Cluster id '%s' must be a natural number", clusterIDString)
		}
		logging.Info("Puting cluster '%d'\n", clusterID)
		core, err := core.NewCore(Config)
		if err != nil {
			return err
		}

		err = core.PutReleaseAgentConfig(clusterID, natsChannelName, optimiserNatsChannelName, natsToken)
		if err != nil {
			return err
		}

		fmt.Printf("Cluster '%d' has been put\n", clusterID)

		return nil
	},
}

func init() {
	putCmd.AddCommand(putClusterCmd)

	putClusterCmd.Flags().StringVar(&natsChannelName, "nats-channel-name", "", "NATS channel name")
	putClusterCmd.MarkFlagRequired("nats-channel-name")

	putClusterCmd.Flags().StringVar(&optimiserNatsChannelName, "optimiser-nats-channel-name", "", "optimiser NATS channel name")
	putClusterCmd.MarkFlagRequired("optimiser-nats-channel-name")

	putClusterCmd.Flags().StringVar(&natsToken, "nats-token", "", "NATS token")
}
