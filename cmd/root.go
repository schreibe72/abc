// Copyright © 2016 Manfred Schreiber <software@manfredschreiber.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/schreibe72/abc/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verbose        bool
	key            string
	account        string
	blobname       string
	workercount    int
	big            bool
	pipe           bool
	filename       string
	container      string
	prefix         string
	cfgFile        string
	logger         *log.Logger
	contentSetting storage.ContentSetting
)

type Config struct {
	Key     string
	Account string
}

var RootCmd = &cobra.Command{
	Use:   "abc",
	Short: "Azure Blob Command to up and download files",
	Long: `This Tool can upload files to a azure storage account
	and sets the content type of this file. It also can create,
	delete and list container.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose Output")
	RootCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "Azure Storage Account Key")
	RootCmd.PersistentFlags().StringVarP(&account, "account", "a", "", "Azure Storage Account")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	home := os.Getenv("HOME")
	cfgFile = filepath.Join(home, ".azure", "abc-config.json")
	logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgFilename := path.Base(cfgFile)
	cfgExt := path.Ext(cfgFilename)
	cfgName := cfgFilename[0 : len(cfgFilename)-len(cfgExt)]
	cfgExt = cfgExt[1:len(cfgExt)] //Remove leading .

	viper.SetConfigType(cfgExt)
	viper.SetConfigName(cfgName)           // name of config file (without extension)
	viper.AddConfigPath(path.Dir(cfgFile)) // adding home directory as first search path
	viper.AutomaticEnv()                   // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if key == "" {
			key = viper.GetString("Key")
		}
		if account == "" {
			account = viper.GetString("Account")
		}
	}

}
