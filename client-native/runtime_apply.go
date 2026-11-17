// Copyright 2026 HAProxy Technologies
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

package cn

import (
	"sync"
	"time"

	clientnative "github.com/haproxytech/client-native/v6"
	configuration_options "github.com/haproxytech/client-native/v6/configuration/options"
	runtime_api "github.com/haproxytech/client-native/v6/runtime"

	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/reload_agent"
)

// maxOutcomes bounds the remembered commit outcomes; a lost entry only means "reload".
const maxOutcomes = 256

var runtimeApplyState struct {
	client       clientnative.HAProxyClient
	ra           reload_agent.IReloadAgent
	needReloadAt time.Time
	outcomes     map[int64]bool
	paused       int
	mu           sync.Mutex
}

// BindRuntimeApply lets committed changes be applied through the runtime API
// when client-native finds they do not need a reload. Until it is called every
// commit needs a reload.
func BindRuntimeApply(client clientnative.HAProxyClient, ra reload_agent.IReloadAgent) {
	runtimeApplyState.mu.Lock()
	defer runtimeApplyState.mu.Unlock()
	runtimeApplyState.client = client
	runtimeApplyState.ra = ra
}

// PauseRuntimeApply makes commits need a reload until the returned function is
// called. It is for requests that must not change the running HAProxy.
func PauseRuntimeApply() (resume func()) {
	runtimeApplyState.mu.Lock()
	runtimeApplyState.paused++
	runtimeApplyState.mu.Unlock()
	var once sync.Once
	return func() {
		once.Do(func() {
			runtimeApplyState.mu.Lock()
			runtimeApplyState.paused--
			runtimeApplyState.mu.Unlock()
		})
	}
}

// NeedsReload reports whether the commit made on top of the given
// configuration version still needs a reload, and forgets it. A version it
// knows nothing about needs one.
func NeedsReload(version int64) bool {
	runtimeApplyState.mu.Lock()
	defer runtimeApplyState.mu.Unlock()
	reload, found := runtimeApplyState.outcomes[version]
	delete(runtimeApplyState.outcomes, version)
	return !found || reload
}

// applyRuntime returns the runtime to apply a commit with, or nil when runtime
// apply is paused or HAProxy may not run the configuration file as it was
// before the commit.
func applyRuntime() runtime_api.Runtime {
	runtimeApplyState.mu.Lock()
	client, ra, needReloadAt, paused := runtimeApplyState.client, runtimeApplyState.ra, runtimeApplyState.needReloadAt, runtimeApplyState.paused > 0
	runtimeApplyState.mu.Unlock()
	if paused || client == nil || ra == nil || !ra.Synced(needReloadAt) {
		return nil
	}
	rt, err := client.Runtime()
	if err != nil {
		return nil
	}
	return rt
}

// commitDone records the outcome of a commit.
func commitDone(res configuration_options.CommitResult) {
	if res.Err != nil {
		log.Warningf("runtime apply failed, reloading: %v", res.Err)
	}
	runtimeApplyState.mu.Lock()
	defer runtimeApplyState.mu.Unlock()
	if res.Reload {
		runtimeApplyState.needReloadAt = time.Now()
	}
	if len(runtimeApplyState.outcomes) >= maxOutcomes {
		clear(runtimeApplyState.outcomes)
	}
	if runtimeApplyState.outcomes == nil {
		runtimeApplyState.outcomes = map[int64]bool{}
	}
	runtimeApplyState.outcomes[res.Version] = res.Reload
}
