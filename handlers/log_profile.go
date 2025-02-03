// Copyright 2025 HAProxy Technologies
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
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/log_profile"
)

type GetLogProfilesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Get all log-profiles
func (h *GetLogProfilesHandlerImpl) Handle(params log_profile.GetLogProfilesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewGetLogProfilesDefault(int(*e.Code)).WithPayload(e)
	}

	_, profiles, err := configuration.GetLogProfiles(t)
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewGetLogProfilesDefault(int(*e.Code)).WithPayload(e)
	}

	return log_profile.NewGetLogProfilesOK().WithPayload(profiles)
}

type GetLogProfileHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Get one log-profile
func (h *GetLogProfileHandlerImpl) Handle(params log_profile.GetLogProfileParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewGetLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	_, profile, err := configuration.GetLogProfile(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewGetLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	return log_profile.NewGetLogProfileOK().WithPayload(profile)
}

type CreateLogProfileHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *CreateLogProfileHandlerImpl) Handle(params log_profile.CreateLogProfileParams, principal interface{}) middleware.Responder {
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
		return log_profile.NewCreateLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewCreateLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateLogProfile(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return log_profile.NewCreateLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_profile.NewCreateLogProfileDefault(int(*e.Code)).WithPayload(e)
			}
			return log_profile.NewCreateLogProfileCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return log_profile.NewCreateLogProfileAccepted().WithReloadID(rID).WithPayload(params.Data)
	}

	return log_profile.NewCreateLogProfileAccepted().WithPayload(params.Data)
}

type EditLogProfileHandler struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *EditLogProfileHandler) Handle(params log_profile.EditLogProfileParams, principal interface{}) middleware.Responder {
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
		return log_profile.NewEditLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewEditLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.EditLogProfile(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewEditLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	return log_profile.NewEditLogProfileOK().WithPayload(params.Data)
}

type DeleteLogProfileHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *DeleteLogProfileHandlerImpl) Handle(params log_profile.DeleteLogProfileParams, principal interface{}) middleware.Responder {
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
		return log_profile.NewDeleteLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewDeleteLogProfileDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteLogProfile(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return log_profile.NewDeleteLogProfileDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return log_profile.NewDeleteLogProfileDefault(int(*e.Code)).WithPayload(e)
			}
			return log_profile.NewDeleteLogProfileNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return log_profile.NewDeleteLogProfileAccepted().WithReloadID(rID)
	}

	return log_profile.NewDeleteLogProfileAccepted()
}
