package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/server_switching_rule"
)

//CreateServerSwitchingRuleHandlerImpl implementation of the CreateServerSwitchingRuleHandler interface using client-native client
type CreateServerSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteServerSwitchingRuleHandlerImpl implementation of the DeleteServerSwitchingRuleHandler interface using client-native client
type DeleteServerSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetServerSwitchingRuleHandlerImpl implementation of the GetServerSwitchingRuleHandler interface using client-native client
type GetServerSwitchingRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetServerSwitchingRulesHandlerImpl implementation of the GetServerSwitchingRulesHandler interface using client-native client
type GetServerSwitchingRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceServerSwitchingRuleHandlerImpl implementation of the ReplaceServerSwitchingRuleHandler interface using client-native client
type ReplaceServerSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.CreateServerSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateServerSwitchingRule(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewCreateServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return server_switching_rule.NewCreateServerSwitchingRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.DeleteServerSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteServerSwitchingRule(params.ID, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewDeleteServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return server_switching_rule.NewDeleteServerSwitchingRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.GetServerSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rule, err := h.Client.Configuration.GetServerSwitchingRule(params.ID, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewGetServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return server_switching_rule.NewGetServerSwitchingRuleOK().WithPayload(rule)
}

//Handle executing the request and returning a response
func (h *GetServerSwitchingRulesHandlerImpl) Handle(params server_switching_rule.GetServerSwitchingRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rules, err := h.Client.Configuration.GetServerSwitchingRules(params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewGetServerSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return server_switching_rule.NewGetServerSwitchingRulesOK().WithPayload(rules)
}

//Handle executing the request and returning a response
func (h *ReplaceServerSwitchingRuleHandlerImpl) Handle(params server_switching_rule.ReplaceServerSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditServerSwitchingRule(params.ID, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return server_switching_rule.NewReplaceServerSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return server_switching_rule.NewReplaceServerSwitchingRuleOK().WithPayload(params.Data)
}
