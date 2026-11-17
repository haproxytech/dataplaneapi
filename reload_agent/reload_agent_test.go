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
	"os"
	"testing"
	"time"

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

func newSyncedTestAgent(t *testing.T, reloadCmd string) *ReloadAgent {
	t.Helper()
	t.Setenv("CI_DATAPLANE_RELOAD_DELAY_OVERRIDE", "50")
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	f, err := os.CreateTemp(t.TempDir(), "config.cfg")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	ra, err := NewReloadAgent(ReloadAgentParams{
		ReloadCmd:  reloadCmd,
		RestartCmd: "echo ok",
		ConfigFile: f.Name(),
		Retention:  1,
		Ctx:        ctx,
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
