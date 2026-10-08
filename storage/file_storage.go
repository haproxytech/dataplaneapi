package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/google/renameio"
	native_misc "github.com/haproxytech/client-native/v6/misc"
	jsoniter "github.com/json-iterator/go"
)

type fileStorage[T Storable] struct {
	filePath string
}

func (f *fileStorage[T]) Get() (T, error) {
	returnData := new(T)

	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return *returnData, err
	}

	if len(data) == 0 {
		return *returnData, nil
	}
	json := jsoniter.ConfigFastest
	if err := json.Unmarshal(data, returnData); err != nil {
		return *returnData, err
	}

	return *returnData, nil
}

func (f *fileStorage[T]) Store(data T) error {
	if err := f.initFile(); err != nil {
		return err
	}
	json := jsoniter.ConfigFastest
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// may hold the cluster token and service discovery credentials
	if err = renameio.WriteFile(f.filePath, jsonData, 0o600); err != nil {
		return err
	}
	return nil
}

func (f *fileStorage[T]) initFile() error {
	if _, err := native_misc.CheckOrCreateWritableDirectory(filepath.Dir(f.filePath)); err != nil {
		log.Fatalf("error initializing dataplane internal storage: %v", err)
	}
	if _, err := os.Stat(f.filePath); err != nil {
		if os.IsNotExist(err) {
			file, errCreate := os.OpenFile(f.filePath, os.O_CREATE|os.O_WRONLY, 0o600)
			if errCreate != nil {
				return errCreate
			}
			return file.Close()
		}
		return err
	}
	return nil
}
