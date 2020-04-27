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

package haproxy

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/haproxytech/models/v2"

	log "github.com/sirupsen/logrus"
)

type IReloadAgent interface {
	Init(delay int, reloadCmd string, restartCmd string, configFile string, retention int) error
	Reload() string
	Restart() error
	ForceReload() error
	GetReloads() models.Reloads
	GetReload(id string) *models.Reload
}

type reloadCache struct {
	failedReloads map[string]*models.Reload
	lastSuccess   *models.Reload
	next          string
	current       string
	index         int64
	retention     int
	mu            sync.RWMutex
}

// ReloadAgent handles all reloads, scheduled or forced
type ReloadAgent struct {
	delay         int
	reloadCmd     string
	restartCmd    string
	configFile    string
	lkgConfigFile string
	cache         reloadCache
}

// Init a new reload agent
func (ra *ReloadAgent) Init(delay int, reloadCmd string, restartCmd string, configFile string, retention int) error {
	ra.reloadCmd = reloadCmd
	ra.restartCmd = restartCmd
	ra.configFile = configFile
	if delay == 0 {
		delay = 5
	}
	ra.delay = delay
	ra.lkgConfigFile = configFile + ".lkg"

	// create last known good file, assume it is valid when starting
	if err := copyFile(ra.configFile, ra.lkgConfigFile); err != nil {
		return err
	}
	ra.cache.Init(retention)
	go ra.handleReloads()
	return nil
}

func (ra *ReloadAgent) handleReloads() {
	//nolint:gosimple
	for {
		select {
		case <-time.After(time.Duration(ra.delay) * time.Second):
			if ra.cache.next != "" {
				ra.cache.mu.Lock()
				ra.cache.current = ra.cache.next
				ra.cache.next = ""
				ra.cache.mu.Unlock()
				response, err := ra.reloadHAProxy()
				if err != nil {
					ra.cache.failReload(response)
					log.Warning("Reload failed " + err.Error())
				} else {
					ra.cache.succeedReload(response)
				}
			}
		}
	}
}

func (ra *ReloadAgent) reloadHAProxy() (string, error) {
	// try the reload
	log.Debug("Reload started...")
	t := time.Now()
	output, err := execCmd(ra.reloadCmd)
	log.Debug("Reload finished.")
	log.Debug("Time elapsed: ", time.Since(t))
	if err != nil {
		reloadFailedError := err
		// if failed, return to last known good file and restart and return the original file
		log.Info("Reload failed, restarting with last known good config...")
		if err := copyFile(ra.configFile, ra.configFile+".bck"); err != nil {
			return fmt.Sprintf("Reload failed: %s, failed to backup original config file for restart.", output), err
		}
		defer func() {
			// nolint:errcheck
			copyFile(ra.configFile+".bck", ra.configFile)
			os.Remove(ra.configFile + ".bck")
		}()
		if err := copyFile(ra.lkgConfigFile, ra.configFile); err != nil {
			return fmt.Sprintf("Reload failed: %s, failed to revert to last known good config file", output), err
		}
		if err := ra.restartHAProxy(); err != nil {
			log.Warn("Restart failed, please check the reason and restart manually: ", err)
			return fmt.Sprintf("Reload failed: %s, failed to restart HAProxy, please check and start manually", output), err
		}
		log.Debug("HAProxy restarted with last known good config.")
		return output, reloadFailedError
	}
	log.Debug("Reload successful")
	// if success, replace last known good file
	// nolint:errcheck
	copyFile(ra.configFile, ra.lkgConfigFile)
	return output, nil
}

func (ra *ReloadAgent) restartHAProxy() error {
	_, err := execCmd(ra.restartCmd)
	if err != nil {
		return err
	}
	return nil
}

func execCmd(cmd string) (string, error) {
	strArr := strings.Split(cmd, " ")
	var c *exec.Cmd
	if len(strArr) == 1 {
		//nolint:gosec
		c = exec.Command(strArr[0])
	} else {
		//nolint:gosec
		c = exec.Command(strArr[0], strArr[1:]...)
	}
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		return stderr.String(), fmt.Errorf("executing %s failed: %s", cmd, err)
	}
	return stdout.String(), nil
}

// Reload schedules a reload
func (ra *ReloadAgent) Reload() string {
	if ra.cache.next == "" {
		ra.cache.newReload()
	}
	return ra.cache.next
}

// ForceReload calls reload directly
func (ra *ReloadAgent) ForceReload() error {
	r, err := ra.reloadHAProxy()
	if err != nil {
		return NewReloadError(fmt.Sprintf("Reload failed: %v, %v", err, r))
	}
	return nil
}

