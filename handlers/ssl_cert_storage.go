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
	"path/filepath"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/storage"
	models "github.com/haproxytech/models/v2"
)

// StorageGetAllStorageSSLCertificatesHandlerImpl implementation of the StorageGetAllStorageSSLCertificatesHandler interface
type StorageGetAllStorageSSLCertificatesHandlerImpl struct {
	Client *client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *StorageGetAllStorageSSLCertificatesHandlerImpl) Handle(params storage.GetAllStorageSSLCertificatesParams, principal interface{}) middleware.Responder {
	filelist, err := h.Client.SSLCertStorage.GetAll()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetAllStorageSSLCertificatesDefault(int(*e.Code)).WithPayload(e)
	}

	retFiles := []*models.SslCertificate{}
	for _, f := range filelist {
		retFiles = append(retFiles, &models.SslCertificate{
			File:        f,
			Description: "managed SSL file",
			StorageName: filepath.Base(f),
		})
	}
	return &storage.GetAllStorageSSLCertificatesOK{Payload: retFiles}
}

// StorageGetOneStorageMapHandlerImpl implementation of the StorageGetOneStorageSSLCertificateHandler interface
type StorageGetOneStorageSSLCertificateHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *StorageGetOneStorageSSLCertificateHandlerImpl) Handle(params storage.GetOneStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	filename, err := h.Client.SSLCertStorage.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	if filename == "" {
		return storage.NewGetOneStorageSSLCertificateNotFound()
	}
	retf := &models.SslCertificate{
		File:        filename,
		Description: "managed SSL file",
		StorageName: filepath.Base(filename),
	}
	return storage.NewGetOneStorageSSLCertificateOK().WithPayload(retf)
}

// StorageDeleteStorageSSLCertificateHandlerImpl implementation of the StorageDeleteStorageSSLCertificateHandler interface
type StorageDeleteStorageSSLCertificateHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *StorageDeleteStorageSSLCertificateHandlerImpl) Handle(params storage.DeleteStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	err := h.Client.SSLCertStorage.Delete(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	return storage.NewDeleteStorageSSLCertificateNoContent()
}

// StorageReplaceStorageSSLCertificateHandlerImpl implementation of the StorageReplaceStorageSSLCertificateHandler interface
type StorageReplaceStorageSSLCertificateHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *StorageReplaceStorageSSLCertificateHandlerImpl) Handle(params storage.ReplaceStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	filename, err := h.Client.SSLCertStorage.Replace(params.Name, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	retf := &models.SslCertificate{
		File:        filename,
		Description: "managed SSL file",
		StorageName: filepath.Base(filename),
	}
	return storage.NewReplaceStorageSSLCertificateAccepted().WithPayload(retf)
}

// StorageCreateStorageSSLCertificateHandlerImpl implementation of the StorageCreateStorageSSLCertificateHandler interface
type StorageCreateStorageSSLCertificateHandlerImpl struct {
	Client *client_native.HAProxyClient
}

func (h *StorageCreateStorageSSLCertificateHandlerImpl) Handle(params storage.CreateStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	file, ok := params.FileUpload.(*runtime.File)
	if !ok {
		return storage.NewCreateStorageSSLCertificateBadRequest()
	}
	filename, err := h.Client.SSLCertStorage.Create(file.Header.Filename, params.FileUpload)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	retf := &models.SslCertificate{
		File:        filename,
		Description: "managed SSL file",
		StorageName: filepath.Base(filename),
	}
	return storage.NewCreateStorageSSLCertificateCreated().WithPayload(retf)
}
