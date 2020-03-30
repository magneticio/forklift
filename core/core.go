package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/sql"
	"github.com/magneticio/forklift/util"
	policies "github.com/magneticio/vamp-policies"
)

type Core struct {
	Conf models.ForkliftConfiguration
}

func NewCore(conf models.ForkliftConfiguration) (*Core, error) {

	return &Core{
		Conf: conf,
	}, nil
}

func (c *Core) GetArtifact(organization string, environment string, name string, kind string) (*models.SqlElement, error) {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(environment)
	if err != nil {
		return nil, err
	}

	organizationConfig, err := c.GetNamespaceDatabaseConfiguration(organization)
	if err != nil {
		return nil, err
	}

	namespacedOrganizationName := organizationConfig.Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
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

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return err
	}

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

	encodedPassword := util.EncodeString(password, configuration.Vamp.Security.PasswordHashAlgorithm, configuration.Vamp.Security.PasswordHashSalt)

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

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	return client.Insert(databaseConfig.Sql.Database, databaseConfig.Sql.Table, sqlElementAsJsonString)

}

func (c *Core) UpdateUser(namespace string, name string, role string, password string) error {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return err
	}

	userElement, validationError := c.GetUser(namespace, name)
	if validationError != nil {

		return validationError
	}
	if userElement == nil {

		return errors.New(fmt.Sprintf("%v %v does not exist", models.UsersKind, name))
	}

	// get organization Configuration using namespace
	configuration, configurationError := c.getConfig(namespace)
	if configurationError != nil {
		return configurationError
	}

	encodedPassword := util.EncodeString(password, configuration.Vamp.Security.PasswordHashAlgorithm, configuration.Vamp.Security.PasswordHashSalt)

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

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	return client.InsertOrReplace(databaseConfig.Sql.Database, databaseConfig.Sql.Table, sqlElement.Name, sqlElement.Kind, sqlElementAsJsonString)

}

func (c *Core) DeleteUser(namespace string, user string) error {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return err
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	return client.DeleteByNameAndKind(databaseConfig.Sql.Database, databaseConfig.Sql.Table, user, models.UsersKind) //TODO admin should be a constant

}

func (c *Core) AddUser(namespace string, user string) error {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return err
	}

	// get organization Configuration using namespace
	configuration, configurationError := c.getConfig(namespace)
	if configurationError != nil {
		return configurationError
	}

	var userArtifact models.Artifact

	marshallError := json.Unmarshal([]byte(user), &userArtifact)
	if marshallError != nil {
		return marshallError
	}

	userArtifact.Password = util.EncodeString(userArtifact.Password, configuration.Vamp.Security.PasswordHashAlgorithm, configuration.Vamp.Security.PasswordHashSalt)

	userJson, marshalError := json.Marshal(userArtifact)
	if marshalError != nil {
		return marshalError
	}

	sqlElement, convertError := ConvertToSqlElement(string(userJson))
	if convertError != nil {
		return convertError
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	sqlElementString, jsonMarshallError := json.Marshal(sqlElement)
	if jsonMarshallError != nil {
		return jsonMarshallError
	}

	return client.InsertOrReplace(databaseConfig.Sql.Database, databaseConfig.Sql.Table, sqlElement.Name, sqlElement.Kind, string(sqlElementString))

}

func (c *Core) GetUser(namespace string, name string) (*models.SqlElement, error) {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return nil, err
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
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

func (c *Core) ListArtifacts(organization string, environment string, kind string) ([]string, error) {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(environment)
	if err != nil {
		return nil, err
	}

	organizationConfig, err := c.GetNamespaceDatabaseConfiguration(organization)
	if err != nil {
		return nil, err
	}

	namespacedOrganizationName := organizationConfig.Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return nil, clientError
	}

	result, queryError := client.List(namespacedOrganizationName, databaseConfig.Sql.Table, kind)
	if queryError != nil {
		return nil, queryError
	}

	if result == nil {
		return nil, nil
	}

	var names []string

	for _, element := range result {

		var sqlElement models.SqlElement

		jsonUnmarshallError := json.Unmarshal([]byte(element.Record), &sqlElement)
		if jsonUnmarshallError != nil {
			return nil, jsonUnmarshallError
		}

		names = append(names, sqlElement.Name)
	}

	return names, nil
}

