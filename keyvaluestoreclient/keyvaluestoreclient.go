package keyvaluestoreclient

import (
	"errors"

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

func NewKeyValueStoreClient(config models.KeyValueStoreConfiguration) (KeyValueStoreClient, error) {
	if config.Type == "vault" {
		params := map[string]string{
			"cert":   config.Vault.ClientTlsCert,
			"key":    config.Vault.ClientTlsKey,
			"caCert": config.Vault.ServerTlsCert,
		}

		vaultKVclient, vaultKVclientError := kvstore.NewVaultKeyValueStore(config.Vault.Url, config.Vault.Token, params)
		if vaultKVclientError != nil {
			return nil, vaultKVclientError
		}

		return vaultKVclient, nil
	}
	return nil, errors.New("Unsupported Key Value Store Client: " + config.Type)
}
