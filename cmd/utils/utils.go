package utils

import (
	"encoding/json"
	"errors"
	"os"
)

const StoragePath string = "./storage"
const StorageFile string = "/data.json"

func WriteStorageData(data map[string]string) error {
	if _, err := os.Stat(StoragePath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(StoragePath, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = os.Create(StoragePath + StorageFile)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(StoragePath+StorageFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {

		}
	}(f)

	jsonInd, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonInd)
	if err != nil {
		return err
	}
	return nil
}

func CheckStorageEntryExists() bool {
	if _, err := os.Stat(StoragePath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