func (c *Core) ListUsers(namespace string) ([]string, error) {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return nil, err
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return nil, clientError
	}

	result, queryError := client.List(databaseConfig.Sql.Database, databaseConfig.Sql.Table, "users")
	if queryError != nil {
		return nil, queryError
	}

	if result == nil {
		return nil, nil
	}

	var names []string

	for _, element := range result {

		var sqlElement models.SqlElement

		jsonUnmarshallError := json.Unmarshal([]byte(element.Record), &sqlElement)
		if jsonUnmarshallError != nil {
			return nil, jsonUnmarshallError
		}

		names = append(names, sqlElement.Name)
	}

	return names, nil
}

func (c *Core) AddPolicy(organization string, environment string, policyContent string) error {
	logging.Info("Adding policy:\n")
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(environment)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	policyAPI := policies.NewPolicyAPI(keyValueStoreClient, c.Conf.ReleaseAgentKeyValueStorePath)
	return policyAPI.Save(policyContent)
}

func (c *Core) DeleteReleasePolicy(organization string, environment string, policyName string) error {
	logging.Info("Deleting policy: %v\n", policyName)
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(environment)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	policyAPI := policies.NewPolicyAPI(keyValueStoreClient, c.Conf.ReleaseAgentKeyValueStorePath)
	return policyAPI.Delete(policyName)
}

// AddReleasePlan - adds release plan to key value store
func (c *Core) AddReleasePlan(name string, content string) error {
	keyValueStoreConfig := c.GetKeyValueStoreConfiguration()
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := path.Join(c.Conf.ReleasePlansKeyValueStorePath, name)
	logging.Info("Storing Release Plan Under Key: %v\n", key)
	keyValueStoreClientPutError := keyValueStoreClient.Put(key, content)
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}
	return nil
}

// DeleteReleasePlan - deletes release plan from key value store
func (c *Core) DeleteReleasePlan(name string) error {
	keyValueStoreConfig := c.GetKeyValueStoreConfiguration()
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := path.Join(c.Conf.ReleasePlansKeyValueStorePath, name)
	logging.Info("Deleting Release Plan Under Key: %v\n", key)
	keyValueStoreClientDeleteError := keyValueStoreClient.Delete(key)
	if keyValueStoreClientDeleteError != nil {
		return keyValueStoreClientDeleteError
	}
	return nil
}

func (c *Core) addArtifactToDatabase(organization string, environment string, content string) error {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(environment)
	if err != nil {
		return err
	}

	organizationDbConfiguration, err := c.GetNamespaceDatabaseConfiguration(organization)
	if err != nil {
		return err
	}

	sqlElement, convertError := ConvertToSqlElement(content)
	if convertError != nil {
		return convertError
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	sqlElementString, jsonMarshallError := json.Marshal(sqlElement)
	if jsonMarshallError != nil {
		return jsonMarshallError
	}

	if sqlElement.Kind == "workflows" {
		tokenSqlElementAsString, generateTokenError := c.GenerateTokenForWorkflow(environment, sqlElement.Name, "admin")
		if generateTokenError != nil {
			return generateTokenError
		}
		insertTokenError := client.InsertOrReplace(organizationDbConfiguration.Sql.Database, organizationDbConfiguration.Sql.Table, sqlElement.Name, "tokens", tokenSqlElementAsString)
		if insertTokenError != nil {
			return insertTokenError
		}
	}

	return client.InsertOrReplace(organizationDbConfiguration.Sql.Database, databaseConfig.Sql.Table, sqlElement.Name, sqlElement.Kind, string(sqlElementString))
}

func (c *Core) addArtifactToVault(organization string, environment string, content string) error {
	var artifact models.Artifact
	unmarshallError := json.Unmarshal([]byte(content), &artifact)
	if unmarshallError != nil {
		logging.Error("Unmarshalling error : %v\n", unmarshallError)
		return unmarshallError
	}

	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(environment)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := path.Join(keyValueStoreConfig.BasePath, c.Conf.ReleaseAgentKeyValueStorePath, artifact.Kind, artifact.Name)
	logging.Info("Storing Artifact Under Key: %v\n", key)
	keyValueStoreClientPutError := keyValueStoreClient.Put(key, content)
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}
	return nil
}

