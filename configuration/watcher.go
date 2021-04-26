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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

type ConfigWatcherParams struct {
	FilePath string
	Callback func()
	Version  string
	Ctx      context.Context
}

type ConfigWatcherUpdate struct {
	Version string
	Content string
}

type ConfigWatcher struct {
	configFile string
	update     chan string
	callback   func()
	done       <-chan struct{}
	wa         *fsnotify.Watcher
}

func NewConfigWatcher(params ConfigWatcherParams) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.WithError(err).Info("Failed to initialize watcher")
		return nil, err
	}
	if err := watcher.Add(filepath.Dir(params.FilePath)); err != nil {
		log.WithError(err).Info(fmt.Sprintf("Failed to watch config file: %s", params.FilePath))
		return nil, err
	}

	ctx := params.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	cw := &ConfigWatcher{
		wa:         watcher,
		configFile: params.FilePath,
		update:     make(chan string),
		callback:   params.Callback,
		done:       ctx.Done(),
	}
	return cw, nil
}

func (w *ConfigWatcher) Update(hash string) {
	w.update <- hash
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
			log.WithError(err).Info("Closing config file watcher watcher")
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
	content, err := ioutil.ReadFile(w.configFile)
	if err != nil {
		log.Warningf("Watcher: error reading config file: %s", err.Error())
		return false
	}
	lines := strings.Split(string(content), "\n")
	parts := strings.Split(lines[0], "=")
	if len(parts) != 2 || parts[0] != "# _md5hash" {
		return true
	}
	//nolint:gosec
	bHash := md5.Sum([]byte(strings.Join(lines[1:], "\n")))
	hash := fmt.Sprintf("%x", bHash)
	return parts[1] != hash
}
