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
	// fmt.Printf("Vault store at key: %v\n", key)
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

func (c *Core) ListOrganizations(baseNamespace string) ([]string, error) {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration("")
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return nil, keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath
	list, keyValueStoreClientListError := keyValueStoreClient.List(key)
	if keyValueStoreClientListError != nil {
		return nil, keyValueStoreClientListError
	}
	filteredMap := make(map[string]bool)
	for _, name := range list {
		if strings.HasPrefix(name, baseNamespace) {
			filteredName := strings.Split(name, "-")
			if len(filteredName) == 2 {
				filteredMap[filteredName[1]] = true
			}
		}
	}
	filteredReducedList := make([]string, len(filteredMap))
	i := 0
	for k, _ := range filteredMap {
		filteredReducedList[i] = k
		i++
	}
	return filteredReducedList, nil
}

func (c *Core) ListEnvironments(baseNamespace string, organization string) ([]string, error) {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration("")
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return nil, keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath
	list, keyValueStoreClientListError := keyValueStoreClient.List(key)
	if keyValueStoreClientListError != nil {
		return nil, keyValueStoreClientListError
	}
	filteredMap := make(map[string]bool)
	filterPrefix := baseNamespace + "-" + organization
	for _, name := range list {
		if strings.HasPrefix(name, filterPrefix) {
			filteredName := strings.Split(name, "-")
			if len(filteredName) == 3 {
				filteredMap[filteredName[2]] = true
			}
		}
	}
	filteredReducedList := make([]string, len(filteredMap))
	i := 0
	for k, _ := range filteredMap {
		filteredReducedList[i] = k
		i++
	}
	return filteredReducedList, nil
}

func (c *Core) CreateEnvironment(namespace string, organization string, elements map[string]string, configuration Configuration) error {

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
	// fmt.Printf("Vault store at key: %v\n", key)
	keyValueStoreClientPutError := keyValueStoreClient.PutValue(key, string(value))
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	namespacedOrganizationName := c.GetNamespaceDatabaseConfiguration(organization).Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.SetupEnvironment(namespacedOrganizationName, databaseConfig.Sql.Table, elements)

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