// AddArtifact : adds artifact to sql database and Vault
func (c *Core) AddArtifact(organization string, environment string, content string) error {

	if c.Conf.DatabaseEnabled {

		dbError := c.addArtifactToDatabase(organization, environment, content)
		if dbError != nil {
			return dbError
		}

	}

	return c.addArtifactToVault(organization, environment, content)
}

func (c *Core) deleteArtifactFromDatabase(organization string, environment string, name string, kind string) error {

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(environment)
	if err != nil {
		return err
	}

	organizationDatabaseConfig, err := c.GetNamespaceDatabaseConfiguration(organization)
	if err != nil {
		return err
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	tokenName := GenerateTokenName(environment, name, kind)

	deleteTokenError := client.DeleteByNameAndKind(organizationDatabaseConfig.Sql.Database, organizationDatabaseConfig.Sql.Table, tokenName, "tokens")
	if deleteTokenError != nil {
		return deleteTokenError
	}

	return client.DeleteByNameAndKind(organizationDatabaseConfig.Sql.Database, databaseConfig.Sql.Table, name, kind)
}

func (c *Core) deleteArtifactFromVault(organization string, environment string, name string, kind string) error {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(environment)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := path.Join(keyValueStoreConfig.BasePath, c.Conf.ReleaseAgentKeyValueStorePath, kind, name)
	logging.Info("Deleting Artifact Under Key: %v\n", key)
	keyValueStoreClientDeleteError := keyValueStoreClient.Delete(key)
	if keyValueStoreClientDeleteError != nil {
		return keyValueStoreClientDeleteError
	}
	return nil
}

// DeleteArtifact : deletes artifact from sql database and Vault
func (c *Core) DeleteArtifact(organization string, environment string, name string, kind string) error {

	if c.Conf.DatabaseEnabled {

		if dbError := c.deleteArtifactFromDatabase(organization, environment, name, kind); dbError != nil {
			return dbError
		}

	}

	return c.deleteArtifactFromVault(organization, environment, name, kind)
}

func (c *Core) CreateOrganization(namespace string, configuration models.VampConfiguration) error {

	putConfigError := c.putConfig(namespace, configuration)
	if putConfigError != nil {
		return putConfigError
	}

	if !c.Conf.DatabaseEnabled {
		//If database is not enabled we return
		return nil
	}

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return err
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	return client.SetupOrganization(databaseConfig.Sql.Database, databaseConfig.Sql.Table)

}

func (c *Core) UpdateOrganization(namespace string, configuration models.VampConfiguration) error {
	return c.putConfig(namespace, configuration)
}

func (c *Core) ListOrganizations(baseNamespace string) ([]string, error) {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration("")
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return nil, keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath
	logging.Info("Listing Values in Key Value Store with namespace: %v\n", baseNamespace)
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
	logging.Info("Listing Values in Key Value Store Under Key: %v\n", key)
	list, keyValueStoreClientListError := keyValueStoreClient.List(key)
	if keyValueStoreClientListError != nil {
		return nil, keyValueStoreClientListError
	}
	filteredMap := make(map[string]bool)
	filterPrefix := baseNamespace + "-" + organization + "-"
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

func (c *Core) CreateEnvironment(namespace string, organization string, elements []string, configuration models.VampConfiguration) error {
	putConfigError := c.putConfig(namespace, configuration)
	if putConfigError != nil {
		return putConfigError
	}

	if !c.Conf.DatabaseEnabled {
		//If database is not enabled we return
		return nil
	}

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return err
	}

	organizationDatabaseConfig, err := c.GetNamespaceDatabaseConfiguration(organization)
	if err != nil {
		return err
	}

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	sqlElements := make([]string, len(elements))
	for i, element := range elements {
		sqlElement, convertError := ConvertToSqlElementJson(element)
		if convertError != nil {
			return convertError
		}
		sqlElements[i] = sqlElement

		sqlElementAsStruct, conversionError := ConvertToSqlElement(element)
		if conversionError != nil {
			return conversionError
		}
		if sqlElementAsStruct.Kind == "workflows" {
			tokenSqlElementAsString, generateTokenError := c.GenerateTokenForWorkflow(namespace, sqlElementAsStruct.Name, "admin")
			if generateTokenError != nil {
				return generateTokenError
			}
			insertTokenError := client.InsertOrReplace(organizationDatabaseConfig.Sql.Database, organizationDatabaseConfig.Sql.Table, sqlElementAsStruct.Name, "tokens", tokenSqlElementAsString)
			if insertTokenError != nil {
				return insertTokenError
			}
		}
	}

	return client.SetupEnvironment(organizationDatabaseConfig.Sql.Database, databaseConfig.Sql.Table, sqlElements)

}

func (c *Core) UpdateEnvironment(namespace string, organization string, elements []string, configuration models.VampConfiguration) error {

	putConfigError := c.putConfig(namespace, configuration)
	if putConfigError != nil {
		return putConfigError
	}

	if !c.Conf.DatabaseEnabled {
		//If database is not enabled we return
		return nil
	}

	databaseConfig, err := c.GetNamespaceDatabaseConfiguration(namespace)
	if err != nil {
		return err
	}

	organizationConfig, err := c.GetNamespaceDatabaseConfiguration(organization)
	if err != nil {
		return err
	}

	namespacedOrganizationName := organizationConfig.Sql.Database

	client, clientError := sql.NewSqlClient(databaseConfig)
	if clientError != nil {
		logging.Error("Client error: %v\n", clientError.Error())
		return clientError
	}

	return client.UpdateEnvironment(namespacedOrganizationName, databaseConfig.Sql.Table, elements)

}

// GetNamespaceDatabaseConfiguration retrieves the database configuration by namespace
func (c *Core) GetNamespaceDatabaseConfiguration(namespace string) (models.Database, error) {

	if !c.Conf.DatabaseEnabled {
		return models.Database{}, errors.New("Database is not enabled")
	}

	return models.Database{
		Sql: models.SqlConfiguration{
			Database: Namespaced(namespace, c.Conf.DatabaseName),
			Table:    Namespaced(namespace, c.Conf.DatabaseTable),
			User:     c.Conf.DatabaseUser,
			Password: c.Conf.DatabasePassword,
			Url:      Namespaced(namespace, c.Conf.DatabaseURL),
		},
		Type: c.Conf.DatabaseType,
	}, nil
}

func (c *Core) GetNamespaceKeyValueStoreConfiguration(namespace string) *models.KeyValueStoreConfiguration {
	return &models.KeyValueStoreConfiguration{
		Type:     c.Conf.KeyValueStoreType,
		BasePath: Namespaced(namespace, c.Conf.KeyValueStoreBasePath),
		Vault: models.VaultKeyValueStoreConfiguration{
			Url:               c.Conf.KeyValueStoreUrL,
			Token:             c.Conf.KeyValueStoreToken,
			ServerTlsCert:     c.Conf.KeyValueStoreServerTlsCert,
			ClientTlsCert:     c.Conf.KeyValueStoreClientTlsCert,
			ClientTlsKey:      c.Conf.KeyValueStoreClientTlsKey,
			KvMode:            c.Conf.KeyValueStoreKvMode,
			FallbackKvVersion: c.Conf.KeyValueStoreFallbackKvVersion,
		},
	}
}

// GetKeyValueStoreConfiguration - gets key value store configuration
func (c *Core) GetKeyValueStoreConfiguration() *models.KeyValueStoreConfiguration {

	return &models.KeyValueStoreConfiguration{
		Type: c.Conf.KeyValueStoreType,
		Vault: models.VaultKeyValueStoreConfiguration{
			Url:               c.Conf.KeyValueStoreUrL,
			Token:             c.Conf.KeyValueStoreToken,
			ServerTlsCert:     c.Conf.KeyValueStoreServerTlsCert,
			ClientTlsCert:     c.Conf.KeyValueStoreClientTlsCert,
			ClientTlsKey:      c.Conf.KeyValueStoreClientTlsKey,
			KvMode:            c.Conf.KeyValueStoreKvMode,
			FallbackKvVersion: c.Conf.KeyValueStoreFallbackKvVersion,
		},
	}

}

func (c *Core) DeleteOrganization(namespace string) error {
	return c.deleteConfig(namespace)
}

func (c *Core) DeleteEnvironment(namespace string) error {
	return c.deleteConfig(namespace)
}

func (c *Core) ShowOrganization(namespace string) (*models.VampConfiguration, error) {
	conf, err := c.getConfig(namespace)
	if err != nil {
		if strings.HasPrefix(err.Error(), "No Values") {
			return nil, errors.New("No organization found")
		}
		return nil, err
	}
	return conf, err
}

func (c *Core) ShowEnvironment(namespace string) (*models.VampConfiguration, error) {
	conf, err := c.getConfig(namespace)
	if err != nil {
		if strings.HasPrefix(err.Error(), "No Values") {
			return nil, errors.New("No environment found")
		}
		return nil, err
	}
	return conf, err
}

func (c *Core) putConfig(namespace string, configuration models.VampConfiguration) error {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(namespace)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath + "/configuration/applied"
	logging.Info("Storing Config Under Key: %v\n", key)
	value, jsonMarshallError := json.Marshal(configuration)
	if jsonMarshallError != nil {
		return jsonMarshallError
	}
	keyValueStoreClientPutError := keyValueStoreClient.Put(key, string(value))
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}
	return nil
}

