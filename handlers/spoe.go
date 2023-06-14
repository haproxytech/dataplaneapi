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
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/spoe"
)

// SpoeCreateSpoeHandlerImpl implementation of the SpoeCreateSpoeAgentHandler interface
type SpoeCreateSpoeHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeCreateSpoeHandlerImpl) Handle(params spoe.CreateSpoeParams, principal interface{}) middleware.Responder {
	file, ok := params.FileUpload.(*runtime.File)
	if !ok {
		return spoe.NewCreateSpoeBadRequest()
	}

	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeDefault(int(*e.Code)).WithPayload(e)
	}

	path, err := spoeStorage.Create(file.Header.Filename, params.FileUpload)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewCreateSpoeCreated().WithPayload(path)
}

// SpoeDeleteSpoeFileHandlerImpl implementation of the SpoeDeleteSpoeFileHandler interface
type SpoeDeleteSpoeFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeDeleteSpoeFileHandlerImpl) Handle(params spoe.DeleteSpoeFileParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeFileDefault(int(*e.Code)).WithPayload(e)
	}

	err = spoeStorage.Delete(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeFileDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewDeleteSpoeFileNoContent()
}

// SpoeGetAllSpoeFilesHandlerImpl implementation of the SpoeGetAllSpoeFilesHandler interface
type SpoeGetAllSpoeFilesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// SpoeGetAllSpoeFilesHandlerImpl implementation of the SpoeGetAllSpoeFilesHandler
func (h *SpoeGetAllSpoeFilesHandlerImpl) Handle(params spoe.GetAllSpoeFilesParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetAllSpoeFilesDefault(int(*e.Code)).WithPayload(e)
	}

	files, err := spoeStorage.GetAll()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetAllSpoeFilesDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewGetAllSpoeFilesOK().WithPayload(files)
}

// SpoeGetOneSpoeFileHandlerImpl implementation of the MapsGetOneRuntimeMapHandler interface
type SpoeGetOneSpoeFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeGetOneSpoeFileHandlerImpl) Handle(params spoe.GetOneSpoeFileParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetOneSpoeFileDefault(int(*e.Code)).WithPayload(e)
	}

	path, err := spoeStorage.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetOneSpoeFileDefault(int(*e.Code)).WithPayload(e)
	}
	if path == "" {
		return spoe.NewGetOneSpoeFileNotFound()
	}
	return spoe.NewGetOneSpoeFileOK().WithPayload(&spoe.GetOneSpoeFileOKBody{Data: path})
}
