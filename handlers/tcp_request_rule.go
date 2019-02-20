package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/tcp_request_rule"
)

//CreateTCPRequestRuleHandlerImpl implementation of the CreateTCPRequestRuleHandler interface using client-native client
type CreateTCPRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteTCPRequestRuleHandlerImpl implementation of the DeleteTCPRequestRuleHandler interface using client-native client
type DeleteTCPRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetTCPRequestRuleHandlerImpl implementation of the GetTCPRequestRuleHandler interface using client-native client
type GetTCPRequestRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetTCPRequestRulesHandlerImpl implementation of the GetTCPRequestRulesHandler interface using client-native client
type GetTCPRequestRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceTCPRequestRuleHandlerImpl implementation of the ReplaceTCPRequestRuleHandler interface using client-native client
type ReplaceTCPRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.CreateTCPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateTCPRequestRule(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewCreateTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_request_rule.NewCreateTCPRequestRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.DeleteTCPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteTCPRequestRule(params.ID, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewDeleteTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_request_rule.NewDeleteTCPRequestRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.GetTCPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rule, err := h.Client.Configuration.GetTCPRequestRule(params.ID, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_request_rule.NewGetTCPRequestRuleOK().WithPayload(rule)
}

//Handle executing the request and returning a response
func (h *GetTCPRequestRulesHandlerImpl) Handle(params tcp_request_rule.GetTCPRequestRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rules, err := h.Client.Configuration.GetTCPRequestRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewGetTCPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_request_rule.NewGetTCPRequestRulesOK().WithPayload(rules)
}

//Handle executing the request and returning a response
func (h *ReplaceTCPRequestRuleHandlerImpl) Handle(params tcp_request_rule.ReplaceTCPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditTCPRequestRule(params.ID, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_request_rule.NewReplaceTCPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_request_rule.NewReplaceTCPRequestRuleOK().WithPayload(params.Data)
}
