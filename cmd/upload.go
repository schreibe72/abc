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
	"mime"
	"os"
	"path/filepath"

	"github.com/schreibe72/abc/storage"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "uploads a file to a container in your storage account",
	Long:  `uploads a file to a selected container and stets the contentType`,
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
		if workercount == 0 {
			return errors.New("Workercount shoud not be empty!")
		}
		if !pipe && filename == "" {
			return errors.New("You need an input! Pipe or Filename")
		}
		if pipe && blobname == "" {
			return errors.New("You need an Blobname!")
		}
		return nil
	},
	RunE: upload,
}

func init() {
	RootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringVarP(&blobname, "blobname", "n", "", "The Blob File Name (required for pipe)")
	uploadCmd.Flags().IntVarP(&workercount, "worker", "w", 10, "download Worker Count")
	uploadCmd.Flags().BoolVarP(&pipe, "pipe", "p", false, "incoming Pipe")
	uploadCmd.Flags().BoolVarP(&big, "big", "b", false, "spilt file which are bigger than 195GB in part Blockblobs")
	uploadCmd.Flags().StringVarP(&contentSetting.ContentType, "contentType", "T", "", "Contenttype for the uploaded file")
	uploadCmd.Flags().StringVarP(&contentSetting.CacheControl, "cacheControl", "C", "", "CacheControl for the uploaded file")
	uploadCmd.Flags().StringVarP(&contentSetting.ContentLanguage, "contentLanguage", "L", "", "ContentLanguage for the uploaded file")
	uploadCmd.Flags().StringVarP(&contentSetting.ContentEncoding, "contentEncoding", "E", "", "ContentEncoding for the uploaded file")
	uploadCmd.Flags().StringVarP(&filename, "filename", "f", "", "Filename to upload (required if no pipe)")
	uploadCmd.Flags().StringVarP(&container, "container", "c", "", "a Azure Container (required)")
	uploadCmd.MarkFlagFilename("filename", "")
	uploadCmd.MarkFlagRequired("container")
}

func upload(cmd *cobra.Command, args []string) error {

	if verbose {
		logger.Printf("Account: %s\nContainer: %s\n", account, container)
	}

	s, err := storage.NewStorageClient(key, account, workercount)
	if err != nil {
		return err
	}
	s.Verbose = verbose
	s.Logger = logger
	s.CreateContainer(container)
	if pipe {
		if verbose {
			logger.Printf("Blobname: %s\n", blobname)
			logger.Printf("Big: %t\n", big)
		}
		err := s.SaveBlob(os.Stdin, container, blobname, big, contentSetting)
		check(err)
	} else {
		f, err := os.Open(filename)
		check(err)
		defer f.Close()
		if contentSetting.ContentType == "" {
			contentSetting.ContentType = mime.TypeByExtension(filepath.Ext(filename))
		}
		if blobname == "" {
			blobname = filepath.Base(filename)
		}

		if verbose {
			logger.Printf("Upload %s as ContentType %s\n", blobname, contentSetting.ContentType)
		}
		big, err = FileIsTooBig(filename)
		check(err)
		if verbose {
			logger.Printf("Filename: %s\n", filename)
			logger.Printf("Blobname: %s\n", blobname)
			logger.Printf("Big: %t\n", big)
		}
		err = s.SaveBlob(f, container, blobname, big, contentSetting)
		check(err)
	}
	return nil
}
