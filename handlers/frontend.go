package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/controller/haproxy"
	"github.com/haproxytech/controller/misc"
	"github.com/haproxytech/controller/operations/frontend"
)

//CreateFrontendHandlerImpl implementation of the CreateFrontendHandler interface using client-native client
type CreateFrontendHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteFrontendHandlerImpl implementation of the DeleteFrontendHandler interface using client-native client
type DeleteFrontendHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetFrontendHandlerImpl implementation of the GetFrontendHandler interface using client-native client
type GetFrontendHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetFrontendsHandlerImpl implementation of the GetFrontendsHandler interface using client-native client
type GetFrontendsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceFrontendHandlerImpl implementation of the ReplaceFrontendHandler interface using client-native client
type ReplaceFrontendHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateFrontendHandlerImpl) Handle(params frontend.CreateFrontendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateFrontend(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewCreateFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return frontend.NewCreateFrontendCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteFrontendHandlerImpl) Handle(params frontend.DeleteFrontendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteFrontend(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewDeleteFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return frontend.NewDeleteFrontendNoContent()
}

//Handle executing the request and returning a response
func (h *GetFrontendHandlerImpl) Handle(params frontend.GetFrontendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetFrontend(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewGetFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	return frontend.NewGetFrontendOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetFrontendsHandlerImpl) Handle(params frontend.GetFrontendsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetFrontends(t)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewGetFrontendsDefault(int(*e.Code)).WithPayload(e)
	}
	return frontend.NewGetFrontendsOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceFrontendHandlerImpl) Handle(params frontend.ReplaceFrontendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	ondisk, err := h.Client.Configuration.GetFrontend(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewReplaceFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	reload := changeThroughRuntimeAPI(*params.Data, *ondisk.Data, "", "", h.Client)

	err = h.Client.Configuration.EditFrontend(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return frontend.NewReplaceFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if reload {
		h.ReloadAgent.Reload()
	}

	return frontend.NewReplaceFrontendOK().WithPayload(params.Data)
}
