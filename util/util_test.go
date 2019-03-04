package util_test

import (
	"testing"

	"github.com/magneticio/forklift/util"
	"github.com/stretchr/testify/assert"
)

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
}
