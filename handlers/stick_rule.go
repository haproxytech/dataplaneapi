package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/stick_rule"
)

//CreateStickRuleHandlerImpl implementation of the CreateStickRuleHandler interface using client-native client
type CreateStickRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteStickRuleHandlerImpl implementation of the DeleteStickRuleHandler interface using client-native client
type DeleteStickRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetStickRuleHandlerImpl implementation of the GetStickRuleHandler interface using client-native client
type GetStickRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetStickRulesHandlerImpl implementation of the GetStickRulesHandler interface using client-native client
type GetStickRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceStickRuleHandlerImpl implementation of the ReplaceStickRuleHandler interface using client-native client
type ReplaceStickRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateStickRuleHandlerImpl) Handle(params stick_rule.CreateStickRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateStickRule(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewCreateStickRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_rule.NewCreateStickRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteStickRuleHandlerImpl) Handle(params stick_rule.DeleteStickRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteStickRule(params.ID, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewDeleteStickRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_rule.NewDeleteStickRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetStickRuleHandlerImpl) Handle(params stick_rule.GetStickRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rule, err := h.Client.Configuration.GetStickRule(params.ID, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewGetStickRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_rule.NewGetStickRuleOK().WithPayload(rule)
}

//Handle executing the request and returning a response
func (h *GetStickRulesHandlerImpl) Handle(params stick_rule.GetStickRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rules, err := h.Client.Configuration.GetStickRules(params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewGetStickRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_rule.NewGetStickRulesOK().WithPayload(rules)
}

//Handle executing the request and returning a response
func (h *ReplaceStickRuleHandlerImpl) Handle(params stick_rule.ReplaceStickRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditStickRule(params.ID, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_rule.NewReplaceStickRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_rule.NewReplaceStickRuleOK().WithPayload(params.Data)
}
