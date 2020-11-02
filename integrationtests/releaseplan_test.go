// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegrationReleasePlanCommands(t *testing.T) {
	var releasePlanPath = "./resources/releaseplan.json"
	var clusterID = uint64(7777)
	var applicationID = uint64(9999)
	var serviceID = uint64(56789)
	var serviceVersion = "1.0.5"

	Convey("When executing put release plan command with all necessary flags", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --cluster %d --application %d --service %d --file %s",
			serviceVersion,
			clusterID,
			applicationID,
			serviceID,
			releasePlanPath,
		)
		stdoutLines, err := runCommand(command)

		Convey("error should not be thrown", func() {
			So(err, ShouldBeNil)
		})

		Convey("response should contain information that release plan has been put", func() {
			So(stdoutLines[0], ShouldEqual, "Release plan for service version '1.0.5' has been put")
		})

		Convey("release plan should be saved to Vault", func() {
			vaultPolicy, _, err := readValueFromVault("/v1/secret/vamp/projects/1/clusters/7777/applications/9999/release-plans/56789/1.0.1")
			snapshot, _ := readSnapshot("./snapshots/releaseplan.txt")
			So(err, ShouldBeNil)
			So(vaultPolicy, ShouldEqual, snapshot)
		})

		Convey("and showing release plan", func() {
			showReleasePlanCommand := fmt.Sprintf(
				"show releaseplan %s --cluster %d --application %d --service %d",
				serviceVersion,
				clusterID,
				applicationID,
				serviceID,
			)

			stdoutLines, err := runCommand(showReleasePlanCommand)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain release plan", func() {
				snapshot, _ := readSnapshot("./snapshots/releaseplan_show.txt")
				So(toText(stdoutLines), ShouldEqual, snapshot)
			})
		})

		Convey("and listing release plans", func() {
			listReleasePlansCommand := fmt.Sprintf(
				"list releaseplans --cluster %d --application %d --service %d",
				clusterID,
				applicationID,
				serviceID,
			)
			stdoutLines, err := runCommand(listReleasePlansCommand)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain release plans list", func() {
				So(stdoutLines[0], ShouldEqual, "- 1.0.5")
			})
		})

		Convey("and deleting it afterwards", func() {
			deleteReleasePlanCommand := fmt.Sprintf(
				"delete releaseplan %s --cluster %d --application %d --service %d",
				serviceVersion,
				clusterID,
				applicationID,
				serviceID,
			)
			stdoutLines, err := runCommand(deleteReleasePlanCommand)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain information that release plan has been deleted", func() {
				So(stdoutLines[0], ShouldEqual, "Release plan for service version '1.0.5' has been deleted")
			})

			Convey("release plan should be removed from Vault", func() {
				_, exists, _ := readValueFromVault("/v1/secret/vamp/projects/1/clusters/7777/applications/9999/release-plans/56789/1.0.1")
				So(exists, ShouldEqual, false)
			})
		})
	})

	Convey("When executing put release plan command without service version value", t, func() {
		command := fmt.Sprintf(
			"put releaseplan --cluster %d --application %d --service %d --file %s",
			clusterID,
			applicationID,
			serviceID,
			releasePlanPath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `Not enough arguments, service version needed`)
		})
	})

	Convey("When executing put release plan command without cluster flag", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --application %d --service %d --file %s",
			serviceVersion,
			applicationID,
			serviceID,
			releasePlanPath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `cluster id must be provided`)
		})
	})

	Convey("When executing put release plan command without application flag", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --cluster %d --service %d --file %s",
			serviceVersion,
			clusterID,
			serviceID,
			releasePlanPath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "application" not set`)
		})
	})

	Convey("When executing put release plan command without service flag", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --cluster %d --application %d --file %s",
			serviceVersion,
			clusterID,
			applicationID,
			releasePlanPath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "service" not set`)
		})
	})

	Convey("When executing put release plan command without file flag", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --cluster %d --application %d --service %d",
			serviceVersion,
			clusterID,
			applicationID,
			serviceID,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "file" not set`)
		})
	})
}
