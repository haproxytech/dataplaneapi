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
	"github.com/haproxytech/dataplaneapi/operations/process_manager"
)

type DeleteProgramHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (d DeleteProgramHandlerImpl) Handle(params process_manager.DeleteProgramParams, _ interface{}) middleware.Responder {
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

		return process_manager.NewDeleteProgramDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := d.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return process_manager.NewDeleteProgramDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.DeleteProgram(params.Name, t, v); err != nil {
		e := misc.HandleError(err)
		return process_manager.NewDeleteProgramDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = d.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)

				return process_manager.NewDeleteProgramDefault(int(*e.Code)).WithPayload(e)
			}

			return process_manager.NewDeleteProgramNoContent()
		}

		return process_manager.NewDeleteProgramAccepted().WithReloadID(d.ReloadAgent.Reload())
	}

	return process_manager.NewDeleteProgramAccepted()
}

type CreateProgramHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (c CreateProgramHandlerImpl) Handle(params process_manager.CreateProgramParams, _ interface{}) middleware.Responder {
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

		return process_manager.NewCreateProgramDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := c.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return process_manager.NewCreateProgramDefault(int(*e.Code)).WithPayload(e)
	}

	if err = configuration.CreateProgram(params.Data, t, v); err != nil {
		e := misc.HandleError(err)
		return process_manager.NewCreateProgramDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = c.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)

				return process_manager.NewCreateProgramDefault(int(*e.Code)).WithPayload(e)
			}

			return process_manager.NewCreateProgramCreated().WithPayload(params.Data)
		}

		return process_manager.NewCreateProgramAccepted().WithReloadID(c.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return process_manager.NewCreateProgramAccepted().WithPayload(params.Data)
}

type GetProgramHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetProgramHandlerImpl) Handle(params process_manager.GetProgramParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)

		return process_manager.NewGetProgramDefault(int(*e.Code)).WithPayload(e)
	}

	_, r, err := configuration.GetProgram(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)

		return process_manager.NewGetProgramDefault(int(*e.Code)).WithPayload(e)
	}

	return process_manager.NewGetProgramOK().WithPayload(r)
}

type GetProgramsHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (g GetProgramsHandlerImpl) Handle(params process_manager.GetProgramsParams, _ interface{}) middleware.Responder {
	var t string

	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := g.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)

		return process_manager.NewGetProgramsDefault(int(*e.Code)).WithPayload(e)
	}

	_, r, err := configuration.GetPrograms(t)
	if err != nil {
		e := misc.HandleError(err)

		return process_manager.NewGetProgramsDefault(int(*e.Code)).WithPayload(e)
	}

	return process_manager.NewGetProgramsOK().WithPayload(r)
}

type ReplaceProgramHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (r ReplaceProgramHandlerImpl) Handle(params process_manager.ReplaceProgramParams, _ interface{}) middleware.Responder {
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

		return process_manager.NewReplaceProgramDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := r.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)

		return process_manager.NewReplaceProgramDefault(int(*e.Code)).WithPayload(e)
	}

	params.Data.Name = params.Name

	if err = configuration.EditProgram(params.Name, params.Data, t, v); err != nil {
		e := misc.HandleError(err)

		return process_manager.NewReplaceProgramDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			if err = r.ReloadAgent.ForceReload(); err != nil {
				e := misc.HandleError(err)

				return process_manager.NewReplaceProgramDefault(int(*e.Code)).WithPayload(e)
			}

			return process_manager.NewReplaceProgramOK().WithPayload(params.Data)
		}

		return process_manager.NewReplaceProgramAccepted().WithReloadID(r.ReloadAgent.Reload()).WithPayload(params.Data)
	}

	return process_manager.NewReplaceProgramAccepted().WithPayload(params.Data)
}
