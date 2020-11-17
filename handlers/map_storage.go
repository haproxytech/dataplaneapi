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
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/storage"
)

//GetMapStorageHandlerImpl implementation of the StorageGetAllStorageMapFilesHandler interface
type GetAllStorageMapFilesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//Handle executing the request and returning a response
func (h *GetAllStorageMapFilesHandlerImpl) Handle(params storage.GetAllStorageMapFilesParams, principal interface{}) middleware.Responder {
	files, err := h.Client.MapStorage.GetAll()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetAllStorageMapFilesDefault(int(*e.Code)).WithPayload(e)
	}
	return &storage.GetAllStorageMapFilesOK{Payload: files}
}

// StorageGetOneStorageMapHandlerImpl implementation of the StorageGetOneStorageMapHandler interface
type GetOneStorageMapHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *GetOneStorageMapHandlerImpl) Handle(params storage.GetOneStorageMapParams, principal interface{}) middleware.Responder {
	m, err := h.Client.MapStorage.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageMapDefault(int(*e.Code)).WithPayload(e)
	}
	if m == "" {
		return storage.NewGetOneStorageMapNotFound()
	}
	return storage.NewGetOneStorageMapOK().WithPayload(m)
}

//StorageDeleteStorageMapHandlerImpl implementation of the StorageDeleteStorageMapHandler interface
type StorageDeleteStorageMapHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *StorageDeleteStorageMapHandlerImpl) Handle(params storage.DeleteStorageMapParams, principal interface{}) middleware.Responder {
	err := h.Client.MapStorage.Delete(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageMapDefault(int(*e.Code)).WithPayload(e)
	}
	return storage.NewDeleteStorageMapNoContent()
}

//StorageReplaceStorageMapFileHandlerImpl implementation of the StorageReplaceStorageMapFileHandler interface
type StorageReplaceStorageMapFileHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *StorageReplaceStorageMapFileHandlerImpl) Handle(params storage.ReplaceStorageMapFileParams, principal interface{}) middleware.Responder {
	f, err := h.Client.MapStorage.Replace(params.Name, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageMapFileDefault(int(*e.Code)).WithPayload(e)
	}
	return storage.NewReplaceStorageMapFileOK().WithPayload(f)
}
