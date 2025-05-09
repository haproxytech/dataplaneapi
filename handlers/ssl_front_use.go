// Copyright 2025 HAProxy Technologies
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
	cn "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/configuration/parents"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	sfu "github.com/haproxytech/dataplaneapi/operations/s_s_l_front_use"
)

type GetAllSSLFrontUsesHandlerImpl struct {
	Client cn.HAProxyClient
}

func (h GetAllSSLFrontUsesHandlerImpl) Handle(params sfu.GetAllSSLFrontUsesParams, i interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewGetAllSSLFrontUsesDefault(int(*e.Code)).WithPayload(e)
	}

	_, uses, err := configuration.GetSSLFrontUses(string(parents.FrontendParentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewGetAllSSLFrontUsesDefault(int(*e.Code)).WithPayload(e)
	}

	return sfu.NewGetAllSSLFrontUsesOK().WithPayload(uses)
}

type CreateSSLFrontUseHandlerImpl struct {
	Client      cn.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h CreateSSLFrontUseHandlerImpl) Handle(params sfu.CreateSSLFrontUseParams, i interface{}) middleware.Responder {
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
		return sfu.NewCreateSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewCreateSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateSSLFrontUse(string(parents.FrontendParentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewCreateSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return sfu.NewCreateSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
			}
			return sfu.NewCreateSSLFrontUseCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return sfu.NewCreateSSLFrontUseAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return sfu.NewCreateSSLFrontUseAccepted().WithPayload(params.Data)
}

type GetSSLFrontUseHandlerImpl struct {
	Client cn.HAProxyClient
}

func (h GetSSLFrontUseHandlerImpl) Handle(params sfu.GetSSLFrontUseParams, i interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewGetSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	_, use, err := configuration.GetSSLFrontUse(params.Index, string(parents.FrontendParentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewGetSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	return sfu.NewGetSSLFrontUseOK().WithPayload(use)
}

type ReplaceSSLFrontUseHandlerImpl struct {
	Client      cn.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h ReplaceSSLFrontUseHandlerImpl) Handle(params sfu.ReplaceSSLFrontUseParams, i interface{}) middleware.Responder {
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
		return sfu.NewReplaceSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewReplaceSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditSSLFrontUse(params.Index, string(parents.FrontendParentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewReplaceSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return sfu.NewReplaceSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
			}
			return sfu.NewReplaceSSLFrontUseOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return sfu.NewReplaceSSLFrontUseAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return sfu.NewReplaceSSLFrontUseOK().WithPayload(params.Data)
}

type DeleteSSLFrontUseHandlerImpl struct {
	Client      cn.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h DeleteSSLFrontUseHandlerImpl) Handle(params sfu.DeleteSSLFrontUseParams, i interface{}) middleware.Responder {
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
		return sfu.NewDeleteSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewDeleteSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteSSLFrontUse(params.Index, string(parents.FrontendParentType), params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sfu.NewDeleteSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return sfu.NewDeleteSSLFrontUseDefault(int(*e.Code)).WithPayload(e)
			}
			return sfu.NewDeleteSSLFrontUseNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return sfu.NewDeleteSSLFrontUseAccepted().WithReloadID(rID)
	}

	return sfu.NewDeleteSSLFrontUseAccepted()
}
