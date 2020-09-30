package core

import (
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
	return policyAPI.Delete(strconv.FormatUint(policyID, 10))
}

// UpsertReleasePlan - upserts release plan in key value store
func (c *Core) UpsertReleasePlan(serviceID uint64, serviceVersion string, releasePlanContent string) error {
	releasePlanKey := c.getReleasePlansPath(serviceID, serviceVersion)
	return c.kvClient.Put(releasePlanKey, releasePlanContent)
}

// DeleteReleasePlan - deletes release plan from key value store
func (c *Core) DeleteReleasePlan(serviceID uint64, serviceVersion string) error {
	releasePlanKey := c.getReleasePlansPath(serviceID, serviceVersion)
	return c.kvClient.Delete(releasePlanKey)
}

func (c *Core) getReleasePlansPath(serviceID uint64, serviceVersion string) string {
	return path.Join(c.projectPath, "release-plans", strconv.FormatUint(serviceID, 10), serviceVersion)
}
