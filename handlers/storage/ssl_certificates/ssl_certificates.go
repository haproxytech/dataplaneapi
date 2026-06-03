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

package ssl_certificates

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/strfmt"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all SSL certificate storage routes onto r using spec-based request validation.
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

// HandlerImpl implements ServerInterface for HAProxy SSL certificate storage.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func (h *HandlerImpl) GetAllStorageSSLCertificates(w http.ResponseWriter, r *http.Request) {
	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filelist, err := sslStorage.GetAll()
	if err != nil {
		respond.Error(w, err)
		return
	}

	retFiles := models.SslCertificates{}
	for _, f := range filelist {
		retFiles = append(retFiles, &models.SslCertificate{
			File:        f,
			Description: "managed SSL file",
			StorageName: filepath.Base(f),
		})
	}

	respond.JSON(w, http.StatusOK, retFiles)
}

func (h *HandlerImpl) CreateStorageSSLCertificate(w http.ResponseWriter, r *http.Request, params CreateStorageSSLCertificateParams) {
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

	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	// We need to read the cert here because we are going to write it twice:
	// once to the filesystem, and (optionally) to HAProxy's runtime socket.
	certBody, err := io.ReadAll(file)
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, size, err := sslStorage.Create(header.Filename, io.NopCloser(bytes.NewReader(certBody)))
	if err != nil {
		respond.Error(w, err)
		return
	}

	info, err := sslStorage.GetCertificatesInfo(filename)
	if err != nil {
		respond.Error(w, err)
		return
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

	if params.SkipReload {
		respond.JSON(w, http.StatusCreated, retf)
		return
	}

	if params.ForceReload {
		if err = h.ReloadAgent.ForceReload(); err != nil {
			respond.Error(w, err)
			return
		}
		respond.JSON(w, http.StatusCreated, retf)
		return
	}

	// Try to push the new cert to HAProxy using the runtime socket.
	if err = pushCertToRuntime(h.Client, filename, string(certBody), true); err != nil {
		log.Debugf("failed to push certificate via runtime, reloading instead: %s", err.Error())
		respond.Accepted(w, h.ReloadAgent.Reload(), retf)
		return
	}

	respond.JSON(w, http.StatusCreated, retf)
}

func (h *HandlerImpl) DeleteStorageSSLCertificate(w http.ResponseWriter, r *http.Request, name string, params DeleteStorageSSLCertificateParams) {
	configuration, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	runningConf := strings.NewReader(configuration.Parser().String())

	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, _, err := sslStorage.Get(name)
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

	if err = sslStorage.Delete(name); err != nil {
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

	if rt, rtErr := h.Client.Runtime(); rtErr == nil {
		if err = rt.DeleteCertEntry(name); err == nil {
			respond.NoContent(w)
			return
		}
		log.Debugf("failed to delete certificate via runtime, reloading instead: %s", err.Error())
	}

	respond.Accepted(w, h.ReloadAgent.Reload(), nil)
}

func (h *HandlerImpl) GetOneStorageSSLCertificate(w http.ResponseWriter, r *http.Request, name string) {
	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, size, err := sslStorage.Get(name)
	if err != nil {
		respond.Error(w, err)
		return
	}

	info, err := sslStorage.GetCertificatesInfo(name)
	if err != nil {
		respond.Error(w, err)
		return
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

	respond.JSON(w, http.StatusOK, retf)
}

func (h *HandlerImpl) ReplaceStorageSSLCertificate(w http.ResponseWriter, r *http.Request, name string, params ReplaceStorageSSLCertificateParams) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		respond.BadRequest(w, err.Error())
		return
	}

	sslStorage, err := h.Client.SSLCertStorage()
	if err != nil {
		respond.Error(w, err)
		return
	}

	filename, err := sslStorage.Replace(name, string(data))
	if err != nil {
		respond.Error(w, err)
		return
	}

	info, err := sslStorage.GetCertificatesInfo(filename)
	if err != nil {
		respond.Error(w, err)
		return
	}

	size := int64(len(data))
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

	// Try to push the new cert to HAProxy using the runtime socket.
	if err = pushCertToRuntime(h.Client, filename, string(data), false); err != nil {
		log.Debugf("failed to push certificate via runtime, reloading instead: %s", err.Error())
		respond.Accepted(w, h.ReloadAgent.Reload(), retf)
		return
	}

	respond.JSON(w, http.StatusOK, retf)
}

func pushCertToRuntime(c client_native.HAProxyClient, filename, body string, newCert bool) error {
	rt, err := c.Runtime()
	if err != nil {
		return err
	}

	if newCert {
		if err = rt.NewCertEntry(filename); err != nil {
			return err
		}
	}

	if err = rt.SetCertEntry(filename, body); err != nil {
		return err
	}

	return rt.CommitCertEntry(filename)
}
