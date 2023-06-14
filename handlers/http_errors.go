// Copyright 2022 HAProxy Technologies
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
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/http_errors"
)

// CreateHTTPErrorsSectionHandlerImpl implementation of the CreateHTTPErrorsSectionHandler interface using client-native client
type CreateHTTPErrorsSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteHTTPErrorsSectionHandlerImpl implementation of the DeleteHTTPErrorsSectionHandler interface using client-native client
type DeleteHTTPErrorsSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetHTTPErrorsSectionHandlerImpl implementation of the GetHTTPErrorsSectionHandler interface using client-native client
type GetHTTPErrorsSectionHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetHTTPErrorsSectionsHandlerImpl implementation of the GetHTTPErrorsSectionsHandler interface using client-native client
type GetHTTPErrorsSectionsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPErrorsSectionHandlerImpl implementation of the ReplaceHTTPErrorsSectionHandler interface using client-native client
type ReplaceHTTPErrorsSectionHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPErrorsSectionHandlerImpl) Handle(params http_errors.CreateHTTPErrorsSectionParams, principal interface{}) middleware.Responder {
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
		return http_errors.NewCreateHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewCreateHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateHTTPErrorsSection(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return http_errors.NewCreateHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_errors.NewCreateHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return http_errors.NewCreateHTTPErrorsSectionCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_errors.NewCreateHTTPErrorsSectionAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return http_errors.NewCreateHTTPErrorsSectionAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPErrorsSectionHandlerImpl) Handle(params http_errors.DeleteHTTPErrorsSectionParams, principal interface{}) middleware.Responder {
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
		return http_errors.NewDeleteHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewDeleteHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPErrorsSection(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewDeleteHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_errors.NewDeleteHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
			}
			return http_errors.NewDeleteHTTPErrorsSectionNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_errors.NewDeleteHTTPErrorsSectionAccepted().WithReloadID(rID)
	}

	return http_errors.NewDeleteHTTPErrorsSectionAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPErrorsSectionHandlerImpl) Handle(params http_errors.GetHTTPErrorsSectionParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewGetHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	v, ms, err := configuration.GetHTTPErrorsSection(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewGetHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	return http_errors.NewGetHTTPErrorsSectionOK().WithPayload(&http_errors.GetHTTPErrorsSectionOKBody{Version: v, Data: ms})
}

// Handle executing the request and returning a response
func (h *GetHTTPErrorsSectionsHandlerImpl) Handle(params http_errors.GetHTTPErrorsSectionsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewGetHTTPErrorsSectionsDefault(int(*e.Code)).WithPayload(e)
	}

	v, ms, err := configuration.GetHTTPErrorsSections(t)
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewGetHTTPErrorsSectionsDefault(int(*e.Code)).WithPayload(e)
	}

	return http_errors.NewGetHTTPErrorsSectionsOK().WithPayload(&http_errors.GetHTTPErrorsSectionsOKBody{Version: v, Data: ms})
}

func (h *ReplaceHTTPErrorsSectionHandlerImpl) Handle(params http_errors.ReplaceHTTPErrorsSectionParams, principal interface{}) middleware.Responder {
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
		return http_errors.NewReplaceHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewReplaceHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPErrorsSection(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_errors.NewReplaceHTTPErrorsSectionDefault(int(*e.Code)).WithPayload(e)
	}

	return http_errors.NewReplaceHTTPErrorsSectionOK().WithPayload(params.Data)
}
