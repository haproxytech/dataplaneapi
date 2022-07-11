// Copyright 2019 HAProxy Technologies
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
//

package configuration

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/storage"
	jsoniter "github.com/json-iterator/go"
)

func DecodeBootstrapKey(key string) (map[string]string, error) {
	raw, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", key, err)
	}
	var decodedKey map[string]string
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(raw, &decodedKey)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", key, err)
	}
	return decodedKey, nil
}

func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CheckIfStorageDirIsOK(storageDir string, config *Configuration) error {
	if storageDir == "" {
		return errors.New("storage-dir in bootstrap key is empty")
	}
	_, errStorage := misc.CheckOrCreateWritableDirectory(storageDir)
	if errStorage != nil {
		return errStorage
	}
	dirs := []storage.FileType{
		storage.BackupsType, storage.MapsType, storage.SSLType,
		storage.SpoeTransactionsType, storage.SpoeType,
		storage.TransactionsType,
		storage.FileType("certs-cluster"),
	}
	for _, dir := range dirs {
		_, errStorage := misc.CheckOrCreateWritableDirectory(path.Join(storageDir, string(dir)))
		if errStorage != nil {
			return errStorage
		}
	}
	config.Cluster.StorageDir.Store(storageDir)
	config.HAProxy.ClusterTLSCertDir = path.Join(storageDir, "certs-cluster")
	config.Cluster.CertificateDir.Store(path.Join(storageDir, "certs-cluster"))
	return nil
}
