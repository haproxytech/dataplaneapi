// Copyright 2020 HAProxy Technologies
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
	"context"
	"crypto/md5"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/haproxytech/dataplaneapi/log"
)

type ConfigWatcherParams struct {
	Ctx      context.Context
	Callback func()
	FilePath string
	Version  string
}

type ConfigWatcher struct {
	callback   func()
	done       <-chan struct{}
	wa         *fsnotify.Watcher
	configFile string
}

func NewConfigWatcher(params ConfigWatcherParams) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Infof("Failed to initialize watcher: %s", err.Error())
		return nil, err
	}
	if err := watcher.Add(filepath.Dir(params.FilePath)); err != nil {
		log.Infof("Failed to watch config file: %s: %s", params.FilePath, err.Error())
		return nil, err
	}

	ctx := params.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cw := &ConfigWatcher{
		wa:         watcher,
		configFile: params.FilePath,
		callback:   params.Callback,
		done:       ctx.Done(),
	}
	return cw, nil
}

func (w *ConfigWatcher) Listen() {
	defer w.wa.Close()
	for {
		select {
		case event, ok := <-w.wa.Events:
			if !ok {
				return
			}
			if event.Name != w.configFile {
				continue
			}
			if w.checkFlags(event) {
				if w.invalidHash() {
					w.callback()
				}
			}
		case err, ok := <-w.wa.Errors:
			log.Infof("Closing config file watcher watcher: %s", err.Error())
			if !ok {
				return
			}
		case <-w.done:
			return
		}
	}
}

func (w *ConfigWatcher) checkFlags(event fsnotify.Event) bool {
	return event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create
}

func (w *ConfigWatcher) invalidHash() bool {
	content, err := os.ReadFile(w.configFile)
	if err != nil {
		log.Warningf("Watcher: error reading config file: %s", err.Error())
		return false
	}
	lines := strings.Split(string(content), "\n")
	parts := strings.Split(lines[0], "=")
	if len(parts) != 2 || parts[0] != "# _md5hash" {
		return true
	}
	bHash := md5.Sum([]byte(strings.Join(lines[1:], "\n")))
	hash := hex.EncodeToString(bHash[:])
	return parts[1] != hash
}
