package core

import (
	"fmt"
	"strings"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/sql"
	"github.com/magneticio/forklift/util"
)

type Configuration struct {
	VampConfiguration models.VampConfiguration          `yaml:"vamp,omitempty" json:"vamp,omitempty"`
	Sql               models.SqlConfiguration           `yaml:"sql,omitempty" json:"sql,omitempty"`
	KeyValueStore     models.KeyValueStoreConfiguration `yaml:"key-value-store,omitempty" json:"key-value-store,omitempty"`
	Hocon             string                            `hocon:"sql,omitempty" json:"hocon,omitempty"`
}

type Core struct {
	Conf Configuration `yaml:"forklift,omitempty" json:"forklift,omitempty"`
}

func NewCore(conf Configuration) (*Core, error) {

	return &Core{
		Conf: conf,
	}, nil
}

func (c *Core) CreateOrganization(namespace string) error {

	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(namespace)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath
	value := c.Conf.Hocon
	keyValueStoreClientPutError := keyValueStoreClient.PutValue(key, value)
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}

	sqlConfig := c.GetNamespaceSqlConfiguration(namespace)

	host, hostError := util.GetHostFromUrl(sqlConfig.Url)
	if hostError != nil {
		fmt.Printf("Error: %v\n", hostError.Error())
		return hostError
	}

	client, clientError := sql.NewMySqlClient(sqlConfig.User, sqlConfig.Password, host, sqlConfig.Database)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.SetupOrganization(sqlConfig.Database, sqlConfig.Table)

}

func (c *Core) GetNamespaceSqlConfiguration(namespace string) *models.SqlConfiguration {
	return &models.SqlConfiguration{
		Database:          Namespaced(namespace, c.Conf.Sql.Database),
		Table:             Namespaced(namespace, c.Conf.Sql.Table),
		User:              Namespaced(namespace, c.Conf.Sql.User),
		Password:          Namespaced(namespace, c.Conf.Sql.Password),
		Url:               Namespaced(namespace, c.Conf.Sql.Url),
		DatabaseServerUrl: Namespaced(namespace, c.Conf.Sql.DatabaseServerUrl),
	}
}

func (c *Core) GetNamespaceKeyValueStoreConfiguration(namespace string) *models.KeyValueStoreConfiguration {
	return &models.KeyValueStoreConfiguration{
		Type:     c.Conf.KeyValueStore.Type,
		BasePath: Namespaced(namespace, c.Conf.KeyValueStore.BasePath),
		Vault: models.VaultKeyValueStoreConfiguration{
			Url:   Namespaced(namespace, c.Conf.KeyValueStore.Vault.Url),
			Token: Namespaced(namespace, c.Conf.KeyValueStore.Vault.Token),
		},
	}
}

func Namespaced(namespace string, text string) string {
	return strings.Replace(text, "${namespace}", namespace, -1)
}
