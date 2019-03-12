package keyvaluestoreclient

import (
	"errors"
	"fmt"

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
