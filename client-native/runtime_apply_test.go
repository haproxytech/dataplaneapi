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
	"errors"
	"testing"
	"time"

	clientnative "github.com/haproxytech/client-native/v6"
	configuration_options "github.com/haproxytech/client-native/v6/configuration/options"
	runtime_api "github.com/haproxytech/client-native/v6/runtime"
	"github.com/stretchr/testify/assert"

	"github.com/haproxytech/dataplaneapi/reload_agent"
)

type fakeClient struct {
	clientnative.HAProxyClient
	rt  runtime_api.Runtime
	err error
}

func (c fakeClient) Runtime() (runtime_api.Runtime, error) { return c.rt, c.err }

type fakeRuntime struct{ runtime_api.Runtime }

type fakeAgent struct {
	reload_agent.IReloadAgent
	asked  time.Time
	synced bool
}

func (a *fakeAgent) Synced(t time.Time) bool {
	a.asked = t
	return a.synced
}

func resetRuntimeApply(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		runtimeApplyState.mu.Lock()
		defer runtimeApplyState.mu.Unlock()
		runtimeApplyState.client = nil
		runtimeApplyState.ra = nil
		runtimeApplyState.needReloadAt = time.Time{}
		runtimeApplyState.outcomes = nil
	})
}

func TestApplyRuntimeUnbound(t *testing.T) {
	resetRuntimeApply(t)
	assert.True(t, applyRuntime() == nil)
}

func TestApplyRuntimeSynced(t *testing.T) {
	resetRuntimeApply(t)
	rt := &fakeRuntime{}
	BindRuntimeApply(fakeClient{rt: rt}, &fakeAgent{synced: true})
	assert.Same(t, rt, applyRuntime())
}

func TestApplyRuntimeNotSynced(t *testing.T) {
	resetRuntimeApply(t)
	BindRuntimeApply(fakeClient{rt: &fakeRuntime{}}, &fakeAgent{})
	assert.True(t, applyRuntime() == nil)
}

func TestApplyRuntimeClientError(t *testing.T) {
	resetRuntimeApply(t)
	BindRuntimeApply(fakeClient{rt: &fakeRuntime{}, err: errors.New("no runtime")}, &fakeAgent{synced: true})
	assert.True(t, applyRuntime() == nil)
}

func TestCommitDoneReload(t *testing.T) {
	resetRuntimeApply(t)
	ra := &fakeAgent{synced: true}
	BindRuntimeApply(fakeClient{rt: &fakeRuntime{}}, ra)

	before := time.Now()
	commitDone(configuration_options.CommitResult{Version: 7, Reload: true, Err: errors.New("boom")})
	applyRuntime()
	assert.False(t, ra.asked.Before(before))
	assert.True(t, NeedsReload(7))
}

func TestCommitDoneNoReload(t *testing.T) {
	resetRuntimeApply(t)
	commitDone(configuration_options.CommitResult{Version: 3})
	assert.False(t, NeedsReload(3))
	assert.True(t, NeedsReload(3))
}

func TestNeedsReloadUnknownVersion(t *testing.T) {
	resetRuntimeApply(t)
	assert.True(t, NeedsReload(42))
}

func TestCommitDoneBounded(t *testing.T) {
	resetRuntimeApply(t)
	for v := range int64(maxOutcomes) {
		commitDone(configuration_options.CommitResult{Version: v})
	}
	assert.Len(t, runtimeApplyState.outcomes, maxOutcomes)
	commitDone(configuration_options.CommitResult{Version: maxOutcomes})
	assert.Len(t, runtimeApplyState.outcomes, 1)
	assert.True(t, NeedsReload(0))
	assert.False(t, NeedsReload(maxOutcomes))
}
