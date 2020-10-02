package core

import (
	"encoding/json"
	"fmt"
	"path"
	"strconv"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/models"
	policies "github.com/magneticio/vamp-policies"
)

type Core struct {
	kvClient    keyvaluestoreclient.KeyValueStoreClient
	projectPath string
	clusterID   *uint64
}

func NewCore(conf models.ForkliftConfiguration) (*Core, error) {
	if conf.ProjectID == nil {
		return nil, fmt.Errorf("project id must be provided")
	}
	projectPath := path.Join(conf.KeyValueStoreBasePath, "projects", strconv.FormatUint(*conf.ProjectID, 10))
	config := models.VaultKeyValueStoreConfiguration{
		URL:               conf.KeyValueStoreURL,
		Token:             conf.KeyValueStoreToken,
		ServerTLSCert:     conf.KeyValueStoreServerTLSCert,
		ClientTLSCert:     conf.KeyValueStoreClientTLSCert,
		ClientTLSKey:      conf.KeyValueStoreClientTLSKey,
		KvMode:            conf.KeyValueStoreKvMode,
		FallbackKvVersion: conf.KeyValueStoreFallbackKvVersion,
	}
	kvClient, err := keyvaluestoreclient.NewKeyValueStoreClient(config)
	if err != nil {
		return nil, err
	}

	return &Core{
		kvClient:    kvClient,
		projectPath: projectPath,
		clusterID:   conf.ClusterID,
	}, nil
}

// PutPolicy - puts policy to key value store
func (c *Core) PutPolicy(policyID uint64, policyContent string) error {
	policyAPI := policies.NewPolicyAPI(c.kvClient, c.projectPath)
	return policyAPI.Save(strconv.FormatUint(policyID, 10), policyContent)
}

// DeletePolicy - deletes policy from key value store
func (c *Core) DeletePolicy(policyID uint64) error {
	policyAPI := policies.NewPolicyAPI(c.kvClient, c.projectPath)
	policyKey := strconv.FormatUint(policyID, 10)
	_, err := policyAPI.Find(policyKey)
	if err != nil {
		return fmt.Errorf("cannot find policy: %v", err)
	}
	return policyAPI.Delete(policyKey)
}

// PutReleasePlan - puts release plan to key value store
func (c *Core) PutReleasePlan(serviceID uint64, serviceVersion string, releasePlanContent string) error {
	releasePlanKey := c.getReleasePlanKey(serviceID, serviceVersion)
	return c.kvClient.Put(releasePlanKey, releasePlanContent)
}

// DeleteReleasePlan - deletes release plan from key value store
func (c *Core) DeleteReleasePlan(serviceID uint64, serviceVersion string) error {
	releasePlanKey := c.getReleasePlanKey(serviceID, serviceVersion)
	exists, err := c.kvClient.Exists(releasePlanKey)
	if err != nil {
		return fmt.Errorf("cannot find release plan: %v", err)
	}
	if !exists {
		return fmt.Errorf("release plan does not exist")
	}
	return c.kvClient.Delete(releasePlanKey)
}

// PutReleaseAgentConfig - puts Release Agent config to key value store
func (c *Core) PutReleaseAgentConfig(clusterID uint64, natsChannelName, optimiserNatsChannelName, natsToken string) error {
	if natsChannelName == "" {
		return fmt.Errorf("NATS channel name must not be empty")
	}
	releaseAgentConfigKey := c.getReleaseAgentConfigKey(clusterID)
	existingReleaseAgentConfig, exists, err := c.getReleaseAgentConfig(releaseAgentConfigKey)
	if err != nil {
		return err
	}

	var releaseAgentConfig models.ReleaseAgentConfig
	if exists {
		releaseAgentConfig = models.ReleaseAgentConfig{
			NatsChannel:                 natsChannelName,
			NatsToken:                   natsToken,
			OptimiserNatsChannel:        optimiserNatsChannelName,
			K8SNamespaceToApplicationID: existingReleaseAgentConfig.K8SNamespaceToApplicationID,
		}
	} else {
		releaseAgentConfig = models.ReleaseAgentConfig{
			NatsChannel:                 natsChannelName,
			NatsToken:                   natsToken,
			OptimiserNatsChannel:        optimiserNatsChannelName,
			K8SNamespaceToApplicationID: make(map[string]uint64),
		}
	}

	return c.saveReleaseAgentConfig(releaseAgentConfigKey, releaseAgentConfig)
}

// DeleteReleaseAgentConfig - deletes Release Agent config from key value store
func (c *Core) DeleteReleaseAgentConfig(clusterID uint64) error {
	releaseAgentConfigKey := c.getReleaseAgentConfigKey(clusterID)
	_, exists, err := c.getReleaseAgentConfig(releaseAgentConfigKey)
	if err != nil {
		return fmt.Errorf("cannot find Release Agent config: %v", err)
	}
	if !exists {
		return fmt.Errorf("Release Agent config does not exist")
	}

	return c.kvClient.Delete(releaseAgentConfigKey)
}

