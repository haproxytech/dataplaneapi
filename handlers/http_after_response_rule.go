package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/http_after_response_rule"
)

type CreateHTTPAfterResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (c CreateHTTPAfterResponseRuleHandlerImpl) Handle(params http_after_response_rule.CreateHTTPAfterResponseRuleParams, _ interface{}) middleware.Responder {
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
		return http_after_response_rule.NewCreateHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := c.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewCreateHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateHTTPAfterResponseRule(params.ParentType, params.ParentName, params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewCreateHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = c.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return http_after_response_rule.NewCreateHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}

			return http_after_response_rule.NewCreateHTTPAfterResponseRuleCreated().WithPayload(params.Data)
		}

		return http_after_response_rule.NewCreateHTTPAfterResponseRuleAccepted().WithReloadID(c.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return http_after_response_rule.NewCreateHTTPAfterResponseRuleAccepted().WithPayload(params.Data)
}

type DeleteHTTPAfterResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (d DeleteHTTPAfterResponseRuleHandlerImpl) Handle(params http_after_response_rule.DeleteHTTPAfterResponseRuleParams, _ interface{}) middleware.Responder {
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
		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := d.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.DeleteHTTPAfterResponseRule(params.Index, params.ParentType, params.ParentName, t, v); err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = d.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return http_after_response_rule.NewDeleteHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}

			return http_after_response_rule.NewDeleteHTTPAfterResponseRuleNoContent()
		}

		return http_after_response_rule.NewDeleteHTTPAfterResponseRuleAccepted().WithReloadID(d.ReloadAgent.Reload())
	}

	return http_after_response_rule.NewDeleteHTTPAfterResponseRuleAccepted()
}

type GetHTTPAfterResponseRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetHTTPAfterResponseRuleHandlerImpl) Handle(params http_after_response_rule.GetHTTPAfterResponseRuleParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	v, rule, err := configuration.GetHTTPAfterResponseRule(params.Index, params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	return http_after_response_rule.NewGetHTTPAfterResponseRuleOK().WithPayload(&http_after_response_rule.GetHTTPAfterResponseRuleOKBody{Version: v, Data: rule})
}

type GetHTTPAfterResponseRulesHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetHTTPAfterResponseRulesHandlerImpl) Handle(params http_after_response_rule.GetHTTPAfterResponseRulesParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetHTTPAfterResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}

	v, rules, err := configuration.GetHTTPAfterResponseRules(params.ParentType, params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewGetHTTPAfterResponseRulesDefault(int(*e.Code)).WithPayload(e)
	}

	return http_after_response_rule.NewGetHTTPAfterResponseRulesOK().WithPayload(&http_after_response_rule.GetHTTPAfterResponseRulesOKBody{Version: v, Data: rules})
}

type ReplaceHTTPAfterResponseRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (r ReplaceHTTPAfterResponseRuleHandlerImpl) Handle(params http_after_response_rule.ReplaceHTTPAfterResponseRuleParams, _ interface{}) middleware.Responder {
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
		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := r.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.EditHTTPAfterResponseRule(params.Index, params.ParentType, params.ParentName, params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = r.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return http_after_response_rule.NewReplaceHTTPAfterResponseRuleDefault(int(*e.Code)).WithPayload(e)
			}

			return http_after_response_rule.NewReplaceHTTPAfterResponseRuleOK().WithPayload(params.Data)
		}

		return http_after_response_rule.NewReplaceHTTPAfterResponseRuleAccepted().WithReloadID(r.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return http_after_response_rule.NewReplaceHTTPAfterResponseRuleAccepted().WithPayload(params.Data)
}
