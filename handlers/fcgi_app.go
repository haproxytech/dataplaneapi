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

package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/fcgi_app"
)

type CreateFCGIAppHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (c CreateFCGIAppHandlerImpl) Handle(params fcgi_app.CreateFCGIAppParams, _ interface{}) middleware.Responder {
	var t string
	var v int64

	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		code := misc.ErrHTTPBadRequest

		e := &models.Error{
			Message: misc.StringP("Both force_reload and transaction specified, specify only one"),
			Code:    &code,
		}

		return fcgi_app.NewCreateFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := c.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return fcgi_app.NewCreateFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateFCGIApplication(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return fcgi_app.NewCreateFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = c.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)

				return fcgi_app.NewCreateFCGIAppDefault(int(*e.Code)).WithPayload(e)
			}

			return fcgi_app.NewCreateFCGIAppCreated().WithPayload(params.Data)
		}

		return fcgi_app.NewCreateFCGIAppAccepted().WithReloadID(c.ReloadAgent.Reload()).WithPayload(params.Data)
	}
	return fcgi_app.NewCreateFCGIAppAccepted().WithPayload(params.Data)
}

type DeleteFCGIAppHandlerImpl struct {
	ReloadAgent haproxy.IReloadAgent
	Client      client_native.HAProxyClient
}

func (d DeleteFCGIAppHandlerImpl) Handle(params fcgi_app.DeleteFCGIAppParams, _ interface{}) middleware.Responder {
	var t string
	var v int64

	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		code := misc.ErrHTTPBadRequest

		e := &models.Error{
			Message: misc.StringP("Both force_reload and transaction specified, specify only one"),
			Code:    &code,
		}

		return fcgi_app.NewDeleteFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := d.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return fcgi_app.NewDeleteFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.DeleteFCGIApplication(params.Name, t, v); err != nil {
		e := misc.HandleError(err)
		return fcgi_app.NewDeleteFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = d.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)

				return fcgi_app.NewDeleteFCGIAppDefault(int(*e.Code)).WithPayload(e)
			}

			return fcgi_app.NewDeleteFCGIAppNoContent()
		}

		return fcgi_app.NewDeleteFCGIAppAccepted().WithReloadID(d.ReloadAgent.Reload())
	}

	return fcgi_app.NewDeleteFCGIAppAccepted()
}

type GetFCGIAppHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetFCGIAppHandlerImpl) Handle(params fcgi_app.GetFCGIAppParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewGetFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	v, r, err := configuration.GetFCGIApplication(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewGetFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	return fcgi_app.NewGetFCGIAppOK().WithPayload(&fcgi_app.GetFCGIAppOKBody{Version: v, Data: r})
}

type GetFCGIAppsHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetFCGIAppsHandlerImpl) Handle(params fcgi_app.GetFCGIAppsParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewGetFCGIAppsDefault(int(*e.Code)).WithPayload(e)
	}

	v, r, err := configuration.GetFCGIApplications(t)
	if err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewGetFCGIAppsDefault(int(*e.Code)).WithPayload(e)
	}

	return fcgi_app.NewGetFCGIAppsOK().WithPayload(&fcgi_app.GetFCGIAppsOKBody{Version: v, Data: r})
}

type ReplaceFCGIAppHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (r ReplaceFCGIAppHandlerImpl) Handle(params fcgi_app.ReplaceFCGIAppParams, _ interface{}) middleware.Responder {
	var t string
	var v int64

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		code := misc.ErrHTTPBadRequest

		e := &models.Error{
			Message: misc.StringP("Both force_reload and transaction specified, specify only one"),
			Code:    &code,
		}

		return fcgi_app.NewReplaceFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := r.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewReplaceFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	params.Data.Name = params.Name

	if err = configuration.EditFCGIApplication(params.Name, params.Data, t, v); err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewReplaceFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = r.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)

				return fcgi_app.NewReplaceFCGIAppDefault(int(*e.Code)).WithPayload(e)
			}

			return fcgi_app.NewReplaceFCGIAppOK().WithPayload(params.Data)
		}

		return fcgi_app.NewReplaceFCGIAppAccepted().WithReloadID(r.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return fcgi_app.NewReplaceFCGIAppAccepted().WithPayload(params.Data)
}
