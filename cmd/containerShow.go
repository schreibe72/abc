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
	"errors"
	"fmt"

	"github.com/schreibe72/abc/storage"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var containerShowCmd = &cobra.Command{
	Use:   "show",
	Short: "show the containing blobs in your storage account",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if container == "" {
			return errors.New("Container shoud not be empty!")
		}
		return nil
	},
	RunE: containerShow,
}

func init() {
	containerCmd.AddCommand(containerShowCmd)
	containerShowCmd.Flags().StringVarP(&container, "container", "c", "", "a Azure Container (required)")
	containerShowCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "a Azure Blob Prefix")
}

func containerShow(cmd *cobra.Command, args []string) error {
	s, err := storage.NewStorageClient(key, account, workercount)
	if err != nil {
		return err
	}

	s.Verbose = verbose
	s.Logger = logger

	l, err := s.ShowContainer(container, prefix)
	if err != nil {
		return err
	}
	for _, b := range l {
		fmt.Println(b)
	}
	return nil
}
