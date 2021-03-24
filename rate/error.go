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

import "fmt"

func NewTransactionLimitReachedError(limit uint64) error {
	return &TransactionLimitReachedErr{limit: limit}
}

type TransactionLimitReachedErr struct {
	limit uint64
}

func (l TransactionLimitReachedErr) Error() string {
	return fmt.Sprintf("cannot start a new transaction, reached the maximum amount of %d active transactions available", l.limit)
}
