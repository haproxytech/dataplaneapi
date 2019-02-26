package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/server"
	"github.com/haproxytech/models"
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

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.CreateServer(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server.NewCreateServerDefault(int(*e.Code)).WithPayload(e)
			}
			return server.NewCreateServerCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return server.NewCreateServerAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return server.NewCreateServerAccepted().WithPayload(params.Data)
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

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.DeleteServer(params.Name, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return server.NewDeleteServerDefault(int(*e.Code)).WithPayload(e)
			}
			return server.NewDeleteServerNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return server.NewDeleteServerAccepted().WithReloadID(rID)
	}
	return server.NewDeleteServerAccepted()
}

//Handle executing the request and returning a response
func (h *GetServerHandlerImpl) Handle(params server.GetServerParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	srv, err := h.Client.Configuration.GetServer(params.Name, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServerDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetServerOK().WithPayload(srv)
}

//Handle executing the request and returning a response
func (h *GetServersHandlerImpl) Handle(params server.GetServersParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	srvs, err := h.Client.Configuration.GetServers(params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewGetServersDefault(int(*e.Code)).WithPayload(e)
	}
	return server.NewGetServersOK().WithPayload(srvs)
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

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}

	ondisk, err := h.Client.Configuration.GetServer(params.Name, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}

	err = h.Client.Configuration.EditServer(params.Name, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		reload := changeThroughRuntimeAPI(*params.Data, *ondisk.Data, params.Backend, "", h.Client)
		if reload {
			if *params.ForceReload {
				err := h.ReloadAgent.ForceReload()
				if err != nil {
					e := misc.HandleError(err)
					return server.NewReplaceServerDefault(int(*e.Code)).WithPayload(e)
				}
				return server.NewReplaceServerOK().WithPayload(params.Data)
			}
			rID := h.ReloadAgent.Reload()
			return server.NewReplaceServerAccepted().WithReloadID(rID).WithPayload(params.Data)
		}
		return server.NewReplaceServerOK().WithPayload(params.Data)
	}
	return server.NewReplaceServerAccepted().WithPayload(params.Data)
}
