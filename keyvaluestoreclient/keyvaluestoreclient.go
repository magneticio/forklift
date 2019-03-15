package keyvaluestoreclient

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/models"
)

type KeyValueStoreClient interface {
	GetValue(string) (string, error)
	PutValue(string, string) error
	Delete(string) error
	List(string) ([]string, error)
}

type VaultKeyValueStoreClient struct {
	Url    string
	Token  string
	Params map[string]string
	Client *vaultapi.Client
}

func NewKeyValueStoreClient(config models.KeyValueStoreConfiguration) (KeyValueStoreClient, error) {
	if config.Type == "vault" {
		params := map[string]string{
			"cert":   config.Vault.ClientTlsCert,
			"key":    config.Vault.ClientTlsKey,
			"caCert": config.Vault.ServerTlsCert,
		}
		vaultKVclient, vaultKVclientError := NewVaultKeyValueStoreClient(config.Vault.Url, config.Vault.Token, params)
		if vaultKVclientError != nil {
			return nil, vaultKVclientError
		}
		return vaultKVclient, nil
	}
	return nil, errors.New("Unsupported Key Value Store Client: " + config.Type)
}

func NewVaultKeyValueStoreClient(address string, token string, params map[string]string) (*VaultKeyValueStoreClient, error) {

	logging.Info("Initialising Vault Client with address %v\n", address)

	config, configErr := getConfig(address, params["cert"], params["key"], params["caCert"])
	if configErr != nil {
		logging.Error("Error getting config %v\n", configErr.Error())
		return nil, configErr
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
		logging.Error("Error initialising client %v\n", err.Error())
		return nil, err
	}

	client.SetToken(token)

	return &VaultKeyValueStoreClient{
		Url:    address,
		Token:  token,
		Params: params,
		Client: client,
	}, nil
}

func (c *VaultKeyValueStoreClient) getClient() (*vaultapi.Client, error) {
	// TODO: This will check for token renewal

	logging.Info("Retrievng token")

	token := c.Client.Auth().Token()

	logging.Info("Looking up token")

	tokenSecret, err := token.LookupSelf()
	if err != nil {
		logging.Error("Could not lookup token due to %v", err.Error())
		return c.Client, nil
	}

	logging.Info("Checking if token is renewable")

	renewable, _ := tokenSecret.TokenIsRenewable()
	if renewable {

		ttl, _ := tokenSecret.Data["creation_ttl"].(json.Number).Int64()

		renewPeriod := ttl / 2

		if renewPeriod < 1 {
			return nil, errors.New("Token renew period is invalid")
		}

		logging.Info("Attempting token renewal with ttl %v seconds", renewPeriod)

		_, err := c.Client.Auth().Token().RenewSelf(int(renewPeriod))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to renew token - %v", err))
		}

	} else {
		logging.Info("Token is not renewable")
	}

	return c.Client, nil
}

func (c *VaultKeyValueStoreClient) Delete(keyName string) error {
	logging.Info("Deleting from Vault key %v\n", keyName)
	err := c.DeleteData(fixPath(keyName), nil) // nil mean versions are not defined
	if err != nil {
		logging.Error("Error while deleting from Vault key %v - %v\n", keyName, err.Error())
		return err
	}
	return nil
}

func getConfig(address, cert, key, caCert string) (*vaultapi.Config, error) {
	conf := vaultapi.DefaultConfig()
	conf.Address = address

	tlsConfig := &tls.Config{}
	if cert != "" && key != "" {
		clientCert, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
		tlsConfig.BuildNameToCertificate()
	}

	if caCert != "" {
		ca, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(ca)
		tlsConfig.RootCAs = caCertPool
	}

	conf.HttpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return conf, nil
}

// TODO: use trimPrefix
func fixPath(path string) string {
	if len(path) > 0 && string(path[0]) == "/" {
		return strings.Replace(path, "/", "", 1)
	}
	return path
}

func fixPathSuffix(path string) string {
	return strings.TrimSuffix(path, "/")
}

func secret(keyName string) string {
	return fmt.Sprintf("%s/%s", "secret", keyName)
}

func valueMap(value string) map[string]interface{} {
	return map[string]interface{}{
		"value": value,
	}
}

func (c *VaultKeyValueStoreClient) PutValue(key string, value string) error {
	return c.PutData(fixPath(key), valueMap(value), -1) // -1 means new version
}

func (c *VaultKeyValueStoreClient) GetValue(key string) (string, error) {
	secretValues, err := c.GetData(fixPath(key), 0) //0 means lastest version
	if err != nil {
		return "", err
	}
	value, ok := secretValues["value"].(string)
	if !ok {
		return "", errors.New("Value is not available for key: " + key)
	}
	return value, nil
}

func (c *VaultKeyValueStoreClient) List(key string) ([]string, error) {
	logging.Info("Getting list from Vault with key %v\n", key)
	secretData, listErr := c.ListData(fixPath(key))
	if listErr != nil {
		logging.Error("Error while getting list from Vault with key %v - %v\n", key, listErr.Error())
		return nil, listErr
	}
	if secretData == nil {
		return nil, errors.New("List is not available for path: " + key)
	}
	if val, ok := secretData["keys"]; ok {
		if keysTemp, castOk := val.([]interface{}); castOk {
			keys := make([]string, len(keysTemp))
			for index, k := range keysTemp {
				if str, strCastOk := k.(string); strCastOk {
					keys[index] = fixPathSuffix(str)
				}
			}
			return keys, nil
		}
	}
	return nil, errors.New("List is not available for path: " + key)

}
