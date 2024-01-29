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
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
)

type awsServiceDiscovery struct {
	services    Store
	client      configuration.Configuration
	reloadAgent haproxy.IReloadAgent
	context     context.Context
}

// NewAWSDiscoveryService creates a new ServiceDiscovery that connects to AWS
func NewAWSDiscoveryService(params ServiceDiscoveriesParams) ServiceDiscovery {
	return &awsServiceDiscovery{
		services:    NewInstanceStore(),
		client:      params.Client,
		reloadAgent: params.ReloadAgent,
		context:     params.Context,
	}
}

func (a awsServiceDiscovery) AddNode(id string, params ServiceDiscoveryParams) (err error) {
	aParams, ok := params.(*models.AwsRegion)
	if !ok {
		return errors.New("expected *models.AwsRegion")
	}

	var instance *awsInstance
	instance, err = newAWSRegionInstance(a.context, aParams, a.client, a.reloadAgent)
	if err != nil {
		return
	}

	if err = a.services.Create(id, instance); err != nil {
		return
	}

	if *aParams.Enabled {
		instance.start()
	}
	return
}

func (a awsServiceDiscovery) GetNode(id string) (params ServiceDiscoveryParams, err error) {
	var i interface{}
	if i, err = a.services.Read(id); err != nil {
		return
	}
	return i.(*awsInstance).params, nil
}

func (a awsServiceDiscovery) GetNodes() (ServiceDiscoveryParams, error) {
	var awsRegions models.AwsRegions
	for _, as := range a.services.List() {
		awsRegions = append(awsRegions, as.(*awsInstance).params)
	}
	return awsRegions, nil
}

func (a awsServiceDiscovery) RemoveNode(id string) error {
	return a.services.Delete(id)
}

func (a awsServiceDiscovery) UpdateNode(id string, params ServiceDiscoveryParams) error {
	newParams, ok := params.(*models.AwsRegion)
	if !ok {
		return errors.New("expected *models.AwsRegion")
	}
	return a.services.Update(id, func(item interface{}) (err error) {
		ai := item.(*awsInstance)

		if err = ai.updateTimeout(*newParams.RetryTimeout); err != nil {
			if *ai.params.Enabled {
				ai.stop()
			}
			return errors.New("invalid retry_timeout")
		}

		switch {
		case *newParams.Enabled == *ai.params.Enabled:
			break
		case *newParams.Enabled && !*ai.params.Enabled:
			ai.start()
		default:
			ai.stop()
		}
		ai.params = newParams

		if *ai.params.Enabled {
			ai.update <- struct{}{}
		}

		return
	})
}
