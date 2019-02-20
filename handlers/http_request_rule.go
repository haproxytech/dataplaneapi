package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/http_request_rule"
)

//CreateHTTPRequestRuleHandlerImpl implementation of the CreateHTTPRequestRuleHandler interface using client-native client
type CreateHTTPRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteHTTPRequestRuleHandlerImpl implementation of the DeleteHTTPRequestRuleHandler interface using client-native client
type DeleteHTTPRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetHTTPRequestRuleHandlerImpl implementation of the GetHTTPRequestRuleHandler interface using client-native client
type GetHTTPRequestRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetHTTPRequestRulesHandlerImpl implementation of the GetHTTPRequestRulesHandler interface using client-native client
type GetHTTPRequestRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceHTTPRequestRuleHandlerImpl implementation of the ReplaceHTTPRequestRuleHandler interface using client-native client
type ReplaceHTTPRequestRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.CreateHTTPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateHTTPRequestRule(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewCreateHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return http_request_rule.NewCreateHTTPRequestRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.DeleteHTTPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteHTTPRequestRule(params.ID, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewDeleteHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return http_request_rule.NewDeleteHTTPRequestRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.GetHTTPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rule, err := h.Client.Configuration.GetHTTPRequestRule(params.ID, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return http_request_rule.NewGetHTTPRequestRuleOK().WithPayload(rule)
}

//Handle executing the request and returning a response
func (h *GetHTTPRequestRulesHandlerImpl) Handle(params http_request_rule.GetHTTPRequestRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	rules, err := h.Client.Configuration.GetHTTPRequestRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewGetHTTPRequestRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return http_request_rule.NewGetHTTPRequestRulesOK().WithPayload(rules)
}

//Handle executing the request and returning a response
func (h *ReplaceHTTPRequestRuleHandlerImpl) Handle(params http_request_rule.ReplaceHTTPRequestRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditHTTPRequestRule(params.ID, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_request_rule.NewReplaceHTTPRequestRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return http_request_rule.NewReplaceHTTPRequestRuleOK().WithPayload(params.Data)
}
