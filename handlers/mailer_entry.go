// Copyright 2022 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/mailer_entry"
)

type CreateMailerEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

type DeleteMailerEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

type GetMailerEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

type GetMailerEntriesHandlerImpl struct {
	Client client_native.HAProxyClient
}

type ReplaceMailerEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Handle executing the request and returning a response
func (h *CreateMailerEntryHandlerImpl) Handle(params mailer_entry.CreateMailerEntryParams, principal interface{}) middleware.Responder {
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
		return mailer_entry.NewCreateMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewCreateMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.CreateMailerEntry(params.MailersSection, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewCreateMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return mailer_entry.NewCreateMailerEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return mailer_entry.NewCreateMailerEntryCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return mailer_entry.NewCreateMailerEntryAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return mailer_entry.NewCreateMailerEntryAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteMailerEntryHandlerImpl) Handle(params mailer_entry.DeleteMailerEntryParams, principal interface{}) middleware.Responder {
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
		return mailer_entry.NewDeleteMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewDeleteMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteMailerEntry(params.Name, params.MailersSection, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewDeleteMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return mailer_entry.NewDeleteMailerEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return mailer_entry.NewDeleteMailerEntryNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return mailer_entry.NewDeleteMailerEntryAccepted().WithReloadID(rID)
	}
	return mailer_entry.NewDeleteMailerEntryAccepted()
}

// Handle executing the request and returning a response
func (h *GetMailerEntryHandlerImpl) Handle(params mailer_entry.GetMailerEntryParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewGetMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	v, me, err := configuration.GetMailerEntry(params.Name, params.MailersSection, t)
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewGetMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	return mailer_entry.NewGetMailerEntryOK().WithPayload(&mailer_entry.GetMailerEntryOKBody{Version: v, Data: me})
}

// Handle executing the request and returning a response
func (h *GetMailerEntriesHandlerImpl) Handle(params mailer_entry.GetMailerEntriesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewGetMailerEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	v, mes, err := configuration.GetMailerEntries(params.MailersSection, t)
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewGetMailerEntriesDefault(int(*e.Code)).WithPayload(e)
	}
	return mailer_entry.NewGetMailerEntriesOK().WithPayload(&mailer_entry.GetMailerEntriesOKBody{Version: v, Data: mes})
}

// Handle executing the request and returning a response
func (h *ReplaceMailerEntryHandlerImpl) Handle(params mailer_entry.ReplaceMailerEntryParams, principal interface{}) middleware.Responder {
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
		return mailer_entry.NewReplaceMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewReplaceMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditMailerEntry(params.Name, params.MailersSection, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return mailer_entry.NewReplaceMailerEntryDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return mailer_entry.NewReplaceMailerEntryDefault(int(*e.Code)).WithPayload(e)
			}
			return mailer_entry.NewReplaceMailerEntryOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return mailer_entry.NewReplaceMailerEntryAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return mailer_entry.NewReplaceMailerEntryAccepted().WithPayload(params.Data)
}
