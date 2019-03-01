package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/sql"
	"github.com/magneticio/forklift/util"
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

func (c *Core) CreateAdmin(namespace string, name string, password string) error {

	version := "1.0.4" // TODO: this should be a constant
	kind := "users"    // TODO: this should be a constant

	encodedPassword := util.EncodeString(password, c.Conf.VampConfiguration.Security.PasswordHashAlgorithm, c.Conf.VampConfiguration.Security.PasswordHashSalt)

	artifact := models.Artifact{
		Name:     name,
		Password: encodedPassword,
		Kind:     kind,
		Roles:    []string{"admin"},
		Metadata: map[string]string{},
	}

	artifactAsJson, artifactJsonError := json.Marshal(artifact)
	if artifactJsonError != nil {
		return artifactJsonError
	}

	artifactAsJsonString := string(artifactAsJson)

	sqlElement := models.SqlElement{
		Version:   version,
		Instance:  util.UUID(),
		Timestamp: util.Timestamp(),
		Name:      name,
		Kind:      kind,
		Artifact:  artifactAsJsonString,
	}

	sqlElementAsJson, sqlElementJsonError := json.Marshal(sqlElement)
	if sqlElementJsonError != nil {
		return sqlElementJsonError
	}

	sqlElementAsJsonString := string(sqlElementAsJson)

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.Insert(databaseConfig.Sql.Database, databaseConfig.Sql.Table, sqlElementAsJsonString)

}

func (c *Core) AddAdmin(namespace string, admin string) error {

	sqlElement, convertError := ConvertToSqlElement(admin)
	if convertError != nil {
		return convertError
	}

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.Insert(databaseConfig.Sql.Database, databaseConfig.Sql.Table, sqlElement)

}

func (c *Core) DeleteAdmin(namespace string, admin string) error {

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.DeleteByNameAndKind(databaseConfig.Sql.Database, databaseConfig.Sql.Table, admin, "users") //TODO admin should be a constant

}

func (c *Core) CreateOrganization(namespace string, configuration Configuration) error {
	putConfigError := c.putConfig(namespace, configuration)
	if putConfigError != nil {
		return putConfigError
	}

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.SetupOrganization(databaseConfig.Sql.Database, databaseConfig.Sql.Table)

}

func (c *Core) UpdateOrganization(namespace string, configuration Configuration) error {
	return c.putConfig(namespace, configuration)
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

func (c *Core) CreateEnvironment(namespace string, organization string, elements []string, configuration Configuration) error {
	putConfigError := c.putConfig(namespace, configuration)
	if putConfigError != nil {
		return putConfigError
	}

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	namespacedOrganizationName := c.GetNamespaceDatabaseConfiguration(organization).Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	sqlElements := make([]string, len(elements))
	for i, element := range elements {
		sqlElement, convertError := ConvertToSqlElement(element)
		if convertError != nil {
			return convertError
		}
		sqlElements[i] = sqlElement
	}

	return client.SetupEnvironment(namespacedOrganizationName, databaseConfig.Sql.Table, sqlElements)

}

func (c *Core) UpdateEnvironment(namespace string, organization string, elements []string, configuration Configuration) error {

	putConfigError := c.putConfig(namespace, configuration)
	if putConfigError != nil {
		return putConfigError
	}

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	namespacedOrganizationName := c.GetNamespaceDatabaseConfiguration(organization).Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.UpdateEnvironment(namespacedOrganizationName, databaseConfig.Sql.Table, elements)

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

func (c *Core) DeleteOrganization(namespace string) error {
	return c.deleteConfig(namespace)
}

func (c *Core) DeleteEnvironment(namespace string) error {
	return c.deleteConfig(namespace)
}

func (c *Core) putConfig(namespace string, configuration Configuration) error {
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
	keyValueStoreClientPutError := keyValueStoreClient.PutValue(key, string(value))
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}
	return nil
}

func (c *Core) deleteConfig(namespace string) error {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(namespace)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath + "/configuration/applied"
	keyValueStoreClientPutError := keyValueStoreClient.Delete(key)
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}
	return nil
}

func Namespaced(namespace string, text string) string {
	return strings.Replace(text, "${namespace}", namespace, -1)
}

func ConvertToSqlElement(artifactAsJsonString string) (string, error) {
	var artifact models.Artifact
	unmarshallError := json.Unmarshal([]byte(artifactAsJsonString), &artifact)
	if unmarshallError != nil {
		return "", unmarshallError
	}
	version := "1.0.4" // TODO: this should be a constant
	sqlElement := models.SqlElement{
		Version:   version,
		Instance:  util.UUID(),
		Timestamp: util.Timestamp(),
		Name:      artifact.Name,
		Kind:      artifact.Kind,
		Artifact:  artifactAsJsonString,
	}
	sqlElementString, jsonMarshallError := json.Marshal(sqlElement)
	if jsonMarshallError != nil {
		return "", jsonMarshallError
	}
	return string(sqlElementString), nil
}
