package keyvaluestoreclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	vaultapi "github.com/hashicorp/vault/api"
)

type KeyValueStoreClient interface {
	Get(string) (map[string]interface{}, error)
	Put(string, map[string]interface{}) error
	Delete(string) error
}

type VaultKeyValueStoreClient struct {
	Url    string
	Token  string
	Params map[string]string
	Client *vaultapi.Client
}

func NewVaultKeyValueStoreClient(address string, token string, params map[string]string) (*VaultKeyValueStoreClient, error) {

	config, configErr := getConfig(address, params["cert"], params["key"], params["caCert"])
	if configErr != nil {
		return nil, configErr
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
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

func (c *VaultKeyValueStoreClient) GetClient() *vaultapi.Client {
	// This will check for token renewal
	return c.Client
}

func (c *VaultKeyValueStoreClient) Put(keyName string, secretData map[string]interface{}) error {
	_, err := c.GetClient().Logical().Write(secret(keyName), secretData)
	if err != nil {
		return err
	}
	return nil
}

func (c *VaultKeyValueStoreClient) Get(keyName string) (map[string]interface{}, error) {
	secretValues, err := c.GetClient().Logical().Read(secret(keyName))
	if err != nil {
		return nil, nil
	}
	return secretValues.Data, nil
}

func (c *VaultKeyValueStoreClient) Delete(keyName string) error {
	_, err := c.GetClient().Logical().Delete(secret(keyName))
	if err != nil {
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

func secret(keyName string) string {
	return fmt.Sprintf("%s/%s", "secret", keyName)
}
