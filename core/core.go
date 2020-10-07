package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strconv"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/models"
	policies "github.com/magneticio/vamp-policies"
	policiesModel "github.com/magneticio/vamp-policies/policy/domain/model/policy"
	"github.com/magneticio/vamp-policies/policy/interface/api"
	policiesDTO "github.com/magneticio/vamp-policies/policy/interface/persistence/vault/dto"
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

// ListPolicies - lists existing policies
func (c *Core) ListPolicies() ([]models.PolicyView, error) {
	policyAPI := policies.NewPolicyAPI(c.kvClient, c.projectPath)
	apiPolicyViews, err := policyAPI.FindAll()
	if err != nil {
		return nil, fmt.Errorf("cannot list policies: %v", err)
	}

	policyViews := make([]models.PolicyView, len(apiPolicyViews))
	for i, apiPolicyView := range apiPolicyViews {
		policyViews[i] = models.PolicyView{
			ID:   apiPolicyView.PolicyID,
			Name: apiPolicyView.PolicyName,
			Type: string(apiPolicyView.PolicyType),
		}
	}

	return policyViews, nil
}

// GetPolicyString - gets exisiting policy string
func (c *Core) GetPolicyString(policyID uint64) (string, error) {
	policyAPI := policies.NewPolicyAPI(c.kvClient, c.projectPath)
	policyView, err := policyAPI.FindByID(policyID)
	if err != nil {
		return "", fmt.Errorf("cannot get policy: %v", err)
	}
	switch policyView.PolicyType {
	case api.ReleasePolicyType:
		policy, err := policyAPI.GetReleasePolicyByID(policyID)
		if err != nil {
			return "", fmt.Errorf("cannot get release policy: %v", err)
		}
		return getReleasePolicyString(policy)
	case api.ValidationPolicyType:
		policy, err := policyAPI.GetValidationPolicyByID(policyID)
		if err != nil {
			return "", fmt.Errorf("cannot get validation policy: %v", err)
		}
		return getValidationPolicyString(policy)
	}
	return "", fmt.Errorf("unsupported policy type: %v", policyView.PolicyType)
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

// ListClusters - lists existing clusters
func (c *Core) ListClusters() ([]models.ClusterView, error) {
	clustersPath := path.Join(c.projectPath, "clusters")
	clusterIDStrings, err := c.kvClient.List(clustersPath)
	if err != nil {
		return nil, fmt.Errorf("cannot list clusters: %v", err)
	}
	clusterIDs := make([]uint64, len(clusterIDStrings))
	for i, clusterIDString := range clusterIDStrings {
		clusterID, err := strconv.ParseUint(clusterIDString, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("found cluster with invalid id '%s'", clusterIDString)
		}
		clusterIDs[i] = clusterID
	}

	clusters := make([]models.ClusterView, len(clusterIDs))

	for i, clusterID := range clusterIDs {
		releaseAgentConfigKey := c.getReleaseAgentConfigKey(clusterID)
		releaseAgentConfig, exists, err := c.getReleaseAgentConfig(releaseAgentConfigKey)
		if err != nil {
			return nil, fmt.Errorf("cannot get cluster '%d': %v", clusterID, err)
		}
		if !exists {
			return nil, fmt.Errorf("cluster config for cluster '%d' does not exist", clusterID)
		}
		clusters[i] = models.ClusterView{
			ID:                   clusterID,
			NatsChannel:          releaseAgentConfig.NatsChannel,
			OptimiserNatsChannel: releaseAgentConfig.OptimiserNatsChannel,
		}
	}

	return clusters, nil
}

// GetCluster - gets existing cluster
func (c *Core) GetCluster(clusterID uint64) (*models.ClusterView, error) {
	releaseAgentConfigKey := c.getReleaseAgentConfigKey(clusterID)
	releaseAgentConfig, exists, err := c.getReleaseAgentConfig(releaseAgentConfigKey)
	if err != nil {
		return nil, fmt.Errorf("cannot get cluster: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("cluster config does not exist")
	}
	return &models.ClusterView{
		ID:                   clusterID,
		NatsChannel:          releaseAgentConfig.NatsChannel,
		OptimiserNatsChannel: releaseAgentConfig.OptimiserNatsChannel,
		NatsToken:            releaseAgentConfig.NatsToken,
	}, nil
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

// ListApplications - lists existing applications
func (c *Core) ListApplications() ([]models.ApplicationView, error) {
	if c.clusterID == nil {
		return nil, fmt.Errorf("cluster id must be provided")
	}
	releaseAgentConfigKey := c.getReleaseAgentConfigKey(uint64(*c.clusterID))
	releaseAgentConfig, exists, err := c.getReleaseAgentConfig(releaseAgentConfigKey)
	if err != nil {
		return nil, fmt.Errorf("cannot find Release Agent config: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("Release Agent config does not exist")
	}
	applications := make([]models.ApplicationView, len(releaseAgentConfig.K8SNamespaceToApplicationID))
	i := 0
	for namespace, applicationID := range releaseAgentConfig.K8SNamespaceToApplicationID {
		applications[i] = models.ApplicationView{
			ID:        applicationID,
			Namespace: namespace,
		}
		i++
	}

	return applications, nil
}

// GetApplication - gets existing application
func (c *Core) GetApplication(applicationID uint64) (*models.ApplicationView, error) {
	if c.clusterID == nil {
		return nil, fmt.Errorf("cluster id must be provided")
	}
	releaseAgentConfigKey := c.getReleaseAgentConfigKey(uint64(*c.clusterID))
	releaseAgentConfig, exists, err := c.getReleaseAgentConfig(releaseAgentConfigKey)
	if err != nil {
		return nil, fmt.Errorf("cannot find Release Agent config: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("Release Agent config does not exist")
	}

	for namespace, configApplicationID := range releaseAgentConfig.K8SNamespaceToApplicationID {
		if applicationID == configApplicationID {
			return &models.ApplicationView{
				ID:        applicationID,
				Namespace: namespace,
			}, nil
		}
	}

	return nil, fmt.Errorf("application '%d' not found", applicationID)
}

// PutServiceConfig - puts service to key value store
func (c *Core) PutServiceConfig(serviceID, applicationID uint64, serviceConfigText string) error {
	if c.clusterID == nil {
		return fmt.Errorf("cluster id must be provided")
	}

	var serviceConfig models.ServiceConfig
	if err := json.Unmarshal([]byte(serviceConfigText), &serviceConfig); err != nil {
		return fmt.Errorf("cannot deserialize service config: %v", err)
	}
	if err := models.NewValidateDTO()(serviceConfig); err != nil {
		return fmt.Errorf("service config validation failed: %v", err)
	}
	if err := serviceConfig.Validate(); err != nil {
		return fmt.Errorf("service config validation failed: %v", err)
	}

	serviceConfigKey := c.getServiceConfigKey(*c.clusterID, applicationID, serviceID)

	return c.kvClient.Put(serviceConfigKey, serviceConfigText)
}

// DeleteServiceConfig - deletes service config from key value store
func (c *Core) DeleteServiceConfig(serviceID, applicationID uint64) error {
	if c.clusterID == nil {
		return fmt.Errorf("cluster id must be provided")
	}

	serviceConfigKey := c.getServiceConfigKey(*c.clusterID, applicationID, serviceID)
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

func (c *Core) getServiceConfigKey(clusterID, applicationID, serviceID uint64) string {
	return path.Join(
		c.getClusterPath(clusterID),
		"applications",
		strconv.FormatUint(applicationID, 10),
		"service-configs",
		strconv.FormatUint(serviceID, 10),
	)
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

func getReleasePolicyString(policy *policiesModel.Policy) (string, error) {
	policyDTO, err := policiesDTO.ToPolicyDTO(policy)
	if err != nil {
		return "", fmt.Errorf("cannot get release policy string: %v", err)
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(policyDTO)
	if err != nil {
		return "", fmt.Errorf("cannot marshal release policy: %v", err)
	}
	return buf.String(), nil
}

func getValidationPolicyString(policy *policiesModel.ValidationPolicy) (string, error) {
	policyDTO, err := policiesDTO.ToValidationPolicyDTO(policy)
	if err != nil {
		return "", fmt.Errorf("cannot get validation policy string: %v", err)
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(policyDTO)
	if err != nil {
		return "", fmt.Errorf("cannot marshal validation policy: %v", err)
	}
	return buf.String(), nil
}
