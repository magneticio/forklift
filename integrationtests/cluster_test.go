// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegrationClusterCommands(t *testing.T) {
	var natsChannelName = "nats-channel"
	var optimiserChannelName = "optimiser-channel"
	var natsToken = "nats-token"

	Convey("When executing put cluster command with all necessary flags", t, func() {
		command := fmt.Sprintf(
			"put cluster 1 --nats-channel-name %s --optimiser-nats-channel-name %s --nats-token %s",
			natsChannelName,
			optimiserChannelName,
			natsToken,
		)

		stdoutLines, err := runCommand(command)

		Convey("error should not be thrown", func() {
			So(err, ShouldBeNil)
		})

		Convey("response should contain information that cluster has been put", func() {
			So(stdoutLines[0], ShouldEqual, "Cluster '1' has been put")
		})

		Convey("cluster should be saved to Vault", func() {
			clusterConfig, _, err := readValueFromVault("/v1/secret/vamp/projects/1/clusters/1/release-agent-config")
			snapshot, _ := readSnapshot("./snapshots/clusterconfig.txt")
			So(err, ShouldBeNil)
			So(clusterConfig, ShouldEqual, snapshot)
		})

		Convey("and deleting it afterwards", func() {
			stdoutLines, err := runCommand("delete cluster 1")

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain information that cluster has been deleted", func() {
				So(stdoutLines[0], ShouldEqual, "Cluster '1' has been deleted")
			})

			Convey("cluster should be removed from Vault", func() {
				_, exists, _ := readValueFromVault("/v1/secret/vamp/projects/1/clusters/1/release-agent-config")
				So(exists, ShouldEqual, false)
			})
		})
	})

	Convey("When executing put cluster command without nats-channel-name flag", t, func() {
		command := fmt.Sprintf(
			"put cluster 1 --optimiser-nats-channel-name %s --nats-token %s",
			optimiserChannelName,
			natsToken,
		)

		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "nats-channel-name" not set`)
		})
	})

	Convey("When executing put cluster command without optimiser-nats-channel-name flag", t, func() {
		command := fmt.Sprintf(
			"put cluster 1 --nats-channel-name %s --nats-token %s",
			natsChannelName,
			natsToken,
		)

		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "optimiser-nats-channel-name" not set`)
		})
	})
}
