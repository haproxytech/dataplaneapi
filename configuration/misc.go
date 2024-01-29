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
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/renameio"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/storage"
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

	var keySummary string
	if len(key) > 10 {
		keySummary = key[:4] + "..." + key[len(key)-5:]
	} else {
		keySummary = key
	}

	if expiryUnixTS, ok := decodedKey["expiring-time"]; ok {
		tUnix, ok2 := strconv.ParseInt(expiryUnixTS, 10, 64)
		if ok2 != nil {
			return nil, fmt.Errorf("bootstrap key %s error, decoding expiry to int: %s", keySummary, expiryUnixTS)
		}
		expiryTime := time.Unix(tUnix, 0)
		if expiryTime.Before(time.Now()) {
			return nil, fmt.Errorf("refusing to use expired bootstrap key: %s expired on: %s", keySummary, strfmt.DateTime(expiryTime))
		}
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
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func RemoveStorageFolder(storageDir string) error {
	return os.RemoveAll(storageDir)
}

func InitStorageNoticeFile(storageDir string) error {
	content := strings.Builder{}

	_, _ = fmt.Fprintf(&content, "# *********************************************************************************\n")
	_, _ = fmt.Fprintf(&content, "# NOTE: This storage folder contains files managed by HAProxy Fusion Control Plane:\n")
	_, _ = fmt.Fprintf(&content, "#       manual edits may cause issues and misconfigurations.\n")

	return renameio.WriteFile(path.Join(storageDir, "NOTICE"), []byte(content.String()), os.ModePerm)
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
	return nil
}

func removeFromSlice[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}
