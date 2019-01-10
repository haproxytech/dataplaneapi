package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/sites"
)

//CreateSiteHandlerImpl implementation of the CreateSiteHandler interface using client-native client
type CreateSiteHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteSiteHandlerImpl implementation of the DeleteSiteHandler interface using client-native client
type DeleteSiteHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetSiteHandlerImpl implementation of the GetSiteHandler interface using client-native client
type GetSiteHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetSitesHandlerImpl implementation of the GetSitesHandler interface using client-native client
type GetSitesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceSiteHandlerImpl implementation of the ReplaceSiteHandler interface using client-native client
type ReplaceSiteHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateSiteHandlerImpl) Handle(params sites.CreateSiteParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateSite(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewCreateSiteDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return sites.NewCreateSiteCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteSiteHandlerImpl) Handle(params sites.DeleteSiteParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteSite(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewDeleteSiteDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return sites.NewDeleteSiteNoContent()
}

//Handle executing the request and returning a response
func (h *GetSiteHandlerImpl) Handle(params sites.GetSiteParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetSite(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewGetSiteDefault(int(*e.Code)).WithPayload(e)
	}
	return sites.NewGetSiteOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetSitesHandlerImpl) Handle(params sites.GetSitesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetSites(t)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewGetSitesDefault(int(*e.Code)).WithPayload(e)
	}
	return sites.NewGetSitesOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceSiteHandlerImpl) Handle(params sites.ReplaceSiteParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditSite(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return sites.NewReplaceSiteDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return sites.NewReplaceSiteOK().WithPayload(params.Data)
}
