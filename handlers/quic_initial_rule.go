package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	cnconstants "github.com/haproxytech/client-native/v6/configuration/parents"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/quic_initial_rule"
)

type CreateQUICInitialRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (c CreateQUICInitialRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params quic_initial_rule.CreateQUICInitialRuleFrontendParams, _ interface{}) middleware.Responder {
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
		return quic_initial_rule.NewCreateQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := c.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewCreateQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateQUICInitialRule(params.Index, string(parentType), params.ParentName, params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewCreateQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = c.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return quic_initial_rule.NewCreateQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
			}

			return quic_initial_rule.NewCreateQUICInitialRuleFrontendCreated().WithPayload(params.Data)
		}

		return quic_initial_rule.NewCreateQUICInitialRuleFrontendAccepted().WithReloadID(c.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return quic_initial_rule.NewCreateQUICInitialRuleFrontendAccepted().WithPayload(params.Data)
}

type DeleteQUICInitialRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (d DeleteQUICInitialRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params quic_initial_rule.DeleteQUICInitialRuleFrontendParams, _ interface{}) middleware.Responder {
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
		return quic_initial_rule.NewDeleteQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := d.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewDeleteQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.DeleteQUICInitialRule(params.Index, string(parentType), params.ParentName, t, v); err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewDeleteQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = d.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return quic_initial_rule.NewDeleteQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
			}

			return quic_initial_rule.NewDeleteQUICInitialRuleFrontendNoContent()
		}

		return quic_initial_rule.NewDeleteQUICInitialRuleFrontendAccepted().WithReloadID(d.ReloadAgent.Reload())
	}

	return quic_initial_rule.NewDeleteQUICInitialRuleFrontendAccepted()
}

type GetQUICInitialRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetQUICInitialRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params quic_initial_rule.GetQUICInitialRuleFrontendParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewGetQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rule, err := configuration.GetQUICInitialRule(params.Index, string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewGetQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	return quic_initial_rule.NewGetQUICInitialRuleFrontendOK().WithPayload(rule)
}

type GetAllQUICInitialRuleHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetAllQUICInitialRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params quic_initial_rule.GetAllQUICInitialRuleFrontendParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewGetAllQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	_, rules, err := configuration.GetQUICInitialRules(string(parentType), params.ParentName, t)
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewGetAllQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	return quic_initial_rule.NewGetAllQUICInitialRuleFrontendOK().WithPayload(rules)
}

type ReplaceQUICInitialRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (r ReplaceQUICInitialRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params quic_initial_rule.ReplaceQUICInitialRuleFrontendParams, _ interface{}) middleware.Responder {
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
		return quic_initial_rule.NewReplaceQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := r.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewReplaceQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.EditQUICInitialRule(params.Index, string(parentType), params.ParentName, params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewReplaceQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = r.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)
				return quic_initial_rule.NewReplaceQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
			}

			return quic_initial_rule.NewReplaceQUICInitialRuleFrontendOK().WithPayload(params.Data)
		}

		return quic_initial_rule.NewReplaceQUICInitialRuleFrontendAccepted().WithReloadID(r.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return quic_initial_rule.NewReplaceQUICInitialRuleFrontendAccepted().WithPayload(params.Data)
}

type ReplaceAllQUICInitialRuleHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *ReplaceAllQUICInitialRuleHandlerImpl) Handle(parentType cnconstants.CnParentType, params quic_initial_rule.ReplaceAllQUICInitialRuleFrontendParams, principal interface{}) middleware.Responder {
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
		return quic_initial_rule.NewReplaceAllQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewReplaceAllQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.ReplaceQUICInitialRules(string(parentType), params.ParentName, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return quic_initial_rule.NewReplaceAllQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return quic_initial_rule.NewReplaceAllQUICInitialRuleFrontendDefault(int(*e.Code)).WithPayload(e)
			}
			return quic_initial_rule.NewReplaceAllQUICInitialRuleFrontendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return quic_initial_rule.NewReplaceAllQUICInitialRuleFrontendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return quic_initial_rule.NewReplaceAllQUICInitialRuleFrontendAccepted().WithPayload(params.Data)
}
