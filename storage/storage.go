package storage

import (
	"errors"
	"log"

	as "github.com/Azure/azure-sdk-for-go/storage"
)

//MaxBlobBlockSize
//
const (
	MaxBlobBockCount = 50000
	MaxBlobBlockSize = as.MaxBlobBlockSize
)

type ContentSetting as.BlobHeaders

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

func NewStorageClient(key string, account string, worker int) (StorageAttributes, error) {
	s := StorageAttributes{Key: key,
		Account:     account,
		WorkerCount: worker,
	}
	if err := validateAccountCredentials(s.Account, s.Key); err != nil {
		return StorageAttributes{}, err
	}
	client, err := as.NewBasicClient(s.Account, s.Key)
	if err != nil {
		return StorageAttributes{}, err
	}
	s.storageClient = client
	s.blobService = s.storageClient.GetBlobService()
	return s, nil
}

func validateBlobName(container string, name string) error {
	switch {
	case container == "":
		return errors.New("no container provided")
	case name == "":
		return errors.New("no blob name provided")
	}
	return nil
}

func validateAccountCredentials(account string, key string) error {
	switch {
	case account == "":
		return errors.New("no Storage Account provided")
	case key == "":
		return errors.New("no Azure Storage Access Key provided")
	}
	return nil
}
