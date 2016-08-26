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

	"github.com/schreibe72/abc/storage"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var containerListCmd = &cobra.Command{
	Use:   "list",
	Short: "list all containers in your storage account",
	Long: `Here you can list all containers in your storage account. You can also list
	all container by a certain prefix.`,
	RunE: containerList,
}

func init() {
	containerCmd.AddCommand(containerListCmd)
	containerListCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "a Azure Container Prefix")
}

func containerList(cmd *cobra.Command, args []string) error {
	s, err := storage.NewStorageClient(key, account, workercount)
	if err != nil {
		return err
	}

	s.Verbose = verbose
	s.Logger = logger

	l, err := s.ListContainer(prefix)
	if err != nil {
		return err
	}
	for _, b := range l {
		fmt.Println(b)
	}
	return nil
}
