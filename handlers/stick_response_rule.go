package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/stick_response_rule"
)

//CreateStickResponseRuleHandlerImpl implementation of the CreateStickResponseRuleHandler interface using client-native client
type CreateStickResponseRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteStickResponseRuleHandlerImpl implementation of the DeleteStickResponseRuleHandler interface using client-native client
type DeleteStickResponseRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetStickResponseRuleHandlerImpl implementation of the GetStickResponseRuleHandler interface using client-native client
type GetStickResponseRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetStickResponseRulesHandlerImpl implementation of the GetStickResponseRulesHandler interface using client-native client
type GetStickResponseRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceStickResponseRuleHandlerImpl implementation of the ReplaceStickResponseRuleHandler interface using client-native client
type ReplaceStickResponseRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateStickResponseRuleHandlerImpl) Handle(params stick_response_rule.CreateStickResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateStickResponseRule(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_response_rule.NewCreateStickResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_response_rule.NewCreateStickResponseRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteStickResponseRuleHandlerImpl) Handle(params stick_response_rule.DeleteStickResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteStickResponseRule(params.ID, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_response_rule.NewDeleteStickResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_response_rule.NewDeleteStickResponseRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetStickResponseRuleHandlerImpl) Handle(params stick_response_rule.GetStickResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetStickResponseRule(params.ID, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_response_rule.NewGetStickResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_response_rule.NewGetStickResponseRuleOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetStickResponseRulesHandlerImpl) Handle(params stick_response_rule.GetStickResponseRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetStickResponseRules(params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_response_rule.NewGetStickResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_response_rule.NewGetStickResponseRulesOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceStickResponseRuleHandlerImpl) Handle(params stick_response_rule.ReplaceStickResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditStickResponseRule(params.ID, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_response_rule.NewReplaceStickResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_response_rule.NewReplaceStickResponseRuleOK().WithPayload(params.Data)
}
