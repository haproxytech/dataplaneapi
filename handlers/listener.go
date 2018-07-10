package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/controller/haproxy"
	"github.com/haproxytech/controller/misc"
	"github.com/haproxytech/controller/operations/listener"
)

//CreateListenerHandlerImpl implementation of the CreateListenerHandler interface using client-native client
type CreateListenerHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteListenerHandlerImpl implementation of the DeleteListenerHandler interface using client-native client
type DeleteListenerHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetListenerHandlerImpl implementation of the GetListenerHandler interface using client-native client
type GetListenerHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetListenersHandlerImpl implementation of the GetListenersHandler interface using client-native client
type GetListenersHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceListenerHandlerImpl implementation of the ReplaceListenerHandler interface using client-native client
type ReplaceListenerHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateListenerHandlerImpl) Handle(params listener.CreateListenerParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateListener(params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return listener.NewCreateListenerDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return listener.NewCreateListenerCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteListenerHandlerImpl) Handle(params listener.DeleteListenerParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteListener(params.Name, params.Frontend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return listener.NewDeleteListenerDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return listener.NewDeleteListenerNoContent()
}

//Handle executing the request and returning a response
func (h *GetListenerHandlerImpl) Handle(params listener.GetListenerParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetListener(params.Name, params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return listener.NewGetListenerDefault(int(*e.Code)).WithPayload(e)
	}
	return listener.NewGetListenerOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetListenersHandlerImpl) Handle(params listener.GetListenersParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetListeners(params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return listener.NewGetListenersDefault(int(*e.Code)).WithPayload(e)
	}
	return listener.NewGetListenersOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceListenerHandlerImpl) Handle(params listener.ReplaceListenerParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditListener(params.Name, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return listener.NewReplaceListenerDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return listener.NewReplaceListenerOK().WithPayload(params.Data)
}
