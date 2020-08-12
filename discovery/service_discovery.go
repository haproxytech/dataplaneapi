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
	"sync"

	"github.com/haproxytech/client-native/v2/configuration"
)

//ServiceDiscoveryParams configuration for a specific service discovery
type ServiceDiscoveryParams interface{}

//ServiceDiscovery represents the required methods for a service discovery
type ServiceDiscovery interface {
	AddNode(id string, params ServiceDiscoveryParams) error
	GetNode(id string) (ServiceDiscoveryParams, error)
	GetNodes() (ServiceDiscoveryParams, error)
	RemoveNode(id string) error
	UpdateNode(id string, params ServiceDiscoveryParams) error
}

// ServiceDiscoveries manages all registered service discovery services
type ServiceDiscoveries interface {
	AddNode(serviceName string, id string, params ServiceDiscoveryParams) error
	AddService(serviceName string, serviceImpl ServiceDiscovery) error
	GetNode(serviceName string, id string) (ServiceDiscoveryParams, error)
	GetNodes(serviceName string) (ServiceDiscoveryParams, error)
	RemoveNode(serviceName string, id string) error
	RemoveService(serviceName string) error
	UpdateNode(serviceName string, id string, params ServiceDiscoveryParams) error
}

//NewServiceDiscoveries creates a new ServiceDiscoveries instance
func NewServiceDiscoveries(client *configuration.Client) ServiceDiscoveries {
	sd := &serviceDiscoveryImpl{
		services: make(map[string]ServiceDiscovery),
		client:   client,
	}
	//nolint
	sd.AddService("consul", NewConsulDiscoveryService(client))
	return sd
}

type serviceDiscoveryImpl struct {
	client   *configuration.Client
	services map[string]ServiceDiscovery
	mu       sync.RWMutex
}

func (s *serviceDiscoveryImpl) AddNode(serviceName string, id string, params ServiceDiscoveryParams) error {
	s.mu.RLock()
	sd, ok := s.services[serviceName]
	s.mu.RUnlock()
	if !ok {
		return errors.New("service not found")
	}
	return sd.AddNode(id, params)
}

func (s *serviceDiscoveryImpl) AddService(serviceName string, serviceImpl ServiceDiscovery) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.services[serviceName]
	if ok {
		return errors.New("service already exists")
	}
	s.services[serviceName] = serviceImpl
	return nil
}

func (s *serviceDiscoveryImpl) GetNode(serviceName string, id string) (ServiceDiscoveryParams, error) {
	s.mu.RLock()
	sd, ok := s.services[serviceName]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("service not found")
	}
	return sd.GetNode(id)
}

func (s *serviceDiscoveryImpl) GetNodes(serviceName string) (ServiceDiscoveryParams, error) {
	s.mu.RLock()
	sd, ok := s.services[serviceName]
	s.mu.RUnlock()
	if !ok {
		return nil, errors.New("service not found")
	}
	return sd.GetNodes()
}

func (s *serviceDiscoveryImpl) RemoveNode(serviceName string, id string) error {
	s.mu.RLock()
	sd, ok := s.services[serviceName]
	s.mu.RUnlock()
	if !ok {
		return errors.New("service not found")
	}
	return sd.RemoveNode(id)
}

func (s *serviceDiscoveryImpl) RemoveService(serviceName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.services[serviceName]
	if !ok {
		return errors.New("service not found")
	}
	delete(s.services, serviceName)
	return nil
}

func (s *serviceDiscoveryImpl) UpdateNode(serviceName string, id string, params ServiceDiscoveryParams) error {
	s.mu.RLock()
	sd, ok := s.services[serviceName]
	s.mu.RUnlock()
	if !ok {
		return errors.New("service not found")
	}
	return sd.UpdateNode(id, params)
}
