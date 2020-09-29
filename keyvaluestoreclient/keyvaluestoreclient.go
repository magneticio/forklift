package keyvaluestoreclient

import (
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/vamp-sdk-go/kvstore"
)

type KeyValueStoreClient interface {
	Get(string) (string, error)
	Exists(key string) (bool, error)
	Put(string, string) error
	Delete(string) error
	List(string) ([]string, error)
}

func NewKeyValueStoreClient(config models.VaultKeyValueStoreConfiguration) (KeyValueStoreClient, error) {
	params := map[string]string{
		"cert":   config.ClientTlsCert,
		"key":    config.ClientTlsKey,
		"caCert": config.ServerTlsCert,
	}

	vaultKVclient, vaultKVclientError := kvstore.NewVaultKeyValueStore(config.Url, config.Token, params)
	if vaultKVclientError != nil {
		return nil, vaultKVclientError
	}

	return vaultKVclient, nil
}
