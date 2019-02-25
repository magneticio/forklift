package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/sql"
)

type Configuration struct {
	VampConfiguration models.VampConfiguration `yaml:"vamp,omitempty" json:"vamp,omitempty"`
}

type Core struct {
	Conf Configuration `yaml:"forklift,omitempty" json:"forklift,omitempty"`
}

func NewCore(conf Configuration) (*Core, error) {

	return &Core{
		Conf: conf,
	}, nil
}

func (c *Core) CreateOrganization(namespace string, configuration Configuration) error {

	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(namespace)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath + "/configuration/applied"
	value, jsonMarshallError := json.Marshal(configuration)
	if jsonMarshallError != nil {
		return jsonMarshallError
	}
	fmt.Printf("Vault store at key: %v\n", key)
	keyValueStoreClientPutError := keyValueStoreClient.PutValue(key, string(value))
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.SetupOrganization(databaseConfig.Sql.Database, databaseConfig.Sql.Table)

}

func (c *Core) GetNamespaceDatabaseConfiguration(namespace string) models.Database {
	databaseConf := c.Conf.VampConfiguration.Persistence.Database

	return models.Database{
		Sql: models.SqlConfiguration{
			Database:          Namespaced(namespace, databaseConf.Sql.Database),
			Table:             Namespaced(namespace, databaseConf.Sql.Table),
			User:              Namespaced(namespace, databaseConf.Sql.User),
			Password:          Namespaced(namespace, databaseConf.Sql.Password),
			Url:               Namespaced(namespace, databaseConf.Sql.Url),
			DatabaseServerUrl: Namespaced(namespace, databaseConf.Sql.DatabaseServerUrl),
		},
		Type: databaseConf.Type,
	}
}

func (c *Core) GetNamespaceKeyValueStoreConfiguration(namespace string) *models.KeyValueStoreConfiguration {
	return &models.KeyValueStoreConfiguration{
		Type:     c.Conf.VampConfiguration.Persistence.KeyValueStore.Type,
		BasePath: Namespaced(namespace, c.Conf.VampConfiguration.Persistence.KeyValueStore.BasePath),
		Vault: models.VaultKeyValueStoreConfiguration{
			Url:   Namespaced(namespace, c.Conf.VampConfiguration.Persistence.KeyValueStore.Vault.Url),
			Token: Namespaced(namespace, c.Conf.VampConfiguration.Persistence.KeyValueStore.Vault.Token),
		},
	}
}

func Namespaced(namespace string, text string) string {
	return strings.Replace(text, "${namespace}", namespace, -1)
}
