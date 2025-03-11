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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	cnconf "github.com/haproxytech/client-native/v6/configuration"
	models "github.com/haproxytech/client-native/v6/models"
	cnruntime "github.com/haproxytech/client-native/v6/runtime"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/storage"
)

type StorageGetAllStorageSSLCrtListFilesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *StorageGetAllStorageSSLCrtListFilesHandlerImpl) Handle(params storage.GetAllStorageSSLCrtListFilesParams, principal interface{}) middleware.Responder {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetAllStorageSSLCrtListFilesDefault(int(*e.Code)).WithPayload(e)
	}

	filelist, err := crtListStorage.GetAll()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetAllStorageSSLCrtListFilesDefault(int(*e.Code)).WithPayload(e)
	}

	retFiles := models.SslCrtListFiles{}
	for _, f := range filelist {
		retFiles = append(retFiles, &models.SslCrtListFile{
			File:        f,
			Description: "managed certificate list",
			StorageName: filepath.Base(f),
		})
	}
	return &storage.GetAllStorageSSLCrtListFilesOK{Payload: retFiles}
}

type StorageGetOneStorageSSLCrtListFileHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *StorageGetOneStorageSSLCrtListFileHandlerImpl) Handle(params storage.GetOneStorageSSLCrtListFileParams, principal interface{}) middleware.Responder {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	filename, _, err := crtListStorage.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}
	if filename == "" {
		return storage.NewGetOneStorageSSLCrtListFileNotFound()
	}

	f, err := os.Open(filename)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetOneStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	return storage.NewGetOneStorageSSLCrtListFileOK().WithPayload(f)
}

type StorageDeleteStorageSSLCrtListFileHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageDeleteStorageSSLCrtListFileHandlerImpl) Handle(params storage.DeleteStorageSSLCrtListFileParams, principal interface{}) middleware.Responder {
	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}
	runningConf := strings.NewReader(configuration.Parser().String())

	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	filename, _, err := crtListStorage.Get(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	// this is far from perfect but should provide a basic level of protection
	scanner := bufio.NewScanner(runningConf)

	lineNr := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' && strings.Contains(line, filename) {
			errCode := misc.ErrHTTPConflict
			errMsg := fmt.Sprintf("rejecting attempt to delete file %s referenced in haproxy conf at line %d: %s", filename, lineNr-1, line)
			e := &models.Error{Code: &errCode, Message: &errMsg}
			return storage.NewDeleteStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
		}
		lineNr++
	}

	err = crtListStorage.Delete(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
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
		return storage.NewDeleteStorageSSLCrtListFileNoContent()
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewReplaceStorageMapFileDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewDeleteStorageSSLCrtListFileNoContent()
	}

	rID := h.ReloadAgent.Reload()
	return storage.NewDeleteStorageSSLCrtListFileAccepted().WithReloadID(rID)
}

type StorageReplaceStorageSSLCrtListFileHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageReplaceStorageSSLCrtListFileHandlerImpl) Handle(params storage.ReplaceStorageSSLCrtListFileParams, principal interface{}) middleware.Responder {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	filename, err := crtListStorage.Replace(params.Name, params.Data)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewReplaceStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	retf := &models.SslCrtListFile{
		File:        filename,
		Description: "managed certificate list",
		StorageName: filepath.Base(filename),
		Size:        misc.Int64P(len(params.Data)),
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
		return storage.NewReplaceStorageSSLCrtListFileOK().WithPayload(retf)
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewReplaceStorageMapFileDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewReplaceStorageSSLCrtListFileOK().WithPayload(retf)
	}

	rID := h.ReloadAgent.Reload()
	return storage.NewReplaceStorageSSLCrtListFileAccepted().WithReloadID(rID).WithPayload(retf)
}

type StorageCreateStorageSSLCrtListFileHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageCreateStorageSSLCrtListFileHandlerImpl) Handle(params storage.CreateStorageSSLCrtListFileParams, principal interface{}) middleware.Responder {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	file, ok := params.FileUpload.(*runtime.File)
	if !ok {
		return storage.NewCreateStorageSSLCrtListFileBadRequest()
	}
	filename, size, err := crtListStorage.Create(file.Header.Filename, params.FileUpload)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCrtListFileDefault(int(*e.Code)).WithPayload(e)
	}

	retf := &models.SslCrtListFile{
		File:        filename,
		Description: "managed certificate list",
		StorageName: filepath.Base(filename),
		Size:        &size,
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
		return storage.NewCreateStorageSSLCrtListFileCreated().WithPayload(retf)
	}
	rID := h.ReloadAgent.Reload()
	return storage.NewCreateStorageSSLCrtListFileAccepted().WithReloadID(rID).WithPayload(retf)
}

//
// crt-list entries
//

type StorageGetStorageSSLCrtListEntriesHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *StorageGetStorageSSLCrtListEntriesHandlerImpl) Handle(params storage.GetStorageSSLCrtListEntriesParams, principal interface{}) middleware.Responder {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetStorageSSLCrtListEntriesDefault(int(*e.Code)).WithPayload(e)
	}
	contents, err := crtListStorage.GetContents(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetStorageSSLCrtListEntriesDefault(int(*e.Code)).WithPayload(e)
	}
	entries, err := cnruntime.ParseCrtListEntries(contents)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewGetStorageSSLCrtListEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return storage.NewGetStorageSSLCrtListEntriesOK().WithPayload(entries)
}

type StorageCreateStorageSSLCrtListEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *StorageCreateStorageSSLCrtListEntryHandlerImpl) Handle(params storage.CreateStorageSSLCrtListEntryParams, principal interface{}) middleware.Responder {
	entry := params.Data
	if entry.File == "" {
		e := misc.SetError(400, "missing certificate file name")
		return storage.NewCreateStorageSSLCrtListEntryBadRequest().WithPayload(e)
	}

	// Serialize the entry to a single line.
	var sb strings.Builder
	sb.Grow(len(entry.File) + len(entry.SSLBindConfig) + 32)
	sb.WriteString(entry.File)
	if entry.SSLBindConfig != "" {
		sb.WriteString(" [")
		sb.WriteString(entry.SSLBindConfig)
		sb.WriteByte(']')
	}
	if len(entry.SNIFilter) > 0 {
		sb.WriteByte(' ')
		sb.WriteString(strings.Join(entry.SNIFilter, " "))
	}
	sb.WriteByte('\n')
	line := sb.String()

	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	// Append the entry to the crt-list file.
	content, err := crtListStorage.GetContents(params.Name)
	if err == nil {
		_, err = crtListStorage.Replace(params.Name, content+line)
	} else if errors.Is(err, cnconf.ErrObjectDoesNotExist) {
		_, _, err = crtListStorage.Create(params.Name, io.NopCloser(strings.NewReader(line)))
	}
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewCreateStorageSSLCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	// Try to add the entry to the HAProxy runtime. Force the reload if that fails.
	forceReload := false
	if params.ForceReload != nil {
		forceReload = *params.ForceReload
	}
	rt, err := h.Client.Runtime()
	if err != nil {
		forceReload = true
	} else {
		if errAdd := rt.AddCrtListEntry(params.Name, *entry); errAdd != nil {
			log.Warning("failed to add crt-list entry via runtime: ", errAdd)
			forceReload = true
		}
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewCreateStorageSSLCrtListEntryDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewCreateStorageSSLCrtListEntryCreated().WithPayload(entry)
	}
	rID := h.ReloadAgent.Reload()
	return storage.NewCreateStorageSSLCrtListEntryAccepted().WithReloadID(rID)
}

type StorageDeleteStorageSSLCrtListEntryHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// Delete an entry in a crt-list. If the line_number is not provided,
// delete the first match, which is what the runtime API does.
func (h *StorageDeleteStorageSSLCrtListEntryHandlerImpl) Handle(params storage.DeleteStorageSSLCrtListEntryParams, principal interface{}) middleware.Responder {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	// Modify the crt-file on storage.
	content, err := crtListStorage.GetContents(params.Name)
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}
	lineno := int64(0)
	firstMatch := true
	if params.LineNumber < 1 {
		firstMatch = false
	}
	var sb strings.Builder // the modified contents
	sb.Grow(len(content))
	strings.SplitSeq(content, "\n")(func(line string) bool {
		lineno++
		if strings.HasPrefix(line, params.Certificate) {
			if params.LineNumber == lineno || firstMatch {
				// skip that line
				firstMatch = false
				return true
			}
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
		return true
	})
	_, err = crtListStorage.Replace(params.Name, sb.String())
	if err != nil {
		e := misc.HandleError(err)
		return storage.NewDeleteStorageSSLCrtListEntryDefault(int(*e.Code)).WithPayload(e)
	}

	// Try to delete the entry with the HAProxy runtime. Force the reload if that fails.
	forceReload := false
	if params.ForceReload != nil {
		forceReload = *params.ForceReload
	}
	rt, err := h.Client.Runtime()
	if err != nil {
		forceReload = true
	} else {
		num := &params.LineNumber
		if params.LineNumber == 0 {
			num = nil
		}
		if errDel := rt.DeleteCrtListEntry(params.Name, params.Certificate, num); errDel != nil {
			log.Warning("failed to delete crt-list entry via runtime: ", errDel)
			forceReload = true
		}
	}

	if forceReload {
		err := h.ReloadAgent.ForceReload()
		if err != nil {
			e := misc.HandleError(err)
			return storage.NewDeleteStorageSSLCrtListEntryDefault(int(*e.Code)).WithPayload(e)
		}
		return storage.NewDeleteStorageSSLCrtListEntryNoContent()
	}
	rID := h.ReloadAgent.Reload()
	return storage.NewDeleteStorageSSLCrtListEntryAccepted().WithReloadID(rID)
}
