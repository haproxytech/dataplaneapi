package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/stick_request_rule"
)

//CreateStickRequestRuleHandlerImpl implementation of the CreateStickRequestRuleHandler interface using client-native client
type CreateStickRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteStickRequestRuleHandlerImpl implementation of the DeleteStickRequestRuleHandler interface using client-native client
type DeleteStickRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetStickRequestRuleHandlerImpl implementation of the GetStickRequestRuleHandler interface using client-native client
type GetStickRequestRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetStickRequestRulesHandlerImpl implementation of the GetStickRequestRulesHandler interface using client-native client
type GetStickRequestRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceStickRequestRuleHandlerImpl implementation of the ReplaceStickRequestRuleHandler interface using client-native client
type ReplaceStickRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateStickRequestRuleHandlerImpl) Handle(params stick_request_rule.CreateStickRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateStickRequestRule(params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_request_rule.NewCreateStickRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_request_rule.NewCreateStickRequestRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteStickRequestRuleHandlerImpl) Handle(params stick_request_rule.DeleteStickRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteStickRequestRule(params.ID, params.Backend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_request_rule.NewDeleteStickRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_request_rule.NewDeleteStickRequestRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetStickRequestRuleHandlerImpl) Handle(params stick_request_rule.GetStickRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetStickRequestRule(params.ID, params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_request_rule.NewGetStickRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_request_rule.NewGetStickRequestRuleOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetStickRequestRulesHandlerImpl) Handle(params stick_request_rule.GetStickRequestRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetStickRequestRules(params.Backend, t)
	if err != nil {
		e := misc.HandleError(err)
		return stick_request_rule.NewGetStickRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return stick_request_rule.NewGetStickRequestRulesOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceStickRequestRuleHandlerImpl) Handle(params stick_request_rule.ReplaceStickRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditStickRequestRule(params.ID, params.Backend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return stick_request_rule.NewReplaceStickRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return stick_request_rule.NewReplaceStickRequestRuleOK().WithPayload(params.Data)
}
