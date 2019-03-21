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
	contents, error := util.ReadFilesIndirectory("../resources/artifacts")

	assert.Nil(t, error)
	/* Example print
	  for file, content := range contents {
			fmt.Printf("File :%v size: %v\n", file, len(content))
		}
	*/
	assert.Equal(t, 1761, len(contents["../resources/artifacts/breeds/quantification.yml"]))
	assert.Equal(t, 194, len(contents["../resources/artifacts/workflows/quantification.yml"]))

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
}
