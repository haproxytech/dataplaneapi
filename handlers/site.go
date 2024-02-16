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
	"github.com/haproxytech/dataplaneapi/operations/sites"
)

// CreateSiteHandlerImpl implementation of the CreateSiteHandler interface using client-native client
type CreateSiteHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteSiteHandlerImpl implementation of the DeleteSiteHandler interface using client-native client
type DeleteSiteHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetSiteHandlerImpl implementation of the GetSiteHandler interface using client-native client
type GetSiteHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetSitesHandlerImpl implementation of the GetSitesHandler interface using client-native client
type GetSitesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceSiteHandlerImpl implementation of the ReplaceSiteHandler interface using client-native client
type ReplaceSiteHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateSiteHandlerImpl) Handle(params sites.CreateSiteParams, principal interface{}) middleware.Responder {
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
		return sites.NewCreateSiteDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewCreateSiteDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateSite(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewCreateSiteDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return sites.NewCreateSiteDefault(int(*e.Code)).WithPayload(e)
			}
			return sites.NewCreateSiteCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return sites.NewCreateSiteAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return sites.NewCreateSiteAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteSiteHandlerImpl) Handle(params sites.DeleteSiteParams, principal interface{}) middleware.Responder {
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
		return sites.NewDeleteSiteDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewCreateSiteDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteSite(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewDeleteSiteDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return sites.NewDeleteSiteDefault(int(*e.Code)).WithPayload(e)
			}
			return sites.NewDeleteSiteNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return sites.NewDeleteSiteAccepted().WithReloadID(rID)
	}
	return sites.NewDeleteSiteAccepted()
}

// Handle executing the request and returning a response
func (h *GetSiteHandlerImpl) Handle(params sites.GetSiteParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewGetSiteDefault(int(*e.Code)).WithPayload(e)
	}

	_, site, err := configuration.GetSite(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewGetSiteDefault(int(*e.Code)).WithPayload(e)
	}
	return sites.NewGetSiteOK().WithPayload(site)
}

// Handle executing the request and returning a response
func (h *GetSitesHandlerImpl) Handle(params sites.GetSitesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewGetSitesDefault(int(*e.Code)).WithPayload(e)
	}

	_, s, err := configuration.GetSites(t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return sites.NewGetSitesOK().WithPayload(models.Sites{})
		}
		return sites.NewGetSitesDefault(int(*e.Code)).WithPayload(e)
	}
	return sites.NewGetSitesOK().WithPayload(s)
}

// Handle executing the request and returning a response
func (h *ReplaceSiteHandlerImpl) Handle(params sites.ReplaceSiteParams, principal interface{}) middleware.Responder {
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
		return sites.NewReplaceSiteDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewReplaceSiteDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditSite(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewReplaceSiteDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return sites.NewReplaceSiteDefault(int(*e.Code)).WithPayload(e)
			}
			return sites.NewReplaceSiteOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return sites.NewReplaceSiteAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return sites.NewReplaceSiteAccepted().WithPayload(params.Data)
}
