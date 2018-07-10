package haproxy

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/haproxytech/client-native/runtime"
)

// ReloadAgent handles all reloads, scheduled or forced
type ReloadAgent struct {
	reloadDelay   int
	reloads       chan int
	needsReload   bool
	reloadCmd     string
	runtimeClient *runtime.Client
}

// Init a new reload agent
func (ra *ReloadAgent) Init(delay int, cmd string, cli *runtime.Client) {
	ra.reloadDelay = delay
	ra.needsReload = false
	ra.reloads = make(chan int)
	ra.reloadCmd = cmd
	ra.runtimeClient = cli
	go ra.handleReloads()
}

func (ra *ReloadAgent) handleReloads() {
	for {
		select {
		case <-ra.reloads:
			ra.needsReload = true
		case <-time.After(time.Duration(10) * time.Second):
			if ra.needsReload {
				err := ra.reloadHAProxy()
				if err != nil {
					log.Warning("Reload failed " + err.Error())
				} else {
					ra.needsReload = false
				}
			}
		}
	}
}

func (ra *ReloadAgent) reloadHAProxy() error {
	strArr := strings.Split(ra.reloadCmd, " ")
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
		return fmt.Errorf("%v: %v", err.Error(), string(stderr.Bytes()))
	}
	return nil
}

// Reload schedules a reload
func (ra *ReloadAgent) Reload() {
	ra.reloads <- 0
}

// ForceReload calls reload directly
func (ra *ReloadAgent) ForceReload() error {
	err := ra.reloadHAProxy()
	if err != nil {
		return err
	}
	return nil
}
