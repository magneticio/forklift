package core

import (
	"fmt"
	"strings"

	"github.com/magneticio/forklift/sql"
	"github.com/magneticio/forklift/util"
)

type SqlConfiguration struct {
	Database          string `yaml:"database,omitempty" json:"database,omitempty"`
	Table             string `yaml:"table,omitempty" json:"table,omitempty"`
	User              string `yaml:"user,omitempty" json:"user,omitempty"`
	Password          string `yaml:"password,omitempty" json:"password,omitempty"`
	Url               string `yaml:"url,omitempty" json:"url,omitempty"`
	DatabaseServerUrl string `yaml:"database-server-url,omitempty" json:"database-server-url,omitempty"`
}

type Configuration struct {
	Sql SqlConfiguration `yaml:"sql,omitempty" json:"sql,omitempty"`
}

type Core struct {
	Conf Configuration
}

func NewCore(conf Configuration) (*Core, error) {

	return &Core{
		Conf: conf,
	}, nil
}

func (c *Core) CreateOrganisation(namespacedOrganisation string) error {

	sqlConfig := c.GetNamespaceSqlConfiguration(namespacedOrganisation)

	host, hostError := util.GetHostFromUrl(sqlConfig.Url)
	if hostError != nil {
		fmt.Printf("Error: %v\n", hostError.Error())
		return hostError
	}

	client, clientError := sql.NewMySqlClient(sqlConfig.User, sqlConfig.Password, host, sqlConfig.Database)
	if clientError != nil {
		fmt.Printf("Error: %v\n", clientError.Error())
		return clientError
	}

	return client.SetupOrganisation(namespacedOrganisation)

}

func (c *Core) GetNamespaceSqlConfiguration(namespace string) *SqlConfiguration {
	return &SqlConfiguration{
		Database:          Namespaced(namespace, c.Conf.Sql.Database),
		Table:             Namespaced(namespace, c.Conf.Sql.Table),
		User:              Namespaced(namespace, c.Conf.Sql.User),
		Password:          Namespaced(namespace, c.Conf.Sql.Password),
		Url:               Namespaced(namespace, c.Conf.Sql.Url),
		DatabaseServerUrl: Namespaced(namespace, c.Conf.Sql.DatabaseServerUrl),
	}
}

func Namespaced(namespace string, text string) string {
	return strings.Replace(text, "${namespace}", namespace, -1)
}
