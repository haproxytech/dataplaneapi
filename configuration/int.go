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
	"fmt"
	"sync/atomic"
)

type AtomicInt struct {
	value atomic.Value
}

func (s *AtomicInt) Load() int {
	v := s.value.Load()
	if v == nil {
		return 0
	}
	return v.(int)
}

func (s *AtomicInt) Store(str int) {
	s.value.Store(str)
}

func (s *AtomicInt) String() string {
	return fmt.Sprintf("%d", s.Load())
}

func (s *AtomicInt) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var buf int
	err := unmarshal(&buf)
	if err != nil {
		return err
	}

	s.Store(buf)
	return nil
}

func (s AtomicInt) MarshalYAML() (interface{}, error) {
	return s.Load(), nil
}
