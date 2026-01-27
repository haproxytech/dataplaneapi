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

package configuration

import (
	"sync/atomic"
)

type AtomicBool struct {
	value atomic.Value
}

func (s *AtomicBool) Load() bool {
	v := s.value.Load()
	if v == nil {
		return false
	}
	return v.(bool)
}

func (s *AtomicBool) Store(str bool) {
	s.value.Store(str)
}

func (s *AtomicBool) String() string {
	if s.Load() {
		return "true"
	}
	return "false"
}

func (s *AtomicBool) UnmarshalYAML(unmarshal func(any) error) error {
	var buf bool
	err := unmarshal(&buf)
	if err != nil {
		return err
	}

	s.Store(buf)
	return nil
}

func (s AtomicBool) MarshalYAML() (any, error) {
	return s.Load(), nil
}
