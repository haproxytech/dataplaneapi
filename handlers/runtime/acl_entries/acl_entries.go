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

package acl_entries

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all acl_entries routes onto r using spec-based request validation.
func RegisterRouter(r chi.Router, client client_native.HAProxyClient) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{Client: client}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for HAProxy runtime ACL entries.
type HandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *HandlerImpl) GetAllRuntimeAclFileEntries(w http.ResponseWriter, r *http.Request, parentName string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	files, err := rt.GetACLFilesEntries(parentName)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, files)
}

func (h *HandlerImpl) CreateRuntimeAclFileEntry(w http.ResponseWriter, r *http.Request, parentName string) {
	var data CreateRuntimeAclFileEntryJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.AddACLFileEntry(parentName, data.Value); err != nil {
		respond.Error(w, err)
		return
	}

	fileEntry, err := rt.GetACLFileEntry(parentName, data.Value)
	if err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, fileEntry)
}

func (h *HandlerImpl) AddPayloadRuntimeACL(w http.ResponseWriter, r *http.Request, parentName string) {
	var data AddPayloadRuntimeACLJSONRequestBody
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = rt.AddACLAtomic(parentName, data); err != nil {
		respond.RuntimeError(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, data)
}

func (h *HandlerImpl) DeleteRuntimeAclFileEntry(w http.ResponseWriter, r *http.Request, parentName string, id string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	if err = rt.DeleteACLFileEntry(parentName, "#"+id); err != nil {
		respond.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HandlerImpl) GetRuntimeAclFileEntry(w http.ResponseWriter, r *http.Request, parentName string, id string) {
	rt, err := h.Client.Runtime()
	if err != nil {
		respond.Error(w, err)
		return
	}
	fileEntry, err := rt.GetACLFileEntry(parentName, id)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, fileEntry)
}
