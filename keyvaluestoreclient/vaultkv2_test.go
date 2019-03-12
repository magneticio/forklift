package keyvaluestoreclient_test

import (
	"fmt"
	"testing"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/stretchr/testify/assert"
)

func valueMap(value string) map[string]interface{} {
	return map[string]interface{}{
		"value": value,
	}
}

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

	key := "secret/vamp"
	valueExpected := map[string]interface{}(map[string]interface{}{"value": "kv2test"})

	putErr := vaultKeyValueStoreClient.PutData(key, valueExpected, -1) // 0 means no cas
	assert.Nil(t, putErr)

	valueActual, getErr := vaultKeyValueStoreClient.GetData(key, 0) // 0 means latest version
	fmt.Printf("valueActual %v\n", valueActual)
	assert.Nil(t, getErr)
	assert.Equal(t, valueExpected, valueActual)

	vaultKeyValueStoreClient.PutData(key+"/key1", valueMap("value1"), -1)
	vaultKeyValueStoreClient.PutData(key+"/key2", valueMap("value2"), -1)
	vaultKeyValueStoreClient.PutData(key+"/key3", valueMap("value3"), -1)

	keys := []interface{}{"key1", "key2", "key3"}
	expectedListData := map[string]interface{}{
		"keys": keys,
	}
	actualListData, listDataError := vaultKeyValueStoreClient.ListData(key)
	assert.Nil(t, listDataError)
	assert.Equal(t, expectedListData, actualListData)
}
