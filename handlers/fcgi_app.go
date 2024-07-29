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
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

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

	var err error
	if err = c.createFCGIApplication(params, t, v); err != nil {
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

func (c CreateFCGIAppHandlerImpl) createFCGIApplication(params fcgi_app.CreateFCGIAppParams, t string, v int64) error {
	configuration, err := c.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.CreateStructuredFCGIApplication(params.Data, t, v)
	}
	return configuration.CreateFCGIApplication(params.Data, t, v)
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

	_, r, err := g.getFCGIApplication(params, t)
	if err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewGetFCGIAppDefault(int(*e.Code)).WithPayload(e)
	}

	return fcgi_app.NewGetFCGIAppOK().WithPayload(r)
}

func (g GetFCGIAppHandlerImpl) getFCGIApplication(params fcgi_app.GetFCGIAppParams, t string) (int64, *models.FCGIApp, error) {
	configuration, err := g.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredFCGIApplication(params.Name, t)
	}
	return configuration.GetFCGIApplication(params.Name, t)
}

type GetFCGIAppsHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetFCGIAppsHandlerImpl) Handle(params fcgi_app.GetFCGIAppsParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	_, r, err := g.getFCGIApplications(params, t)
	if err != nil {
		e := misc.HandleError(err)

		return fcgi_app.NewGetFCGIAppsDefault(int(*e.Code)).WithPayload(e)
	}

	return fcgi_app.NewGetFCGIAppsOK().WithPayload(r)
}

func (g GetFCGIAppsHandlerImpl) getFCGIApplications(params fcgi_app.GetFCGIAppsParams, t string) (int64, models.FCGIApps, error) {
	configuration, err := g.Client.Configuration()
	if err != nil {
		return 0, nil, err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.GetStructuredFCGIApplications(t)
	}
	return configuration.GetFCGIApplications(t)
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

	params.Data.Name = params.Name

	var err error
	if err = r.editFCGIApplication(params, t, v); err != nil {
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

func (r ReplaceFCGIAppHandlerImpl) editFCGIApplication(params fcgi_app.ReplaceFCGIAppParams, t string, v int64) error {
	configuration, err := r.Client.Configuration()
	if err != nil {
		return err
	}
	if params.FullSection != nil && *params.FullSection {
		return configuration.EditStructuredFCGIApplication(params.Name, params.Data, t, v)
	}
	return configuration.EditFCGIApplication(params.Name, params.Data, t, v)
}
