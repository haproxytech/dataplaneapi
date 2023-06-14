// Copyright 2019 HAProxy Technologies
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
	capture "github.com/haproxytech/dataplaneapi/operations/declare_capture"
)

type CreateDeclareCaptureHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

type DeleteDeclareCaptureHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

type GetDeclareCaptureHandlerImpl struct {
	Client client_native.HAProxyClient
}

type GetDeclareCapturesHandlerImpl struct {
	Client client_native.HAProxyClient
}

type ReplaceDeclareCaptureHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *CreateDeclareCaptureHandlerImpl) Handle(params capture.CreateDeclareCaptureParams, principal interface{}) middleware.Responder {
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
		return capture.NewCreateDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewGetDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}
	_, frontend, err := configuration.GetFrontend(params.Frontend, t)
	if frontend == nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	if err != nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	err = configuration.CreateDeclareCapture(params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewCreateDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return capture.NewCreateDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
			}
			return capture.NewCreateDeclareCaptureCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return capture.NewCreateDeclareCaptureAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return capture.NewCreateDeclareCaptureAccepted().WithPayload(params.Data)
}

func (h *DeleteDeclareCaptureHandlerImpl) Handle(params capture.DeleteDeclareCaptureParams, principal interface{}) middleware.Responder {
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
		return capture.NewCreateDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewDeleteDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}
	err = configuration.DeleteDeclareCapture(params.Index, params.Frontend, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewDeleteDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return capture.NewDeleteDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
			}
			return capture.NewDeleteDeclareCaptureNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return capture.NewCreateDeclareCaptureAccepted().WithReloadID(rID)
	}
	return capture.NewDeleteDeclareCaptureAccepted()
}

func (h *GetDeclareCaptureHandlerImpl) Handle(params capture.GetDeclareCaptureParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewGetDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}

	_, frontend, err := configuration.GetFrontend(params.Frontend, t)
	if frontend == nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	if err != nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	v, data, err := configuration.GetDeclareCapture(params.Index, params.Frontend, t)
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewGetDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}
	return capture.NewGetDeclareCaptureOK().WithPayload(&capture.GetDeclareCaptureOKBody{Version: v, Data: data})
}

func (h *GetDeclareCapturesHandlerImpl) Handle(params capture.GetDeclareCapturesParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewGetDeclareCapturesDefault(int(*e.Code)).WithPayload(e)
	}

	_, frontend, err := configuration.GetFrontend(params.Frontend, t)
	if frontend == nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	if err != nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	v, data, err := configuration.GetDeclareCaptures(params.Frontend, t)
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			return capture.NewGetDeclareCapturesOK().WithPayload(&capture.GetDeclareCapturesOKBody{Version: v, Data: data})
		}
		return capture.NewGetDeclareCapturesDefault(int(*e.Code)).WithPayload(e)
	}
	return capture.NewGetDeclareCapturesOK().WithPayload(&capture.GetDeclareCapturesOKBody{Version: v, Data: data})
}

func (h *ReplaceDeclareCaptureHandlerImpl) Handle(params capture.ReplaceDeclareCaptureParams, principal interface{}) middleware.Responder {
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
		return capture.NewReplaceDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewReplaceDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}

	_, frontend, err := configuration.GetFrontend(params.Frontend, t)
	if frontend == nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	if err != nil {
		return capture.NewGetDeclareCaptureNotFound()
	}
	err = configuration.EditDeclareCapture(params.Index, params.Frontend, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return capture.NewReplaceDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return capture.NewReplaceDeclareCaptureDefault(int(*e.Code)).WithPayload(e)
			}
			return capture.NewReplaceDeclareCaptureOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return capture.NewReplaceDeclareCaptureAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return capture.NewReplaceDeclareCaptureAccepted().WithPayload(params.Data)
}
