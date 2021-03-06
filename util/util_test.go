package util_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/magneticio/forklift/util"
	"github.com/stretchr/testify/assert"
)

//Utility test to check that the url query parameters are parsed correctly
func TestUrlParsing(t *testing.T) {

	u, _ := url.ParseRequestURI(strings.TrimPrefix(
		"jdbc:mysql://mysql.default.svc.cluster.local:3306/vamp-${parent}?test=1&test2=2&useSSL=false", "jdbc:"))

	assert.Equal(t, "test=1&test2=2&useSSL=false", u.Query().Encode())

	u, _ = url.ParseRequestURI(strings.TrimPrefix(
		"jdbc:mysql://mysql.default.svc.cluster.local:3306/vamp-${parent}", "jdbc:"))

	assert.Equal(t, "", u.Query().Encode())

}

func TestValidateName(t *testing.T) {

	assert.True(t, util.ValidateName("organization1212"))

	assert.False(t, util.ValidateName("organization-1212"))

	assert.False(t, util.ValidateName("ORGANIZATION1212"))

}

func TestReadFilesIndirectory(t *testing.T) {
	contents, error := util.ReadFilesIndirectory("resources")

	assert.Nil(t, error)
	assert.Equal(t, 4386, len(contents["resources/testpolicy1.json"]))
	assert.Equal(t, 1031, len(contents["resources/testpolicy2.json"]))
}

func TestUUID(t *testing.T) {
	uuid := util.UUID()
	assert.Equal(t, 36, len(uuid))
}

func TestEncodeString(t *testing.T) {

	result := util.EncodeString("password", "SHA-512", "d4f22852-e281-428f-8968-1265b1c5a1b0")

	expected := "2b6bd58fa3c7412421821e8a567f0fc958727090cfc08b0cc0f4349d642f30505997a9a199a3bfe87f3388d5d8828d5cc8094c383aebc5054e421b656f8515fa"

	assert.Equal(t, expected, result)

	expected2 := "cf72ed5638f7e2629c86e1552a3bd0b6c852d048"

	text2 := "class io.vamp.common.Namespace@vampio-testorg2-testenv"
	result2 := util.EncodeString(text2, "SHA-1", "v1")
	assert.Equal(t, expected2, result2)

	expected3 := "6d1339c7c7a1ac54246a57320bb1dd15176ce29"

	text3 := "class io.vamp.common.Namespace@vampio-organization-environment"
	result3 := util.EncodeString(text3, "SHA-1", "v1")
	assert.Equal(t, expected3, result3)
}

func TestConvertJSONToJSON(t *testing.T) {
	jsonWithEscapedChars := `{"value": "health \u003e= baselines.minHealth"}`
	prettyJSON, err := util.Convert("json", "json", jsonWithEscapedChars)

	assert.NoError(t, err)

	expected := "{\n    \"value\": \"health >= baselines.minHealth\"\n}"

	assert.Equal(t, expected, prettyJSON)
}
