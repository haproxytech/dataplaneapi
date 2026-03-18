// Copyright 2025 HAProxy Technologies
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

package cn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	clientnative "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/runtime"
	"github.com/haproxytech/dataplaneapi/log"
)

// Listen for HAProxy's event on the master socket.

// Events categgories
const (
	EventAcme = "acme"
)

type HAProxyEventListener struct {
	mu       sync.Mutex
	listener *runtime.EventListener
	client   clientnative.HAProxyClient // for storage only
	rt       runtime.Runtime
	// done is closed to signal the listen goroutine to exit.
	// A new channel is created each time a new goroutine is started.
	done      chan struct{}
	stop      atomic.Bool
	lastEvent time.Time
}

var (
	ErrNoMasterSocket = errors.New("master socket not configured")
	ErrOldVersion     = errors.New("this version of HAProxy does not support event sinks")
)

func ListenHAProxyEvents(ctx context.Context, client clientnative.HAProxyClient) (*HAProxyEventListener, error) {
	rt, err := client.Runtime()
	if err != nil {
		return nil, err
	}
	if rt == nil || rt.IsStatsSocket() {
		return nil, ErrNoMasterSocket
	}

	version, err := rt.GetVersion()
	if err != nil {
		return nil, err
	}
	// v3.2+
	if version.Major < 3 || (version.Major == 3 && version.Minor < 2) {
		return nil, ErrOldVersion
	}

	el, err := newHAProxyEventListener(rt.SocketPath())
	if err != nil {
		return nil, err
	}

	h := &HAProxyEventListener{
		client:   client,
		listener: el,
		rt:       rt,
		done:     make(chan struct{}),
	}

	go h.listen(ctx)

	log.Debugf("listening for HAProxy events on: %s", rt.SocketPath())

	return h, nil
}

// Reconfigure a running listener with a new Runtime.
func (h *HAProxyEventListener) Reconfigure(ctx context.Context, rt runtime.Runtime) error {
	if rt == nil || rt.IsStatsSocket() {
		return ErrNoMasterSocket
	}

	if rt.SocketPath() == h.rt.SocketPath() {
		h.mu.Lock()
		h.rt = rt
		h.mu.Unlock()
		return nil
	}

	// Signal the old goroutine to stop and wait for it to exit.
	h.stopAndWait()

	h.mu.Lock()
	h.rt = rt
	h.done = make(chan struct{})
	h.stop.Store(false)
	h.mu.Unlock()

	go h.listen(ctx)

	return nil
}

func (h *HAProxyEventListener) resetLocked() {
	if h.listener != nil {
		if err := h.listener.Close(); err != nil {
			log.Warning(err)
		}
		h.listener = nil
	}
}

// stopAndWait signals the listen goroutine to stop and waits for it to exit.
func (h *HAProxyEventListener) stopAndWait() {
	h.stop.Store(true)
	h.mu.Lock()
	h.resetLocked()
	h.mu.Unlock()
	// Wait for the goroutine to acknowledge the stop.
	<-h.done
}

func (h *HAProxyEventListener) Stop() error {
	h.stopAndWait()
	return nil
}

func newHAProxyEventListener(socketPath string) (*runtime.EventListener, error) {
	// This is both the connect and write timeout.
	// Use a small value here since at this point dataplane is supposed
	// to be already connected to the master socket.
	timeout := 3 * time.Second

	el, err := runtime.NewEventListener("unix", socketPath, "dpapi", timeout, "-w", "-0")
	if err != nil {
		return nil, fmt.Errorf("could not listen to HAProxy's events: %w", err)
	}

	return el, nil
}

const maxRetryBackoff = 60 * time.Second

func (h *HAProxyEventListener) listen(ctx context.Context) {
	defer close(h.done)

	var err error
	retryAfter := 100 * time.Millisecond
	loggedDisconnect := false

	for {
		if h.stop.Load() {
			h.mu.Lock()
			h.resetLocked()
			h.mu.Unlock()
			return
		}

		h.mu.Lock()
		needsConnect := h.listener == nil
		h.mu.Unlock()

		if needsConnect {
			var el *runtime.EventListener
			el, err = newHAProxyEventListener(h.rt.SocketPath())
			if err != nil {
				if !loggedDisconnect {
					log.Warningf("event listener disconnected, reconnecting: %v", err)
					loggedDisconnect = true
				} else {
					log.Debugf("event listener reconnect attempt: %v", err)
				}
				select {
				case <-time.After(retryAfter):
				case <-ctx.Done():
					return
				}
				retryAfter *= 2
				if retryAfter >= maxRetryBackoff {
					log.Warning("event listener giving up reconnection attempts")
					h.stop.Store(true)
				}
				continue
			}

			h.mu.Lock()
			h.listener = el
			h.mu.Unlock()

			retryAfter = 100 * time.Millisecond
			loggedDisconnect = false
			log.Debugf("event listener reconnected to: %s", h.rt.SocketPath())
		}

		for {
			ev, err := h.listener.Listen(ctx)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					log.Debugf("event listener error: %v", err)
				}
				h.mu.Lock()
				h.resetLocked()
				h.mu.Unlock()
				break
			}

			h.handle(ev)

			if h.listener == nil { // just in case
				break
			}
		}
	}
}

func (h *HAProxyEventListener) handle(ev runtime.Event) {
	if !ev.Timestamp.After(h.lastEvent) {
		// Event already seen! Skip.
		log.Debugf("events: skipping already seen: '%s'", ev.String())
		return
	}
	h.lastEvent = ev.Timestamp

	log.Debugf("events: new: '%s'", ev.String())

	category, rest, ok := strings.Cut(ev.Message, " ")
	if !ok {
		log.Warningf("failed to parse HAProxy Event: '%s'", ev.Message)
		return
	}

	if category == EventAcme {
		h.handleAcmeEvent(rest)
		return
	}

	// Do not expect dataplaneapi to be able to handle all event types.
	log.Debugf("unknown HAProxy Event type: '%s'", ev.Message)
}
