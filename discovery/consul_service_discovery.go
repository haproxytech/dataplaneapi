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
	"time"

	"github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/models/v2"
)

type consulServiceDiscovery struct {
	consulServices map[string]*consulInstance
	client         *configuration.Client
	mu             sync.RWMutex
}

//NewConsulDiscoveryService creates a new ServiceDiscovery that connects to consul
func NewConsulDiscoveryService(client *configuration.Client) ServiceDiscovery {
	return &consulServiceDiscovery{
		consulServices: make(map[string]*consulInstance),
		client:         client,
	}
}

func (c *consulServiceDiscovery) AddNode(id string, params ServiceDiscoveryParams) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.consulServices[id]
	if ok {
		return fmt.Errorf("instance already exists for: %s", id)
	}
	cParams, ok := params.(*models.Consul)
	if !ok {
		return errors.New("expected models.Consuls")
	}
	timeout, err := time.ParseDuration(fmt.Sprintf("%ds", *cParams.RetryTimeout))
	if err != nil {
		return err
	}
	instance := &consulInstance{
		params:  cParams,
		timeout: timeout,
		discoveryConfig: NewServiceDiscoveryInstance(c.client, discoveryInstanceParams{
			Whitelist:       cParams.ServiceWhitelist,
			Blacklist:       cParams.ServiceBlacklist,
			ServerSlotsBase: int(*cParams.ServerSlotsBase),
			SlotsGrowthType: *cParams.ServerSlotsGrowthType,
			SlotsIncrement:  int(cParams.ServerSlotsGrowthIncrement),
		}),
		prevIndexes: make(map[string]uint64),
	}
	c.consulServices[id] = instance
	instance.prevEnabled = *cParams.Enabled

	if *cParams.Enabled {
		return instance.start()
	}
	return nil
}

func (c *consulServiceDiscovery) GetNode(id string) (ServiceDiscoveryParams, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ci, ok := c.consulServices[id]
	if !ok {
		return nil, errors.New("instance not found")
	}
	return ci.params, nil
}

func (c *consulServiceDiscovery) GetNodes() (ServiceDiscoveryParams, error) {
	c.mu.RLock()
	var consuls models.Consuls
	for _, ci := range c.consulServices {
		consuls = append(consuls, ci.params)
	}
	return consuls, nil
}

func (c *consulServiceDiscovery) RemoveNode(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.consulServices[id]
	if !ok {
		return errors.New("instance not found")
	}
	delete(c.consulServices, id)
	return nil
}

func (c *consulServiceDiscovery) UpdateNode(id string, params ServiceDiscoveryParams) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ci, ok := c.consulServices[id]
	if !ok {
		return errors.New("instance not found")
	}
	cParams, ok := params.(*models.Consul)
	if !ok {
		return errors.New("expected models.Consuls")
	}
	ci.params = cParams
	err := ci.updateTimeout(int(*cParams.RetryTimeout))
	if err != nil {
		ci.stop()
		return errors.New("invalid retry_timeout")
	}
	return ci.handelStateChange()
}
