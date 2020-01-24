package core

import (
	"testing"

	"github.com/magneticio/forklift/models"
	"github.com/stretchr/testify/assert"
)

func TestGetNamespaceSqlConfiguration(t *testing.T) {

	c := Core{
		Conf: models.ForkliftConfiguration{
			Namespace:             "namespace",
			DatabaseEnabled:       true,
			DatabaseType:          "database-type-${namespace}",
			DatabaseName:          "database-name-${namespace}",
			DatabaseURL:           "database-url-${namespace}",
			DatabaseUser:          "database-user-${namespace}",
			DatabasePassword:      "database-password-${namespace}",
			DatabaseTable:         "database-table-${namespace}",
			KeyValueStoreUrL:      "kv-url-${namespace}",
			KeyValueStoreToken:    "kv-token-${namespace}",
			KeyValueStoreBasePath: "kv-bas-path-${namespace}",
			KeyValueStoreType:     "kv-type-${namespace}",
		},
	}

	result, err := c.GetNamespaceDatabaseConfiguration("namespace")

	assert.Nil(t, err)

	assert.Equal(t, result.Sql.Database, "database-name-namespace")
	assert.Equal(t, result.Sql.Table, "database-table-namespace")
	assert.Equal(t, result.Sql.User, "database-user-${namespace}")
	assert.Equal(t, result.Sql.Password, "database-password-${namespace}")
	assert.Equal(t, result.Sql.Url, "database-url-namespace")
	assert.Equal(t, result.Type, "database-type-${namespace}")
}
