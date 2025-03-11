// Copyright 2025 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handlers

import (
	"io"

	oapi "github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/dataplaneapi/misc"
	ssl_runtime "github.com/haproxytech/dataplaneapi/operations/s_s_l_runtime"
)

type GetAllCrlHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetAllCrlHandlerImpl) Handle(params ssl_runtime.GetAllCrlParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCrlDefault(int(*e.Code)).WithPayload(e)
	}

	files, err := runtime.ShowCrlFiles()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCrlDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetAllCrlOK().WithPayload(files)
}

type CreateCrlHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h CreateCrlHandlerImpl) Handle(params ssl_runtime.CreateCrlParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCrlDefault(int(*e.Code)).WithPayload(e)
	}

	file, ok := params.FileUpload.(*oapi.File)
	if !ok {
		return ssl_runtime.NewCreateCaFileBadRequest()
	}

	payload, err := io.ReadAll(file)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCrlDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.NewCrlFile(file.Header.Filename)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCrlDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.SetCrlFile(file.Header.Filename, string(payload))
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCrlDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.CommitCrlFile(file.Header.Filename)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCrlDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewCreateCrlCreated()
}

type GetCrlHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetCrlHandlerImpl) Handle(params ssl_runtime.GetCrlParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCrlDefault(int(*e.Code)).WithPayload(e)
	}

	entries, err := runtime.ShowCrlFile(params.Name, params.Index)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCrlDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetCrlOK().WithPayload(*entries)
}

type ReplaceCrlHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h ReplaceCrlHandlerImpl) Handle(params ssl_runtime.ReplaceCrlParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCrlDefault(int(*e.Code)).WithPayload(e)
	}

	payload, err := io.ReadAll(params.FileUpload)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCrlDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.SetCrlFile(params.Name, string(payload))
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCrlDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.CommitCrlFile(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCrlDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewReplaceCrlOK()
}

type DeleteCrlHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h DeleteCrlHandlerImpl) Handle(params ssl_runtime.DeleteCrlParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCrlDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.DeleteCrlFile(params.Name); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCrlDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewDeleteCrlNoContent()
}
