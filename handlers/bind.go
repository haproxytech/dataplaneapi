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
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/bind"
)

// CreateBindHandlerImpl implementation of the CreateBindHandler interface using client-native client
type CreateBindHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteBindHandlerImpl implementation of the DeleteBindHandler interface using client-native client
type DeleteBindHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetBindHandlerImpl implementation of the GetBindHandler interface using client-native client
type GetBindHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetBindsHandlerImpl implementation of the GetBindsHandler interface using client-native client
type GetBindsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceBindHandlerImpl implementation of the ReplaceBindHandler interface using client-native client
type ReplaceBindHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func bindTypeParams(frontend *string, parentType *string, parentName *string) (pType string, pName string, err error) {
	if frontend != nil && *frontend != "" {
		return "frontend", *frontend, nil
	}
	if parentType == nil || *parentType == "" {
		return "", "", fmt.Errorf("parentType empty")
	}
	if parentName == nil || *parentName == "" {
		return "", "", fmt.Errorf("parentName empty")
	}
	return *parentType, *parentName, nil
}

// Handle executing the request and returning a response
func (h *CreateBindHandlerImpl) Handle(params bind.CreateBindParams, principal interface{}) middleware.Responder {
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
		return bind.NewCreateBindDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewCreateBindDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := bindTypeParams(params.Frontend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewCreateBindDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.CreateBind(pType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewCreateBindDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return bind.NewCreateBindDefault(int(*e.Code)).WithPayload(e)
			}
			return bind.NewCreateBindCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return bind.NewCreateBindAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return bind.NewCreateBindAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteBindHandlerImpl) Handle(params bind.DeleteBindParams, principal interface{}) middleware.Responder {
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
		return bind.NewDeleteBindDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewDeleteBindDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := bindTypeParams(params.Frontend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewDeleteBindDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.DeleteBind(params.Name, pType, pName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewDeleteBindDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return bind.NewDeleteBindDefault(int(*e.Code)).WithPayload(e)
			}
			return bind.NewDeleteBindNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return bind.NewDeleteBindAccepted().WithReloadID(rID)
	}
	return bind.NewDeleteBindAccepted()
}

// Handle executing the request and returning a response
func (h *GetBindHandlerImpl) Handle(params bind.GetBindParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewGetBindDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := bindTypeParams(params.Frontend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewGetBindDefault(int(*e.Code)).WithPayload(e)
	}
	_, b, err := configuration.GetBind(params.Name, pType, pName, t)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewGetBindDefault(int(*e.Code)).WithPayload(e)
	}
	return bind.NewGetBindOK().WithPayload(b)
}

// Handle executing the request and returning a response
func (h *GetBindsHandlerImpl) Handle(params bind.GetBindsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewGetBindsDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := bindTypeParams(params.Frontend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewGetBindsDefault(int(*e.Code)).WithPayload(e)
	}
	_, bs, err := configuration.GetBinds(pType, pName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return bind.NewGetBindsOK().WithPayload(models.Binds{})
		}
		return bind.NewGetBindsDefault(int(*e.Code)).WithPayload(e)
	}
	return bind.NewGetBindsOK().WithPayload(bs)
}

// Handle executing the request and returning a response
func (h *ReplaceBindHandlerImpl) Handle(params bind.ReplaceBindParams, principal interface{}) middleware.Responder {
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
		return bind.NewReplaceBindDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewReplaceBindDefault(int(*e.Code)).WithPayload(e)
	}

	pType, pName, err := bindTypeParams(params.Frontend, params.ParentType, params.ParentName)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewReplaceBindDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.EditBind(params.Name, pType, pName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewReplaceBindDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return bind.NewReplaceBindDefault(int(*e.Code)).WithPayload(e)
			}
			return bind.NewReplaceBindOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return bind.NewReplaceBindAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return bind.NewReplaceBindAccepted().WithPayload(params.Data)
}
