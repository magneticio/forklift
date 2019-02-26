package keyvaluestoreclient_test

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/magneticio/forklift/keyvaluestoreclient"
)

// testVaultServer creates a test vault cluster and returns a configured API
// client and closer function.
func testVaultServer(t testing.TB) (*api.Client, map[string]string, func()) {
	t.Helper()

	client, _, params, closer := testVaultServerUnseal(t)
	return client, params, closer
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnseal(t testing.TB) (*api.Client, []string, map[string]string, func()) {
	t.Helper()

	return testVaultServerCoreConfig(t, &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		// Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
		LogicalBackends: map[string]logical.Factory{
			"database":       database.Factory,
			"generic-leased": vault.LeasedPassthroughBackendFactory,
			"pki":            pki.Factory,
			"transit":        transit.Factory,
		},
		BuiltinRegistry: builtinplugins.Registry,
	})
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper.
func testVaultServerCoreConfig(t testing.TB, coreConfig *vault.CoreConfig) (*api.Client, []string, map[string]string, func()) {
	t.Helper()

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Convert the unseal keys to base64 encoded, since these are how the user
	// will get them.
	unsealKeys := make([]string, len(cluster.BarrierKeys))
	for i := range unsealKeys {
		unsealKeys[i] = base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i])
	}

	params := map[string]string{
		"cert":   cluster.TempDir + "/ca_cert.pem",
		"key":    cluster.TempDir + "/ca_key.pem",
		"caCert": cluster.TempDir + "/ca_cert.pem",
	}

	return client, unsealKeys, params, func() { defer cluster.Cleanup() }
}

// testVaultServerBad creates an http server that returns a 500 on each request
// to simulate failures.
func testVaultServerBad(t testing.TB) (*api.Client, func()) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	server := &http.Server{
		Addr: "127.0.0.1:0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
		}),
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	client, err := api.NewClient(&api.Config{
		Address: "http://" + listener.Addr().String(),
	})
	if err != nil {
		t.Fatal(err)
	}

	return client, func() {
		ctx, done := context.WithTimeout(context.Background(), 5*time.Second)
		defer done()

		server.Shutdown(ctx)
	}
}

func TestVaultKeyVauleStoreClientValueWrappers(t *testing.T) {
	client, params, closer := testVaultServer(t)
	defer closer()
	url := client.Address()
	token := client.Token()

	// TODO: Update this test when certs is added to the config

	vaultKeyValueStoreClient, clientErr := keyvaluestoreclient.NewVaultKeyValueStoreClient(url, token, params)
	assert.Nil(t, clientErr)
	assert.NotNil(t, vaultKeyValueStoreClient)

	key := "secret/test2"
	value := "value2"
	putErr := vaultKeyValueStoreClient.PutValue(key, value)
	assert.Nil(t, putErr)

	value2, getErr := vaultKeyValueStoreClient.GetValue(key)
	assert.Nil(t, getErr)
	assert.Equal(t, value, value2)

	deleteErr := vaultKeyValueStoreClient.Delete(key)
	assert.Nil(t, deleteErr)

	basePath := "/secret/vamp"
	key1 := basePath + "/key1/subpath"
	valueKey1 := "test1"
	putErr1 := vaultKeyValueStoreClient.PutValue(key1, valueKey1)
	assert.Nil(t, putErr1)
	key2 := basePath + "/key2/subpath"
	valueKey2 := "test2"
	putErr2 := vaultKeyValueStoreClient.PutValue(key2, valueKey2)
	assert.Nil(t, putErr2)

	list, listError := vaultKeyValueStoreClient.List(basePath)
	assert.Nil(t, listError)
	listExpected := []string{"key1", "key2"}
	assert.Equal(t, listExpected, list)
}
