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

	"github.com/haproxytech/client-native/v4/configuration"
	"github.com/haproxytech/client-native/v4/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
)

type nomadServiceDiscovery struct {
	nomadServices Store
	client        configuration.Configuration
	reloadAgent   haproxy.IReloadAgent
	context       context.Context
}

// NewNomadDiscoveryService creates a new ServiceDiscovery that connects to Nomad.
func NewNomadDiscoveryService(params ServiceDiscoveriesParams) ServiceDiscovery {
	return &nomadServiceDiscovery{
		nomadServices: NewInstanceStore(),
		client:        params.Client,
		reloadAgent:   params.ReloadAgent,
		context:       params.Context,
	}
}

func (n *nomadServiceDiscovery) AddNode(id string, params ServiceDiscoveryParams) (err error) {
	nParams, ok := params.(*models.Nomad)
	if !ok {
		return errors.New("expected models.Nomad")
	}

	var timeout time.Duration
	timeout, err = time.ParseDuration(fmt.Sprintf("%ds", *nParams.RetryTimeout))
	if err != nil {
		return err
	}

	logFields := map[string]interface{}{"ServiceDiscovery": "Nomad", "ID": *nParams.ID}

	instance := &nomadInstance{
		params:  nParams,
		ctx:     n.context,
		timeout: timeout,
		discoveryConfig: NewServiceDiscoveryInstance(n.client, n.reloadAgent, discoveryInstanceParams{
			Allowlist:       nParams.Allowlist,
			Denylist:        nParams.Denylist,
			LogFields:       logFields,
			ServerSlotsBase: int(*nParams.ServerSlotsBase),
			SlotsGrowthType: *nParams.ServerSlotsGrowthType,
			SlotsIncrement:  int(nParams.ServerSlotsGrowthIncrement),
		}),
		prevIndexes: make(map[string]uint64),
		logFields:   logFields,
	}

	if err = n.nomadServices.Create(id, instance); err != nil {
		return
	}

	instance.prevEnabled = *nParams.Enabled

	if *nParams.Enabled {
		return instance.start()
	}
	return nil
}

func (n *nomadServiceDiscovery) GetNode(id string) (p ServiceDiscoveryParams, err error) {
	var i interface{}
	if i, err = n.nomadServices.Read(id); err != nil {
		return
	}
	p = i.(*nomadInstance).params
	return
}

func (n *nomadServiceDiscovery) GetNodes() (ServiceDiscoveryParams, error) {
	var nomads models.Nomads
	for _, ni := range n.nomadServices.List() {
		nomads = append(nomads, ni.(*nomadInstance).params)
	}
	return nomads, nil
}

func (n *nomadServiceDiscovery) RemoveNode(id string) error {
	return n.nomadServices.Delete(id)
}

func (n *nomadServiceDiscovery) UpdateNode(id string, params ServiceDiscoveryParams) (err error) {
	nParams, ok := params.(*models.Nomad)
	if !ok {
		return errors.New("expected models.Nomads")
	}
	return n.nomadServices.Update(id, func(item interface{}) error {
		ni := item.(*nomadInstance)
		ni.params = nParams
		if err = ni.updateTimeout(int(*nParams.RetryTimeout)); err != nil {
			ni.stop()
			return errors.New("invalid retry_timeout")
		}
		return ni.handleStateChange()
	})
}
