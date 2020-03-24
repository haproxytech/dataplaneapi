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

type AtomicString struct {
	value atomic.Value
}

func (s *AtomicString) Load() string {
	v := s.value.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

func (s *AtomicString) Store(str string) {
	s.value.Store(str)
}

func (s *AtomicString) String() string {
	return s.Load()
}

func (s *AtomicString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf string
	err := unmarshal(&buf)
	if err != nil {
		return err
	}

	s.Store(buf)
	return nil
}

func (s AtomicString) MarshalYAML() (interface{}, error) {
	return s.Load(), nil
}
