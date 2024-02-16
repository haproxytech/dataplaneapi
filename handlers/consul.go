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

package handlers

import (
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native/v6/models"

	sc "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/service_discovery"
)

// CreateConsulHandlerImpl implementation of the CreateConsulHandler interface using client-native client
type CreateConsulHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.Consul) error
	UseValidation   bool
}

// DeleteConsulHandlerImpl implementation of the DeleteConsulHandler interface using client-native client
type DeleteConsulHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.Consul) error
}

// GetConsulHandlerImpl implementation of the GetConsulHandler interface using client-native client
type GetConsulHandlerImpl struct {
	Discovery sc.ServiceDiscoveries
}

// GetConsulsHandlerImpl implementation of the GetConsulsHandler interface using client-native client
type GetConsulsHandlerImpl struct {
	Discovery sc.ServiceDiscoveries
}

// ReplaceConsulHandlerImpl implementation of the ReplaceConsulHandler interface using client-native client
type ReplaceConsulHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.Consul) error
	UseValidation   bool
}

// Handle executing the request and returning a response
func (c *CreateConsulHandlerImpl) Handle(params service_discovery.CreateConsulParams, principal interface{}) middleware.Responder {
	params.Data.ID = sc.NewServiceDiscoveryUUID()
	if err := sc.ValidateConsulData(params.Data, c.UseValidation); err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateConsulDefault(int(*e.Code)).WithPayload(e)
	}
	setFilters(params.Data)
	if params.Data.HealthCheckPolicy != nil && *params.Data.HealthCheckPolicy == models.ConsulHealthCheckPolicyMin && params.Data.HealthCheckPolicyMin <= 0 {
		e := &models.Error{
			Message: misc.StringP("health_check_policy_min is required for 'min' health_check_policy"),
			Code:    misc.Int64P(int(misc.ErrHTTPBadRequest)),
		}
		return service_discovery.NewCreateConsulDefault(int(*e.Code)).WithPayload(e)
	}
	err := c.Discovery.AddNode("consul", *params.Data.ID, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateConsulDefault(int(*e.Code)).WithPayload(e)
	}
	consuls, err := getConsuls(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateConsulDefault(int(*e.Code)).WithPayload(e)
	}
	err = c.PersistCallback(consuls)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateConsulDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewCreateConsulCreated().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (c *DeleteConsulHandlerImpl) Handle(params service_discovery.DeleteConsulParams, principal interface{}) middleware.Responder {
	err := c.Discovery.RemoveNode("consul", params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceConsulDefault(int(*e.Code)).WithPayload(e)
	}
	consuls, err := getConsuls(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewDeleteConsulDefault(int(*e.Code)).WithPayload(e)
	}
	err = c.PersistCallback(consuls)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewDeleteConsulDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewDeleteConsulNoContent()
}

// Handle executing the request and returning a response
func (c *GetConsulHandlerImpl) Handle(params service_discovery.GetConsulParams, principal interface{}) middleware.Responder {
	nodes, err := c.Discovery.GetNode("consul", params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewGetConsulsDefault(int(*e.Code)).WithPayload(e)
	}
	consul, ok := nodes.(*models.Consul)
	if !ok {
		e := misc.HandleError(errors.New("expected *models.Consul"))
		return service_discovery.NewGetConsulsDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewGetConsulOK().WithPayload(consul)
}

// Handle executing the request and returning a response
func (c *GetConsulsHandlerImpl) Handle(params service_discovery.GetConsulsParams, principal interface{}) middleware.Responder {
	consuls, err := getConsuls(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewGetConsulDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewGetConsulsOK().WithPayload(consuls)
}

// Handle executing the request and returning a response
func (c *ReplaceConsulHandlerImpl) Handle(params service_discovery.ReplaceConsulParams, principal interface{}) middleware.Responder {
	if err := sc.ValidateConsulData(params.Data, c.UseValidation); err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceConsulDefault(int(*e.Code)).WithPayload(e)
	}
	setFilters(params.Data)
	if *params.Data.HealthCheckPolicy == models.ConsulHealthCheckPolicyMin && params.Data.HealthCheckPolicyMin <= 0 {
		e := &models.Error{
			Message: misc.StringP("health_check_policy_min is required for 'min' health_check_policy"),
			Code:    misc.Int64P(int(misc.ErrHTTPBadRequest)),
		}
		return service_discovery.NewCreateConsulDefault(int(*e.Code)).WithPayload(e)
	}
	err := c.Discovery.UpdateNode("consul", *params.Data.ID, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceConsulDefault(int(*e.Code)).WithPayload(e)
	}
	consuls, err := getConsuls(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceConsulDefault(int(*e.Code)).WithPayload(e)
	}
	err = c.PersistCallback(consuls)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewDeleteConsulDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewReplaceConsulOK().WithPayload(params.Data)
}

func setFilters(data *models.Consul) {
	if len(data.ServiceAllowlist) == 0 && len(data.ServiceWhitelist) > 0 {
		data.ServiceAllowlist = data.ServiceWhitelist
	}

	if len(data.ServiceDenylist) == 0 && len(data.ServiceBlacklist) > 0 {
		data.ServiceDenylist = data.ServiceBlacklist
	}
}

func getConsuls(discovery sc.ServiceDiscoveries) (models.Consuls, error) {
	nodes, err := discovery.GetNodes("consul")
	if err != nil {
		return nil, err
	}
	consuls, ok := nodes.(models.Consuls)
	if !ok {
		return nil, errors.New("expected models.Consuls")
	}
	return consuls, nil
}
