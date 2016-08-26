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
	"os"

	"github.com/schreibe72/abc/storage"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download a blob",
	Long:  `Downloads a blob file selected in you storage account and selected container`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if key == "" {
			return errors.New("Key should not be empty!")
		}
		if account == "" {
			return errors.New("Account shoud not be empty!")
		}
		if container == "" {
			return errors.New("Container shoud not be empty!")
		}
		return nil
	},
	RunE: download,
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&blobname, "blobname", "n", "", "The Blob File Name (required)")
	downloadCmd.Flags().IntVarP(&workercount, "worker", "w", 10, "download Worker Count")
	downloadCmd.Flags().BoolVarP(&pipe, "pipe", "p", false, "outgoing Pipe")
	downloadCmd.Flags().StringVarP(&filename, "filename", "f", "", "Filename to download")
	downloadCmd.Flags().StringVarP(&container, "container", "c", "", "a Azure Container (required)")
}

func download(cmd *cobra.Command, args []string) error {

	s, err := storage.NewStorageClient(key, account, workercount)
	if err != nil {
		return err
	}

	s.Verbose = verbose
	s.Logger = logger

	if verbose {
		logger.Printf("Account: %s\nContainer: %s\n", account, container)

	}

	if pipe {
		if verbose {
			logger.Println("Output to Pipe")
		}
		check(s.LoadBlob(os.Stdout, container, blobname))
	} else {
		if verbose {
			logger.Println("Output to File")
		}
		if filename == "" {
			filename = blobname
		}
		f, err := os.Create(filename)
		check(err)
		defer f.Close()
		check(s.LoadBlob(f, container, blobname))
		f.Close()
	}
	return nil
}
