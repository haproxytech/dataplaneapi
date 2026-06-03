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

package acl

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
)

// RegisterRouter registers all ACL routes onto r using spec-based request validation
// and a shared error handler.
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

// HandlerImpl implements ServerInterface for HAProxy ACL configuration.
type HandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// --- Backend ---

func (h *HandlerImpl) GetAllAclBackend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllAclBackendParams) {
	h.getAllAcl(w, r, "backend", parentName, params.AclName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllAclBackend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllAclBackendParams) {
	h.replaceAllAcl(w, r, "backend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteAclBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteAclBackendParams) {
	h.deleteAcl(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetAclBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetAclBackendParams) {
	h.getAcl(w, r, "backend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateAclBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateAclBackendParams) {
	h.createAcl(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceAclBackend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceAclBackendParams) {
	h.replaceAcl(w, r, "backend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Defaults ---

func (h *HandlerImpl) GetAllAclDefaults(w http.ResponseWriter, r *http.Request, parentName string, params GetAllAclDefaultsParams) {
	h.getAllAcl(w, r, "defaults", parentName, params.AclName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllAclDefaults(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllAclDefaultsParams) {
	h.replaceAllAcl(w, r, "defaults", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteAclDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteAclDefaultsParams) {
	h.deleteAcl(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetAclDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetAclDefaultsParams) {
	h.getAcl(w, r, "defaults", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateAclDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateAclDefaultsParams) {
	h.createAcl(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceAclDefaults(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceAclDefaultsParams) {
	h.replaceAcl(w, r, "defaults", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- FCGIApp ---

func (h *HandlerImpl) GetAllAclFCGIApp(w http.ResponseWriter, r *http.Request, parentName string, params GetAllAclFCGIAppParams) {
	h.getAllAcl(w, r, "fcgi-app", parentName, params.AclName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllAclFCGIApp(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllAclFCGIAppParams) {
	h.replaceAllAcl(w, r, "fcgi-app", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteAclFCGIApp(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteAclFCGIAppParams) {
	h.deleteAcl(w, r, "fcgi-app", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetAclFCGIApp(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetAclFCGIAppParams) {
	h.getAcl(w, r, "fcgi-app", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateAclFCGIApp(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateAclFCGIAppParams) {
	h.createAcl(w, r, "fcgi-app", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceAclFCGIApp(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceAclFCGIAppParams) {
	h.replaceAcl(w, r, "fcgi-app", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Frontend ---

func (h *HandlerImpl) GetAllAclFrontend(w http.ResponseWriter, r *http.Request, parentName string, params GetAllAclFrontendParams) {
	h.getAllAcl(w, r, "frontend", parentName, params.AclName, params.TransactionId)
}

func (h *HandlerImpl) ReplaceAllAclFrontend(w http.ResponseWriter, r *http.Request, parentName string, params ReplaceAllAclFrontendParams) {
	h.replaceAllAcl(w, r, "frontend", parentName, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) DeleteAclFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params DeleteAclFrontendParams) {
	h.deleteAcl(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) GetAclFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params GetAclFrontendParams) {
	h.getAcl(w, r, "frontend", parentName, index, params.TransactionId)
}

func (h *HandlerImpl) CreateAclFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params CreateAclFrontendParams) {
	h.createAcl(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

func (h *HandlerImpl) ReplaceAclFrontend(w http.ResponseWriter, r *http.Request, parentName string, index int, params ReplaceAclFrontendParams) {
	h.replaceAcl(w, r, "frontend", parentName, index, params.TransactionId, int64(params.Version), params.ForceReload)
}

// --- Shared implementations ---

func (h *HandlerImpl) getAllAcl(w http.ResponseWriter, r *http.Request, parentType, parentName, aclName, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	var aclNamePtr *string
	if aclName != "" {
		aclNamePtr = &aclName
	}
	var acls Acls
	if aclNamePtr != nil {
		_, acls, err = cfg.GetACLs(parentType, parentName, txID, *aclNamePtr)
	} else {
		_, acls, err = cfg.GetACLs(parentType, parentName, txID)
	}
	if err != nil {
		e := misc.HandleContainerGetError(err)
		if *e.Code == misc.ErrHTTPOk {
			respond.JSON(w, http.StatusOK, Acls{})
			return
		}
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, acls)
}

func (h *HandlerImpl) getAcl(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string) {
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	_, data, err := cfg.GetACL(int64(index), parentType, parentName, txID)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, data)
}

func (h *HandlerImpl) createAcl(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Acl
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.CreateACL(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusCreated, &data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) deleteAcl(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.DeleteACL(int64(index), parentType, parentName, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.NoContent(w)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), nil)
		return
	}
	respond.Accepted(w, "", nil)
}

func (h *HandlerImpl) replaceAcl(w http.ResponseWriter, r *http.Request, parentType, parentName string, index int, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Acl
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.EditACL(int64(index), parentType, parentName, &data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusOK, &data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), &data)
		return
	}
	respond.JSON(w, http.StatusAccepted, &data)
}

func (h *HandlerImpl) replaceAllAcl(w http.ResponseWriter, r *http.Request, parentType, parentName, txID string, version int64, forceReload bool) {
	if txID != "" && forceReload {
		respond.BadRequest(w, "Both force_reload and transaction specified, specify only one")
		return
	}
	var data Acls
	if !respond.DecodeBody(r, w, &data) {
		return
	}
	cfg, err := h.Client.Configuration()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = cfg.ReplaceAcls(parentType, parentName, data, txID, version); err != nil {
		respond.Error(w, err)
		return
	}
	if txID == "" {
		if forceReload {
			if err = h.ReloadAgent.ForceReload(); err != nil {
				respond.Error(w, err)
				return
			}
			respond.JSON(w, http.StatusOK, data)
			return
		}
		respond.Accepted(w, h.ReloadAgent.Reload(), data)
		return
	}
	respond.JSON(w, http.StatusAccepted, data)
}
