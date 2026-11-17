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

package reload_agent

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReloadAgentDoesntMissReloads(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	f, err := os.CreateTemp(t.TempDir(), "config.cfg")
	require.NoError(t, err)
	assert.NotNil(t, f)
	t.Cleanup(func() {
		cancel()
		assert.NoError(t, os.Remove(f.Name()))
	})

	reloadAgentParams := ReloadAgentParams{
		Delay:      1,
		ReloadCmd:  `echo "systemctl reload haproxy"`,
		RestartCmd: `echo "systemctl restart haproxy"`,
		ConfigFile: f.Name(),
		BackupDir:  "",
		Retention:  1,
		Ctx:        ctx,
	}

	ra, err := NewReloadAgent(reloadAgentParams)
	require.NoError(t, err)
	assert.NotNil(t, ra)

	var reloadID, firstReloadID, secondReloadID string

	// trigger a reload
	reloadID = ra.Reload()
	assert.NotEmpty(t, reloadID)
	firstReloadID = reloadID

	// trigger another reload shortly after the first one but before the
	// delay has elapsed which should yield the first reload ID
	time.Sleep(10 * time.Millisecond)
	reloadID = ra.Reload()
	assert.EqualValues(t, firstReloadID, reloadID)

	// sleep for as long as the delay duration to mimic a slightly
	// slower DataplaneAPI operation
	time.Sleep(time.Duration(reloadAgentParams.Delay) * time.Second)

	// Since this is happening after the delay has elapsed, it should create
	// a new reload ID
	reloadID = ra.Reload()
	assert.NotEmpty(t, reloadID)
	secondReloadID = reloadID
	assert.NotEqualValues(t, firstReloadID, secondReloadID)
}

func newSyncedTestAgent(t *testing.T, reloadCmd string, workerPID ...func() (int64, error)) *ReloadAgent {
	t.Helper()
	t.Setenv("CI_DATAPLANE_RELOAD_DELAY_OVERRIDE", "50")
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	f, err := os.CreateTemp(t.TempDir(), "config.cfg")
	require.NoError(t, err)
	require.NoError(t, f.Close())
	var wp func() (int64, error)
	if len(workerPID) > 0 {
		wp = workerPID[0]
	}

	ra, err := NewReloadAgent(ReloadAgentParams{
		ReloadCmd:  reloadCmd,
		RestartCmd: "echo ok",
		ConfigFile: f.Name(),
		Retention:  1,
		Ctx:        ctx,
		WorkerPID:  wp,
	})
	require.NoError(t, err)
	return ra
}

func TestSyncedFreshAgent(t *testing.T) {
	ra := newSyncedTestAgent(t, "echo ok")
	assert.True(t, ra.Synced(time.Now().Add(-time.Minute)))
	assert.True(t, ra.Synced(time.Time{}))
}

func TestSyncedScheduledReload(t *testing.T) {
	ra := newSyncedTestAgent(t, "echo ok")
	before := time.Now()
	ra.Reload()
	assert.False(t, ra.Synced(time.Time{}))
	assert.False(t, ra.Synced(before))

	require.Eventually(t, func() bool { return ra.Synced(time.Time{}) }, 2*time.Second, 10*time.Millisecond)
	assert.True(t, ra.Synced(before))
	assert.False(t, ra.Synced(time.Now()))
}

func TestSyncedFailedReload(t *testing.T) {
	ra := newSyncedTestAgent(t, "false")
	created := time.Now()
	require.Error(t, ra.ForceReload())
	assert.False(t, ra.Synced(created))
	assert.True(t, ra.Synced(time.Time{}))
}

func TestSyncedDoesNotBlockDuringReload(t *testing.T) {
	ra := newSyncedTestAgent(t, "sleep 0.3")
	before := time.Now()
	done := make(chan error, 1)
	go func() { done <- ra.ForceReload() }()
	require.Eventually(t, ra.running.Load, time.Second, time.Millisecond)

	start := time.Now()
	assert.False(t, ra.Synced(time.Time{}))
	assert.Less(t, time.Since(start), 50*time.Millisecond)
	require.NoError(t, <-done)
	assert.True(t, ra.Synced(before))
	assert.False(t, ra.Synced(start))
}

func TestSyncedForceReloadWithoutSchedule(t *testing.T) {
	ra := newSyncedTestAgent(t, "echo ok")
	before := time.Now()
	assert.False(t, ra.Synced(before))
	require.NoError(t, ra.ForceReload())
	assert.True(t, ra.Synced(before))
}

