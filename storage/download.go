package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func (a *StorageAttributes) loadBlob(w io.Writer, container string, name string) error {

	if err := validateBlobName(container, name); err != nil {
		return err
	}

	r, err := a.blobService.GetBlobRange(container, name, "0-", map[string]string{})
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	r.Close()
	return err
}

// LoadBlob downloads the blob. It desides wether this is a bundle with spilt files or not.
func (a *StorageAttributes) LoadBlob(w io.Writer, container string, name string) error {

	if err := validateBlobName(container, name); err != nil {
		return err
	}

	bundleFileName := fmt.Sprintf("%s-bundle.json", name)
	blobExist, err := a.blobService.BlobExists(container, bundleFileName)
	if err != nil {
		return err
	}
	if blobExist {
		if a.Verbose {
			a.Logger.Println("Bundle Exists")
		}
		return a.loadBlobBundle(w, container, bundleFileName)
	}
	if a.Verbose {
		a.Logger.Println("Normal File")
	}
	return a.loadBlob(w, container, name)
}

func (a *StorageAttributes) loadBlobBundle(w io.Writer, container string, name string) error {

	if err := validateBlobName(container, name); err != nil {
		return err
	}

	var data []byte
	var b bundle
	r, err := a.blobService.GetBlobRange(container, name, "0-", map[string]string{})
	if err != nil {
		return err
	}
	data, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &b)
	if err != nil {
		a.Logger.Println(string(data) + "\n\n")
		return err
	}
	for _, item := range b.Bundle {
		if a.Verbose {
			a.Logger.Printf("Downloading Part: %s\n", item.BlobName)
			err = a.loadBlob(w, container, item.BlobName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
