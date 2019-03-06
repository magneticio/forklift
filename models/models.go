package models

type VampConfiguration struct {
	Persistence     Persistence     `yaml:"persistence,omitempty" json:"persistence,omitempty"`
	Model           Model           `yaml:"model,omitempty" json:"model,omitempty"`
	Security        Security        `yaml:"security,omitempty" json:"security,omitempty"`
	Pulse           Pulse           `yaml:"pulse,omitempty" json:"pulse,omitempty"`
	Metadata        Metadata        `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	ContainerDriver ContainerDriver `yaml:"container-driver,omitempty" json:"container-driver,omitempty"`
	Lifter          Lifter          `yaml:"lifter,omitempty" json:"lifter,omitempty"`
	GatewayDriver   GatewayDriver   `yaml:"gateway-driver,omitempty" json:"gateway-driver,omitempty"`
	WorkflowDriver  WorkflowDriver  `yaml:"workflow-driver,omitempty" json:"workflow-driver,omitempty"`
	Operation       Operation       `yaml:"operation,omitempty" json:"operation,omitempty"`
	Namespace       string          `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}

type WorkflowDriver struct {
	Workflow Workflow `yaml:"workflow,omitempty" json:"workflow,omitempty"`
	Type     string   `yaml:"type,omitempty" json:"type,omitempty"`
}

type Workflow struct {
	VampKeyValueStoreType        string       `yaml:"vamp-key-value-store-type,omitempty" json:"vamp-key-value-store-type,omitempty"`
	Deployables                  []Deployable `yaml:"deployables,omitempty" json:"deployables,omitempty"`
	Scale                        Scale        `yaml:"scale,omitempty" json:"scale,omitempty"`
	VampKeyValueStoreConnection  string       `yaml:"vamp-key-value-store-connection,omitempty" json:"vamp-key-value-store-connection,omitempty"`
	VampWorkflowExecutionPeriod  int          `yaml:"vamp-workflow-execution-period,omitempty" json:"vamp-workflow-execution-period,omitempty"`
	VampKeyValueStoreToken       string       `yaml:"vamp-key-value-store-token,omitempty" json:"vamp-key-value-store-token,omitempty"`
	VampWorkflowExecutionTimeout int          `yaml:"vamp-workflow-execution-timeout,omitempty" json:"vamp-workflow-execution-timeout,omitempty"`
	VampElasticsearchUrl         string       `yaml:"vamp-elasticsearch-url,omitempty" json:"vamp-elasticsearch-url,omitempty"`
	VampKeyValueStorePath        string       `yaml:"vamp-key-value-store-path,omitempty" json:"vamp-key-value-store-path,omitempty"`
	VampUrl                      string       `yaml:"vamp-url,omitempty" json:"vamp-url,omitempty"`
}

type Operation struct {
	Synchronization Synchronization `yaml:"synchronization,omitempty" json:"synchronization,omitempty"`
	Deployment      Deployment      `yaml:"deployment,omitempty" json:"deployment,omitempty"`
	Gateway         Gateway         `yaml:"gateway,omitempty" json:"gateway,omitempty"`
}

type Synchronization struct {
	Period     string                    `yaml:"period,omitempty" json:"period,omitempty"`
	Check      Check                     `yaml:"check,omitempty" json:"check,omitempty"`
	Deployment SynchronizationDeployment `yaml:"deployment,omitempty" json:"deployment,omitempty"`
}

type Check struct {
	HealthChecks         bool `yaml:"health-checks,omitempty" json:"health-checks,omitempty"`
	Deployable           bool `yaml:"deployable,omitempty" json:"deployable,omitempty"`
	Instances            bool `yaml:"instances,omitempty" json:"instances,omitempty"`
	Ports                bool `yaml:"ports,omitempty" json:"ports,omitempty"`
	Cpu                  bool `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	EnvironmentVariables bool `yaml:"environment-variables,omitempty" json:"environment-variables,omitempty"`
	Memory               bool `yaml:"memory,omitempty" json:"memory,omitempty"`
}

type SynchronizationDeployment struct {
	RefetchBreedOnUpdate bool `yaml:"refetch-breed-on-update,omitempty" json:"refetch-breed-on-update,omitempty"`
}

type Deployment struct {
	Scale     Scale    `yaml:"scale,omitempty" json:"scale,omitempty"`
	Arguments []string `yaml:"arguments" json:"arguments"`
}

type Scale struct {
	CPU       float32 `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Instances int     `yaml:"instances,omitempty" json:"instances,omitempty"`
	Memory    string  `yaml:"memory,omitempty" json:"memory,omitempty"`
}

type Gateway struct {
	VirtualHosts VirtualHosts `yaml:"virtual-hosts,omitempty" json:"virtual-hosts,omitempty"`
	Selector     string       `yaml:"selector,omitempty" json:"selector,omitempty"`
}

type VirtualHosts struct {
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty"`
}

type Deployable struct {
	Type  string `yaml:"type,omitempty" json:"type,omitempty"`
	Breed string `yaml:"breed,omitempty" json:"breed,omitempty"`
}

