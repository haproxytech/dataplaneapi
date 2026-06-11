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

package maps

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all map storage routes onto r using spec-based request validation.
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

// HandlerImpl implements ServerInterface for HAProxy map storage.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetAllStorageMapFiles(w http.ResponseWriter, r *http.Request) {
	st, err := h.Client.MapStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filenames, err := st.GetAll()
	if err != nil {
		respond.Error(w, err)
		return
	}

	retMaps := models.Maps{}
	for _, f := range filenames {
		retMaps = append(retMaps, &models.Map{
			Description: "managed map file",
			File:        f,
			ID:          "",
			StorageName: filepath.Base(f),
		})
	}

	respond.JSON(w, http.StatusOK, retMaps)
}

func (h *HandlerImpl) CreateStorageMapFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respond.MultipartError(w, err)
		return
	}
	file, header, err := r.FormFile("file_upload")
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}
	defer file.Close()

	st, err := h.Client.MapStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, size, err := st.Create(header.Filename, io.NopCloser(file))
	if err != nil {
		respond.RuntimeError(w, err)
		return
	}

	me := &models.Map{
		Description: "managed but not loaded map file (no runtime ID)",
		File:        filename,
		StorageName: filepath.Base(filename),
		Size:        &size,
	}
	// no reload or force reload since this is just a file upload,
	// haproxy configuration has not been changed
	respond.JSON(w, http.StatusCreated, me)
}

func (h *HandlerImpl) DeleteStorageMap(w http.ResponseWriter, r *http.Request, name string) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}

	st, err := h.Client.MapStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	runningConf := strings.NewReader(configuration.Parser().String())

	filename, _, err := st.Get(name)
	if err != nil {
		respond.Error(w, err)
		return
	}

	// this is far from perfect but should provide a basic level of protection
	scanner := bufio.NewScanner(runningConf)

	lineNr := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, filename) && !strings.HasPrefix(line, "#") {
			errCode := misc.ErrHTTPConflict
			errMsg := fmt.Sprintf("rejecting attempt to delete file %s referenced in haproxy conf at line %d: %s", filename, lineNr-1, line)
			respond.JSON(w, int(errCode), &models.Error{Code: &errCode, Message: &errMsg})
			return
		}
		lineNr++
	}

	if err = st.Delete(name); err != nil {
		respond.Error(w, err)
		return
	}

	respond.NoContent(w)
}

func (h *HandlerImpl) GetOneStorageMap(w http.ResponseWriter, r *http.Request, name string) {
	st, err := h.Client.MapStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	rc, err := st.GetRawContents(name)
	if err != nil {
		respond.Error(w, err)
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	respond.Copy(w, rc)
}

func (h *HandlerImpl) ReplaceStorageMapFile(w http.ResponseWriter, r *http.Request, name string, params ReplaceStorageMapFileParams) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	st, err := h.Client.MapStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if _, err = st.Replace(name, string(data)); err != nil {
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
