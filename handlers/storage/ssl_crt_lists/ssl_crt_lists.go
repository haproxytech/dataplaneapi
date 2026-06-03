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

package ssl_crt_lists

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	cnconf "github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	cnruntime "github.com/haproxytech/client-native/v6/runtime"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all SSL crt-list storage routes onto r using spec-based request validation.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient, ra haproxy.IReloadAgent) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client, ReloadAgent: ra}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy SSL crt-list storage.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetAllStorageSSLCrtListFiles(w http.ResponseWriter, r *http.Request) {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filelist, err := crtListStorage.GetAll()
	if err != nil {
		respond.Error(w, err)
		return
	}

	retFiles := models.SslCrtListFiles{}
	for _, f := range filelist {
		retFiles = append(retFiles, &models.SslCrtListFile{
			File:        f,
			Description: "managed certificate list",
			StorageName: filepath.Base(f),
		})
	}

	respond.JSON(w, http.StatusOK, retFiles)
}

func (h *HandlerImpl) CreateStorageSSLCrtListFile(w http.ResponseWriter, r *http.Request, params CreateStorageSSLCrtListFileParams) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	file, header, err := r.FormFile("file_upload")
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	defer file.Close()

	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, size, err := crtListStorage.Create(header.Filename, io.NopCloser(file))
	if err != nil {
		respond.Error(w, err)
		return
	}

	retf := &models.SslCrtListFile{
		File:        filename,
		Description: "managed certificate list",
		StorageName: filepath.Base(filename),
		Size:        &size,
	}

	if params.ForceReload {
		if err = h.ReloadAgent.ForceReload(); err != nil {
			respond.Error(w, err)
			return
		}
		respond.JSON(w, http.StatusCreated, retf)
		return
	}

	respond.Accepted(w, h.ReloadAgent.Reload(), retf)
}

func (h *HandlerImpl) DeleteStorageSSLCrtListFile(w http.ResponseWriter, r *http.Request, name string, params DeleteStorageSSLCrtListFileParams) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	runningConf := strings.NewReader(configuration.Parser().String())

	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, _, err := crtListStorage.Get(name)
	if err != nil {
		respond.Error(w, err)
		return
	}

	// this is far from perfect but should provide a basic level of protection
	scanner := bufio.NewScanner(runningConf)

	lineNr := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' && strings.Contains(line, filename) {
			errCode := misc.ErrHTTPConflict
			errMsg := fmt.Sprintf("rejecting attempt to delete file %s referenced in haproxy conf at line %d: %s", filename, lineNr-1, line)
			respond.JSON(w, int(errCode), &models.Error{Code: &errCode, Message: &errMsg})
			return
		}
		lineNr++
	}

	if err = crtListStorage.Delete(name); err != nil {
		respond.Error(w, err)
		return
	}

	if params.SkipReload {
		respond.NoContent(w)
		return
	}

	if params.ForceReload {
		if err = h.ReloadAgent.ForceReload(); err != nil {
			respond.Error(w, err)
			return
		}
		respond.NoContent(w)
		return
	}

	respond.Accepted(w, h.ReloadAgent.Reload(), nil)
}

func (h *HandlerImpl) GetOneStorageSSLCrtListFile(w http.ResponseWriter, r *http.Request, name string) {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	rc, err := crtListStorage.GetRawContents(name)
	if err != nil {
		respond.Error(w, err)
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	respond.Copy(w, rc)
}

func (h *HandlerImpl) ReplaceStorageSSLCrtListFile(w http.ResponseWriter, r *http.Request, name string, params ReplaceStorageSSLCrtListFileParams) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, err := crtListStorage.Replace(name, string(data))
	if err != nil {
		respond.Error(w, err)
		return
	}

	size := int64(len(data))
	retf := &models.SslCrtListFile{
		File:        filename,
		Description: "managed certificate list",
		StorageName: filepath.Base(filename),
		Size:        &size,
	}

	if params.SkipReload {
		respond.JSON(w, http.StatusOK, retf)
		return
	}

	if params.ForceReload {
		if err = h.ReloadAgent.ForceReload(); err != nil {
			respond.Error(w, err)
			return
		}
		respond.JSON(w, http.StatusOK, retf)
		return
	}

	respond.Accepted(w, h.ReloadAgent.Reload(), retf)
}

