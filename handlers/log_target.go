package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/log_target"
)

//CreateLogTargetHandlerImpl implementation of the CreateLogTargetHandler interface using client-native client
type CreateLogTargetHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteLogTargetHandlerImpl implementation of the DeleteLogTargetHandler interface using client-native client
type DeleteLogTargetHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetLogTargetHandlerImpl implementation of the GetLogTargetHandler interface using client-native client
type GetLogTargetHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetLogTargetsHandlerImpl implementation of the GetLogTargetsHandler interface using client-native client
type GetLogTargetsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceLogTargetHandlerImpl implementation of the ReplaceLogTargetHandler interface using client-native client
type ReplaceLogTargetHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateLogTargetHandlerImpl) Handle(params log_target.CreateLogTargetParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateLogTarget(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewCreateLogTargetDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return log_target.NewCreateLogTargetCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteLogTargetHandlerImpl) Handle(params log_target.DeleteLogTargetParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteLogTarget(params.ID, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewDeleteLogTargetDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return log_target.NewDeleteLogTargetNoContent()
}

//Handle executing the request and returning a response
func (h *GetLogTargetHandlerImpl) Handle(params log_target.GetLogTargetParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	logTarget, err := h.Client.Configuration.GetLogTarget(params.ID, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewGetLogTargetDefault(int(*e.Code)).WithPayload(e)
	}
	return log_target.NewGetLogTargetOK().WithPayload(logTarget)
}

//Handle executing the request and returning a response
func (h *GetLogTargetsHandlerImpl) Handle(params log_target.GetLogTargetsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	logTargets, err := h.Client.Configuration.GetLogTargets(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewGetLogTargetsDefault(int(*e.Code)).WithPayload(e)
	}
	return log_target.NewGetLogTargetsOK().WithPayload(logTargets)
}

//Handle executing the request and returning a response
func (h *ReplaceLogTargetHandlerImpl) Handle(params log_target.ReplaceLogTargetParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditLogTarget(params.ID, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_target.NewReplaceLogTargetDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return log_target.NewReplaceLogTargetOK().WithPayload(params.Data)
}
