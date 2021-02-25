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
	"fmt"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/client-native/v2/models"
	log "github.com/sirupsen/logrus"
)

type MapSync struct {
	mapQuitChan chan struct{}
	mu          sync.RWMutex
}

func NewMapSync() *MapSync {
	return &MapSync{
		mapQuitChan: make(chan struct{}),
	}
}

// Stop stops maps syncing
func (ms *MapSync) Stop() {
	ms.mapQuitChan <- struct{}{}
}

// SyncAll sync maps file entries with runtime maps entries for all configured files.
// Missing runtime entries are appended to the map file
func (ms *MapSync) SyncAll(client client_native.IHAProxyClient) {

	haproxyOptions := Get().HAProxy

	d := time.Duration(haproxyOptions.UpdateMapFilesPeriod)
	ticker := time.NewTicker(d * time.Second)

	for {
		select {
		case <-ticker.C:
			maps, err := client.GetRuntime().ShowMaps()
			if err != nil {
				log.Warning("show maps sync error: ", err.Error())
				continue
			}
			for _, mp := range maps {
				go func(mp *models.Map) {
					_, err := ms.Sync(mp, client)
					if err != nil {
						log.Warning(err.Error())
					}
				}(mp)
			}
		case <-ms.mapQuitChan:
			return
		}
	}
}

// Sync syncs one map file to runtime entries
func (ms *MapSync) Sync(mp *models.Map, client client_native.IHAProxyClient) (bool, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	rawFile, err := os.Open(mp.File)
	if err != nil {
		return false, fmt.Errorf("error reading map file: %s %s", mp.File, err.Error())
	}
	fileEntries := client.GetRuntime().ParseMapEntriesFromFile(rawFile, false)
	sort.Slice(fileEntries, func(i, j int) bool { return fileEntries[i].Key < fileEntries[j].Key })

	// runtime map entries
	id := fmt.Sprintf("#%s", mp.ID)
	runtimeEntries, err := client.GetRuntime().ShowMapEntries(id)
	if err != nil {
		return false, fmt.Errorf("getting runtime entries error: id: %s %s", id, err.Error())
	}
	sort.Slice(runtimeEntries, func(i, j int) bool { return runtimeEntries[i].Key < runtimeEntries[j].Key })

	if len(fileEntries) != len(runtimeEntries) {
		return dumpRuntimeEntries(mp.File, runtimeEntries)
	}

	if !equalSomeEntries(fileEntries, runtimeEntries) {
		return dumpRuntimeEntries(mp.File, runtimeEntries)
	}

	if !equal(fileEntries, runtimeEntries) {
		return dumpRuntimeEntries(mp.File, runtimeEntries)
	}
	return true, nil
}

// equalSomeEntries compares last few runtime entries with file entries
// if records differs, check is run against random entries
func equalSomeEntries(fEntries, rEntries models.MapEntries, index ...int) bool {
	if len(fEntries) != len(rEntries) {
		return false
	}

	max := 0
	switch l := len(rEntries); {
	case l > 19:
		for i := l - 20; i < l; i++ {
			if rEntries[i].Key != fEntries[i].Key || rEntries[i].Value != fEntries[i].Value {
				return false
			}
		}
		max = l - 19
	case l == 0:
		return true
	default:
		max = l
	}

	maxRandom := 10
	if max < 10 {
		maxRandom = max
	}

	for i := 0; i < maxRandom; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		// There's no need for strong number generation, here, just need for performance
		r := rand.Intn(max) // nolint:gosec
		if len(index) > 0 {
			r = index[0]
		}
		if rEntries[r].Key != fEntries[r].Key || rEntries[r].Value != fEntries[r].Value {
			return false
		}
	}
	return true
}

// equal compares runtime and map entries
// Returns true if all entries are same, otherwise returns false
func equal(a, b models.MapEntries) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		// ID should not be compared, since it doesn't exists in file
		if v.Key != b[i].Key || v.Value != b[i].Value {
			return false
		}
	}
	return true
}

// dumpRuntimeEntries dumps runtime entries into map file
// Returns true,nil if succeed, otherwise retuns false,error
func dumpRuntimeEntries(file string, me models.MapEntries) (bool, error) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return false, fmt.Errorf("error opening map file: %s %s", file, err.Error())
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return false, fmt.Errorf("error truncating map file: %s %s", file, err.Error())
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return false, fmt.Errorf("error setting file to offset: %s %s", file, err.Error())
	}

	for _, e := range me {
		line := fmt.Sprintf("%s %s%s", e.Key, e.Value, "\n")
		_, err = f.WriteString(line)
		if err != nil {
			return false, fmt.Errorf("error writing map file: %s %s", file, err.Error())
		}
	}
	log.Infof("map file %s synced with runtime entries", file)
	return true, nil
}
