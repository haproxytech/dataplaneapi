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
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/dataplaneapi/misc"
	ssl_runtime "github.com/haproxytech/dataplaneapi/operations/s_s_l_runtime"
)

type GetAllCrtListsHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetAllCrtListsHandlerImpl) Handle(params ssl_runtime.GetAllCrtListsParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCrtListsDefault(int(*e.Code)).WithPayload(e)
	}

	files, err := runtime.ShowCrtLists()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCrtListsDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetAllCrtListsOK().WithPayload(files)
}

type GetAllCrtListEntriesHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h GetAllCrtListEntriesHandlerImpl) Handle(params ssl_runtime.GetAllCrtListEntriesParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCrtListEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	crtList, err := runtime.ShowCrtListEntries(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewGetAllCrtListEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewGetAllCrtListEntriesOK().WithPayload(crtList)
}

type AddCrtListEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h AddCrtListEntryHandlerImpl) Handle(params ssl_runtime.AddCrtListEntryParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewAddCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.AddCrtListEntry(params.Name, *params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewAddCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewAddCrtListEntryCreated()
}

type DeleteCrtListEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h DeleteCrtListEntryHandlerImpl) Handle(params ssl_runtime.DeleteCrtListEntryParams, i interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	if err = runtime.DeleteCrtListEntry(params.Name, params.CertFile, params.LineNumber); err != nil {
		e := misc.HandleError(err)
		return ssl_runtime.NewDeleteCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	return ssl_runtime.NewDeleteCrtListEntryNoContent()
}
