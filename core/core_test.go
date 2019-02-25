package core

import (
	"testing"

	"github.com/magneticio/forklift/models"
	"github.com/stretchr/testify/assert"
)

func TestGetNamespaceSqlConfiguration(t *testing.T) {

	c := Core{
		Conf: Configuration{
			VampConfiguration: models.VampConfiguration{
				Persistence: models.Persistence{
					Database: models.Database{
						Sql: models.SqlConfiguration{
							Database:          "database-${namespace}",
							Table:             "table-${namespace}",
							User:              "user-${namespace}",
							Password:          "password-${namespace}",
							Url:               "url-${namespace}",
							DatabaseServerUrl: "databaseServerUrl-${namespace}",
						},
						Type: "type",
					},
				},
			},
		},
	}

	result := c.GetNamespaceDatabaseConfiguration("namespace")

	assert.Equal(t, result.Sql.Database, "database-namespace")
	assert.Equal(t, result.Sql.Table, "table-namespace")
	assert.Equal(t, result.Sql.User, "user-namespace")
	assert.Equal(t, result.Sql.Password, "password-namespace")
	assert.Equal(t, result.Sql.Url, "url-namespace")
	assert.Equal(t, result.Sql.DatabaseServerUrl, "databaseServerUrl-namespace")
	assert.Equal(t, result.Type, "type")
}
