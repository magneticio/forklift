package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"unsafe"

	"github.com/hashicorp/vault/api"
	"github.com/magneticio/forklift/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	vaultAddress = "http://vault:8200"
	vaultToken   = "vamp"
)

// setEnvVariables - set environment variables needed by Forklift tests
func setEnvVariables() {
	os.Setenv("VAMP_FORKLIFT_VAULT_ADDR", vaultAddress)
	os.Setenv("VAMP_FORKLIFT_VAULT_TOKEN", vaultToken)
}

// deleteEmpty - deletes empty strings from an array of strings
func deleteEmpty(s []string) []string {
	var ret []string
	for _, str := range s {
		if str != "" {
			ret = append(ret, str)
		}
	}
	return ret
}

// getTestVaultClient - returns test Vault client
func getTestVaultClient() *api.Client {
	config := &api.Config{
		Address: vaultAddress,
	}
	client, _ := api.NewClient(config)
	client.SetToken(vaultToken)

	return client
}

// readValueFromVault - reads value from Vault based on provided path
func readValueFromVault(path string) (string, bool, error) {
	vaultClient := getTestVaultClient()
	req := vaultClient.NewRequest("GET", path)
	resp, err := vaultClient.RawRequest(req)

	if resp.StatusCode == 404 {
		return "", false, nil
	}

	if err != nil {
		return "", false, err
	}

	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)

	var respJSON map[string]interface{}

	err = json.Unmarshal(buf.Bytes(), &respJSON)
	if err != nil {
		return "", true, fmt.Errorf("cannot unmarshal Vault response: %v", err)
	}

	data, ok := respJSON["data"]
	if !ok {
		return "", true, fmt.Errorf("data not found in Vault response")
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return "", true, fmt.Errorf("data in Vault response is of invalid type")
	}

	dataData, ok := dataMap["data"]
	if !ok {
		return "", true, fmt.Errorf("data not found in Vault response")
	}

	dataDataMap, ok := dataData.(map[string]interface{})
	if !ok {
		return "", true, fmt.Errorf("data in Vault response is of invalid type")
	}

	value, ok := dataDataMap["value"]
	if !ok {
		return "", true, fmt.Errorf("value not found in Vault response")
	}

	valueText, err := json.Marshal(value)
	if err != nil {
		return "", true, fmt.Errorf("cannot get value text: %v", err)
	}

	return string(valueText), true, nil
}

// readSnapshot - reads test snapshot
func readSnapshot(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read snapshot: %v", err)
	}

	return string(data), nil
}

// runCommand - runs Forklift command
func runCommand(commandText string) ([]string, error) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	command := cmd.RootCmd()
	defer resetFlags(command)
	defer func() {
		os.Stdout = oldStdout
	}()

	args := strings.Split(commandText, " ")
	command.SetArgs(args)
	setEnvVariables()

	err := command.Execute()
	if err != nil {
		return nil, err
	}

	w.Close()
	stdoutBytes, _ := ioutil.ReadAll(r)
	stdoutLines := strings.Split(string(stdoutBytes), "\n")

	return deleteEmpty(stdoutLines), nil
}

// resetFlags - resets Cobra flags to initial values
// if it's not executed, errors for missing mandatory flags are not being thrown
func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Value.Type() == "stringSlice" {
			// XXX: unfortunately, flag.Value.Set() appends to original
			// slice, not resets it, so we retrieve pointer to the slice here
			// and set it to new empty slice manually
			value := reflect.ValueOf(flag.Value).Elem().FieldByName("value")
			ptr := (*[]string)(unsafe.Pointer(value.Pointer()))
			*ptr = make([]string, 0)
		}
		flag.Value.Set(flag.DefValue)
		flag.Changed = false
	})
	for _, cmd := range cmd.Commands() {
		resetFlags(cmd)
	}
}
