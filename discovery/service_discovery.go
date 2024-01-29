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
	"context"
	"errors"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/dataplaneapi/haproxy"
)

// ServiceDiscoveryParams configuration for a specific service discovery
type ServiceDiscoveryParams interface{}

// ServiceDiscovery represents the required methods for a service discovery
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

// ServiceDiscoveriesParams contain the parameters for the service discovery initialization
type ServiceDiscoveriesParams struct {
	Client      configuration.Configuration
	ReloadAgent haproxy.IReloadAgent
	Context     context.Context
}

// NewServiceDiscoveries creates a new ServiceDiscoveries instance
func NewServiceDiscoveries(params ServiceDiscoveriesParams) ServiceDiscoveries {
	sd := &serviceDiscoveryImpl{
		services: NewInstanceStore(),
	}
	sd.AddService("consul", NewConsulDiscoveryService(params))
	_ = sd.AddService("aws", NewAWSDiscoveryService(params))
	return sd
}

type serviceDiscoveryImpl struct {
	services Store
}

func (s *serviceDiscoveryImpl) AddNode(serviceName string, id string, params ServiceDiscoveryParams) error {
	sd, err := s.services.Read(serviceName)
	if err != nil {
		return errors.New("service not found")
	}
	return sd.(ServiceDiscovery).AddNode(id, params)
}

func (s *serviceDiscoveryImpl) AddService(serviceName string, serviceImpl ServiceDiscovery) error {
	if err := s.services.Create(serviceName, serviceImpl); err != nil {
		return errors.New("service already exists")
	}
	return nil
}

func (s *serviceDiscoveryImpl) GetNode(serviceName string, id string) (ServiceDiscoveryParams, error) {
	sd, err := s.services.Read(serviceName)
	if err != nil {
		return nil, errors.New("service not found")
	}
	return sd.(ServiceDiscovery).GetNode(id)
}

func (s *serviceDiscoveryImpl) GetNodes(serviceName string) (ServiceDiscoveryParams, error) {
	sd, err := s.services.Read(serviceName)
	if err != nil {
		return nil, errors.New("service not found")
	}
	return sd.(ServiceDiscovery).GetNodes()
}

func (s *serviceDiscoveryImpl) RemoveNode(serviceName string, id string) error {
	sd, err := s.services.Read(serviceName)
	if err != nil {
		return errors.New("service not found")
	}
	return sd.(ServiceDiscovery).RemoveNode(id)
}

func (s *serviceDiscoveryImpl) RemoveService(serviceName string) error {
	if err := s.services.Delete(serviceName); err != nil {
		return errors.New("service not found")
	}
	return nil
}

func (s *serviceDiscoveryImpl) UpdateNode(serviceName string, id string, params ServiceDiscoveryParams) error {
	sd, err := s.services.Read(serviceName)
	if err != nil {
		return errors.New("service not found")
	}
	return sd.(ServiceDiscovery).UpdateNode(id, params)
}
