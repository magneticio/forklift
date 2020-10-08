// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegrationServiceCommands(t *testing.T) {
	var clusterID = uint64(889)
	var applicationID = uint64(112)
	var serviceID = uint64(4555)
	var validServicePath = "./resources/validservice.json"
	var invalidServicePath = "./resources/invalidservice.json"

	Convey("When executing put service command with valid service config", t, func() {
		command := fmt.Sprintf(
			"put service --cluster %d --file %s",
			clusterID,
			validServicePath,
		)
		stdoutLines, err := runCommand(command)

		Convey("error should not be thrown", func() {
			So(err, ShouldBeNil)
		})

		Convey("response should contain information that service has been put", func() {
			So(stdoutLines[0], ShouldEqual, "Service has been put")
		})

		Convey("service config should be saved to Vault", func() {
			vaultService, _, err := readValueFromVault("/v1/secret/vamp/projects/1/clusters/889/applications/112/services/4555")
			snapshot, _ := readSnapshot("./snapshots/service.txt")
			So(err, ShouldBeNil)
			So(vaultService, ShouldEqual, snapshot)
		})

		Convey("and showing service", func() {
			showServiceCommand := fmt.Sprintf(
				"show service %d --cluster %d --application %d",
				serviceID,
				clusterID,
				applicationID,
			)

			stdoutLines, err := runCommand(showServiceCommand)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain service", func() {
				snapshot, _ := readSnapshot("./snapshots/service_show.txt")
				So(toText(stdoutLines), ShouldEqual, snapshot)
			})
		})

		Convey("and listing services", func() {
			listServicesCommand := fmt.Sprintf(
				"list services --cluster %d --application %d",
				clusterID,
				applicationID,
			)
			stdoutLines, err := runCommand(listServicesCommand)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain services list", func() {
				So(stdoutLines[0], ShouldEqual, "- 4555")
			})
		})

		Convey("and deleting it afterwards", func() {
			deleteServiceCommand := fmt.Sprintf(
				"delete service %d --cluster %d --application %d",
				serviceID,
				clusterID,
				applicationID,
			)
			stdoutLines, err := runCommand(deleteServiceCommand)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain information that service has been deleted", func() {
				So(stdoutLines[0], ShouldEqual, "Service '4555' has been deleted")
			})

			Convey("service config should be removed from Vault", func() {
				_, exists, _ := readValueFromVault("/v1/secret/vamp/projects/1/clusters/889/applications/112/services/4555")
				So(exists, ShouldEqual, false)
			})
		})
	})

	Convey("When executing put service command with invalid service config", t, func() {
		command := fmt.Sprintf(
			"put service --cluster %d --file %s",
			clusterID,
			invalidServicePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, "service config validation failed: k8s_namespace field is required")
		})
	})

	Convey("When executing put service command without cluster", t, func() {
		command := fmt.Sprintf(
			"put service --file %s",
			validServicePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `cluster id must be provided`)
		})
	})

	Convey("When executing put service command without file flag", t, func() {
		command := fmt.Sprintf(
			"put service --cluster %d",
			clusterID,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "file" not set`)
		})
	})
}
