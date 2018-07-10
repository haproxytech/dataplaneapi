package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/controller/haproxy"
	"github.com/haproxytech/controller/misc"
	"github.com/haproxytech/controller/operations/tcp_connection_rule"
)

//CreateTCPConnectionRuleHandlerImpl implementation of the CreateTCPConnectionRuleHandler interface using client-native client
type CreateTCPConnectionRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteTCPConnectionRuleHandlerImpl implementation of the DeleteTCPConnectionRuleHandler interface using client-native client
type DeleteTCPConnectionRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetTCPConnectionRuleHandlerImpl implementation of the GetTCPConnectionRuleHandler interface using client-native client
type GetTCPConnectionRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetTCPConnectionRulesHandlerImpl implementation of the GetTCPConnectionRulesHandler interface using client-native client
type GetTCPConnectionRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceTCPConnectionRuleHandlerImpl implementation of the ReplaceTCPConnectionRuleHandler interface using client-native client
type ReplaceTCPConnectionRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateTCPConnectionRuleHandlerImpl) Handle(params tcp_connection_rule.CreateTCPConnectionRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateTCPConnectionRule(params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_connection_rule.NewCreateTCPConnectionRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_connection_rule.NewCreateTCPConnectionRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteTCPConnectionRuleHandlerImpl) Handle(params tcp_connection_rule.DeleteTCPConnectionRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteTCPConnectionRule(params.ID, params.Frontend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_connection_rule.NewDeleteTCPConnectionRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_connection_rule.NewDeleteTCPConnectionRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetTCPConnectionRuleHandlerImpl) Handle(params tcp_connection_rule.GetTCPConnectionRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetTCPConnectionRule(params.ID, params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_connection_rule.NewGetTCPConnectionRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_connection_rule.NewGetTCPConnectionRuleOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetTCPConnectionRulesHandlerImpl) Handle(params tcp_connection_rule.GetTCPConnectionRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetTCPConnectionRules(params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_connection_rule.NewGetTCPConnectionRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_connection_rule.NewGetTCPConnectionRulesOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceTCPConnectionRuleHandlerImpl) Handle(params tcp_connection_rule.ReplaceTCPConnectionRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditTCPConnectionRule(params.ID, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_connection_rule.NewReplaceTCPConnectionRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_connection_rule.NewReplaceTCPConnectionRuleOK().WithPayload(params.Data)
}
