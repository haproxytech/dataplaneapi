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
	"fmt"
	"time"

	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
)

type consulServiceDiscovery struct {
	consulServices Store
	client         configuration.Configuration
	reloadAgent    haproxy.IReloadAgent
	context        context.Context
}

// NewConsulDiscoveryService creates a new ServiceDiscovery that connects to consul
func NewConsulDiscoveryService(params ServiceDiscoveriesParams) ServiceDiscovery {
	return &consulServiceDiscovery{
		consulServices: NewInstanceStore(),
		client:         params.Client,
		reloadAgent:    params.ReloadAgent,
		context:        params.Context,
	}
}

func (c *consulServiceDiscovery) AddNode(id string, params ServiceDiscoveryParams) (err error) {
	cParams, ok := params.(*models.Consul)
	if !ok {
		return errors.New("expected models.Consuls")
	}

	var timeout time.Duration
	timeout, err = time.ParseDuration(fmt.Sprintf("%ds", *cParams.RetryTimeout))
	if err != nil {
		return err
	}

	logFields := map[string]interface{}{"ServiceDiscovery": "Consul", "ID": *cParams.ID}

	instance := &consulInstance{
		params:  cParams,
		ctx:     c.context,
		timeout: timeout,
		discoveryConfig: NewServiceDiscoveryInstance(c.client, c.reloadAgent, discoveryInstanceParams{
			Allowlist:       cParams.ServiceAllowlist,
			Denylist:        cParams.ServiceDenylist,
			LogFields:       logFields,
			ServerSlotsBase: int(*cParams.ServerSlotsBase),
			SlotsGrowthType: *cParams.ServerSlotsGrowthType,
			SlotsIncrement:  int(cParams.ServerSlotsGrowthIncrement),
		}),
		prevIndexes: make(map[string]uint64),
		logFields:   logFields,
	}

	if err = c.consulServices.Create(id, instance); err != nil {
		return
	}

	instance.prevEnabled = *cParams.Enabled

	if *cParams.Enabled {
		return instance.start()
	}
	return nil
}

func (c *consulServiceDiscovery) GetNode(id string) (p ServiceDiscoveryParams, err error) {
	var i interface{}
	if i, err = c.consulServices.Read(id); err != nil {
		return
	}
	p = i.(*consulInstance).params
	return
}

func (c *consulServiceDiscovery) GetNodes() (ServiceDiscoveryParams, error) {
	var consuls models.Consuls
	for _, ci := range c.consulServices.List() {
		consuls = append(consuls, ci.(*consulInstance).params)
	}
	return consuls, nil
}

func (c *consulServiceDiscovery) RemoveNode(id string) error {
	return c.consulServices.Delete(id)
}

func (c *consulServiceDiscovery) UpdateNode(id string, params ServiceDiscoveryParams) (err error) {
	cParams, ok := params.(*models.Consul)
	if !ok {
		return errors.New("expected models.Consuls")
	}
	return c.consulServices.Update(id, func(item interface{}) error {
		ci := item.(*consulInstance)
		ci.params = cParams
		if err = ci.updateTimeout(int(*cParams.RetryTimeout)); err != nil {
			ci.stop()
			return errors.New("invalid retry_timeout")
		}
		return ci.handleStateChange()
	})
}