func TestSyncedRestart(t *testing.T) {
	ra := newSyncedTestAgent(t, "echo ok")
	before := time.Now()
	require.NoError(t, ra.Restart())
	assert.True(t, ra.Synced(before))
}

// fakeWorker stands for the runtime API: it answers the PID the test sets and
// counts the queries.
type fakeWorker struct {
	pid   atomic.Int64
	fail  atomic.Bool
	calls atomic.Int32
}

func newFakeWorker(pid int64) *fakeWorker {
	w := &fakeWorker{}
	w.pid.Store(pid)
	return w
}

func (w *fakeWorker) PID() (int64, error) {
	w.calls.Add(1)
	if w.fail.Load() {
		return 0, errors.New("runtime API unavailable")
	}
	return w.pid.Load(), nil
}

func TestSyncedWaitsForNewWorker(t *testing.T) {
	w := newFakeWorker(100)
	ra := newSyncedTestAgent(t, "echo ok", w.PID)
	before := time.Now()
	id := ra.Reload()
	require.Eventually(t, func() bool {
		r := ra.GetReload(id)
		return r != nil && r.Status == models.ReloadStatusSucceeded
	}, 2*time.Second, 10*time.Millisecond)
	assert.Equal(t, int32(1), w.calls.Load(), "PID read once before the reload")

	// the old worker still answers
	assert.False(t, ra.Synced(before))
	assert.False(t, ra.Synced(before))
	assert.Equal(t, int32(3), w.calls.Load())
	// earlier checks fail first and do not query
	assert.False(t, ra.Synced(time.Now()))
	assert.Equal(t, int32(3), w.calls.Load())

	w.pid.Store(101)
	assert.True(t, ra.Synced(before))
	assert.Equal(t, int32(4), w.calls.Load())
	assert.True(t, ra.Synced(before))
	assert.True(t, ra.Synced(time.Time{}))
	assert.Equal(t, int32(4), w.calls.Load(), "no query once a new worker answered")
}

func TestSyncedWorkerPIDError(t *testing.T) {
	w := newFakeWorker(100)
	w.fail.Store(true)
	ra := newSyncedTestAgent(t, "echo ok", w.PID)
	before := time.Now()
	// an unknown PID before the reload leaves nothing to wait for
	require.NoError(t, ra.ForceReload())
	assert.True(t, ra.Synced(before))

	w.fail.Store(false)
	before = time.Now()
	require.NoError(t, ra.ForceReload())
	w.fail.Store(true)
	assert.False(t, ra.Synced(before))
	w.fail.Store(false)
	assert.False(t, ra.Synced(before))
	w.pid.Store(101)
	assert.True(t, ra.Synced(before))
}

func TestSyncedNonPositiveWorkerPID(t *testing.T) {
	w := newFakeWorker(0)
	ra := newSyncedTestAgent(t, "echo ok", w.PID)
	before := time.Now()
	require.NoError(t, ra.ForceReload())
	assert.True(t, ra.Synced(before))
	assert.Equal(t, int32(1), w.calls.Load())
}

func TestSyncedFailedReloadKeepsStaleWorker(t *testing.T) {
	w := newFakeWorker(100)
	ra := newSyncedTestAgent(t, "false", w.PID)
	require.Error(t, ra.ForceReload())
	assert.Equal(t, int64(0), ra.staleWorker.Load())
	assert.True(t, ra.Synced(time.Time{}))
	assert.Equal(t, int32(1), w.calls.Load())

	ra.reloadCmd = "echo ok"
	before := time.Now()
	require.NoError(t, ra.ForceReload())
	assert.Equal(t, int64(100), ra.staleWorker.Load())
	ra.reloadCmd = "false"
	w.pid.Store(200)
	require.Error(t, ra.ForceReload())
	assert.Equal(t, int64(100), ra.staleWorker.Load())
	assert.True(t, ra.Synced(before))
}

func TestSyncedForceReloadWithoutScheduleWaitsForNewWorker(t *testing.T) {
	w := newFakeWorker(100)
	ra := newSyncedTestAgent(t, "echo ok", w.PID)
	before := time.Now()
	require.NoError(t, ra.ForceReload())
	assert.False(t, ra.Synced(before))
	w.pid.Store(101)
	assert.True(t, ra.Synced(before))
}

func TestSyncedRestartWaitsForNewWorker(t *testing.T) {
	w := newFakeWorker(100)
	ra := newSyncedTestAgent(t, "echo ok", w.PID)
	before := time.Now()
	require.NoError(t, ra.Restart())
	assert.False(t, ra.Synced(before))
	w.pid.Store(101)
	assert.True(t, ra.Synced(before))
}
