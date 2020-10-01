// Copyright © 2020 Developer <developer@vamp.io>
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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magneticio/forklift/logging"
	"github.com/magneticio/forklift/models"
	"github.com/magneticio/forklift/util"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// AppName - application name
const AppName string = "forklift"

// Version - should be in format d.d.d where d is a decimal number
const Version string = "v1.0.0"

// AddAppName - Application name can change over time so it is made parameteric
func AddAppName(str string) string {
	return strings.Replace(str, "$AppName", AppName, -1)
}

// Config - forklift configuration
var Config models.ForkliftConfiguration

// Common code parameters
var cfgFile string
var configPath string
var configFileType string
var serviceID uint64
var applicationID uint64

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   AddAppName("$AppName"),
	Short: AddAppName("A command line client for $AppName"),
	Long: AddAppName(`$AppName is a setup tool for vamp.
	It is required to have a default config.
	Envrionment variables can be used in combination with the config.
	Environment variables:
		VAMP_FORKLIFT_PROJECT
		VAMP_FORKLIFT_CLUSTER
		VAMP_FORKLIFT_VAULT_ADDR
		VAMP_FORKLIFT_VAULT_TOKEN
		VAMP_FORKLIFT_VAULT_CACERT
		VAMP_FORKLIFT_VAULT_CLIENT_CERT
		VAMP_FORKLIFT_VAULT_CLIENT_KEY`),
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

	rootCmd.PersistentFlags().VarP(Config.ProjectID, "project", "p", "project id")
	rootCmd.PersistentFlags().VarP(Config.ClusterID, "cluster", "c", "cluster id")
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

	setupConfigurationEnvrionmentVariables()

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
	// For Checking during development:
	// fmt.Printf("Config: %v\n", Config)
	jsonConfig, _ := json.Marshal(Config)

	logging.Info("Forklift configuration %v\n", util.PrettyJson(string(jsonConfig)))
}

func setupConfigurationEnvrionmentVariables() {
	viper.BindEnv("project", "VAMP_FORKLIFT_PROJECT")
	viper.BindEnv("cluster", "VAMP_FORKLIFT_CLUSTER")
	viper.BindEnv("key-value-store-url", "VAMP_FORKLIFT_VAULT_ADDR")
	viper.BindEnv("key-value-store-token", "VAMP_FORKLIFT_VAULT_TOKEN")
	viper.BindEnv("key-value-store-server-tls-cert", "VAMP_FORKLIFT_VAULT_CACERT")
	viper.BindEnv("key-value-store-client-tls-cert", "VAMP_FORKLIFT_VAULT_CLIENT_CERT")
	viper.BindEnv("key-value-store-client-tls-key", "VAMP_FORKLIFT_VAULT_CLIENT_KEY")
}
