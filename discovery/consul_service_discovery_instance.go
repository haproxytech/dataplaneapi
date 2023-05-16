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
	"net"
	"strconv"
	"time"

	"github.com/haproxytech/client-native/v4/configuration"
	"github.com/haproxytech/client-native/v4/models"
	"github.com/haproxytech/dataplaneapi/log"

	"github.com/hashicorp/consul/api"
)

type consulService struct {
	params  *models.Consul
	name    string
	servers []configuration.ServiceServer
	changed bool
}

func (c *consulService) GetName() string {
	return c.name
}

func (c *consulService) GetBackendName() string {
	return fmt.Sprintf("consul-backend-%s-%d-%s", *c.params.Address, *c.params.Port, c.name)
}

func (c *consulService) Changed() bool {
	return c.changed
}

func (c *consulService) GetServers() []configuration.ServiceServer {
	return c.servers
}

type consulInstance struct {
	ctx             context.Context
	params          *models.Consul
	api             *api.Client
	discoveryConfig *ServiceDiscoveryInstance
	prevIndexes     map[string]uint64
	update          chan struct{}
	logFields       map[string]interface{}
	timeout         time.Duration
	prevEnabled     bool
}

func (c *consulInstance) start() error {
	c.logDebug("discovery job starting")
	if err := c.setAPIClient(); err != nil {
		return err
	}
	c.update = make(chan struct{}, 1)
	go c.watch()
	return nil
}

func (c *consulInstance) setAPIClient() error {
	address := net.JoinHostPort(*c.params.Address, strconv.FormatInt(*c.params.Port, 10))
	consulConfig := api.DefaultConfig()
	consulConfig.Address = address
	consulConfig.Token = c.params.Token
	cc, err := api.NewClient(consulConfig)
	if err != nil {
		return err
	}
	c.api = cc
	return nil
}

func (c *consulInstance) watch() {
	watchTimer := time.NewTimer(c.timeout)
	defer watchTimer.Stop()

	for {
		select {
		case _, ok := <-c.update:
			if !ok {
				return
			}
			c.logDebug("discovery job update triggered")
			if err := c.setAPIClient(); err != nil {
				c.logErrorf("error while setting up the API client: %s", err.Error())
				c.stop()
				continue
			}
			err := c.discoveryConfig.UpdateParams(discoveryInstanceParams{
				Allowlist:       c.params.ServiceAllowlist,
				Denylist:        c.params.ServiceDenylist,
				LogFields:       c.logFields,
				ServerSlotsBase: int(*c.params.ServerSlotsBase),
				SlotsGrowthType: *c.params.ServerSlotsGrowthType,
				SlotsIncrement:  int(c.params.ServerSlotsGrowthIncrement),
			})
			if err != nil {
				c.stop()
				c.logErrorf("error while updating the instance: %s", err.Error())
			}
		case <-c.ctx.Done():
			c.stop()
		case <-watchTimer.C:
			c.logDebug("discovery job reconciliation started")
			if err := c.updateServices(); err != nil {
				// c.log.Errorf("error while updating service: %w", err)
				c.stop()
			}
			c.logDebug("discovery job reconciliation completed")
			watchTimer.Reset(c.timeout)
		}
	}
}

func (c *consulInstance) stop() {
	c.logDebug("discovery job stopping")
	c.api = nil
	c.prevEnabled = false
	close(c.update)
}

func (c *consulInstance) updateServices() error {
	services := make([]ServiceInstance, 0)
	params := &api.QueryOptions{}
	if c.params.Namespace != "" {
		params.Namespace = c.params.Namespace
	}
	cServices, _, err := c.api.Catalog().Services(params)
	if err != nil {
		return err
	}
	newIndexes := make(map[string]uint64)
	for se := range cServices {
		if se == "consul" {
			continue
		}
		nodes, meta, err := c.api.Health().Service(se, "", false, &api.QueryOptions{})
		if err != nil {
			continue
		}
		newIndexes[se] = meta.LastIndex
		services = append(services, &consulService{
			name:    se,
			params:  c.params,
			servers: c.convertToServers(nodes),
			changed: c.hasServiceChanged(se, meta.LastIndex),
		})
	}
	c.prevIndexes = newIndexes
	return c.discoveryConfig.UpdateServices(services)
}

func (c *consulInstance) convertToServers(nodes []*api.ServiceEntry) []configuration.ServiceServer {
	servers := make([]configuration.ServiceServer, 0)
	for _, node := range nodes {
		if node.Service.Address != "" {
			servers = append(servers, configuration.ServiceServer{
				Address: node.Service.Address,
				Port:    node.Service.Port,
			})
		} else {
			servers = append(servers, configuration.ServiceServer{
				Address: node.Node.Address,
				Port:    node.Service.Port,
			})
		}
	}
	return servers
}

func (c *consulInstance) hasServiceChanged(service string, index uint64) bool {
	prevIndex, ok := c.prevIndexes[service]
	if !ok {
		return true
	}
	return prevIndex != index
}

func (c *consulInstance) updateTimeout(timeoutSeconds int) error {
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", timeoutSeconds))
	if err != nil {
		return err
	}
	c.timeout = timeout
	return nil
}

func (c *consulInstance) handleStateChange() error {
	if c.stateChangedToEnabled() {
		if err := c.start(); err != nil {
			c.prevEnabled = false
			return err
		}
		c.prevEnabled = *c.params.Enabled
		return nil
	}
	if c.stateChangedToDisabled() {
		c.stop()
	}
	if *c.params.Enabled {
		c.update <- struct{}{}
	}
	c.prevEnabled = *c.params.Enabled
	return nil
}

func (c *consulInstance) stateChangedToEnabled() bool {
	return !c.prevEnabled && *c.params.Enabled
}

func (c *consulInstance) stateChangedToDisabled() bool {
	return c.prevEnabled && !*c.params.Enabled
}

func (c *consulInstance) logDebug(message string) {
	log.WithFields(c.logFields, log.DebugLevel, message)
}

func (c *consulInstance) logErrorf(format string, args ...interface{}) {
	log.WithFieldsf(c.logFields, log.ErrorLevel, format, args...)
}
