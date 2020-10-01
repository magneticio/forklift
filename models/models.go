package models

import (
	"fmt"
	"strconv"
)

// NullableUint64 - custom uint64 type to be able to read it's value from both viper config and cobra flag
type NullableUint64 uint64

// Set - sets the value of the NullableUint64 type
func (val *NullableUint64) Set(v string) error {
	if v == "" {
		return nil
	}
	parsedValue, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return fmt.Errorf("value must be a natural number")
	}
	u := NullableUint64(parsedValue)
	val = &u
	return nil
}

// String - gets the string value of the NullableUint64 type
func (val *NullableUint64) String() string {
	if val == nil {
		return ""
	}
	return strconv.FormatUint(uint64(*val), 10)
}

// Type - gets the type text for NullableUint64 type
func (val *NullableUint64) Type() string {
	return "NullableUint64"
}

// UnmarshalYAML - implements the Unmarshaler interface of the yaml pkg
func (val *NullableUint64) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var valueText string
	err := unmarshal(&valueText)
	if err != nil {
		return err
	}
	if valueText == "" {
		return nil
	}
	value, err := strconv.ParseUint(valueText, 10, 64)
	if err != nil {
		return fmt.Errorf("value must be a natural number")
	}
	u := NullableUint64(value)
	val = &u
	fmt.Printf("val3: %v\n", val)
	return nil
}

// ForkliftConfiguration - configuration built from config file, environment variables and flags
type ForkliftConfiguration struct {
	ProjectID                      *NullableUint64 `yaml:"project,omitempty" json:"project,omitempty"`
	ClusterID                      *NullableUint64 `yaml:"cluster,omitempty" json:"cluster,omitempty"`
	KeyValueStoreUrL               string          `yaml:"key-value-store-url,omitempty" json:"key-value-store-url,omitempty"`
	KeyValueStoreToken             string          `yaml:"key-value-store-token,omitempty" json:"key-value-store-token,omitempty"`
	KeyValueStoreBasePath          string          `yaml:"key-value-store-base-path,omitempty" json:"key-value-store-base-path,omitempty"`
	KeyValueStoreServerTlsCert     string          `yaml:"key-value-store-server-tls-cert,omitempty" json:"key-value-store-server-tls-cert,omitempty"`
	KeyValueStoreClientTlsKey      string          `yaml:"key-value-store-client-tls-key,omitempty" json:"key-value-store-client-tls-key,omitempty"`
	KeyValueStoreClientTlsCert     string          `yaml:"key-value-store-client-tls-cert,omitempty" json:"key-value-store-client-tls-cert,omitempty"`
	KeyValueStoreKvMode            string          `yaml:"key-value-store-kv-mode,omitempty" json:"key-value-store-kv-mode,omitempty"`
	KeyValueStoreFallbackKvVersion string          `yaml:"key-value-store-fallback-kv-version,omitempty" json:"key-value-store-fallback-kv-version,omitempty"`
}

// VaultKeyValueStoreConfiguration - Vault configuration
type VaultKeyValueStoreConfiguration struct {
	Url               string `yaml:"url,omitempty" json:"url,omitempty"`
	Token             string `yaml:"token,omitempty" json:"token,omitempty"`
	KvMode            string `yaml:"kv-mode,omitempty" json:"kv-mode,omitempty"`
	FallbackKvVersion string `yaml:"fallback-kv-version,omitempty" json:"fallback-kv-version,omitempty"`
	ServerTlsCert     string `yaml:"server-tls-cert,omitempty" json:"server-tls-cert,omitempty"`
	ClientTlsKey      string `yaml:"client-tls-key,omitempty" json:"client-tls-key,omitempty"`
	ClientTlsCert     string `yaml:"client-tls-cert,omitempty" json:"client-tls-cert,omitempty"`
}

// ReleaseAgentConfig - config for Release Agent
type ReleaseAgentConfig struct {
	NatsChannel                 string            `json:"nats_channel"`
	NatsToken                   string            `json:"nats_token"`
	K8SNamespaceToApplicationID map[string]uint64 `json:"applications"`
	OptimiserNatsChannel        string            `json:"optimiser_nats_channel"`
}

// ServiceConfig - service config for Release Agent
type ServiceConfig struct {
	ApplicationID   *uint64                     `json:"application_id" validate:"required"`
	ServiceID       *uint64                     `json:"service_id" validate:"required"`
	K8SNamespace    string                      `json:"k8s_namespace" validate:"required,min=1"`
	K8sLabels       map[string]string           `json:"k8s_labels" validate:"required,min=1"`
	VersionSelector string                      `json:"version_selector" validate:"required,min=1"`
	DefaultPolicyID *uint64                     `json:"default_policy_id" validate:"required_without_all=PatchPolicyID MinorPolicyID MajorPolicyID,ne=0"`
	PatchPolicyID   *uint64                     `json:"patch_policy_id"`
	MinorPolicyID   *uint64                     `json:"minor_policy_id"`
	MajorPolicyID   *uint64                     `json:"major_policy_id"`
	IngressRules    []*ServiceConfigIngressRule `json:"ingress_rules"`
	IsHeadless      bool                        `json:"headless"`
}

// ServiceConfigIngressRule - service config ingress rule for Release Agent
type ServiceConfigIngressRule struct {
	Domain        string `json:"domain" validate:"required,min=4"`
	TLSSecretName string `json:"tls_secret_name"`
	Path          string `json:"path" validate:"required,min=1"`
	Port          *int64 `json:"port" validate:"required"`
}
