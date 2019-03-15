// Copyright © 2019 Developer <developer@vamp.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/models"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

const AppName string = "forklift"

// version should be in format d.d.d where d is a decimal number
const Version string = "0.1.5"

/*
Application name can change over time so it is made parameteric
*/
func AddAppName(str string) string {
	return strings.Replace(str, "$AppName", AppName, -1)
}

var cfgFile string

var Config ForkliftConfig

// Common code parameters
var configPath string
var configFileType string
var organization string
var environment string
var artficatsPath string
var role string
var kind string

type ForkliftConfig struct {
	Namespace         string                   `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	VampConfiguration models.VampConfiguration `yaml:"forklift,omitempty" json:"forklift,omitempty"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   AddAppName("$AppName"),
	Short: AddAppName("A command line client for $AppName"),
	Long: AddAppName(`$AppName is a organization and environment setup tool for vamp.
    It requires a vamp deployment setup as a backend to work.
    It is required to have a default config.
    Envrionment variables can be used in combination with the config.
    Environment variables:
      VAMP_FORKLIFT_VAULT_ADDR
      VAMP_FORKLIFT_VAULT_TOKEN
      VAMP_FORKLIFT_VAULT_CACERT
      VAMP_FORKLIFT_VAULT_CLIENT_CERT
      VAMP_FORKLIFT_VAULT_CLIENT_KEY
      VAMP_FORKLIFT_MYSQL_HOST
      VAMP_FORKLIFT_MYSQL_CONNECTION_PROPS
      VAMP_FORKLIFT_MYSQL_USER
      VAMP_FORKLIFT_MYSQL_PASSWORD
    `),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	logging.Init(os.Stdout, os.Stderr)

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.forklift/config.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolVarP(&logging.Verbose, "verbose", "v", false, "Verbose")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".$AppName" (without extension).
		path := filepath.FromSlash(home + AddAppName("/.$AppName"))

		if _, pathErr := os.Stat(path); os.IsNotExist(pathErr) {
			// path does not exist
			err = os.MkdirAll(path, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		viper.AddConfigPath(path)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig() // TODO: handle config file autocreation
	// unmarshal config
	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}

	unmarshallError := yaml.Unmarshal(bs, &Config)
	if unmarshallError != nil {
		panic(unmarshallError)
	}

	// TODO: Setup Defaults for Config
	SetupConfigurationEnvrionmentVariables()
	// For Checking during development:
	// fmt.Printf("Config: %v\n", Config)
	// logging.Log("Config: %v\n", Config)
}

func SetupConfigurationEnvrionmentVariables() {
	// VAMP_FORKLIFT_VAULT_ADDR
	viper.BindEnv("vault_addr", "VAMP_FORKLIFT_VAULT_ADDR")
	if viper.GetString("vault_addr") != "" {
		Config.VampConfiguration.Persistence.KeyValueStore.Vault.Url = viper.GetString("vault_addr")
	}
	// VAMP_FORKLIFT_VAULT_TOKEN
	viper.BindEnv("vault_token", "VAMP_FORKLIFT_VAULT_TOKEN")
	if viper.GetString("vault_token") != "" {
		Config.VampConfiguration.Persistence.KeyValueStore.Vault.Token = viper.GetString("vault_token")
	}
	// VAMP_FORKLIFT_VAULT_CACERT
	viper.BindEnv("vault_cacert", "VAMP_FORKLIFT_VAULT_CACERT")
	if viper.GetString("vault_cacert") != "" {
		Config.VampConfiguration.Persistence.KeyValueStore.Vault.ServerTlsCert = viper.GetString("vault_cacert")
	}
	// VAMP_FORKLIFT_VAULT_CLIENT_CERT
	viper.BindEnv("vault_client_cert", "VAMP_FORKLIFT_VAULT_CLIENT_CERT")
	if viper.GetString("vault_client_cert") != "" {
		Config.VampConfiguration.Persistence.KeyValueStore.Vault.ClientTlsCert = viper.GetString("vault_client_cert")
	}
	// VAMP_FORKLIFT_VAULT_CLIENT_KEY
	viper.BindEnv("vault_client_key", "VAMP_FORKLIFT_VAULT_CLIENT_KEY")
	if viper.GetString("vault_client_key") != "" {
		Config.VampConfiguration.Persistence.KeyValueStore.Vault.ClientTlsKey = viper.GetString("vault_client_key")
	}

	// VAMP_FORKLIFT_MYSQL_HOST - for example: mysql://<VAMP_FORKLIFT_MYSQL_HOST>/vamp-${namespace}?useSSL=false
	viper.BindEnv("mysql_host", "VAMP_FORKLIFT_MYSQL_HOST")
	if viper.GetString("mysql_host") != "" {
		url := "mysql://" + viper.GetString("mysql_host") + "/vamp-${namespace}"
		Config.VampConfiguration.Persistence.Database.Sql.Url = url
		databaseServerUrl := "mysql://" + viper.GetString("mysql_host")
		Config.VampConfiguration.Persistence.Database.Sql.DatabaseServerUrl = databaseServerUrl
	}
	// VAMP_FORKLIFT_MYSQL_PARAMS
	viper.BindEnv("mysql_params", "VAMP_FORKLIFT_MYSQL_CONNECTION_PROPS")
	if viper.GetString("mysql_params") != "" {
		Config.VampConfiguration.Persistence.Database.Sql.Url += "?" + viper.GetString("mysql_params")
		Config.VampConfiguration.Persistence.Database.Sql.DatabaseServerUrl += "?" + viper.GetString("mysql_params")
	}
	// VAMP_FORKLIFT_MYSQL_USER
	viper.BindEnv("mysql_user", "VAMP_FORKLIFT_MYSQL_USER")
	if viper.GetString("mysql_user") != "" {
		Config.VampConfiguration.Persistence.Database.Sql.User = viper.GetString("mysql_user")
	}
	// VAMP_FORKLIFT_MYSQL_PASSWORD
	viper.BindEnv("mysql_password", "VAMP_FORKLIFT_MYSQL_PASSWORD")
	if viper.GetString("mysql_password") != "" {
		Config.VampConfiguration.Persistence.Database.Sql.Password = viper.GetString("mysql_password")
	}
}
