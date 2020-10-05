package integrationtests

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/magneticio/forklift/cmd"
)

type restorer func()

// Execute cli and return string lists of standard out, standard error,
// and error if any.
func executeCLI(commandText string) ([]string, []string, error) {
	stdoutText, stderrText, err := executeCLIRaw(commandText)

	return strings.Split(stdoutText, "\n"),
		strings.Split(stderrText, "\n"),
		err
}

// Execute cli and return the strings of standard out, standard error,
// and error.
func executeCLIRaw(commandText string) (string, string, error) {
	oldargs := os.Args
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	os.Args = strings.Split(commandText, " ")
	stdoutR, stdoutW, _ := os.Pipe()
	stderrR, stderrW, _ := os.Pipe()
	os.Stdout = stdoutW
	os.Stderr = stderrW

	err := runForklift()

	stdoutW.Close()
	stderrW.Close()

	stdoutBytes, _ := ioutil.ReadAll(stdoutR)
	stderrBytes, _ := ioutil.ReadAll(stderrR)

	os.Stdout = oldStdout
	os.Stderr = oldStderr
	os.Args = oldargs

	return string(stdoutBytes), string(stderrBytes), err
}

// runForklift runs a command saved in os.Args and returns the error if any
func runForklift() error {
	cmd.Execute()
	return nil
}
