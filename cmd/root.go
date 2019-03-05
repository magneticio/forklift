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
	"strings"

	"github.com/magneticio/forklift/models"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// version should be in format d.d.d where d is a decimal number
const Version string = "0.0.6"
const AppName string = "forklift"

// Backend version is the version this client is tested with
const BackendVersion string = "1.0.4"

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
var artficatsPath string
var role string

type ForkliftConfig struct {
	Namespace         string                   `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	VampConfiguration models.VampConfiguration `yaml:"forklift,omitempty" json:"forklift,omitempty"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   AddAppName("$AppName"),
	Short: AddAppName("A command line client for $AppName"),
	Long: AddAppName(`A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`),
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.forklift.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".forklift" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".forklift")
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
	// fmt.Printf("Config: %v\n", Config)
}
