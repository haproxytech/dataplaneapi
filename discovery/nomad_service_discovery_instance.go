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
	"fmt"
	"time"

	"github.com/haproxytech/client-native/v4/configuration"
	"github.com/haproxytech/client-native/v4/models"
	"github.com/haproxytech/dataplaneapi/log"

	"github.com/hashicorp/nomad/api"
)

type nomadService struct {
	name    string
	changed bool
	params  *models.Nomad
	servers []configuration.ServiceServer
}

func (n *nomadService) GetName() string {
	return n.name
}

func (n *nomadService) GetBackendName() string {
	return fmt.Sprintf("nomad-backend-%s-%s", n.params.Name, n.GetName())
}

func (n *nomadService) Changed() bool {
	return n.changed
}

func (n *nomadService) GetServers() []configuration.ServiceServer {
	return n.servers
}

type nomadInstance struct {
	params          *models.Nomad
	api             *api.Client
	discoveryConfig *ServiceDiscoveryInstance
	prevIndexes     map[string]uint64
	timeout         time.Duration
	prevEnabled     bool
	ctx             context.Context
	update          chan struct{}
	logFields       map[string]interface{}
}

func (n *nomadInstance) start() error {
	n.logDebug("discovery job starting")
	if err := n.setAPIClient(); err != nil {
		return err
	}
	n.update = make(chan struct{}, 1)
	go n.watch()
	return nil
}

func (n *nomadInstance) setAPIClient() error {
	address := fmt.Sprintf("%s:%d", *n.params.Address, *n.params.Port)
	nomadCfg := api.DefaultConfig()
	nomadCfg.Address = address
	nomadCfg.SecretID = n.params.SecretID
	nc, err := api.NewClient(nomadCfg)
	if err != nil {
		return err
	}
	n.api = nc
	return nil
}

func (n *nomadInstance) watch() {
	watchTimer := time.NewTimer(n.timeout)
	defer watchTimer.Stop()

	for {
		select {
		case _, ok := <-n.update:
			if !ok {
				return
			}
			n.logDebug("discovery job update triggered")
			if err := n.setAPIClient(); err != nil {
				n.logErrorf("error while setting up the API client: %s", err.Error())
				n.stop()
				continue
			}
			err := n.discoveryConfig.UpdateParams(discoveryInstanceParams{
				Allowlist:       n.params.Allowlist,
				Denylist:        n.params.Denylist,
				LogFields:       n.logFields,
				ServerSlotsBase: int(*n.params.ServerSlotsBase),
				SlotsGrowthType: *n.params.ServerSlotsGrowthType,
				SlotsIncrement:  int(n.params.ServerSlotsGrowthIncrement),
			})
			if err != nil {
				n.stop()
				n.logErrorf("error while updating the instance: %s", err.Error())
			}
		case <-n.ctx.Done():
			n.stop()
		case <-watchTimer.C:
			n.logDebug("discovery job reconciliation started")
			if err := n.updateServices(); err != nil {
				n.logErrorf("error while updating service: %s", err.Error())
			}
			n.logDebug("discovery job reconciliation completed")
			watchTimer.Reset(n.timeout)
		}
	}
}

func (n *nomadInstance) stop() {
	n.logDebug("discovery job stopping")
	n.api = nil
	n.prevEnabled = false
	close(n.update)
}

func (n *nomadInstance) updateServices() error {
	services := make([]ServiceInstance, 0)

	// Set query options.
	query := &api.QueryOptions{}
	if n.params.Namespace != "" {
		query.Namespace = n.params.Namespace
	}

	// Fetch list of services from Nomad API.
	nServices, _, err := n.api.Services().List(query)
	if err != nil {
		return err
	}

	newIndexes := make(map[string]uint64)

	for _, nSvc := range nServices {
		for _, s := range nSvc.Services {
			svcRegistrations, meta, err := n.api.Services().Get(s.ServiceName, (&api.QueryOptions{Namespace: nSvc.Namespace}))
			if err != nil {
				continue
			}
			if len(svcRegistrations) == 0 {
				continue
			}
			newIndexes[s.ServiceName] = meta.LastIndex
			services = append(services, &nomadService{
				name:    fmt.Sprintf("%s-%s", svcRegistrations[0].Namespace, s.ServiceName),
				params:  n.params,
				servers: n.convertToServers(svcRegistrations),
				changed: n.hasServiceChanged(s.ServiceName, meta.LastIndex),
			})

		}
	}
	n.prevIndexes = newIndexes
	return n.discoveryConfig.UpdateServices(services)
}

func (n *nomadInstance) convertToServers(nodes []*api.ServiceRegistration) []configuration.ServiceServer {
	servers := make([]configuration.ServiceServer, 0)
	for _, node := range nodes {
		servers = append(servers, configuration.ServiceServer{
			Address: node.Address,
			Port:    node.Port,
		})
	}
	return servers
}

func (n *nomadInstance) hasServiceChanged(service string, index uint64) bool {
	prevIndex, ok := n.prevIndexes[service]
	if !ok {
		return true
	}
	return prevIndex == index
}

func (n *nomadInstance) updateTimeout(timeoutSeconds int) error {
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", timeoutSeconds))
	if err != nil {
		return err
	}
	n.timeout = timeout
	return nil
}

func (n *nomadInstance) handleStateChange() error {
	if n.stateChangedToEnabled() {
		if err := n.start(); err != nil {
			n.prevEnabled = false
			return err
		}
		n.prevEnabled = *n.params.Enabled
		return nil
	}
	if n.stateChangedToDisabled() {
		n.stop()
	}
	if *n.params.Enabled {
		n.update <- struct{}{}
	}
	n.prevEnabled = *n.params.Enabled
	return nil
}

func (n *nomadInstance) stateChangedToEnabled() bool {
	return !n.prevEnabled && *n.params.Enabled
}

func (n *nomadInstance) stateChangedToDisabled() bool {
	return n.prevEnabled && !*n.params.Enabled
}

func (n *nomadInstance) logDebug(message string) {
	log.WithFields(n.logFields, log.DebugLevel, message)
}

func (n *nomadInstance) logErrorf(format string, args ...interface{}) {
	log.WithFieldsf(n.logFields, log.ErrorLevel, format, args...)
}
