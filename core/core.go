package core

import (
	"encoding/json"
	"path"
	"strconv"

	"github.com/magneticio/forklift/keyvaluestoreclient"
	"github.com/magneticio/forklift/logging"
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

func (c *Core) UpsertPolicy(policyContent string) error {
	logging.Info("Upserting policy\n")
	policyAPI := policies.NewPolicyAPI(c.kvClient, c.projectPath)
	return policyAPI.Save(policyContent)
}

func (c *Core) DeleteReleasePolicy(policyName string) error {
	logging.Info("Deleting policy: %v\n", policyName)
	policyAPI := policies.NewPolicyAPI(c.kvClient, c.projectPath)
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
