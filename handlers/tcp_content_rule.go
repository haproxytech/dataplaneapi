package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/tcp_content_rule"
)

//CreateTCPContentRuleHandlerImpl implementation of the CreateTCPContentRuleHandler interface using client-native client
type CreateTCPContentRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//DeleteTCPContentRuleHandlerImpl implementation of the DeleteTCPContentRuleHandler interface using client-native client
type DeleteTCPContentRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//GetTCPContentRuleHandlerImpl implementation of the GetTCPContentRuleHandler interface using client-native client
type GetTCPContentRuleHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//GetTCPContentRulesHandlerImpl implementation of the GetTCPContentRulesHandler interface using client-native client
type GetTCPContentRulesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//ReplaceTCPContentRuleHandlerImpl implementation of the ReplaceTCPContentRuleHandler interface using client-native client
type ReplaceTCPContentRuleHandlerImpl struct {
	Client      *client_native.HAProxyClient
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (h *CreateTCPContentRuleHandlerImpl) Handle(params tcp_content_rule.CreateTCPContentRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.CreateTCPContentRule(params.ParentType, params.ParentName, params.Type, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_content_rule.NewCreateTCPContentRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_content_rule.NewCreateTCPContentRuleCreated().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *DeleteTCPContentRuleHandlerImpl) Handle(params tcp_content_rule.DeleteTCPContentRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.DeleteTCPContentRule(params.ID, params.ParentType, params.ParentName, params.Type, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_content_rule.NewDeleteTCPContentRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_content_rule.NewDeleteTCPContentRuleNoContent()
}

//Handle executing the request and returning a response
func (h *GetTCPContentRuleHandlerImpl) Handle(params tcp_content_rule.GetTCPContentRuleParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bck, err := h.Client.Configuration.GetTCPContentRule(params.ID, params.ParentType, params.ParentName, params.Type, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_content_rule.NewGetTCPContentRuleDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_content_rule.NewGetTCPContentRuleOK().WithPayload(bck)
}

//Handle executing the request and returning a response
func (h *GetTCPContentRulesHandlerImpl) Handle(params tcp_content_rule.GetTCPContentRulesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	bcks, err := h.Client.Configuration.GetTCPContentRules(params.ParentType, params.ParentName, params.Type, t)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_content_rule.NewGetTCPContentRulesDefault(int(*e.Code)).WithPayload(e)
	}
	return tcp_content_rule.NewGetTCPContentRulesOK().WithPayload(bcks)
}

//Handle executing the request and returning a response
func (h *ReplaceTCPContentRuleHandlerImpl) Handle(params tcp_content_rule.ReplaceTCPContentRuleParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	err := h.Client.Configuration.EditTCPContentRule(params.ID, params.ParentType, params.ParentName, params.Type, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return tcp_content_rule.NewReplaceTCPContentRuleDefault(int(*e.Code)).WithPayload(e)
	}
	h.ReloadAgent.Reload()
	return tcp_content_rule.NewReplaceTCPContentRuleOK().WithPayload(params.Data)
}
