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
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type ChanNotify struct {
	mu          sync.RWMutex
	subscribers map[string]chan struct{}
}

func NewChanNotify() *ChanNotify {
	cn := &ChanNotify{}
	cn.subscribers = make(map[string]chan struct{})
	return cn
}

func (cn *ChanNotify) Subscribe(name string) chan struct{} {
	cn.mu.Lock()
	defer cn.mu.Unlock()

	c := make(chan struct{}, 1)
	cn.subscribers[name] = c
	return c
}

func (cn *ChanNotify) UnSubscribeAll() {
	cn.mu.Lock()
	defer cn.mu.Unlock()
	cn.subscribers = make(map[string]chan struct{})
}

func (cn *ChanNotify) Notify() {
	cn.notify(0)
}

func (cn *ChanNotify) NotifyWithRetry() {
	cn.notify(3)
}

func (cn *ChanNotify) notify(numTry int) {
	if numTry < 0 {
		return
	}
	cn.mu.RLock()
	defer cn.mu.RUnlock()

	if len(cn.subscribers) == 0 {
		go func() {
			time.Sleep(2 * time.Second)
			numTry--
			cn.notify(numTry)
		}()
		return
	}

	for _, c := range cn.subscribers {
		c <- struct{}{}
	}
}

func (c *Configuration) initSignalHandler() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-osSignals
		log.Println(sig)
		c.Notify.Shutdown.Notify()
	}()

	osSignals2 := make(chan os.Signal, 1)
	signal.Notify(osSignals2, syscall.SIGHUP)

	go func() {
		for {
			<-osSignals2
			c.Notify.Reload.Notify()
		}
	}()
}
