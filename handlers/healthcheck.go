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
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/health_check"
)

// CreateHealthCheckHandlerImpl implementation of the CreateHealthCheckHandler interface using client-native client
type CreateHealthCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteHealthCheckHandlerImpl implementation of the DeleteHealthCheckHandler interface using client-native client
type DeleteHealthCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetHealthCheckHandlerImpl implementation of the GetHealthCheckHandler interface using client-native client
type GetHealthCheckHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetHealthChecksHandlerImpl implementation of the GetHealthChecksHandler interface using client-native client
type GetHealthChecksHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHealthCheckHandlerImpl implementation of the ReplaceHealthCheckHandler interface using client-native client
type ReplaceHealthCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHealthCheckHandlerImpl) Handle(params health_check.CreateHealthCheckParams, principal any) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return health_check.NewCreateHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.createHealthCheck(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return health_check.NewCreateHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return health_check.NewCreateHealthCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return health_check.NewCreateHealthCheckCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return health_check.NewCreateHealthCheckAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return health_check.NewCreateHealthCheckAccepted().WithPayload(params.Data)
}

func (h *CreateHealthCheckHandlerImpl) createHealthCheck(params health_check.CreateHealthCheckParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.CreateStructuredHealthcheck(params.Data, t, v)
	}
	return configuration.CreateHealthcheck(params.Data, t, v)
}

// Handle executing the request and returning a response
func (h *DeleteHealthCheckHandlerImpl) Handle(params health_check.DeleteHealthCheckParams, principal any) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return health_check.NewDeleteHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return health_check.NewDeleteHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHealthcheck(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return health_check.NewDeleteHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return health_check.NewDeleteHealthCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return health_check.NewDeleteHealthCheckNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return health_check.NewDeleteHealthCheckAccepted().WithReloadID(rID)
	}
	return health_check.NewDeleteHealthCheckAccepted()
}

// Handle executing the request and returning a response
func (h *GetHealthCheckHandlerImpl) Handle(params health_check.GetHealthCheckParams, principal any) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, bck, err := h.getHealthCheck(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return health_check.NewGetHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}
	return health_check.NewGetHealthCheckOK().WithPayload(bck)
}

func (h *GetHealthCheckHandlerImpl) getHealthCheck(params health_check.GetHealthCheckParams, t string) (int64, *models.HealthCheck, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredHealthcheck(params.Name, t)
	}
	return configuration.GetHealthcheck(params.Name, t)
}

// Handle executing the request and returning a response
func (h *GetHealthChecksHandlerImpl) Handle(params health_check.GetHealthChecksParams, principal any) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, bcks, err := h.getHealthChecks(params, t)
	if err != nil {
		e := misc.HandleError(err)
		return health_check.NewGetHealthChecksDefault(int(*e.Code)).WithPayload(e)
	}
	return health_check.NewGetHealthChecksOK().WithPayload(bcks)
}

func (h *GetHealthChecksHandlerImpl) getHealthChecks(params health_check.GetHealthChecksParams, t string) (int64, models.Healthchecks, error) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredHealthchecks(t)
	}
	return configuration.GetHealthchecks(t)
}

// Handle executing the request and returning a response
func (h *ReplaceHealthCheckHandlerImpl) Handle(params health_check.ReplaceHealthCheckParams, principal any) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return health_check.NewReplaceHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.editHealthCheck(params, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return health_check.NewReplaceHealthCheckDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return health_check.NewReplaceHealthCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return health_check.NewReplaceHealthCheckOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return health_check.NewReplaceHealthCheckAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return health_check.NewReplaceHealthCheckAccepted().WithPayload(params.Data)
}

func (h *ReplaceHealthCheckHandlerImpl) editHealthCheck(params health_check.ReplaceHealthCheckParams, t string, v int64) error {
	configuration, err := h.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.EditStructuredHealthcheck(params.Name, params.Data, t, v)
	}
	return configuration.EditHealthcheck(params.Name, params.Data, t, v)
}
