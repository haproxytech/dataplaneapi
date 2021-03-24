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
	"sync/atomic"
)

type thresholdLimit struct {
	limit  *uint64
	actual *uint64
}

func (t *thresholdLimit) LimitReached() (err error) {
	if *t.actual >= *t.limit {
		err = NewTransactionLimitReachedError(*t.limit)
	}
	return
}

func NewThresholdLimit(limit uint64, startingFrom uint64) Threshold {
	return &thresholdLimit{
		actual: &startingFrom,
		limit:  &limit,
	}
}

func (t *thresholdLimit) Increase() {
	atomic.AddUint64(t.actual, 1)
}

func (t thresholdLimit) Decrease() {
	atomic.AddUint64(t.actual, ^uint64(0))
}
