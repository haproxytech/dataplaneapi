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

package misc

import (
	"math/rand"
	"testing"
)

func TestRandomString(t *testing.T) {
	for i := 0; i < 1024; i++ {
		size := rand.Intn(512)
		str, err := RandomString(size)
		if err != nil {
			t.Errorf("RandomString returned an error for size %d: %v", size, err)
		}
		if len(str) != size {
			t.Errorf("RandomString returned a string of length %d for size %d", len(str), size)
		}
	}
}