func (h *HandlerImpl) GetStorageSSLCrtListEntries(w http.ResponseWriter, r *http.Request, name string) {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	contents, err := crtListStorage.GetContents(name)
	if err != nil {
		respond.Error(w, err)
		return
	}

	entries, err := cnruntime.ParseCrtListEntries(contents)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, entries)
}

func (h *HandlerImpl) CreateStorageSSLCrtListEntry(w http.ResponseWriter, r *http.Request, name string, params CreateStorageSSLCrtListEntryParams) {
	var entry SslCrtListEntry
	if !respond.DecodeBody(r, w, &entry) {
		return
	}

	if entry.File == "" {
		respond.BadRequest(w, "missing certificate file name")
		return
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
		respond.Error(w, err)
		return
	}

	// Append the entry to the crt-list file.
	content, err := crtListStorage.GetContents(name)
	if err == nil {
		_, err = crtListStorage.Replace(name, content+line)
	} else if errors.Is(err, cnconf.ErrObjectDoesNotExist) {
		_, _, err = crtListStorage.Create(name, io.NopCloser(strings.NewReader(line)))
	}
	if err != nil {
		respond.Error(w, err)
		return
	}

	// Try to add the entry to the HAProxy runtime. Force the reload if that fails.
	forceReload := params.ForceReload
	rt, err := h.Client.Runtime()
	if err != nil {
		forceReload = true
	} else {
		if errAdd := rt.AddCrtListEntry(name, entry); errAdd != nil {
			log.Warning("failed to add crt-list entry via runtime: ", errAdd)
			forceReload = true
		}
	}

	if forceReload {
		if err = h.ReloadAgent.ForceReload(); err != nil {
			respond.Error(w, err)
			return
		}
		respond.JSON(w, http.StatusCreated, &entry)
		return
	}

	respond.Accepted(w, h.ReloadAgent.Reload(), nil)
}

func (h *HandlerImpl) DeleteStorageSSLCrtListEntry(w http.ResponseWriter, r *http.Request, name string, params DeleteStorageSSLCrtListEntryParams) {
	crtListStorage, err := h.Client.CrtListStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	// Modify the crt-file on storage.
	content, err := crtListStorage.GetContents(name)
	if err != nil {
		respond.Error(w, err)
		return
	}

	lineNumber := int64(params.LineNumber)
	lineno := int64(0)
	firstMatch := lineNumber > 0
	var sbOut strings.Builder // the modified contents
	sbOut.Grow(len(content))
	for line := range strings.SplitSeq(content, "\n") {
		lineno++
		if strings.HasPrefix(line, params.Certificate) {
			if lineNumber == lineno || firstMatch {
				// skip that line
				firstMatch = false
				continue
			}
		}
		sbOut.WriteString(line)
		sbOut.WriteByte('\n')
	}

	if _, err = crtListStorage.Replace(name, sbOut.String()); err != nil {
		respond.Error(w, err)
		return
	}

	// Try to delete the entry with the HAProxy runtime. Force the reload if that fails.
	forceReload := params.ForceReload
	rt, err := h.Client.Runtime()
	if err != nil {
		forceReload = true
	} else {
		var num *int64
		if lineNumber != 0 {
			num = &lineNumber
		}
		if errDel := rt.DeleteCrtListEntry(name, params.Certificate, num); errDel != nil {
			log.Warning("failed to delete crt-list entry via runtime: ", errDel)
			forceReload = true
		}
	}

	if forceReload {
		if err = h.ReloadAgent.ForceReload(); err != nil {
			respond.Error(w, err)
			return
		}
		respond.NoContent(w)
		return
	}

	respond.Accepted(w, h.ReloadAgent.Reload(), nil)
}
