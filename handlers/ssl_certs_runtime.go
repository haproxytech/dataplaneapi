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
	"io"

	oapi "github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/dataplaneapi/misc"
	ssl_runtime "github.com/haproxytech/dataplaneapi/operations/s_s_l_runtime"
)

type GetAllCertsHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetAllCertsHandlerImpl) Handle(params ssl_runtime.GetAllCertsParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCertsDefault(int(*e.Code)).WithPayload(e)
	}

	files, err := runtime.ShowCerts()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCertsDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetAllCertsOK().WithPayload(files)
}

type CreateCertHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h CreateCertHandlerImpl) Handle(params ssl_runtime.CreateCertParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCertDefault(int(*e.Code)).WithPayload(e)
	}

	file, ok := params.FileUpload.(*oapi.File)
	if !ok {
		return ssl_runtime.NewCreateCertBadRequest()
	}

	payload, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCertDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.NewCertEntry(file.Header.Filename)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCertDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.SetCertEntry(file.Header.Filename, string(payload))
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCertDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.CommitCertEntry(file.Header.Filename)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCertDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewCreateCertCreated()
}

type GetCertHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetCertHandlerImpl) Handle(params ssl_runtime.GetCertParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCertDefault(int(*e.Code)).WithPayload(e)
	}

	cert, err := runtime.ShowCertificate(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCertDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetCertOK().WithPayload(cert)
}

type ReplaceCertHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h ReplaceCertHandlerImpl) Handle(params ssl_runtime.ReplaceCertParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCertDefault(int(*e.Code)).WithPayload(e)
	}

	payload, err := io.ReadAll(params.FileUpload)
	params.FileUpload.Close()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCertDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.SetCertEntry(params.Name, string(payload))
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCertDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.CommitCertEntry(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewReplaceCertDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewReplaceCertOK()
}

// api.SslRuntimeDeleteCertHandler = nil
type DeleteCertHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h DeleteCertHandlerImpl) Handle(params ssl_runtime.DeleteCertParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCertDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.DeleteCertEntry(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCertDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewDeleteCertNoContent()
}
