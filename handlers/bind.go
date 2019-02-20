package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/bind"
)

//CreateBindHandlerImpl implementation of the CreateBindHandler interface using client-native client
type CreateBindHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteBindHandlerImpl implementation of the DeleteBindHandler interface using client-native client
type DeleteBindHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetBindHandlerImpl implementation of the GetBindHandler interface using client-native client
type GetBindHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetBindsHandlerImpl implementation of the GetBindsHandler interface using client-native client
type GetBindsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceBindHandlerImpl implementation of the ReplaceBindHandler interface using client-native client
type ReplaceBindHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateBindHandlerImpl) Handle(params bind.CreateBindParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateBind(params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewCreateBindDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return bind.NewCreateBindCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteBindHandlerImpl) Handle(params bind.DeleteBindParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteBind(params.Name, params.Frontend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewDeleteBindDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return bind.NewDeleteBindNoContent()
}

//Handle executing the request and returning a response
func (h *GetBindHandlerImpl) Handle(params bind.GetBindParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	b, err := h.Client.Configuration.GetBind(params.Name, params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewGetBindDefault(int(*e.Code)).WithPayload(e)
	}
	return bind.NewGetBindOK().WithPayload(b)
}

//Handle executing the request and returning a response
func (h *GetBindsHandlerImpl) Handle(params bind.GetBindsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bs, err := h.Client.Configuration.GetBinds(params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewGetBindsDefault(int(*e.Code)).WithPayload(e)
	}
	return bind.NewGetBindsOK().WithPayload(bs)
}

//Handle executing the request and returning a response
func (h *ReplaceBindHandlerImpl) Handle(params bind.ReplaceBindParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditBind(params.Name, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return bind.NewReplaceBindDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return bind.NewReplaceBindOK().WithPayload(params.Data)
}
