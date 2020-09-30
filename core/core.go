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
	clusterID   uint64
}

func NewCore(conf models.ForkliftConfiguration) (*Core, error) {
	projectPath := path.Join(conf.KeyValueStoreBasePath, "projects", strconv.FormatUint(conf.ProjectID, 10))
	config := models.VaultKeyValueStoreConfiguration{
		Url:               conf.KeyValueStoreUrL,
		Token:             conf.KeyValueStoreToken,
		ServerTlsCert:     conf.KeyValueStoreServerTlsCert,
		ClientTlsCert:     conf.KeyValueStoreClientTlsCert,
		ClientTlsKey:      conf.KeyValueStoreClientTlsKey,
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

// UpsertPolicy - upserts policy in key value store
func (c *Core) UpsertPolicy(policyID uint64, policyContent string) error {
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

// UpsertReleasePlan - upserts release plan in key value store
func (c *Core) UpsertReleasePlan(serviceID uint64, serviceVersion string, releasePlanContent string) error {
	releasePlanKey := c.getReleasePlansPath(serviceID, serviceVersion)
	return c.kvClient.Put(releasePlanKey, releasePlanContent)
}

// DeleteReleasePlan - deletes release plan from key value store
func (c *Core) DeleteReleasePlan(serviceID uint64, serviceVersion string) error {
	releasePlanKey := c.getReleasePlansPath(serviceID, serviceVersion)
	exists, err := c.kvClient.Exists(releasePlanKey)
	if err != nil {
		return fmt.Errorf("cannot find release plan: %v", err)
	}
	if !exists {
		return fmt.Errorf("release plan does not exist")
	}
	return c.kvClient.Delete(releasePlanKey)
}

// UpsertReleaseAgentConfig - upserts Release Agent config in key value store
func (c *Core) UpsertReleaseAgentConfig(clusterID uint64, natsChannelName, optimiserNatsChannelName, natsToken string) error {
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

	releaseAgentConfigBytes, err := json.Marshal(releaseAgentConfig)
	if err != nil {
		return fmt.Errorf("cannot serialize Release Agent config: %v", err)
	}

	return c.kvClient.Put(releaseAgentConfigKey, string(releaseAgentConfigBytes))
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

func (c *Core) getReleasePlansPath(serviceID uint64, serviceVersion string) string {
	return path.Join(c.projectPath, "release-plans", strconv.FormatUint(serviceID, 10), serviceVersion)
}

func (c *Core) getClusterPath(clusterID uint64) string {
	return path.Join(c.projectPath, "clusters", strconv.FormatUint(clusterID, 10))
}

func (c *Core) getReleaseAgentConfigKey(clusterID uint64) string {
	return path.Join(c.getClusterPath(clusterID), "release-agent-config")
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
