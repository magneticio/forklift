package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/util"
	"github.com/stretchr/testify/assert"
)

func TestOperationTemplateModel(t *testing.T) {

	configPath := "../resources/test/operation-template.yaml"

	configBtye, readErr := util.UseSourceUrl(configPath)
	if readErr != nil {
		fmt.Printf("Error while reading file")
		panic(readErr)
	}

	configText := string(configBtye)
	configJson, convertErr := util.Convert("yaml", "json", configText)
	if convertErr != nil {
		fmt.Printf("Error while converting file")
		panic(convertErr)
	}
	var vampConfig models.VampConfiguration
	unmarshallError := json.Unmarshal([]byte(configJson), &vampConfig)
	if unmarshallError != nil {
		panic(unmarshallError)
	}

	remarshalledJson, marshalError := json.Marshal(vampConfig)
	if marshalError != nil {
		panic(marshalError)
	}

	remarshalledText := string(remarshalledJson)

	remarshalledYaml, convertErr := util.Convert("json", "yaml", remarshalledText)

	result := string(remarshalledYaml)
	assert.Equal(t, result, configText)

}

func TestAdminTemplateModel(t *testing.T) {

	configPath := "../resources/test/admin-template.yaml"

	configBtye, readErr := util.UseSourceUrl(configPath)
	if readErr != nil {
		fmt.Printf("Error while reading file")
		panic(readErr)
	}

	configText := string(configBtye)
	configJson, convertErr := util.Convert("yaml", "json", configText)
	if convertErr != nil {
		fmt.Printf("Error while converting file")
		panic(convertErr)
	}
	var vampConfig models.VampConfiguration
	unmarshallError := json.Unmarshal([]byte(configJson), &vampConfig)
	if unmarshallError != nil {
		panic(unmarshallError)
	}

	expectedFallbackVersion := "${env://VAMP_PERSISTENCE_KEY_VALUE_STORE_VAULT_FALLBACK_KV_VERSION}"
	assert.Equal(t, vampConfig.Vamp.Persistence.KeyValueStore.Vault.FallbackKvVersion, expectedFallbackVersion)

	remarshalledJson, marshalError := json.Marshal(vampConfig)
	if marshalError != nil {
		panic(marshalError)
	}

	remarshalledText := string(remarshalledJson)

	remarshalledYaml, convertErr := util.Convert("json", "yaml", remarshalledText)
	if convertErr != nil {
		panic(convertErr)
	}

	result := string(remarshalledYaml)
	assert.Equal(t, configText, result)
}
