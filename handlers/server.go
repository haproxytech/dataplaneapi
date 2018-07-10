package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/controller/haproxy"
	"github.com/haproxytech/controller/misc"
	"github.com/haproxytech/controller/operations/server"
)

//CreateServerHandlerImpl implementation of the CreateServerHandler interface using client-native client
type CreateServerHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteServerHandlerImpl implementation of the DeleteServerHandler interface using client-native client
type DeleteServerHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetServerHandlerImpl implementation of the GetServerHandler interface using client-native client
type GetServerHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetServersHandlerImpl implementation of the GetServersHandler interface using client-native client
type GetServersHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceServerHandlerImpl implementation of the ReplaceServerHandler interface using client-native client
type ReplaceServerHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateServerHandlerImpl) Handle(params server.CreateServerParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateServer(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return server.NewCreateServerCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteServerHandlerImpl) Handle(params server.DeleteServerParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteServer(params.Name, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return server.NewDeleteServerNoContent()
}

//Handle executing the request and returning a response
func (h *GetServerHandlerImpl) Handle(params server.GetServerParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetServer(params.Name, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServerDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetServerOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetServersHandlerImpl) Handle(params server.GetServersParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetServers(params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServersDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetServersOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceServerHandlerImpl) Handle(params server.ReplaceServerParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditServer(params.Name, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return server.NewReplaceServerOK().WithPayload(params.Data)
}
