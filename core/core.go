package core

import (
	"fmt"
	"strings"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/sql"
	"github.com/magneticio/forklift/util"
)

type SqlConfiguration struct {
	Database          string `yaml:"database,omitempty" json:"database,omitempty"`
	Table             string `yaml:"table,omitempty" json:"table,omitempty"`
	User              string `yaml:"user,omitempty" json:"user,omitempty"`
	Password          string `yaml:"password,omitempty" json:"password,omitempty"`
	Url               string `yaml:"url,omitempty" json:"url,omitempty"`
	DatabaseServerUrl string `yaml:"database-server-url,omitempty" json:"database-server-url,omitempty"`
}

type VaultKeyValueStoreConfiguration struct {
	Url   string `yaml:"url,omitempty" json:"url,omitempty"`
	Token string `yaml:"token,omitempty" json:"token,omitempty"`
}

type KeyValueStoreConfiguration struct {
	Type     string                          `yaml:"type,omitempty" json:"type,omitempty"`
	BasePath string                          `yaml:"base-path,omitempty" json:"base-path,omitempty"`
	Vault    VaultKeyValueStoreConfiguration `yaml:"vault,omitempty" json:"vault,omitempty"`
}

type Configuration struct {
	Sql           SqlConfiguration           `yaml:"sql,omitempty" json:"sql,omitempty"`
	KeyValueStore KeyValueStoreConfiguration `yaml:"key-value-store,omitempty" json:"key-value-store,omitempty"`
	Hocon         string                     `hocon:"sql,omitempty" json:"hocon,omitempty"`
}

type VampConfiguration struct {
	Persistence Persistence `yaml:"persistence,omitempty" json:"persistence,omitempty"`
	Model       Model       `yaml:"model,omitempty" json:"model,omitempty"`
}

type Database struct {
	Type string           `yaml:"type,omitempty" json:"type,omitempty"`
	Sql  SqlConfiguration `yaml:"sql,omitempty" json:"sql,omitempty"`
}

type Persistence struct {
	Database      Database                   `yaml:"database,omitempty" json:"database,omitempty"`
	KeyValueStore KeyValueStoreConfiguration `yaml:"key-value-store,omitempty" json:"key-value-store,omitempty"`
	Transformers  Transformers               `yaml:"transformers,omitempty" json:"transformers,omitempty"`
}

type Resolvers struct {
	namespace []string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}

type Transformers struct {
	classes []string `yaml:"classes,omitempty" json:"classes,omitempty"`
}

type Model struct {
	Resolvers Resolvers `yaml:"resolvers,omitempty" json:"resolvers,omitempty"`
}

type Security struct {
	LookupHashSalt        string `yaml:"lookup-hash-salt,omitempty" json:"lookup-hash-salt,omitempty"`
	LookupHashAlgorithm   string `yaml:"lookup-hash-algorithm,omitempty" json:"lookup-hash-algorithm,omitempty"`
	SessionIdLength       int    `yaml:"session-id-length,omitempty" json:"session-id-length,omitempty"`
	PasswordHashAlgorithm string `yaml:"password-hash-algorithm,omitempty" json:"password-hash-algorithm,omitempty"`
	PasswordHashSalt      string `yaml:"password-hash-salt,omitempty" json:"password-hash-salt,omitempty"`
	TokenValueLength      int    `yaml:"token-value-length,omitempty" json:"token-value-length,omitempty"`
}

type Core struct {
	VampConfiguration vampConfiguration `yaml:"vamp,omitempty" json:"vamp,omitempty"`
	Conf              Configuration     `yaml:"forklift,omitempty" json:"forklift,omitempty"`
}

func NewCore(conf Configuration) (*Core, error) {

	return &Core{
		Conf: conf,
	}, nil
}

func (c *Core) CreateOrganization(namespacedOrganization string) error {

	keyValueStoreConfig := c.GetNamespaceKeyValueStoreConfiguration(namespacedOrganization)
	// TODO: add params
	params := map[string]string{
		"cert":   "???",
		"key":    "???",
		"caCert": "???",
	}
	keyValueStoreClient, keyValueStoreClientError := keyvaluestoreclient.NewVaultKeyValueStoreClient(keyValueStoreConfig.Vault.Url, keyValueStoreConfig.Vault.Token, params)
	if keyValueStoreClientError != nil {
		return keyValueStoreClientError
	}
	key := "vamp/" + namespacedOrganization // this should be fixed
	value := map[string]interface{}{
		"value": c.Conf.Hocon,
	}
	keyValueStoreClientPutError := keyValueStoreClient.Put(key, value)
	if keyValueStoreClientPutError != nil {
		return keyValueStoreClientPutError
	}

	sqlConfig := c.GetNamespaceSqlConfiguration(namespacedOrganization)

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

	return client.SetupOrganization(namespacedOrganization)

}

func (c *Core) GetNamespaceSqlConfiguration(namespace string) *SqlConfiguration {
	return &SqlConfiguration{
		Database:          Namespaced(namespace, c.Conf.Sql.Database),
		Table:             Namespaced(namespace, c.Conf.Sql.Table),
		User:              Namespaced(namespace, c.Conf.Sql.User),
		Password:          Namespaced(namespace, c.Conf.Sql.Password),
		Url:               Namespaced(namespace, c.Conf.Sql.Url),
		DatabaseServerUrl: Namespaced(namespace, c.Conf.Sql.DatabaseServerUrl),
	}
}

func (c *Core) GetNamespaceKeyValueStoreConfiguration(namespace string) *KeyValueStoreConfiguration {
	return &KeyValueStoreConfiguration{
		Type:     c.Conf.KeyValueStore.Type,
		BasePath: Namespaced(namespace, c.Conf.KeyValueStore.BasePath),
		Vault: VaultKeyValueStoreConfiguration{
			Url:   Namespaced(namespace, c.Conf.KeyValueStore.Vault.Url),
			Token: Namespaced(namespace, c.Conf.KeyValueStore.Vault.Token),
		},
	}
}

func Namespaced(namespace string, text string) string {
	return strings.Replace(text, "${namespace}", namespace, -1)
}
