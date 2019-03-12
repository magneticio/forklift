package keyvaluestoreclient_test

import (
	"fmt"
	"testing"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/stretchr/testify/assert"
)

func TestVaultKeyVauleStoreClientKV2(t *testing.T) {

	address := "http://127.0.0.1:8200"
	token := "myroot"
	params := map[string]string{
		"cert":   "",
		"key":    "",
		"caCert": "",
	}

	vaultKeyValueStoreClient, clientErr := keyvaluestoreclient.NewVaultKeyValueStoreClient(address, token, params)
	assert.Nil(t, clientErr)
	assert.NotNil(t, vaultKeyValueStoreClient)

	key := "secret/foo"
	valueExpected := map[string]interface{}(map[string]interface{}{"bar": "baz"})
	valueActual, getErr := vaultKeyValueStoreClient.GetData(key, 2)
	fmt.Printf("valueActual %v\n", valueActual)
	assert.Nil(t, getErr)
	assert.Equal(t, valueExpected, valueActual)

}
