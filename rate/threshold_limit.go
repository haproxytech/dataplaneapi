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

type thresholdLimit struct {
	actual func() uint64
	limit  uint64
}

func (t *thresholdLimit) LimitReached() (err error) {
	if t.actual() >= t.limit {
		err = NewTransactionLimitReachedError(t.limit)
	}
	return
}

func NewThresholdLimit(limit uint64, actual func() uint64) Threshold {
	return &thresholdLimit{
		actual: actual,
		limit:  limit,
	}
}
