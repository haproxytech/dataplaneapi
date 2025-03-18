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
	"bufio"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	client_native "github.com/haproxytech/client-native/v6"
	models "github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/storage"
)

// StorageGetAllStorageSSLCertificatesHandlerImpl implementation of the StorageGetAllStorageSSLCertificatesHandler interface
type StorageGetAllStorageSSLCertificatesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *StorageGetAllStorageSSLCertificatesHandlerImpl) Handle(params storage.GetAllStorageSSLCertificatesParams, principal interface{}) middleware.Responder {
	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetAllStorageSSLCertificatesDefault(int(*e.Code)).WithPayload(e)
	}

	filelist, err := sslStorage.GetAll()
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
	Client client_native.HAProxyClient
}

func (h *StorageGetOneStorageSSLCertificateHandlerImpl) Handle(params storage.GetOneStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	filename, size, err := sslStorage.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	if filename == "" {
		return storage.NewGetOneStorageSSLCertificateNotFound()
	}
	info, err := sslStorage.GetCertificatesInfo(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	retf := &models.SslCertificate{
		File:        filename,
		Description: "managed SSL file",
		StorageName: filepath.Base(filename),
		Size:        &size,
		NotAfter:    (*strfmt.DateTime)(info.NotAfter),
		NotBefore:   (*strfmt.DateTime)(info.NotBefore),
		Issuers:     info.Issuers,
		Domains:     info.DNS,
		IPAddresses: info.IPs,
		Subject:     info.Subject,
		Serial:      info.Serial,
	}
	return storage.NewGetOneStorageSSLCertificateOK().WithPayload(retf)
}

// StorageDeleteStorageSSLCertificateHandlerImpl implementation of the StorageDeleteStorageSSLCertificateHandler interface
type StorageDeleteStorageSSLCertificateHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageDeleteStorageSSLCertificateHandlerImpl) Handle(params storage.DeleteStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	runningConf := strings.NewReader(configuration.Parser().String())

	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	filename, _, err := sslStorage.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	// this is far from perfect but should provide a basic level of protection
	scanner := bufio.NewScanner(runningConf)

	lineNr := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, filename) && !strings.HasPrefix(line, "#") {
			errCode := misc.ErrHTTPConflict
			errMsg := fmt.Sprintf("rejecting attempt to delete file %s referenced in haproxy conf at line %d: %s", filename, lineNr-1, line)
			e := &models.Error{Code: &errCode, Message: &errMsg}
			return storage.NewDeleteStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
		}
		lineNr++
	}

	err = sslStorage.Delete(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	skipReload := false
	if params.SkipReload != nil {
		skipReload = *params.SkipReload
	}
	forceReload := false
	if params.ForceReload != nil {
		forceReload = *params.ForceReload
	}

	if skipReload {
		return storage.NewDeleteStorageSSLCertificateNoContent()
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewReplaceStorageMapFileDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewDeleteStorageSSLCertificateNoContent()
	}

	rID := h.ReloadAgent.Reload()
	return storage.NewDeleteStorageSSLCertificateAccepted().WithReloadID(rID)
}

// StorageReplaceStorageSSLCertificateHandlerImpl implementation of the StorageReplaceStorageSSLCertificateHandler interface
type StorageReplaceStorageSSLCertificateHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageReplaceStorageSSLCertificateHandlerImpl) Handle(params storage.ReplaceStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	filename, err := sslStorage.Replace(params.Name, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	info, err := sslStorage.GetCertificatesInfo(filename)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	retf := &models.SslCertificate{
		File:        filename,
		Description: "managed SSL file",
		StorageName: filepath.Base(filename),
		Size:        misc.Int64P(len(params.Data)),
		NotAfter:    (*strfmt.DateTime)(info.NotAfter),
		NotBefore:   (*strfmt.DateTime)(info.NotBefore),
		Issuers:     info.Issuers,
		Domains:     info.DNS,
		IPAddresses: info.IPs,
		Subject:     info.Subject,
		Serial:      info.Serial,
	}

	skipReload := false
	if params.SkipReload != nil {
		skipReload = *params.SkipReload
	}
	forceReload := false
	if params.ForceReload != nil {
		forceReload = *params.ForceReload
	}

	if skipReload {
		return storage.NewReplaceStorageSSLCertificateOK().WithPayload(retf)
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewReplaceStorageMapFileDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewReplaceStorageSSLCertificateOK().WithPayload(retf)
	}

	rID := h.ReloadAgent.Reload()
	return storage.NewReplaceStorageSSLCertificateAccepted().WithReloadID(rID).WithPayload(retf)
}

// StorageCreateStorageSSLCertificateHandlerImpl implementation of the StorageCreateStorageSSLCertificateHandler interface
type StorageCreateStorageSSLCertificateHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageCreateStorageSSLCertificateHandlerImpl) Handle(params storage.CreateStorageSSLCertificateParams, principal interface{}) middleware.Responder {
	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	file, ok := params.FileUpload.(*runtime.File)
	if !ok {
		return storage.NewCreateStorageSSLCertificateBadRequest()
	}
	filename, size, err := sslStorage.Create(file.Header.Filename, params.FileUpload)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	info, err := sslStorage.GetCertificatesInfo(filename)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCertificateDefault(int(*e.Code)).WithPayload(e)
	}
	retf := &models.SslCertificate{
		File:        filename,
		Description: "managed SSL file",
		StorageName: filepath.Base(filename),
		Size:        &size,
		NotAfter:    (*strfmt.DateTime)(info.NotAfter),
		NotBefore:   (*strfmt.DateTime)(info.NotBefore),
		Issuers:     info.Issuers,
		Domains:     info.DNS,
		IPAddresses: info.IPs,
		Subject:     info.Subject,
		Serial:      info.Serial,
	}

	forceReload := false
	if params.ForceReload != nil {
		forceReload = *params.ForceReload
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewReplaceStorageMapFileDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewCreateStorageSSLCertificateCreated().WithPayload(retf)
	}
	rID := h.ReloadAgent.Reload()
	return storage.NewCreateStorageSSLCertificateAccepted().WithReloadID(rID).WithPayload(retf)
}
