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

package rate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_thresholdLimit_LimitReached(t *testing.T) {
	t.Run("ko", func(t *testing.T) {
		tl := thresholdLimit{
			actual: func() uint64 {
				return 10
			},
			limit: 10,
		}
		assert.Error(t, tl.LimitReached())
	})
	t.Run("ok", func(t *testing.T) {
		tl := thresholdLimit{
			actual: func() uint64 {
				return 0
			},
			limit: 10,
		}
		require.NoError(t, tl.LimitReached())
	})
}
