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
	assert.Equal(t, 1762, len(contents["../resources/artifacts/breeds/quantification.yml"]))
	assert.Equal(t, 395, len(contents["../resources/artifacts/workflows/quantification.yml"]))

}

func TestUUID(t *testing.T) {
	uuid := util.UUID()
	assert.Equal(t, 36, len(uuid))
}
