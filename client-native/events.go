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
	listener  *runtime.EventListener
	client    clientnative.HAProxyClient // for storage only
	rt        runtime.Runtime
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
		// no need to restart the listener
		h.rt = rt
		return nil
	}

	h.Reset()
	h.rt = rt
	h.stop.Store(false)
	go h.listen(ctx)

	return nil
}

func (h *HAProxyEventListener) Reset() {
	if h.listener != nil {
		if err := h.listener.Close(); err != nil {
			log.Warning(err)
		}
		h.listener = nil
	}
}

func (h *HAProxyEventListener) Stop() error {
	h.stop.Store(true)
	return h.listener.Close()
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

func (h *HAProxyEventListener) listen(ctx context.Context) {
	var err error
	retryAfter := 100 * time.Millisecond

	for {
		if h.stop.Load() {
			// Stop requested.
			h.Reset()
			return
		}
		if h.listener == nil {
			h.listener, err = newHAProxyEventListener(h.rt.SocketPath())
			if err != nil {
				// Try again.
				log.Warning(err)
				time.Sleep(retryAfter)
				retryAfter *= 2
				if retryAfter == 51200*time.Millisecond {
					// Give up after 10 iterations.
					h.stop.Store(true)
				}
				continue
			}
		}

		for {
			ev, err := h.listener.Listen(ctx)
			if err != nil {
				// EOF errors usually happen when HAProxy restarts, do not log.
				if !errors.Is(err, io.EOF) {
					log.Warning(err)
				}
				// Reset the connection.
				h.Reset()
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
