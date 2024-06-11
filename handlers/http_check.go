// Copyright 2021 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/http_check"
)

// CreateHTTPCheckHandlerImpl implementation of the CreateHTTPCheckHandler interface using client-native client
type CreateHTTPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteHTTPCheckHandlerImpl implementation of the DeleteHTTPCheckHandler interface using client-native client
type DeleteHTTPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetHTTPCheckHandlerImpl implementation of the GetHTTPCheckHandler interface using client-native client
type GetHTTPCheckHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetHTTPChecksHandlerImpl implementation of the GetHTTPChecksHandler interface using client-native client
type GetHTTPChecksHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceHTTPCheckHandlerImpl implementation of the ReplaceHTTPCheckHandler interface using client-native client
type ReplaceHTTPCheckHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// ReplaceHTTPChecksHandlerImpl implementation of the ReplaceHTTPChecksHandler interface using client-native client
type ReplaceHTTPChecksHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateHTTPCheckHandlerImpl) Handle(params http_check.CreateHTTPCheckParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	pName := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}
	if params.ParentName != nil {
		pName = *params.ParentName
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_check.NewCreateHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewCreateHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateHTTPCheck(params.Index, params.ParentType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewCreateHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewCreateHTTPCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewCreateHTTPCheckCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewCreateHTTPCheckAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_check.NewCreateHTTPCheckAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteHTTPCheckHandlerImpl) Handle(params http_check.DeleteHTTPCheckParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	pName := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}
	if params.ParentName != nil {
		pName = *params.ParentName
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_check.NewDeleteHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewDeleteHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteHTTPCheck(params.Index, params.ParentType, pName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewDeleteHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewDeleteHTTPCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewDeleteHTTPCheckNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewDeleteHTTPCheckAccepted().WithReloadID(rID)
	}
	return http_check.NewDeleteHTTPCheckAccepted()
}

// Handle executing the request and returning a response
func (h *GetHTTPCheckHandlerImpl) Handle(params http_check.GetHTTPCheckParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewGetHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPCheck(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewGetHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}
	return http_check.NewGetHTTPCheckOK().WithPayload(rule)
}

// Handle executing the request and returning a response
func (h *GetHTTPChecksHandlerImpl) Handle(params http_check.GetHTTPChecksParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewGetHTTPChecksDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPChecks(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return http_check.NewGetHTTPChecksOK().WithPayload(models.HTTPChecks{})
		}
		return http_check.NewGetHTTPChecksDefault(int(*e.Code)).WithPayload(e)
	}
	return http_check.NewGetHTTPChecksOK().WithPayload(rules)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPCheckHandlerImpl) Handle(params http_check.ReplaceHTTPCheckParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	pName := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}
	if params.ParentName != nil {
		pName = *params.ParentName
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_check.NewReplaceHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditHTTPCheck(params.Index, params.ParentType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceHTTPCheckDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewReplaceHTTPCheckDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewReplaceHTTPCheckOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewReplaceHTTPCheckAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_check.NewReplaceHTTPCheckAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *ReplaceHTTPChecksHandlerImpl) Handle(params http_check.ReplaceHTTPChecksParams, principal interface{}) middleware.Responder {
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
		return http_check.NewReplaceHTTPChecksDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceHTTPChecksDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceHTTPChecks(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_check.NewReplaceHTTPChecksDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_check.NewReplaceHTTPChecksDefault(int(*e.Code)).WithPayload(e)
			}
			return http_check.NewReplaceHTTPChecksOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_check.NewReplaceHTTPChecksAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_check.NewReplaceHTTPChecksAccepted().WithPayload(params.Data)
}
