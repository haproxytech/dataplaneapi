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
	"github.com/haproxytech/dataplaneapi/operations/server_template"
)

// CreateServerTemplateHandlerImpl implementation of the CreateServerTemplateHandler interface using client-native client
type CreateServerTemplateHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteServerTemplateHandlerImpl implementation of the DeleteServerTemplateHandler interface using client-native client
type DeleteServerTemplateHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetServerTemplateHandlerImpl implementation of the GetServerTemplateHandler interface using client-native client
type GetServerTemplateHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetServerTemplatesHandlerImpl implementation of the GetServerTemplatesHandler interface using client-native client
type GetServerTemplatesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceServerTemplateHandlerImpl implementation of the ReplaceServerTemplateHandler interface using client-native client
type ReplaceServerTemplateHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateServerTemplateHandlerImpl) Handle(params server_template.CreateServerTemplateParams, principal interface{}) middleware.Responder {
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
		return server_template.NewCreateServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewCreateServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateServerTemplate(params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewCreateServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server_template.NewCreateServerTemplateDefault(int(*e.Code)).WithPayload(e)
			}
			return server_template.NewCreateServerTemplateCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return server_template.NewCreateServerTemplateAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return server_template.NewCreateServerTemplateAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteServerTemplateHandlerImpl) Handle(params server_template.DeleteServerTemplateParams, principal interface{}) middleware.Responder {
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
		return server_template.NewDeleteServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewDeleteServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteServerTemplate(params.Prefix, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewDeleteServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server_template.NewDeleteServerTemplateDefault(int(*e.Code)).WithPayload(e)
			}
			return server_template.NewDeleteServerTemplateNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return server_template.NewDeleteServerTemplateAccepted().WithReloadID(rID)
	}
	return server_template.NewDeleteServerTemplateAccepted()
}

// Handle executing the request and returning a response
func (h *GetServerTemplateHandlerImpl) Handle(params server_template.GetServerTemplateParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewGetServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	_, template, err := configuration.GetServerTemplate(params.Prefix, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewGetServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}
	return server_template.NewGetServerTemplateOK().WithPayload(template)
}

// Handle executing the request and returning a response
func (h *GetServerTemplatesHandlerImpl) Handle(params server_template.GetServerTemplatesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewGetServerTemplatesDefault(int(*e.Code)).WithPayload(e)
	}

	_, templates, err := configuration.GetServerTemplates(params.ParentName, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return server_template.NewGetServerTemplatesOK().WithPayload(models.ServerTemplates{})
		}
		return server_template.NewGetServerTemplatesDefault(int(*e.Code)).WithPayload(e)
	}
	return server_template.NewGetServerTemplatesOK().WithPayload(templates)
}

// Handle executing the request and returning a response
func (h *ReplaceServerTemplateHandlerImpl) Handle(params server_template.ReplaceServerTemplateParams, principal interface{}) middleware.Responder {
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
		return server_template.NewReplaceServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewReplaceServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	_, ondisk, err := configuration.GetServerTemplate(params.Prefix, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewReplaceServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditServerTemplate(params.Prefix, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_template.NewReplaceServerTemplateDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		reload := changeThroughRuntimeAPI(*params.Data, *ondisk, params.ParentName, h.Client)
		if reload {
			if *params.ForceReload {
				err := h.ReloadAgent.ForceReload()
				if err != nil {
					e := misc.HandleError(err)
					return server_template.NewReplaceServerTemplateDefault(int(*e.Code)).WithPayload(e)
				}
				return server_template.NewReplaceServerTemplateOK().WithPayload(params.Data)
			}
			rID := h.ReloadAgent.Reload()
			return server_template.NewReplaceServerTemplateAccepted().WithReloadID(rID).WithPayload(params.Data)
		}
		return server_template.NewReplaceServerTemplateOK().WithPayload(params.Data)
	}
	return server_template.NewReplaceServerTemplateAccepted().WithPayload(params.Data)
}
