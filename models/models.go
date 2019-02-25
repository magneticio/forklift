package models

type VampConfiguration struct {
	Persistence Persistence `yaml:"persistence,omitempty" json:"persistence,omitempty"`
	Model       Model       `yaml:"model,omitempty" json:"model,omitempty"`
	Security    Security    `yaml:"security,omitempty" json:"security,omitempty"`
	Pulse       Pulse       `yaml:"pulse,omitempty" json:"pulse,omitempty"`
	Metadata    Metadata    `yaml:"metadata,omitempty" json:"metadata,omitempty"`
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

type Index struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
}
type Elasticsearch struct {
	Index Index  `yaml:"index,omitempty" json:"index,omitempty"`
	Url   string `yaml:"url,omitempty" json:"url,omitempty"`
}

type Pulse struct {
	Type          string        `yaml:"type,omitempty" json:"type,omitempty"`
	Elasticsearch Elasticsearch `yaml:"elasticsearch,omitempty" json:"elasticsearch,omitempty"`
}

type Namespace struct {
	Title string `yaml:"title,omitempty" json:"title,omitempty"`
}

type Metadata struct {
	Namespace Namespace `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}

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
