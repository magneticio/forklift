package models

type ForkliftConfiguration struct {
	ProjectID                      uint64 `yaml:"project-id,omitempty" json:"project-id,omitempty"`
	ClusterID                      uint64 `yaml:"cluster-id,omitempty" json:"cluster-id,omitempty"`
	KeyValueStoreUrL               string `yaml:"key-value-store-url,omitempty" json:"key-value-store-url,omitempty"`
	KeyValueStoreToken             string `yaml:"key-value-store-token,omitempty" json:"key-value-store-token,omitempty"`
	KeyValueStoreBasePath          string `yaml:"key-value-store-base-path,omitempty" json:"key-value-store-base-path,omitempty"`
	KeyValueStoreServerTlsCert     string `yaml:"key-value-store-server-tls-cert,omitempty" json:"key-value-store-server-tls-cert,omitempty"`
	KeyValueStoreClientTlsKey      string `yaml:"key-value-store-client-tls-key,omitempty" json:"key-value-store-client-tls-key,omitempty"`
	KeyValueStoreClientTlsCert     string `yaml:"key-value-store-client-tls-cert,omitempty" json:"key-value-store-client-tls-cert,omitempty"`
	KeyValueStoreKvMode            string `yaml:"key-value-store-kv-mode,omitempty" json:"key-value-store-kv-mode,omitempty"`
	KeyValueStoreFallbackKvVersion string `yaml:"key-value-store-fallback-kv-version,omitempty" json:"key-value-store-fallback-kv-version,omitempty"`
}

type VaultKeyValueStoreConfiguration struct {
	Url               string `yaml:"url,omitempty" json:"url,omitempty"`
	Token             string `yaml:"token,omitempty" json:"token,omitempty"`
	KvMode            string `yaml:"kv-mode,omitempty" json:"kv-mode,omitempty"`
	FallbackKvVersion string `yaml:"fallback-kv-version,omitempty" json:"fallback-kv-version,omitempty"`
	ServerTlsCert     string `yaml:"server-tls-cert,omitempty" json:"server-tls-cert,omitempty"`
	ClientTlsKey      string `yaml:"client-tls-key,omitempty" json:"client-tls-key,omitempty"`
	ClientTlsCert     string `yaml:"client-tls-cert,omitempty" json:"client-tls-cert,omitempty"`
}
