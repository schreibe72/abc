package storage

import (
	"errors"
	"log"

	as "github.com/schreibe72/azure-sdk-for-go/storage"
)

//MaxBlobBlockSize
//
const (
	MaxBlobBockCount = 50000
	MaxBlobBlockSize = as.MaxBlobBlockSize
)

type ContentSetting as.ContentSetting

type StorageAttributes struct {
	Key           string
	Account       string
	storageClient as.Client
	blobService   as.BlobStorageClient
	WorkerCount   int
	Verbose       bool
	Logger        *log.Logger
}

type uploadBlock struct {
	buf []byte
	id  string
}

type bundleItem struct {
	BlobName string `json:"blobName"`
	BlobSize uint64 `json:"blobSize"`
	BlobMD5  string `json:"blobMD5"`
	EOF      bool   `json:"EOF"`
}

type bundle struct {
	Bundle      []bundleItem
	FileName    string
	ContentType string
}

func (a *StorageAttributes) NewStorageClient() error {
	switch {
	case a.Account == "":
		return errors.New("no Storage Account provided")
	case a.Key == "":
		return errors.New("no Azure Storage Access Key provided")
	}
	client, err := as.NewBasicClient(a.Account, a.Key)
	if err != nil {
		return err
	}
	a.storageClient = client
	a.blobService = a.storageClient.GetBlobService()
	return nil
}