type GatewayDriver struct {
	Marshallers []Marshallers `yaml:"marshallers,omitempty" json:"marshallers,omitempty"`
}

type Marshallers struct {
	Type     string   `yaml:"type,omitempty" json:"type,omitempty"`
	Name     string   `yaml:"name,omitempty" json:"name,omitempty"`
	Template Template `yaml:"template,omitempty" json:"template,omitempty"`
}

type Template struct {
	Resource string `yaml:"resource,omitempty" json:"resource,omitempty"`
}

type Lifter struct {
	Artifacts []string `yaml:"artifacts,omitempty" json:"artifacts,omitempty"`
}

type ContainerDriver struct {
	Type       string     `yaml:"type,omitempty" json:"type,omitempty"`
	Kubernetes Kubernetes `yaml:"kubernetes,omitempty" json:"kubernetes,omitempty"`
}

type Kubernetes struct {
	Url                string `yaml:"url,omitempty" json:"url,omitempty"`
	Bearer             string `yaml:"bearer,omitempty" json:"bearer,omitempty"`
	Token              string `yaml:"token,omitempty" json:"token,omitempty"`
	Username           string `yaml:"username,omitempty" json:"username,omitempty"`
	Password           string `yaml:"password,omitempty" json:"password,omitempty"`
	VampGatewayAgentId string `yaml:"vamp-gateway-agent-id,omitempty" json:"vamp-gateway-agent-id,omitempty"`
	TlsCheck           bool   `yaml:"tls-check" json:"tls-check"`
	ServerCaCert       string `yaml:"server-ca-cert,omitempty" json:"server-ca-cert,omitempty"`
	ClientCert         string `yaml:"client-cert,omitempty" json:"client-cert,omitempty"`
	PrivateKey         string `yaml:"private-key,omitempty" json:"private-key,omitempty"`
	Namespace          string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}

type Database struct {
	Type string           `yaml:"type,omitempty" json:"type,omitempty"`
	Sql  SqlConfiguration `yaml:"sql,omitempty" json:"sql,omitempty"`
}

type Persistence struct {
	Database      Database                   `yaml:"database,omitempty" json:"database,omitempty"`
	KeyValueStore KeyValueStoreConfiguration `yaml:"key-value-store,omitempty" json:"key-value-store,omitempty"`
	Transformers  Transformers               `yaml:"transformers" json:"transformers"`
}

type Resolvers struct {
	Deployment []string `yaml:"deployment,omitempty" json:"deployment,omitempty"`
	Namespace  []string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	Workflow   []string `yaml:"workflow,omitempty" json:"workflow,omitempty"`
}

type Transformers struct {
	Classes []string `yaml:"classes" json:"classes"`
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
	Nats          Nats          `yaml:"nats,omitempty" json:"nats,omitempty"`
}

type Nats struct {
	Url       string `yaml:"url,omitempty" json:"url,omitempty"`
	ClusterID string `yaml:"cluster-id,omitempty" json:"cluster-id,omitempty"`
	ClientID  string `yaml:"client-id,omitempty" json:"client-id,omitempty"`
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
	Url               string `yaml:"url,omitempty" json:"url,omitempty"`
	Token             string `yaml:"token,omitempty" json:"token,omitempty"`
	KvMode            string `yaml:"kv-mode,omitempty" json:"kv-mode,omitempty"`
	FallbackKvVersion int    `yaml:"fallback-kv-versione,omitempty" json:"fallback-kv-version,omitempty"`
	ServerTlsCert     string `yaml:"server-tls-cert,omitempty" json:"server-tls-cert,omitempty"`
	ClientTlsKey      string `yaml:"client-tls-key,omitempty" json:"client-tls-key,omitempty"`
	ClientTlsCert     string `yaml:"client-tls-cert,omitempty" json:"client-tls-cert,omitempty"`
}

type KeyValueStoreConfiguration struct {
	Type     string                          `yaml:"type,omitempty" json:"type,omitempty"`
	BasePath string                          `yaml:"base-path,omitempty" json:"base-path,omitempty"`
	Vault    VaultKeyValueStoreConfiguration `yaml:"vault,omitempty" json:"vault,omitempty"`
}

type SqlElement struct {
	Version   string `yaml:"version,omitempty" json:"version,omitempty"`
	Instance  string `yaml:"instance,omitempty" json:"instance,omitempty"`
	Timestamp string `yaml:"timestamp,omitempty" json:"timestamp,omitempty"`
	Name      string `yaml:"name,omitempty" json:"name,omitempty"`
	Kind      string `yaml:"kind,omitempty" json:"kind,omitempty"`
	Artifact  string `yaml:"artifact,omitempty" json:"artifact,omitempty"`
}

type Artifact struct {
	Name     string            `yaml:"name,omitempty" json:"name,omitempty"`
	Password string            `yaml:"password,omitempty" json:"password,omitempty"`
	Kind     string            `yaml:"kind,omitempty" json:"kind,omitempty"`
	Roles    []string          `yaml:"roles,omitempty" json:"roles,omitempty"`
	Metadata map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}
