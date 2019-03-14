package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/magneticio/forklift/core"
	"github.com/magneticio/forklift/util"
	"github.com/stretchr/testify/assert"
)

func TestOperationTemplateModel(t *testing.T) {

	fmt.Println("Reading file")

	configPath := "../resources/test/operation-template.yaml"

	configBtye, readErr := util.UseSourceUrl(configPath)
	if readErr != nil {
		fmt.Printf("Error while reading file")
		panic(readErr)
	}

	configText := string(configBtye)
	configJson, convertErr := util.Convert("yaml", "json", configText)

	//fmt.Printf("Config json %v", configJson)
	if convertErr != nil {
		fmt.Printf("Error while converting file")
		panic(convertErr)
	}
	var vampConfig core.Configuration
	unmarshallError := json.Unmarshal([]byte(configJson), &vampConfig)
	if unmarshallError != nil {
		panic(unmarshallError)
	}
	//fmt.Printf("Config json %v", vampConfig)

	configJson2, marshalError := json.Marshal(vampConfig)
	if marshalError != nil {
		panic(marshalError)
	}

	configText2 := string(configJson2)

	configYaml, convertErr := util.Convert("json", "yaml", configText2)

	result := string(configYaml)
	assert.Equal(t, configText, result)
}

func TestAdminTemplateModel(t *testing.T) {

	fmt.Println("Reading file")

	configPath := "../resources/test/admin-template.yaml"

	configBtye, readErr := util.UseSourceUrl(configPath)
	if readErr != nil {
		fmt.Printf("Error while reading file")
		panic(readErr)
	}

	configText := string(configBtye)
	configJson, convertErr := util.Convert("yaml", "json", configText)

	//fmt.Printf("Config json %v", configJson)
	if convertErr != nil {
		fmt.Printf("Error while converting file")
		panic(convertErr)
	}
	var vampConfig core.Configuration
	unmarshallError := json.Unmarshal([]byte(configJson), &vampConfig)
	if unmarshallError != nil {
		panic(unmarshallError)
	}
	//fmt.Printf("Config json %v", vampConfig)

	configJson2, marshalError := json.Marshal(vampConfig)
	if marshalError != nil {
		panic(marshalError)
	}

	configText2 := string(configJson2)

	configYaml, convertErr := util.Convert("json", "yaml", configText2)

	result := string(configYaml)
	assert.Equal(t, configText, result)
}
