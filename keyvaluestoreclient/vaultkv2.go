package keyvaluestoreclient

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/magneticio/forklift/logging"
)

func (c *VaultKeyValueStoreClient) GetData(key string, version int) (map[string]interface{}, error) {
	client := c.getClient()
	path := sanitizePath(key)
	mountPath, v2, pathError := isKVv2(path, client)
	if pathError != nil {
		logging.Error("Error checking version %s: %s", path, pathError)
		return nil, pathError
	}

	var versionParam map[string]string

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
		logging.Info("Prefix added to the kv path %v", path)
		if version > 0 {
			versionParam = map[string]string{
				"version": fmt.Sprintf("%d", version),
			}
		}
	}

	secret, err := kvReadRequest(client, path, versionParam)
	if err != nil {
		logging.Error("Error reading %s: %s", path, err)
		if secret != nil {
			return secret.Data, nil
		}
		return nil, errors.New(fmt.Sprintf("No value found at %s", path))
	}
	if secret == nil {
		logging.Error("No value found at %s", path)
		return nil, errors.New(fmt.Sprintf("No value found at %s", path))
	}

	data := secret.Data
	if v2 && data != nil {
		data = nil
		dataRaw := secret.Data["data"]
		if dataRaw != nil {
			data = dataRaw.(map[string]interface{})
		}
	}

	if data != nil {
		return data, nil
	}

	return nil, errors.New(fmt.Sprintf("No value found at %s", path))
}

func (c *VaultKeyValueStoreClient) PutData(key string, data map[string]interface{}, cas int) error {
	client := c.getClient()
	path := sanitizePath(key)

	mountPath, v2, pathError := isKVv2(path, client)
	if pathError != nil {
		logging.Error(pathError.Error())
		return pathError
	}

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
		data = map[string]interface{}{
			"data":    data,
			"options": map[string]interface{}{},
		}

		if cas > -1 {
			data["options"].(map[string]interface{})["cas"] = cas
		}
	}

	secret, writeError := client.Logical().Write(path, data)
	if writeError != nil {
		logging.Error("Error writing data to %s: %s", path, writeError)
		if secret != nil {
			logging.Info("Secret: %v\n", secret)
		}
		return writeError
	}
	if secret == nil {
		logging.Info("Success! Data written to: %s", path)
		return nil
	}
	return nil
}

func (c *VaultKeyValueStoreClient) ListData(key string) (map[string]interface{}, error) {
	client := c.getClient()
	path := ensureTrailingSlash(sanitizePath(key))
	mountPath, v2, pathError := isKVv2(path, client)
	if pathError != nil {
		logging.Error(pathError.Error())
		return nil, pathError
	}

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "metadata")
	}

	secret, listError := client.Logical().List(path)
	if listError != nil {
		logging.Error("Error listing %s: %s", path, listError.Error())
		return nil, listError
	}
	if secret == nil || secret.Data == nil {
		logging.Error(fmt.Sprintf("No value found at %s", path))
		return nil, errors.New(fmt.Sprintf("No value found at %s", path))
	}

	// If the secret is wrapped, return the wrapped response.
	if secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		logging.Info("Wrapped Secret %v\n", secret)
		// TODO: handle wrapped secret
	}

	if _, ok := extractListData(secret); !ok {
		logging.Error(fmt.Sprintf("No entries found at %s", path))
		return nil, errors.New(fmt.Sprintf("No entries found at %s", path))
	}

	return secret.Data, nil
}

func (c *VaultKeyValueStoreClient) DeleteData(key string, versions []string) error {
	client := c.getClient()
	path := sanitizePath(key)
	mountPath, v2, pathError := isKVv2(path, client)
	if pathError != nil {
		logging.Error(pathError.Error())
		return pathError
	}

	var secret *api.Secret
	var deleteError error
	if v2 {
		secret, deleteError = c.deleteV2(path, mountPath, versions, true)
	} else {
		secret, deleteError = client.Logical().Delete(path)
	}

	if deleteError != nil {
		logging.Error("Error deleting %s: %s", path, deleteError)
		if secret != nil {
			logging.Info("Secret %v\n", secret)
		}
		return deleteError
	}

	logging.Info("Success! Data deleted (if it existed) at: %s", path)
	return nil
}

func (c *VaultKeyValueStoreClient) deleteV2(path, mountPath string, versions []string, allVersions bool) (*api.Secret, error) {
	client := c.getClient()
	var err error
	var secret *api.Secret
	switch {
	case len(versions) > 0:
		path = addPrefixToVKVPath(path, mountPath, "delete")
		if err != nil {
			return nil, err
		}

		data := map[string]interface{}{
			"versions": kvParseVersionsFlags(versions),
		}

		secret, err = client.Logical().Write(path, data)
	default:
		prefix := "data"
		if allVersions {
			// this deletes all versions of data
			prefix = "metadata"
		}
		path = addPrefixToVKVPath(path, mountPath, prefix)
		if err != nil {
			return nil, err
		}

		secret, err = client.Logical().Delete(path)
	}

	return secret, err
}
