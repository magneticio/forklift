// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegrationApplicationCommands(t *testing.T) {
	var applicationID = uint64(12345)
	var clusterID = uint64(4321)
	var clusterName = "test-cluster"
	var namespace = "test-namespace"
	var natsChannelName = "nats-channel"
	var optimiserChannelName = "optimiser-channel"
	var natsToken = "nats-token"

	Convey("When executing put application command with all necessary flags", t, func() {
		command := fmt.Sprintf(
			"put application %d --namespace %s --cluster %d",
			applicationID,
			namespace,
			clusterID,
		)

		Convey("and cluster is present in Vault", func() {
			putClusterCommand := fmt.Sprintf(
				"put cluster %d --name %s --nats-channel-name %s --optimiser-nats-channel-name %s --nats-token %s",
				clusterID,
				clusterName,
				natsChannelName,
				optimiserChannelName,
				natsToken,
			)
			_, err := runCommand(putClusterCommand)
			So(err, ShouldBeNil)

			stdoutLines, err := runCommand(command)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain information that application has been put", func() {
				So(stdoutLines[0], ShouldEqual, "Application '12345' has been put")
			})

			Convey("application should be saved to Vault", func() {
				clusterConfig, _, err := readValueFromVault("/v1/secret/vamp/projects/1/clusters/4321/release-agent-config")
				snapshot, _ := readSnapshot("./snapshots/applicationconfig.txt")
				So(err, ShouldBeNil)
				So(clusterConfig, ShouldEqual, snapshot)
			})

			Convey("and showing application", func() {
				showApplicationCommand := fmt.Sprintf(
					"show application %d --cluster %d",
					applicationID,
					clusterID,
				)
				stdoutLines, err := runCommand(showApplicationCommand)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})

				Convey("response should contain application definition", func() {
					So(stdoutLines[0], ShouldEqual, "id: 12345")
					So(stdoutLines[1], ShouldEqual, "namespace: test-namespace")
				})
			})

			Convey("and listing applications", func() {
				listApplicationsCommand := fmt.Sprintf(
					"list applications --cluster %d",
					clusterID,
				)
				stdoutLines, err := runCommand(listApplicationsCommand)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})

				Convey("response should contain list of applications", func() {
					So(stdoutLines[0], ShouldEqual, "- id: 12345")
					So(stdoutLines[1], ShouldEqual, "  namespace: test-namespace")
				})
			})

			Convey("and deleting it afterwards", func() {
				deleteApplicationCommand := fmt.Sprintf(
					"delete application %d --cluster %d",
					applicationID,
					clusterID,
				)
				stdoutLines, err := runCommand(deleteApplicationCommand)

				Convey("error should not be thrown", func() {
					So(err, ShouldBeNil)
				})

				Convey("response should contain information that application has been deleted", func() {
					So(stdoutLines[0], ShouldEqual, "Application '12345' has been deleted")
				})

				Convey("application should be removed from Vault", func() {
					clusterConfig, _, err := readValueFromVault("/v1/secret/vamp/projects/1/clusters/4321/release-agent-config")
					snapshot, _ := readSnapshot("./snapshots/clusterconfig.txt")
					So(err, ShouldBeNil)
					So(clusterConfig, ShouldEqual, snapshot)
				})
			})
		})
	})

	Convey("When executing put application command without cluster", t, func() {
		command := fmt.Sprintf(
			"put application %d --namespace %s",
			applicationID,
			namespace,
		)

		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `cluster id must be provided`)
		})
	})

	Convey("When executing put application command without namespace flag", t, func() {
		command := fmt.Sprintf(
			"put application %d",
			applicationID,
		)

		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "namespace" not set`)
		})
	})

	Convey("When executing put application command without application value", t, func() {
		command := fmt.Sprintf(
			"put application --namespace %s",
			namespace,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `Not enough arguments - application id needed`)
		})
	})
}
