// +build integration

package integrationtests

import (
	"testing"

	"github.com/magneticio/forklift/cmd"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegrationVersionCommand(t *testing.T) {
	Convey("Given version command", t, func() {

		Convey("When executing it", func() {
			stdoutLines, err := runCommand("version")

			Convey("version should be returned", func() {
				So(stdoutLines[0], ShouldEqual, "Version: "+cmd.Version)
			})

			Convey("error should not be returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
