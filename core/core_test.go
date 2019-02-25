package core

import (
	"testing"

	"github.com/magneticio/forklift/models"
	"github.com/stretchr/testify/assert"
)

func TestGetNamespaceSqlConfiguration(t *testing.T) {

	c := Core{
		Conf: Configuration{
			Sql: models.SqlConfiguration{
				Database:          "database-${namespace}",
				Table:             "table-${namespace}",
				User:              "user-${namespace}",
				Password:          "password-${namespace}",
				Url:               "url-${namespace}",
				DatabaseServerUrl: "databaseServerUrl-${namespace}",
			},
		},
	}

	result := c.GetNamespaceSqlConfiguration("namespace")

	assert.Equal(t, result.Database, "database-namespace")
	assert.Equal(t, result.Table, "table-namespace")
	assert.Equal(t, result.User, "user-namespace")
	assert.Equal(t, result.Password, "password-namespace")
	assert.Equal(t, result.Url, "url-namespace")
	assert.Equal(t, result.DatabaseServerUrl, "databaseServerUrl-namespace")
}
