package core

import (
	"encoding/json"
	"errors"
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

func (c *Core) GetArtifact(organization string, environment string, name string, kind string) (*models.SqlElement, error) {

	databaseConfig := c.GetNamespaceDatabaseConfiguration(environment)

	namespacedOrganizationName := c.GetNamespaceDatabaseConfiguration(organization).Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return nil, clientError
	}

	result, queryError := client.FindByNameAndKind(namespacedOrganizationName, databaseConfig.Sql.Table, name, kind)
	if queryError != nil {
		return nil, queryError
	}

	if result == nil {
		return nil, nil
	}

	var sqlElement models.SqlElement

	jsonUnmarshallError := json.Unmarshal([]byte(result.Record), &sqlElement)
	if jsonUnmarshallError != nil {
		return nil, jsonUnmarshallError
	}
	return &sqlElement, nil
}

func (c *Core) CreateUser(namespace string, name string, role string, password string) error {

	userElement, validationError := c.GetUser(namespace, name)
	if validationError != nil {

		return validationError
	}
	if userElement != nil {

		return errors.New(fmt.Sprintf("%v %v already exists", models.UsersKind, name))
	}

	// get organization Configuration using namespace
	configuration, configurationError := c.getConfig(namespace)
	if configurationError != nil {
		return configurationError
	}

	encodedPassword := util.EncodeString(password, configuration.VampConfiguration.Security.PasswordHashAlgorithm, configuration.VampConfiguration.Security.PasswordHashSalt)

	artifact := models.Artifact{
		Name:     name,
		Password: encodedPassword,
		Kind:     models.UsersKind,
		Roles:    []string{role},
		Metadata: map[string]string{},
	}

	artifactAsJson, artifactJsonError := json.Marshal(artifact)
	if artifactJsonError != nil {
		return artifactJsonError
	}

	artifactAsJsonString := string(artifactAsJson)

	sqlElement := models.SqlElement{
		Version:   models.BackendVersion,
		Instance:  util.UUID(),
		Timestamp: util.Timestamp(),
		Name:      name,
		Kind:      models.UsersKind,
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

func (c *Core) DeleteUser(namespace string, user string) error {

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.DeleteByNameAndKind(databaseConfig.Sql.Database, databaseConfig.Sql.Table, user, models.UsersKind) //TODO admin should be a constant

}

func (c *Core) AddUser(namespace string, user string) error {

	sqlElement, convertError := ConvertToSqlElement(user)
	if convertError != nil {
		return convertError
	}
	userElement, validationError := c.GetUser(namespace, sqlElement.Name)
	if validationError != nil {
		return validationError
	}
	if userElement != nil {
		return errors.New(fmt.Sprintf("User %v already exists", sqlElement.Name))
	}

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	sqlElementString, jsonMarshallError := json.Marshal(sqlElement)
	if jsonMarshallError != nil {
		return jsonMarshallError
	}

	return client.Insert(databaseConfig.Sql.Database, databaseConfig.Sql.Table, string(sqlElementString))

}

func (c *Core) GetUser(namespace string, name string) (*models.SqlElement, error) {

	databaseConfig := c.GetNamespaceDatabaseConfiguration(namespace)

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return nil, clientError
	}

	result, queryError := client.FindByNameAndKind(databaseConfig.Sql.Database, databaseConfig.Sql.Table, name, models.UsersKind)
	if queryError != nil {
		return nil, queryError
	}

	if result == nil {
		return nil, nil
	}

	var sqlElement models.SqlElement

	jsonUnmarshallError := json.Unmarshal([]byte(result.Record), &sqlElement)
	if jsonUnmarshallError != nil {
		return nil, jsonUnmarshallError
	}
	return &sqlElement, nil
}

func (c *Core) AddArtifact(organization string, environment string, content string) error {

	databaseConfig := c.GetNamespaceDatabaseConfiguration(environment)

	namespacedOrganizationName := c.GetNamespaceDatabaseConfiguration(organization).Sql.Database

	sqlElement, convertError := ConvertToSqlElement(content)
	if convertError != nil {
		return convertError
	}
	element, validationError := c.GetArtifact(organization, environment, sqlElement.Name, sqlElement.Kind)
	if validationError != nil {
		return validationError
	}
	if element != nil {
		return errors.New(fmt.Sprintf("%v %v already exists", sqlElement.Kind, sqlElement.Name))
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	sqlElementString, jsonMarshallError := json.Marshal(sqlElement)
	if jsonMarshallError != nil {
		return jsonMarshallError
	}

	return client.Insert(namespacedOrganizationName, databaseConfig.Sql.Table, string(sqlElementString))

}

func (c *Core) DeleteArtifact(organization string, environment string, name string, kind string) error {

	databaseConfig := c.GetNamespaceDatabaseConfiguration(environment)

	namespacedOrganizationName := c.GetNamespaceDatabaseConfiguration(organization).Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.DeleteByNameAndKind(namespacedOrganizationName, databaseConfig.Sql.Table, name, kind) //TODO admin should be a constant

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
		sqlElement, convertError := ConvertToSqlElementJson(element)
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
	model := c.Conf.VampConfiguration.Persistence.KeyValueStore
	model.BasePath = Namespaced(namespace, c.Conf.VampConfiguration.Persistence.KeyValueStore.BasePath)
	return &model
}

func (c *Core) DeleteOrganization(namespace string) error {
	return c.deleteConfig(namespace)
}

func (c *Core) DeleteEnvironment(namespace string) error {
	return c.deleteConfig(namespace)
}

func (c *Core) ShowOrganization(namespace string) (*Configuration, error) {
	return c.getConfig(namespace)
}

func (c *Core) ShowEnvironment(namespace string) (*Configuration, error) {
	return c.getConfig(namespace)
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

func (c *Core) getConfig(namespace string) (*Configuration, error) {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(namespace)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return nil, keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath + "/configuration/applied"
	configJson, keyValueStoreClientPutError := keyValueStoreClient.GetValue(key)
	if keyValueStoreClientPutError != nil {
		return nil, keyValueStoreClientPutError
	}
	var configuration Configuration
	jsonUnmarshallError := json.Unmarshal([]byte(configJson), &configuration)
	if jsonUnmarshallError != nil {
		return nil, jsonUnmarshallError
	}
	return &configuration, nil
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

func ConvertToSqlElement(artifactAsJsonString string) (*models.SqlElement, error) {
	var artifact models.Artifact
	unmarshallError := json.Unmarshal([]byte(artifactAsJsonString), &artifact)
	if unmarshallError != nil {
		fmt.Printf("Unmarshalling error : %v\n", artifactAsJsonString)
		fmt.Printf("Unmarshalling error : %v\n", unmarshallError.Error())
		return nil, unmarshallError
	}

	return &models.SqlElement{
		Version:   models.BackendVersion,
		Instance:  util.UUID(),
		Timestamp: util.Timestamp(),
		Name:      artifact.Name,
		Kind:      artifact.Kind,
		Artifact:  artifactAsJsonString,
	}, nil
}

func ConvertToSqlElementJson(artifactAsJsonString string) (string, error) {
	sqlElement, conversionError := ConvertToSqlElement(artifactAsJsonString)
	if conversionError != nil {
		return "", conversionError
	}

	sqlElementString, jsonMarshallError := json.Marshal(sqlElement)
	if jsonMarshallError != nil {
		return "", jsonMarshallError
	}
	return string(sqlElementString), nil
}
