package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/http_response_rule"
)

//CreateHTTPResponseRuleHandlerImpl implementation of the CreateHTTPResponseRuleHandler interface using client-native client
type CreateHTTPResponseRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteHTTPResponseRuleHandlerImpl implementation of the DeleteHTTPResponseRuleHandler interface using client-native client
type DeleteHTTPResponseRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetHTTPResponseRuleHandlerImpl implementation of the GetHTTPResponseRuleHandler interface using client-native client
type GetHTTPResponseRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetHTTPResponseRulesHandlerImpl implementation of the GetHTTPResponseRulesHandler interface using client-native client
type GetHTTPResponseRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceHTTPResponseRuleHandlerImpl implementation of the ReplaceHTTPResponseRuleHandler interface using client-native client
type ReplaceHTTPResponseRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.CreateHTTPResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateHTTPResponseRule(params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewCreateHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return http_response_rule.NewCreateHTTPResponseRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.DeleteHTTPResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteHTTPResponseRule(params.ID, params.ParentType, params.ParentName, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewDeleteHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return http_response_rule.NewDeleteHTTPResponseRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.GetHTTPResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetHTTPResponseRule(params.ID, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewGetHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return http_response_rule.NewGetHTTPResponseRuleOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetHTTPResponseRulesHandlerImpl) Handle(params http_response_rule.GetHTTPResponseRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetHTTPResponseRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewGetHTTPResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return http_response_rule.NewGetHTTPResponseRulesOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceHTTPResponseRuleHandlerImpl) Handle(params http_response_rule.ReplaceHTTPResponseRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditHTTPResponseRule(params.ID, params.ParentType, params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_response_rule.NewReplaceHTTPResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return http_response_rule.NewReplaceHTTPResponseRuleOK().WithPayload(params.Data)
}
