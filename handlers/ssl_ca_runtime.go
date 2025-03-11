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

type GetAllCaFilesHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetAllCaFilesHandlerImpl) Handle(params ssl_runtime.GetAllCaFilesParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCaFilesDefault(int(*e.Code)).WithPayload(e)
	}

	files, err := runtime.ShowCAFiles()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCaFilesDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetAllCaFilesOK().WithPayload(files)
}

type CreateCaFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h CreateCaFileHandlerImpl) Handle(params ssl_runtime.CreateCaFileParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	file, ok := params.FileUpload.(*oapi.File)
	if !ok {
		return ssl_runtime.NewCreateCaFileBadRequest()
	}

	// Create a new empty file.
	if err = runtime.NewCAFile(file.Header.Filename); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	// Set its contents.
	payload, err := io.ReadAll(file)
	file.Close()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCaFileDefault(int(*e.Code)).WithPayload(e)
	}
	if err = runtime.SetCAFile(file.Header.Filename, string(payload)); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	// Commit.
	if err = runtime.CommitCAFile(file.Header.Filename); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewCreateCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewCreateCaFileCreated()
}

type GetCaFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetCaFileHandlerImpl) Handle(params ssl_runtime.GetCaFileParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	caFile, err := runtime.GetCAFile(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetCaFileOK().WithPayload(caFile)
}

type SetCaFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h SetCaFileHandlerImpl) Handle(params ssl_runtime.SetCaFileParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewSetCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	payload, err := io.ReadAll(params.FileUpload)
	params.FileUpload.Close()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewSetCaFileDefault(int(*e.Code)).WithPayload(e)
	}
	if err = runtime.SetCAFile(params.Name, string(payload)); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewSetCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.CommitCAFile(params.Name); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewSetCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewCreateCaFileCreated()
}

type DeleteCaFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h DeleteCaFileHandlerImpl) Handle(params ssl_runtime.DeleteCaFileParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.DeleteCAFile(params.Name); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCaFileDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewDeleteCaFileNoContent()
}

type AddCaEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h AddCaEntryHandlerImpl) Handle(params ssl_runtime.AddCaEntryParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewAddCaEntryDefault(int(*e.Code)).WithPayload(e)
	}

	payload, err := io.ReadAll(params.FileUpload)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewAddCaEntryDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.AddCAFileEntry(params.Name, string(payload)); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewAddCaEntryDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.CommitCAFile(params.Name); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewAddCaEntryDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewAddCaEntryCreated()
}

type GetCaEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetCaEntryHandlerImpl) Handle(params ssl_runtime.GetCaEntryParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCaEntryDefault(int(*e.Code)).WithPayload(e)
	}
	entry, err := runtime.ShowCAFile(params.Name, &params.Index)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetCaEntryDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetCaEntryOK().WithPayload(entry)
}
