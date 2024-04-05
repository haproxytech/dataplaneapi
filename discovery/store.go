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

package discovery

import (
	"errors"
	"fmt"
	"sync"
)

type Store interface {
	Create(name string, service interface{}) error
	Read(name string) (interface{}, error)
	Update(name string, mutateFn func(obj interface{}) error) (err error)
	Delete(name string) error
	List() []interface{}
}

type instanceStore struct {
	store map[string]interface{}
	mu    sync.RWMutex
}

func (s *instanceStore) Update(name string, mutateFn func(obj interface{}) error) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var o interface{}
	if o, err = s.get(name); err != nil {
		return fmt.Errorf("cannot update resource: %w", err)
	}

	if err = mutateFn(o); err != nil {
		return fmt.Errorf("cannot update resource: %w", err)
	}

	s.store[name] = o
	return
}

func (s *instanceStore) List() (list []interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, i := range s.store {
		list = append(list, i)
	}
	return
}

func NewInstanceStore() Store {
	return &instanceStore{
		store: map[string]interface{}{},
		mu:    sync.RWMutex{},
	}
}

func (s *instanceStore) Create(name string, service interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.get(name); err == nil {
		return errors.New("instance already exists")
	}

	s.store[name] = service
	return nil
}

func (s *instanceStore) Delete(name string) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err = s.get(name); err != nil {
		return
	}

	delete(s.store, name)
	return nil
}

func (s *instanceStore) Read(name string) (sd interface{}, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.get(name)
}

func (s *instanceStore) get(name string) (sd interface{}, err error) {
	var ok bool
	sd, ok = s.store[name]
	if !ok {
		return nil, errors.New("instance not found")
	}
	return
}
