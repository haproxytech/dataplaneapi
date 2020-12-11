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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewThresholdLimit(t *testing.T) {
	type tc struct {
		limit  uint64
		actual uint64
	}
	for name, tc := range map[string]tc{
		"from zero": {10, 0},
		"reached":   {5, 5},
	} {
		t.Run(name, func(t *testing.T) {
			l := NewThresholdLimit(tc.limit, tc.actual).(*thresholdLimit)
			assert.Equal(t, *l.limit, tc.limit)
			assert.Equal(t, *l.actual, tc.actual)
		})
	}
}

func Test_thresholdLimit_Decrease(t *testing.T) {
	for _, tc := range []uint64{5, 10} {
		t.Run(fmt.Sprintf("%d", tc), func(t *testing.T) {
			l := thresholdLimit{
				limit: func(v uint64) *uint64 {
					return &v
				}(tc),
				actual: func(v uint64) *uint64 {
					return &v
				}(tc),
			}
			var counter uint64
			for *l.actual > 0 {
				l.Decrease()
				counter++
			}
			assert.Equal(t, tc, counter)
		})
	}
}

func Test_thresholdLimit_Increase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		l := thresholdLimit{
			limit: func(v uint64) *uint64 {
				return &v
			}(10),
			actual: func(v uint64) *uint64 {
				return &v
			}(10),
		}
		var counter int
		for *l.actual < *l.limit {
			l.Increase()
			counter++
		}
	})
	t.Run("failure", func(t *testing.T) {
		l := thresholdLimit{
			limit: func(v uint64) *uint64 {
				return &v
			}(10),
			actual: func(v uint64) *uint64 {
				return &v
			}(10),
		}
		assert.NotNil(t, l.LimitReached())
	})
}
