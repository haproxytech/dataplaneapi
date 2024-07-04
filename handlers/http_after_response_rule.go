package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	cnconstants "github.com/haproxytech/client-native/v6/configuration/parents"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/http_after_response_rule"
)

type CreateHTTPAfterResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (c CreateHTTPAfterResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_after_response_rule.CreateHTTPAfterResponseRuleBackendParams, _ interface{}) middleware.Responder {
	t, v := "", int64(0)

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		e := &models.Error{
			Message: misc.StringP("Both force_reload and transaction specified, specify only one"),
			Code:    misc.Int64P(int(misc.ErrHTTPBadRequest)),
		}
		return http_after_response_rule.NewCreateHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := c.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewCreateHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateHTTPAfterResponseRule(params.Index, string(parentType), params.ParentName, params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewCreateHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = c.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return http_after_response_rule.NewCreateHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}

			return http_after_response_rule.NewCreateHTTPAfterResponseRuleBackendCreated().WithPayload(params.Data)
		}

		return http_after_response_rule.NewCreateHTTPAfterResponseRuleBackendAccepted().WithReloadID(c.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return http_after_response_rule.NewCreateHTTPAfterResponseRuleBackendAccepted().WithPayload(params.Data)
}

type DeleteHTTPAfterResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (d DeleteHTTPAfterResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_after_response_rule.DeleteHTTPAfterResponseRuleBackendParams, _ interface{}) middleware.Responder {
	t, v := "", int64(0)

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		e := &models.Error{
			Message: misc.StringP("Both force_reload and transaction specified, specify only one"),
			Code:    misc.Int64P(int(misc.ErrHTTPBadRequest)),
		}
		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := d.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.DeleteHTTPAfterResponseRule(params.Index, string(parentType), params.ParentName, t, v); err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = d.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return http_after_response_rule.NewDeleteHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}

			return http_after_response_rule.NewDeleteHTTPAfterResponseRuleBackendNoContent()
		}

		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleBackendAccepted().WithReloadID(d.ReloadAgent.Reload())
	}

	return http_after_response_rule.NewDeleteHTTPAfterResponseRuleBackendAccepted()
}

type GetHTTPAfterResponseRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetHTTPAfterResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_after_response_rule.GetHTTPAfterResponseRuleBackendParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetHTTPAfterResponseRule(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	return http_after_response_rule.NewGetHTTPAfterResponseRuleBackendOK().WithPayload(rule)
}

type GetAllHTTPAfterResponseRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetAllHTTPAfterResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_after_response_rule.GetAllHTTPAfterResponseRuleBackendParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetAllHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetHTTPAfterResponseRules(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetAllHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	return http_after_response_rule.NewGetAllHTTPAfterResponseRuleBackendOK().WithPayload(rules)
}

type ReplaceHTTPAfterResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (r ReplaceHTTPAfterResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_after_response_rule.ReplaceHTTPAfterResponseRuleBackendParams, _ interface{}) middleware.Responder {
	t, v := "", int64(0)

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		e := &models.Error{
			Message: misc.StringP("Both force_reload and transaction specified, specify only one"),
			Code:    misc.Int64P(int(misc.ErrHTTPBadRequest)),
		}
		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := r.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.EditHTTPAfterResponseRule(params.Index, string(parentType), params.ParentName, params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = r.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return http_after_response_rule.NewReplaceHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}

			return http_after_response_rule.NewReplaceHTTPAfterResponseRuleBackendOK().WithPayload(params.Data)
		}

		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleBackendAccepted().WithReloadID(r.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return http_after_response_rule.NewReplaceHTTPAfterResponseRuleBackendAccepted().WithPayload(params.Data)
}

type ReplaceAllHTTPAfterResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *ReplaceAllHTTPAfterResponseRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params http_after_response_rule.ReplaceAllHTTPAfterResponseRuleBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return http_after_response_rule.NewReplaceAllHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewReplaceAllHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceHTTPAfterResponseRules(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewReplaceAllHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return http_after_response_rule.NewReplaceAllHTTPAfterResponseRuleBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return http_after_response_rule.NewReplaceAllHTTPAfterResponseRuleBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return http_after_response_rule.NewReplaceAllHTTPAfterResponseRuleBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return http_after_response_rule.NewReplaceAllHTTPAfterResponseRuleBackendAccepted().WithPayload(params.Data)
}