// PutApplication - puts application to existing Release Agent config
func (c *Core) PutApplication(applicationID uint64, namespace string) error {
	putApplication := func(releaseAgentConfig *models.ReleaseAgentConfig) {
		for configNamespace, configApplicationID := range releaseAgentConfig.K8SNamespaceToApplicationID {
			if configApplicationID == applicationID {
				delete(releaseAgentConfig.K8SNamespaceToApplicationID, configNamespace)
			}
		}
		releaseAgentConfig.K8SNamespaceToApplicationID[namespace] = applicationID
	}

	return c.onReleaseAgentConfig(putApplication)
}

// DeleteApplication - deletes application from Release Agent config
func (c *Core) DeleteApplication(applicationID uint64) error {
	deleteApplication := func(releaseAgentConfig *models.ReleaseAgentConfig) {
		for configNamespace, configApplicationID := range releaseAgentConfig.K8SNamespaceToApplicationID {
			if configApplicationID == applicationID {
				delete(releaseAgentConfig.K8SNamespaceToApplicationID, configNamespace)
			}
		}
	}

	return c.onReleaseAgentConfig(deleteApplication)
}

// PutServiceConfig - puts service to key value store
func (c *Core) PutServiceConfig(serviceID uint64, serviceConfigText string) error {
	var serviceConfig models.ServiceConfig
	if err := json.Unmarshal([]byte(serviceConfigText), &serviceConfig); err != nil {
		return fmt.Errorf("cannot deserialize service config: %v", err)
	}
	if err := models.NewValidateDTO()(serviceConfig); err != nil {
		return fmt.Errorf("service config validation failed: %v", err)
	}
	serviceConfigKey := c.getServiceConfigKey(serviceID)
	return c.kvClient.Put(serviceConfigKey, serviceConfigText)
}

// DeleteServiceConfig - deletes service config from key value store
func (c *Core) DeleteServiceConfig(serviceID uint64) error {
	serviceConfigKey := c.getServiceConfigKey(serviceID)
	exists, err := c.kvClient.Exists(serviceConfigKey)
	if err != nil {
		return fmt.Errorf("cannot find service config: %v", err)
	}
	if !exists {
		return fmt.Errorf("service config does not exist")
	}
	return c.kvClient.Delete(serviceConfigKey)
}

func (c *Core) onReleaseAgentConfig(apply func(*models.ReleaseAgentConfig)) error {
	if c.clusterID == nil {
		return fmt.Errorf("cluster id must be provided")
	}
	releaseAgentConfigKey := c.getReleaseAgentConfigKey(uint64(*c.clusterID))
	releaseAgentConfig, exists, err := c.getReleaseAgentConfig(releaseAgentConfigKey)
	if err != nil {
		return fmt.Errorf("cannot find Release Agent config: %v", err)
	}
	if !exists {
		return fmt.Errorf("Release Agent config does not exist. Please create cluster first")
	}

	apply(releaseAgentConfig)

	return c.saveReleaseAgentConfig(releaseAgentConfigKey, *releaseAgentConfig)
}

func (c *Core) getClusterPath(clusterID uint64) string {
	return path.Join(c.projectPath, "clusters", strconv.FormatUint(clusterID, 10))
}

func (c *Core) getReleasePlanKey(serviceID uint64, serviceVersion string) string {
	return path.Join(c.projectPath, "release-plans", strconv.FormatUint(serviceID, 10), serviceVersion)
}

func (c *Core) getReleaseAgentConfigKey(clusterID uint64) string {
	return path.Join(c.getClusterPath(clusterID), "release-agent-config")
}

func (c *Core) getServiceConfigKey(serviceID uint64) string {
	return path.Join(c.projectPath, "services", strconv.FormatUint(serviceID, 10))
}

func (c *Core) getReleaseAgentConfig(releaseAgentConfigKey string) (*models.ReleaseAgentConfig, bool, error) {
	configExists, err := c.kvClient.Exists(releaseAgentConfigKey)
	if err != nil {
		return nil, false, fmt.Errorf("cannot check if Release Agent config exists in Vault: %v", err)
	}

	if configExists {
		releaseAgentConfigContent, err := c.kvClient.Get(releaseAgentConfigKey)
		if err != nil {
			return nil, false, fmt.Errorf("cannot get existing Release Agent config from Vault: %v", err)
		}
		var releaseAgentConfig models.ReleaseAgentConfig
		if err = json.Unmarshal([]byte(releaseAgentConfigContent), &releaseAgentConfig); err != nil {
			return nil, false, fmt.Errorf("cannot deserialize existing Release Agent config: %v", err)
		}
		return &releaseAgentConfig, true, nil
	}

	return nil, false, nil
}

func (c *Core) saveReleaseAgentConfig(releaseAgentConfigKey string, releaseAgentConfig models.ReleaseAgentConfig) error {
	releaseAgentConfigBytes, err := json.Marshal(releaseAgentConfig)
	if err != nil {
		return fmt.Errorf("cannot serialize Release Agent config: %v", err)
	}

	return c.kvClient.Put(releaseAgentConfigKey, string(releaseAgentConfigBytes))
}
