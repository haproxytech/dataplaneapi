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
	"github.com/haproxytech/client-native/v4/models"

	sc "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/service_discovery"
)

// CreateNomadHandlerImpl implementation of the CreateNomadHandler interface using client-native client
type CreateNomadHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	UseValidation   bool
	PersistCallback func([]*models.Nomad) error
}

// DeleteNomadHandlerImpl implementation of the DeleteNomadHandler interface using client-native client
type DeleteNomadHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.Nomad) error
}

// GetNomadHandlerImpl implementation of the GetNomadHandler interface using client-native client
type GetNomadHandlerImpl struct {
	Discovery sc.ServiceDiscoveries
}

// GetNomadsHandlerImpl implementation of the GetNomadsHandler interface using client-native client
type GetNomadsHandlerImpl struct {
	Discovery sc.ServiceDiscoveries
}

// ReplaceNomadHandlerImpl implementation of the ReplaceNomadHandler interface using client-native client
type ReplaceNomadHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	UseValidation   bool
	PersistCallback func([]*models.Nomad) error
}

// Handle executing the request and returning a response
func (c *CreateNomadHandlerImpl) Handle(params service_discovery.CreateNomadParams, principal interface{}) middleware.Responder {
	params.Data.ID = sc.NewServiceDiscoveryUUID()
	if err := sc.ValidateNomadData(params.Data, c.UseValidation); err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateNomadDefault(int(*e.Code)).WithPayload(e)
	}
	err := c.Discovery.AddNode("nomad", *params.Data.ID, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateNomadDefault(int(*e.Code)).WithPayload(e)
	}
	nomads, err := getNomads(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateNomadDefault(int(*e.Code)).WithPayload(e)
	}
	err = c.PersistCallback(nomads)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewCreateNomadDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewCreateNomadCreated().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (c *DeleteNomadHandlerImpl) Handle(params service_discovery.DeleteNomadParams, principal interface{}) middleware.Responder {
	err := c.Discovery.RemoveNode("nomad", params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceNomadDefault(int(*e.Code)).WithPayload(e)
	}
	nomads, err := getNomads(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewDeleteNomadDefault(int(*e.Code)).WithPayload(e)
	}
	err = c.PersistCallback(nomads)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewDeleteNomadDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewDeleteNomadNoContent()
}

// Handle executing the request and returning a response
func (c *GetNomadHandlerImpl) Handle(params service_discovery.GetNomadParams, principal interface{}) middleware.Responder {
	nodes, err := c.Discovery.GetNode("nomad", params.ID)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewGetNomadsDefault(int(*e.Code)).WithPayload(e)
	}
	nomad, ok := nodes.(*models.Nomad)
	if !ok {
		e := misc.HandleError(errors.New("expected *models.Nomad"))
		return service_discovery.NewGetNomadsDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewGetNomadOK().WithPayload(&service_discovery.GetNomadOKBody{Data: nomad})
}

// Handle executing the request and returning a response
func (c *GetNomadsHandlerImpl) Handle(params service_discovery.GetNomadsParams, principal interface{}) middleware.Responder {
	nomads, err := getNomads(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewGetNomadDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewGetNomadsOK().WithPayload(&service_discovery.GetNomadsOKBody{Data: nomads})
}

// Handle executing the request and returning a response
func (c *ReplaceNomadHandlerImpl) Handle(params service_discovery.ReplaceNomadParams, principal interface{}) middleware.Responder {
	if err := sc.ValidateNomadData(params.Data, c.UseValidation); err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceNomadDefault(int(*e.Code)).WithPayload(e)
	}
	err := c.Discovery.UpdateNode("nomad", *params.Data.ID, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceNomadDefault(int(*e.Code)).WithPayload(e)
	}
	nomads, err := getNomads(c.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceNomadDefault(int(*e.Code)).WithPayload(e)
	}
	err = c.PersistCallback(nomads)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewDeleteNomadDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewReplaceNomadOK().WithPayload(params.Data)
}

func getNomads(discovery sc.ServiceDiscoveries) (models.Nomads, error) {
	nodes, err := discovery.GetNodes("nomad")
	if err != nil {
		return nil, err
	}
	nomads, ok := nodes.(models.Nomads)
	if !ok {
		return nil, errors.New("expected models.Nomads")
	}
	return nomads, nil
}
