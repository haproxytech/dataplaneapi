package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/operations/global"
	"github.com/haproxytech/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

//GetGlobalHandlerImpl implementation of the GetGlobalHandler interface
type GetGlobalHandlerImpl struct {
	Client *client_native.HAProxyClient
}

// ReplaceGlobalHandlerImpl implementation of the ReplaceGlobalHandler interface
type ReplaceGlobalHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *GetGlobalHandlerImpl) Handle(params global.GetGlobalParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	data, err := h.Client.Configuration.GetGlobalConfiguration(t)
	if err != nil {
		e := misc.HandleError(err)
		return global.NewGetGlobalDefault(int(*e.Code)).WithPayload(e)
	}
	return global.NewGetGlobalOK().WithPayload(data)
}

//Handle executing the request and returning a response
func (h *ReplaceGlobalHandlerImpl) Handle(params global.ReplaceGlobalParams, principal interface{}) middleware.Responder {
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
		return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
	}

	err := h.Client.Configuration.PushGlobalConfiguration(params.Data, t, v)

	if err != nil {
		e := misc.HandleError(err)
		return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
			}
			return global.NewReplaceGlobalOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return global.NewReplaceGlobalAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return global.NewReplaceGlobalAccepted().WithPayload(params.Data)
}
