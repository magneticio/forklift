// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	validServicePath   = "./resources/validservice.json"
	invalidServicePath = "./resources/invalidservice.json"
	serviceID          = uint64(4555)
	clusterID          = uint64(889)
	applicationID      = uint64(112)
)

func TestIntegrationServiceCommands(t *testing.T) {
	Convey("When executing put service command with valid service config", t, func() {
		command := fmt.Sprintf(
			"put service %d --cluster %d --application %d --file %s",
			serviceID,
			clusterID,
			applicationID,
			validServicePath,
		)
		stdoutLines, err := runCommand(command)

		Convey("error should not be thrown", func() {
			So(err, ShouldBeNil)
		})

		Convey("response should contain information that service has been put", func() {
			So(stdoutLines[0], ShouldEqual, "Service '4555' has been put")
		})

		Convey("service config should be saved to Vault", func() {
			vaultService, _, err := readValueFromVault("/v1/secret/data/vamp/projects/1/clusters/889/applications/112/services/4555")
			snapshot, _ := readSnapshot("./snapshots/service.txt")
			So(err, ShouldBeNil)
			So(vaultService, ShouldEqual, snapshot)
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
				_, exists, _ := readValueFromVault("/v1/secret/data/vamp/projects/1/clusters/889/applications/112/services/4555")
				So(exists, ShouldEqual, false)
			})
		})
	})

	Convey("When executing put service command with invalid service config", t, func() {
		command := fmt.Sprintf(
			"put service %d --cluster %d --application %d --file %s",
			serviceID,
			clusterID,
			applicationID,
			invalidServicePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, "service config validation failed: k8s_namespace field is required")
		})
	})

	Convey("When executing put service command without service value", t, func() {
		command := fmt.Sprintf(
			"put service --cluster %d --application %d --file %s",
			clusterID,
			applicationID,
			validServicePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `Not enough arguments - service id needed`)
		})
	})

	Convey("When executing put service command without cluster", t, func() {
		command := fmt.Sprintf(
			"put service %d --application %d --file %s",
			serviceID,
			applicationID,
			validServicePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `cluster id must be provided`)
		})
	})

	Convey("When executing put service command without application flag", t, func() {
		command := fmt.Sprintf(
			"put service %d --cluster %d --file %s",
			serviceID,
			clusterID,
			validServicePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "application" not set`)
		})
	})

	Convey("When executing put service command without file flag", t, func() {
		command := fmt.Sprintf(
			"put service %d --cluster %d --application %d",
			serviceID,
			clusterID,
			applicationID,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "file" not set`)
		})
	})
}
