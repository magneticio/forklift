package models

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

// ForkliftConfiguration - configuration built from config file, environment variables and flags
type ForkliftConfiguration struct {
	ProjectID                      *uint64 `json:"project,omitempty"`
	ClusterID                      *uint64 `json:"cluster,omitempty"`
	KeyValueStoreURL               string  `json:"key-value-store-url,omitempty"`
	KeyValueStoreToken             string  `json:"key-value-store-token,omitempty"`
	KeyValueStoreBasePath          string  `json:"key-value-store-base-path,omitempty"`
	KeyValueStoreServerTLSCert     string  `json:"key-value-store-server-tls-cert,omitempty"`
	KeyValueStoreClientTLSKey      string  `json:"key-value-store-client-tls-key,omitempty"`
	KeyValueStoreClientTLSCert     string  `json:"key-value-store-client-tls-cert,omitempty"`
	KeyValueStoreKvMode            string  `json:"key-value-store-kv-mode,omitempty"`
	KeyValueStoreFallbackKvVersion string  `json:"key-value-store-fallback-kv-version,omitempty"`
}

type tmpForkliftConfiguration struct {
	ProjectID                      string `yaml:"project,omitempty"`
	ClusterID                      string `yaml:"cluster,omitempty"`
	KeyValueStoreURL               string `yaml:"key-value-store-url,omitempty"`
	KeyValueStoreToken             string `yaml:"key-value-store-token,omitempty"`
	KeyValueStoreBasePath          string `yaml:"key-value-store-base-path,omitempty"`
	KeyValueStoreServerTLSCert     string `yaml:"key-value-store-server-tls-cert,omitempty"`
	KeyValueStoreClientTLSKey      string `yaml:"key-value-store-client-tls-key,omitempty"`
	KeyValueStoreClientTLSCert     string `yaml:"key-value-store-client-tls-cert,omitempty"`
	KeyValueStoreKvMode            string `yaml:"key-value-store-kv-mode,omitempty"`
	KeyValueStoreFallbackKvVersion string `yaml:"key-value-store-fallback-kv-version,omitempty"`
}

// UnmarshalYAML - implements the Unmarshaler interface of the yaml pkg
func (conf *ForkliftConfiguration) UnmarshalYAML(node *yaml.Node) error {
	var tmp tmpForkliftConfiguration
	if err := node.Decode(&tmp); err != nil {
		return err
	}

	projectID, err := getUint64FromString(tmp.ProjectID)
	if err != nil {
		return fmt.Errorf("invalid project id: %v", err)
	}

	clusterID, err := getUint64FromString(tmp.ClusterID)
	if err != nil {
		return fmt.Errorf("invalid cluster id: %v", err)
	}

	*conf = ForkliftConfiguration{
		ProjectID:                      projectID,
		ClusterID:                      clusterID,
		KeyValueStoreBasePath:          tmp.KeyValueStoreBasePath,
		KeyValueStoreClientTLSCert:     tmp.KeyValueStoreClientTLSCert,
		KeyValueStoreClientTLSKey:      tmp.KeyValueStoreClientTLSKey,
		KeyValueStoreFallbackKvVersion: tmp.KeyValueStoreFallbackKvVersion,
		KeyValueStoreKvMode:            tmp.KeyValueStoreKvMode,
		KeyValueStoreServerTLSCert:     tmp.KeyValueStoreServerTLSCert,
		KeyValueStoreToken:             tmp.KeyValueStoreToken,
		KeyValueStoreURL:               tmp.KeyValueStoreURL,
	}

	return nil
}

func getUint64FromString(valueText string) (*uint64, error) {
	if valueText == "" {
		return nil, nil
	}
	value, err := strconv.ParseUint(valueText, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("value must be a natural number")
	}
	return &value, nil
}

// VaultKeyValueStoreConfiguration - Vault configuration
type VaultKeyValueStoreConfiguration struct {
	URL               string `yaml:"url,omitempty" json:"url,omitempty"`
	Token             string `yaml:"token,omitempty" json:"token,omitempty"`
	KvMode            string `yaml:"kv-mode,omitempty" json:"kv-mode,omitempty"`
	FallbackKvVersion string `yaml:"fallback-kv-version,omitempty" json:"fallback-kv-version,omitempty"`
	ServerTLSCert     string `yaml:"server-tls-cert,omitempty" json:"server-tls-cert,omitempty"`
	ClientTLSKey      string `yaml:"client-tls-key,omitempty" json:"client-tls-key,omitempty"`
	ClientTLSCert     string `yaml:"client-tls-cert,omitempty" json:"client-tls-cert,omitempty"`
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
	DefaultPolicyID *uint64                     `json:"default_policy_id"`
	PatchPolicyID   *uint64                     `json:"patch_policy_id"`
	MinorPolicyID   *uint64                     `json:"minor_policy_id"`
	MajorPolicyID   *uint64                     `json:"major_policy_id"`
	IngressRules    []*ServiceConfigIngressRule `json:"ingress_rules"`
	IsHeadless      bool                        `json:"headless"`
}

// Validate - additional validation of ServiceConfig structure
// that cannot be achieved using go-playgroud validator
func (sc ServiceConfig) Validate() error {
	if (sc.MajorPolicyID == nil || sc.MinorPolicyID == nil || sc.PatchPolicyID == nil) && sc.DefaultPolicyID == nil {
		return fmt.Errorf("DefaultPolicyID should be defined if any of other policies is not defined")
	}
	return nil
}

// ServiceConfigIngressRule - service config ingress rule for Release Agent
type ServiceConfigIngressRule struct {
	Domain        string `json:"domain" validate:"required,min=4"`
	TLSSecretName string `json:"tls_secret_name"`
	Path          string `json:"path" validate:"required,min=1"`
	Port          *int64 `json:"port" validate:"required"`
}

// ApplicationView - view used as an output for list and show commands
type ApplicationView struct {
	ID        uint64 `yaml:"id"`
	Namespace string `yaml:"namespace"`
}

// ClusterView - view used as an output for list and show commands
type ClusterView struct {
	ID                   uint64 `yaml:"id"`
	NatsChannel          string `yaml:"nats-channel"`
	NatsToken            string `yaml:"nats-token,omitempty"`
	OptimiserNatsChannel string `yaml:"optimiser-nats-channel"`
}
