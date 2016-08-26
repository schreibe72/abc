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

	"github.com/schreibe72/abc/storage"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var blobDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a blob",
	Long:  `delete selected blob files in your storage account`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if blobname == "" {
			return errors.New("Blobname shoud not be empty!")
		}
		return nil
	},
	RunE: blobDelete,
}

func init() {
	blobCmd.AddCommand(blobDeleteCmd)
	blobDeleteCmd.Flags().StringVarP(&container, "container", "c", "", "a Azure Container (required)")
	blobDeleteCmd.Flags().StringVarP(&blobname, "blobname", "n", "", "The Blob File Name (required)")
}

func blobDelete(cmd *cobra.Command, args []string) error {
	s, err := storage.NewStorageClient(key, account, workercount)
	if err != nil {
		return err
	}

	s.Verbose = verbose
	s.Logger = logger

	err = s.DeleteBlob(container, blobname)
	return err
}
