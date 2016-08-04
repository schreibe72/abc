package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

func (a *Arguments) load() error {
	buf, err := ioutil.ReadFile(a.configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, a)
	return err
}

func (a *Arguments) save() error {
	buf, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}

	err = createDirIfNeeded(a.configPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(a.configPath, buf, 0600)
	return err
}

func createDirIfNeeded(configPath string) error {
	dirname := path.Dir(configPath)
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return os.MkdirAll(dirname, 0700)
	}
	return nil
}
