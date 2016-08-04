package storage

import (
	"errors"

	as "github.com/schreibe72/azure-sdk-for-go/storage"
)

func (a *StorageAttributes) ListContainer(prefix string) ([]string, error) {
	var clp as.ListContainersParameters
	var containerNames []string
	if prefix != "" {
		clp.Prefix = prefix
	}
	l, err := a.blobService.ListContainers(clp)
	if err != nil {
		return []string{}, err
	}
	for _, c := range l.Containers {
		containerNames = append(containerNames, c.Name)
	}
	return containerNames, nil
}

func (a *StorageAttributes) CreateContainer(container string) error {
	if container == "" {
		return errors.New("no container provided")
	}
	b, err := a.blobService.CreateContainerIfNotExists(container, "")
	if err != nil {
		return err
	}
	if b && a.Verbose {
		a.Logger.Printf("%s Created", container)
	}
	return nil
}

func (a *StorageAttributes) DeleteContainer(container string) error {
	if container == "" {
		return errors.New("no container provided")
	}
	b, err := a.blobService.DeleteContainerIfExists(container)
	if err != nil {
		return err
	}

	if b && a.Verbose {
		a.Logger.Printf("%s Deleted", container)
	}
	return nil
}

func (a *StorageAttributes) ShowContainer(container string, prefix string) ([]string, error) {
	if container == "" {
		return []string{}, errors.New("no container provided")
	}
	var blp as.ListBlobsParameters
	blobNames := make([]string, 0)
	if prefix != "" {
		blp.Prefix = prefix
	}
	l, err := a.blobService.ListBlobs(container, blp)
	if err != nil {
		return []string{}, err
	}
	for _, b := range l.Blobs {
		blobNames = append(blobNames, b.Name)
	}
	return blobNames, nil
}

func (a *StorageAttributes) DeleteBlob(container string, name string) error {
	switch {
	case container == "":
		return errors.New("no container provided")
	case name == "":
		return errors.New("no blob name provided")
	}
	b, err := a.blobService.DeleteBlobIfExists(container, name, map[string]string{})
	if err != nil {
		return err
	}
	if b && a.Verbose {
		a.Logger.Printf("[%s]%s deleted", container, name)
	}
	return nil
}
