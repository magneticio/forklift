// +build integration

package integrationtests

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	validPolicyFilePath   = "./resources/validpolicy.json"
	invalidPolicyFilePath = "./resources/invalidpolicy.json"
	policyID              = uint64(456)
)

func TestIntegrationPolicyCommands(t *testing.T) {
	Convey("When executing put policy command with valid policy", t, func() {
		command := fmt.Sprintf(
			"put policy %d --file %s",
			policyID,
			validPolicyFilePath,
		)
		stdoutLines, err := runCommand(command)

		Convey("error should not be thrown", func() {
			So(err, ShouldBeNil)
		})

		Convey("response should contain information that policy has been put", func() {
			So(stdoutLines[0], ShouldEqual, "Policy '456' has been put")
		})

		Convey("policy should be saved to Vault", func() {
			vaultPolicy, _, err := readValueFromVault("/v1/secret/data/vamp/projects/1/policies/456")
			snapshot, _ := readSnapshot("./snapshots/policy.txt")
			So(err, ShouldBeNil)
			So(vaultPolicy, ShouldEqual, snapshot)
		})

		Convey("and deleting it afterwards", func() {
			deletePolicyCommand := fmt.Sprintf(
				"delete policy %d",
				policyID,
			)
			stdoutLines, err := runCommand(deletePolicyCommand)

			Convey("error should not be thrown", func() {
				So(err, ShouldBeNil)
			})

			Convey("response should contain information that policy has been deleted", func() {
				So(stdoutLines[0], ShouldEqual, "Policy '456' has been deleted")
			})

			Convey("policy should be removed from Vault", func() {
				_, exists, _ := readValueFromVault("/v1/secret/data/vamp/projects/1/policies/456")
				So(exists, ShouldEqual, false)
			})
		})
	})

	Convey("When executing put policy command with invalid policy", t, func() {
		command := fmt.Sprintf(
			"put policy %d --file %s",
			policyID,
			invalidPolicyFilePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, "Cannot save policy: Found problems with validation policy: Invalid policy type: Unknown policy type: unknown-type")
		})
	})

	Convey("When executing put policy command without policy value", t, func() {
		command := fmt.Sprintf(
			"put policy --file %s",
			validPolicyFilePath,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `Not enough arguments - policy id needed`)
		})
	})

	Convey("When executing put policy command without file flag", t, func() {
		command := fmt.Sprintf(
			"put policy %d",
			policyID,
		)
		_, err := runCommand(command)

		Convey("error should be thrown", func() {
			So(err.Error(), ShouldEqual, `required flag(s) "file" not set`)
		})
	})
}