func (rc *reloadCache) Init(retention int) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.failedReloads = make(map[string]*models.Reload)
	rc.current = ""
	rc.next = ""
	rc.lastSuccess = nil
	rc.index = 0
	rc.retention = retention
}

func (rc *reloadCache) newReload() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.next = rc.generateID()
}

func (rc *reloadCache) failReload(response string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	r := &models.Reload{
		ID:              rc.current,
		Status:          "failed",
		Response:        response,
		ReloadTimestamp: time.Now().Unix(),
	}

	rc.failedReloads[rc.current] = r
	rc.current = ""
	rc.clearReloads()
}

func (rc *reloadCache) succeedReload(response string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	r := &models.Reload{
		ID:              rc.current,
		Status:          "succeeded",
		Response:        response,
		ReloadTimestamp: time.Now().Unix(),
	}

	rc.lastSuccess = r
	rc.current = ""
}

func (rc *reloadCache) clearReloads() {
	now := time.Now().Unix()

	for k, v := range rc.failedReloads {
		if (now - v.ReloadTimestamp) > int64((rc.retention * 86400)) {
			delete(rc.failedReloads, k)
		}
	}
}

func (ra *ReloadAgent) GetReloads() models.Reloads {
	ra.cache.mu.RLock()
	defer ra.cache.mu.RUnlock()

	v := make([]*models.Reload, 0, len(ra.cache.failedReloads))
	for _, value := range ra.cache.failedReloads {
		v = append(v, value)
	}

	if ra.cache.lastSuccess != nil {
		v = append(v, ra.cache.lastSuccess)
	}

	if ra.cache.current != "" {
		r := &models.Reload{
			ID:     ra.cache.current,
			Status: "in_progress",
		}
		v = append(v, r)
	}

	if ra.cache.next != "" {
		r := &models.Reload{
			ID:     ra.cache.next,
			Status: "in_progress",
		}
		v = append(v, r)
	}
	return v
}

func (ra *ReloadAgent) GetReload(id string) *models.Reload {
	ra.cache.mu.RLock()
	defer ra.cache.mu.RUnlock()

	if ra.cache.current == id {
		return &models.Reload{
			ID:     ra.cache.current,
			Status: "in_progress",
		}
	}
	if ra.cache.next == id {
		return &models.Reload{
			ID:     ra.cache.current,
			Status: "in_progress",
		}
	}

	v, ok := ra.cache.failedReloads[id]
	if ok {
		return v
	}
	if ra.cache.lastSuccess != nil {
		if ra.cache.lastSuccess.ID == id {
			return ra.cache.lastSuccess
		}

		// if it is older than last success return success
		sDate, sIndex, err := getTimeIndexFromID(ra.cache.lastSuccess.ID)
		if err != nil {
			return nil
		}
		gDate, gIndex, err := getTimeIndexFromID(id)
		if err != nil {
			return nil
		}

		if gDate.Before(sDate) {
			return &models.Reload{
				ID:     id,
				Status: "succeeded",
			}
		}

		if sIndex > gIndex {
			return &models.Reload{
				ID:     id,
				Status: "succeeded",
			}
		}
	}
	return nil
}

func (ra *ReloadAgent) Restart() error {
	return ra.restartHAProxy()
}

func (rc *reloadCache) generateID() string {
	defer func() {
		rc.index++
	}()
	return fmt.Sprintf("%s-%v", time.Now().Format("2006-01-02"), rc.index)
}

func getTimeIndexFromID(id string) (time.Time, int64, error) {
	data := strings.Split(id, "-")
	index, err := strconv.ParseInt(data[len(data)-1], 10, 64)
	if err != nil {
		return time.Now(), 0, err
	}
	date, err := time.Parse("2006-01-02", strings.Join(data[:len(data)-1], "-"))

	if err != nil {
		return date, 0, err
	}

	return date, index, nil
}

// ReloadError general configuration client error
type ReloadError struct {
	msg string
}

// Error implementation for ConfError
func (e *ReloadError) Error() string {
	return fmt.Sprintf(e.msg)
}

// NewReloadError constructor for ReloadError
func NewReloadError(msg string) *ReloadError {
	return &ReloadError{msg: msg}
}

func copyFile(src, dest string) error {
	srcContent, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcContent.Close()

	destContent, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destContent.Close()

	_, err = io.Copy(destContent, srcContent)
	if err != nil {
		return err
	}
	return destContent.Sync()
}
