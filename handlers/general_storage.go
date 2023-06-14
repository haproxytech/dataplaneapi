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
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"
	models "github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/storage"
)

// StorageCreateStorageGeneralFileHandlerImpl implementation of the StorageCreateStorageGeneralFileHandler interface using client-native client
type StorageCreateStorageGeneralFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *StorageCreateStorageGeneralFileHandlerImpl) Handle(params storage.CreateStorageGeneralFileParams, principal interface{}) middleware.Responder {
	file, ok := params.FileUpload.(*runtime.File)
	if !ok {
		return storage.NewCreateStorageGeneralFileBadRequest()
	}

	gs, err := h.Client.GeneralStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}

	filename, err := gs.Create(file.Header.Filename, params.FileUpload)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return storage.NewCreateStorageGeneralFileDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}

	me := &models.GeneralFile{
		Description: "managed general use file",
		File:        filename,
		StorageName: filepath.Base(filename),
	}

	return storage.NewCreateStorageGeneralFileCreated().WithPayload(me)
}

// StorageGetAllStorageGeneralFilesHandlerImpl implementation of the StorageGetAllStorageGeneralFilesHandler interface
type StorageGetAllStorageGeneralFilesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *StorageGetAllStorageGeneralFilesHandlerImpl) Handle(params storage.GetAllStorageGeneralFilesParams, principal interface{}) middleware.Responder {
	gs, err := h.Client.GeneralStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetAllStorageGeneralFilesDefault(int(*e.Code)).WithPayload(e)
	}

	filenames, err := gs.GetAll()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetAllStorageGeneralFilesDefault(int(*e.Code)).WithPayload(e)
	}

	retFiles := models.GeneralFiles{}
	for _, f := range filenames {
		retFiles = append(retFiles, &models.GeneralFile{
			Description: "managed general use file",
			File:        f,
			ID:          "",
			StorageName: filepath.Base(f),
		})
	}

	return storage.NewGetAllStorageGeneralFilesOK().WithPayload(retFiles)
}

// StorageGetOneStorageGeneralFileHandlerImpl implementation of the StorageGetOneStorageGeneralFileHandler interface
type StorageGetOneStorageGeneralFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *StorageGetOneStorageGeneralFileHandlerImpl) Handle(params storage.GetOneStorageGeneralFileParams, principal interface{}) middleware.Responder {
	gs, err := h.Client.GeneralStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}

	filename, err := gs.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}
	if filename == "" {
		return storage.NewGetOneStorageGeneralFileNotFound()
	}
	f, err := os.Open(filename)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}
	return storage.NewGetOneStorageGeneralFileOK().WithPayload(f)
}

// StorageDeleteStorageGeneralFileHandlerImpl implementation of the StorageDeleteStorageGeneralFileHandler interface
type StorageDeleteStorageGeneralFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *StorageDeleteStorageGeneralFileHandlerImpl) Handle(params storage.DeleteStorageGeneralFileParams, principal interface{}) middleware.Responder {
	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}

	gs, err := h.Client.GeneralStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}

	runningConf := strings.NewReader(configuration.Parser().String())

	filename, err := gs.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
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
			return storage.NewDeleteStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
		}
		lineNr++
	}

	err = gs.Delete(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}
	return storage.NewDeleteStorageGeneralFileNoContent()
}

// StorageReplaceStorageGeneralFileHandlerImpl implementation of the StorageReplaceStorageGeneralFileHandler interface
type StorageReplaceStorageGeneralFileHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageReplaceStorageGeneralFileHandlerImpl) Handle(params storage.ReplaceStorageGeneralFileParams, principal interface{}) middleware.Responder {
	gs, err := h.Client.GeneralStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
	}

	_, err = gs.Replace(params.Name, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
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
		return storage.NewReplaceStorageGeneralFileNoContent()
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewReplaceStorageGeneralFileDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewReplaceStorageGeneralFileNoContent()
	}
	rID := h.ReloadAgent.Reload()
	return storage.NewReplaceStorageGeneralFileAccepted().WithReloadID(rID)
}
