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
	data, err := h.Client.Configuration.GetGlobalConfiguration()
	if err != nil {
		e := misc.HandleError(err)
		return global.NewGetGlobalDefault(int(*e.Code)).WithPayload(e)
	}
	return global.NewGetGlobalOK().WithPayload(data)
}

//Handle executing the request and returning a response
func (h *ReplaceGlobalHandlerImpl) Handle(params global.ReplaceGlobalParams, principal interface{}) middleware.Responder {
	v := int64(0)
	if params.Version == nil {
		msg := "Version not specified"
		code := misc.ErrHTTPBadRequest
		return global.NewReplaceGlobalBadRequest().WithPayload(&models.Error{Code: &code, Message: &msg})
	}

	err := h.Client.Configuration.PushGlobalConfiguration(params.Data, v)

	if err != nil {
		e := misc.HandleError(err)
		return global.NewReplaceGlobalDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return global.NewReplaceGlobalOK().WithPayload(params.Data)
}