func (c *Core) getConfig(namespace string) (*models.VampConfiguration, error) {
	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(namespace)
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewKeyValueStoreClient(*keyValueStoreConfig)
	if keyValueStoreClientError != nil {
		return nil, keyValueStoreClientError
	}
	key := keyValueStoreConfig.BasePath + "/configuration/applied"
	logging.Info("Reading Config Under Key: %v\n", key)
	configJson, keyValueStoreClientGetError := keyValueStoreClient.Get(key)
	if keyValueStoreClientGetError != nil {
		return nil, keyValueStoreClientGetError
	}
	var configuration models.VampConfiguration
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
	logging.Info("Deleting Config Under Key: %v\n", key)
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

func GenerateTokenName(namespace string, workflowName string, kindInTokenName string) string {

	artifactVersion := "v1"
	namespaceReference := "class io.vamp.common.Namespace@" + namespace
	lookupHashAlgorithm := "SHA1" // it is fixed
	logging.Info("namespaceReference %v LookupHashAlgorithm %v, artifactVersion %v\n", namespaceReference, lookupHashAlgorithm, artifactVersion)
	lookupName := util.EncodeString(namespaceReference, lookupHashAlgorithm, artifactVersion)

	return fmt.Sprintf("%s/%s/%s", lookupName, kindInTokenName, workflowName)
}

func (c *Core) GenerateTokenForWorkflow(namespace string, workflowName string, role string) (string, error) {
	// get Configuration using namespace
	s := strings.Split(namespace, "-")
	configuration, configurationError := c.getConfig(s[0] + "-" + s[1])
	if configurationError != nil {
		return "", configurationError
	}
	kind := "tokens"
	kindInTokenName := "workflows"

	tokenName := GenerateTokenName(namespace, workflowName, kindInTokenName)
	//TODO: More meaningful configuration.Vamp.Security.PasswordHashSalt

	encodedValue := util.RandomEncodedString(configuration.Vamp.Security.TokenValueLength)

	artifact := models.Artifact{
		Name:      tokenName,
		Value:     encodedValue,
		Namespace: namespace,
		Kind:      kind,
		Roles:     []string{role},
		Metadata:  map[string]string{},
	}

	artifactAsJson, artifactJsonError := json.Marshal(artifact)
	if artifactJsonError != nil {
		return "", artifactJsonError
	}

	artifactAsJsonString := string(artifactAsJson)
	logging.Info("Token string: %v\n", artifactAsJsonString)

	sqlElement := models.SqlElement{
		Version:   models.BackendVersion,
		Instance:  util.UUID(),
		Timestamp: util.Timestamp(),
		Name:      tokenName,
		Kind:      kind,
		Artifact:  artifactAsJsonString,
	}

	sqlElementAsJson, sqlElementJsonError := json.Marshal(sqlElement)
	if sqlElementJsonError != nil {
		return "", sqlElementJsonError
	}

	sqlElementAsJsonString := string(sqlElementAsJson)

	return sqlElementAsJsonString, nil

}
