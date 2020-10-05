//+build integration

package integrationtests

import (
	"testing"

	"github.com/magneticio/forklift/cmd"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIntegrationVersionCommand(t *testing.T) {
	Convey("Given version command", t, func() {

		command := "forklift version"

		Convey("When executing it", func() {
			lines, _, err := executeCLI(command)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("version should be returned", func() {
				So(lines[0], ShouldEqual, "Version: "+cmd.Version)
			})
		})
	})
}
