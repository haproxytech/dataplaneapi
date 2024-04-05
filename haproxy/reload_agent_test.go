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
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReloadAgentDoesntMissReloads(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	f, err := os.CreateTemp("", "config.cfg")
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
