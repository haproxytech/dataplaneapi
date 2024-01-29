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
	"sync"

	"github.com/haproxytech/client-native/v6/configuration"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
)

// Using a simple mutex to avoid race conditions when multiple Service Discovery instances are trying to commit
// changes at the same time: need to be refactored.
var mutex = &sync.Mutex{}

// ServiceInstance specifies the needed information required from the service to provide for the ServiceDiscoveryInstance.
type ServiceInstance interface {
	GetName() string
	GetBackendName() string
	GetFrom() string
	Changed() bool
	GetServers() []configuration.ServiceServer
}

type confService struct {
	confService *configuration.Service
	cleanup     bool
}

type discoveryInstanceParams struct {
	LogFields       map[string]interface{}
	SlotsGrowthType string
	Allowlist       []string
	Denylist        []string
	ServerSlotsBase int
	SlotsIncrement  int
}

// ServiceDiscoveryInstance manages and updates all services of a single service discovery.
type ServiceDiscoveryInstance struct {
	client        configuration.Configuration
	reloadAgent   haproxy.IReloadAgent
	services      map[string]*confService
	transactionID string
	params        discoveryInstanceParams
}

// NewServiceDiscoveryInstance creates a new ServiceDiscoveryInstance.
func NewServiceDiscoveryInstance(client configuration.Configuration, reloadAgent haproxy.IReloadAgent, params discoveryInstanceParams) *ServiceDiscoveryInstance {
	return &ServiceDiscoveryInstance{
		client:      client,
		reloadAgent: reloadAgent,
		params:      params,
		services:    make(map[string]*confService),
	}
}

// UpdateParams updates the scaling params for each service associated with the service discovery.
func (s *ServiceDiscoveryInstance) UpdateParams(params discoveryInstanceParams) error {
	s.params = params
	for _, se := range s.services {
		err := se.confService.UpdateScalingParams(configuration.ScalingParams{
			BaseSlots:       s.params.ServerSlotsBase,
			SlotsGrowthType: s.params.SlotsGrowthType,
			SlotsIncrement:  s.params.SlotsIncrement,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateServices updates each service and persists the changes inside a single transaction.
func (s *ServiceDiscoveryInstance) UpdateServices(services []ServiceInstance) error {
	mutex.Lock()
	defer mutex.Unlock()

	err := s.startTransaction()
	if err != nil {
		return err
	}
	reload := false
	s.markForCleanUp()
	for _, service := range services {
		if s.serviceNotTracked(service.GetName()) {
			continue
		}
		if !service.Changed() {
			if se, ok := s.services[service.GetName()]; ok {
				se.cleanup = false
			}
			continue
		}
		r, err := s.initService(service)
		if err != nil {
			s.deleteTransaction()
			return err
		}
		reload = reload || r
		se := s.services[service.GetName()]
		r, err = se.confService.Update(service.GetServers())
		if err != nil {
			s.deleteTransaction()
			return err
		}
		reload = reload || r
	}
	r := s.cleanup()
	reload = reload || r
	if reload {
		if err := s.commitTransaction(); err != nil {
			return err
		}
		s.reloadAgent.Reload()
		return nil
	}
	s.deleteTransaction()
	return nil
}

func (s *ServiceDiscoveryInstance) startTransaction() error {
	version, err := s.client.GetVersion("")
	if err != nil {
		return err
	}
	transaction, err := s.client.StartTransaction(version)
	if err != nil {
		return err
	}
	s.transactionID = transaction.ID
	return nil
}

func (s *ServiceDiscoveryInstance) markForCleanUp() {
	for id := range s.services {
		s.services[id].cleanup = true
	}
}

func (s *ServiceDiscoveryInstance) serviceNotTracked(service string) bool {
	if len(s.params.Allowlist) > 0 {
		for _, se := range s.params.Allowlist {
			if se == service {
				return false
			}
		}
		return true
	}
	for _, se := range s.params.Denylist {
		if se == service {
			return true
		}
	}
	return false
}

func (s *ServiceDiscoveryInstance) initService(service ServiceInstance) (bool, error) {
	if se, ok := s.services[service.GetName()]; ok {
		se.confService.SetTransactionID(s.transactionID)
		se.cleanup = false
		return false, nil
	}
	se, err := s.client.NewService(service.GetBackendName(), configuration.ScalingParams{
		BaseSlots:       s.params.ServerSlotsBase,
		SlotsGrowthType: s.params.SlotsGrowthType,
		SlotsIncrement:  s.params.SlotsIncrement,
	})
	if err != nil {
		return false, err
	}
	reload, err := se.Init(s.transactionID, service.GetFrom())
	if err != nil {
		return false, err
	}
	s.services[service.GetName()] = &confService{
		confService: se,
		cleanup:     false,
	}
	return reload, nil
}

func (s *ServiceDiscoveryInstance) cleanup() (reload bool) {
	for service := range s.services {
		if s.services[service].cleanup {
			s.services[service].confService.SetTransactionID(s.transactionID)
			changed, err := s.services[service].confService.Update([]configuration.ServiceServer{})
			if err != nil {
				s.logErrorf("service %s marked for clean-up cannot be updated, %s", service, err.Error())
				continue
			}

			if changed {
				s.logWarningf("service %s marked for clean-up, has not any more backend servers", service)
			}

			reload = reload || changed
		}
	}

	return reload
}

func (s *ServiceDiscoveryInstance) deleteTransaction() {
	if err := s.client.DeleteTransaction(s.transactionID); err != nil {
		s.logWarningf("cannot delete transaction due to an error: %s", err.Error())
	}
	s.transactionID = ""
}

func (s *ServiceDiscoveryInstance) commitTransaction() error {
	_, err := s.client.CommitTransaction(s.transactionID)
	s.transactionID = ""
	return err
}

func (s *ServiceDiscoveryInstance) logWarningf(format string, args ...interface{}) {
	log.WithFieldsf(s.params.LogFields, log.WarnLevel, format, args...)
}

func (s *ServiceDiscoveryInstance) logErrorf(format string, args ...interface{}) {
	log.WithFieldsf(s.params.LogFields, log.ErrorLevel, format, args...)
}
