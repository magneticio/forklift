// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	releasePlanPath = "./resources/releaseplan.json"
	serviceID       = uint64(56789)
	serviceVersion  = "1.0.5"
)

func TestIntegrationReleasePlanCommands(t *testing.T) {
	Convey("When executing put release plan command with all necessary flags", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --service %d --file %s",
			serviceVersion,
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
			vaultPolicy, _, err := readValueFromVault("/v1/secret/data/vamp/projects/1/release-plans/56789/1.0.1")
			snapshot, _ := readSnapshot("./snapshots/releaseplan.txt")
			So(err, ShouldBeNil)
			So(vaultPolicy, ShouldEqual, snapshot)
		})

		Convey("and deleting it afterwards", func() {
			deleteReleasePlanCommand := fmt.Sprintf(
				"delete releaseplan %s --service %d",
				serviceVersion,
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
				_, exists, _ := readValueFromVault("/v1/secret/data/vamp/projects/1/release-plans/56789/1.0.1")
				So(exists, ShouldEqual, false)
			})
		})
	})

	Convey("When executing put release plan command without service version value", t, func() {
		command := fmt.Sprintf(
			"put releaseplan --service %d --file %s",
			serviceID,
			releasePlanPath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `Not enough arguments, service version needed`)
		})
	})

	Convey("When executing put release plan command without service flag", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --file %s",
			serviceVersion,
			releasePlanPath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "service" not set`)
		})
	})

	Convey("When executing put release plan command without file flag", t, func() {
		command := fmt.Sprintf(
			"put releaseplan %s --service %d",
			serviceVersion,
			serviceID,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "file" not set`)
		})
	})
}
