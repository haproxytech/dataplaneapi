package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/backend_switching_rule"
)

//CreateBackendSwitchingRuleHandlerImpl implementation of the CreateBackendSwitchingRuleHandler interface using client-native client
type CreateBackendSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteBackendSwitchingRuleHandlerImpl implementation of the DeleteBackendSwitchingRuleHandler interface using client-native client
type DeleteBackendSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetBackendSwitchingRuleHandlerImpl implementation of the GetBackendSwitchingRuleHandler interface using client-native client
type GetBackendSwitchingRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetBackendSwitchingRulesHandlerImpl implementation of the GetBackendSwitchingRulesHandler interface using client-native client
type GetBackendSwitchingRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceBackendSwitchingRuleHandlerImpl implementation of the ReplaceBackendSwitchingRuleHandler interface using client-native client
type ReplaceBackendSwitchingRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.CreateBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateBackendSwitchingRule(params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewCreateBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return backend_switching_rule.NewCreateBackendSwitchingRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.DeleteBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteBackendSwitchingRule(params.ID, params.Frontend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewDeleteBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return backend_switching_rule.NewDeleteBackendSwitchingRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.GetBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bckRule, err := h.Client.Configuration.GetBackendSwitchingRule(params.ID, params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewGetBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return backend_switching_rule.NewGetBackendSwitchingRuleOK().WithPayload(bckRule)
}

//Handle executing the request and returning a response
func (h *GetBackendSwitchingRulesHandlerImpl) Handle(params backend_switching_rule.GetBackendSwitchingRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bckRules, err := h.Client.Configuration.GetBackendSwitchingRules(params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewGetBackendSwitchingRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return backend_switching_rule.NewGetBackendSwitchingRulesOK().WithPayload(bckRules)
}

//Handle executing the request and returning a response
func (h *ReplaceBackendSwitchingRuleHandlerImpl) Handle(params backend_switching_rule.ReplaceBackendSwitchingRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditBackendSwitchingRule(params.ID, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend_switching_rule.NewReplaceBackendSwitchingRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return backend_switching_rule.NewReplaceBackendSwitchingRuleOK().WithPayload(params.Data)
}
