package cmd

import (
	"os"
	"path"

	"github.com/schreibe72/abc/storage"
)

func check(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func createDirIfNeeded(configPath string) error {
	dirname := path.Dir(configPath)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return os.MkdirAll(dirname, 0700)
	}
	return nil
}

func FileIsTooBig(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	if info.Size() > (storage.MaxBlobBockCount * storage.MaxBlobBlockSize) {
		return true, nil
	}
	return false, nil
}
