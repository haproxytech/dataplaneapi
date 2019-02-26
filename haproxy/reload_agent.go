package haproxy

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/haproxytech/models"

	log "github.com/sirupsen/logrus"
)

type reloadCache struct {
	failedReloads map[string]*models.Reload
	lastSuccess   *models.Reload
	current       string
	index         int64
	retention     int
	mu            sync.RWMutex
}

// ReloadAgent handles all reloads, scheduled or forced
type ReloadAgent struct {
	delay int
	cmd   string
	cache reloadCache
}

// Init a new reload agent
func (ra *ReloadAgent) Init(delay int, cmd string, retention int) {
	ra.cmd = cmd
	ra.delay = delay

	ra.cache.Init(retention)
	go ra.handleReloads()
}

func (ra *ReloadAgent) handleReloads() {
	for {
		select {
		case <-time.After(time.Duration(ra.delay) * time.Second):
			if ra.cache.current != "" {
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
	strArr := strings.Split(ra.cmd, " ")
	var cmd *exec.Cmd
	if len(strArr) == 1 {
		cmd = exec.Command(strArr[0])
	} else {
		cmd = exec.Command(strArr[0], strArr[1:]...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return string(stderr.Bytes()), err
	}
	return string(stdout.Bytes()), nil
}

// Reload schedules a reload
func (ra *ReloadAgent) Reload() string {
	if ra.cache.current == "" {
		ra.cache.newReload()
	}
	return ra.cache.current
}

// ForceReload calls reload directly
func (ra *ReloadAgent) ForceReload() error {
	r, err := ra.reloadHAProxy()
	if err != nil {
		return fmt.Errorf("Reload failed: %v, %v", err, r)
	}
	return nil
}

func (rc *reloadCache) Init(retention int) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.failedReloads = make(map[string]*models.Reload)
	rc.current = ""
	rc.lastSuccess = nil
	rc.index = 0
	rc.retention = retention
}

func (rc *reloadCache) newReload() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.current = rc.generateID()
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
		Status:          "succeded",
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
				Status: "succeded",
			}
		}

		if sIndex > gIndex {
			return &models.Reload{
				ID:     id,
				Status: "succeded",
			}
		}
	}
	return nil
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
