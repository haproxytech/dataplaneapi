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
	"net/http"
	"slices"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/backend"
)

// CreateBackendHandlerImpl implementation of the CreateBackendHandler interface using client-native client
type CreateBackendHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// DeleteBackendHandlerImpl implementation of the DeleteBackendHandler interface using client-native client
type DeleteBackendHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

// GetBackendHandlerImpl implementation of the GetBackendHandler interface using client-native client
type GetBackendHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetBackendsHandlerImpl implementation of the GetBackendsHandler interface using client-native client
type GetBackendsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// ReplaceBackendHandlerImpl implementation of the ReplaceBackendHandler interface using client-native client
type ReplaceBackendHandlerImpl struct {
	Client      client_native.HAProxyClient
	ReloadAgent haproxy.IReloadAgent
}

func logDeprecatedFieldsWarning(b *models.Backend) {
	if b.Httpclose != "" {
		log.Warningf("Field Httpclose is deprecated. Use HTTPConnectionMode.")
	}
	if b.HTTPKeepAlive != "" {
		log.Warningf("Field HTTPKeepAlive is deprecated. Use HTTPConnectionMode.")
	}
	if b.HTTPServerClose != "" {
		log.Warningf("Field HTTPServerClose is deprecated. Use HTTPConnectionMode.")
	}
	if b.IgnorePersist != nil {
		log.Warningf("Field ignore_persist is deprecated. Use ignore_persist_list.")
	}
}

// handleDeprecatedBackendFields adds backward compatibility support for the field ignore_persist,
// which is deprecated in favour of ignore_persist_list.
func handleDeprecatedBackendFields(method string, payload *models.Backend, onDisk *models.Backend) {
	switch method {
	case http.MethodGet:
		// Populate ignore_persist with the first element of ignore_persist_list if it is present.
		if len(payload.IgnorePersistList) > 0 {
			payload.IgnorePersist = &models.BackendIgnorePersist{
				Cond:     payload.IgnorePersistList[0].Cond,
				CondTest: payload.IgnorePersistList[0].CondTest,
			}
		}
	case http.MethodPost, http.MethodPut:
		// Do nothing if ignore_persist is not present in the payload.
		if payload.IgnorePersist == nil {
			return
		}
		// If both ignore_persist and ignore_persist_list are present, ignore and remove ignore_persist.
		if len(payload.IgnorePersistList) > 0 {
			payload.IgnorePersist = nil
			return
		}

		if method == http.MethodPut && onDisk != nil && len(onDisk.IgnorePersistList) > 0 {
			// Preserve ignore_persist_list if it exists prior to modification.
			// Add or reposition the value of ignore_persist unless it is already present as the first element.
			found := -1
			for i, item := range onDisk.IgnorePersistList {
				if *item.Cond == *payload.IgnorePersist.Cond && *item.CondTest == *payload.IgnorePersist.CondTest {
					found = i
					break
				}
			}
			switch found {
			case -1:
				// ignore_persist value is not part of existing ignore_persist_list - insert it in the first position.
				payload.IgnorePersistList = slices.Insert(onDisk.IgnorePersistList, 0, &models.IgnorePersist{
					Cond:     payload.IgnorePersist.Cond,
					CondTest: payload.IgnorePersist.CondTest,
				})
			case 0:
				// ignore_persist value matches the first element of ignore_persist_list - preserve it without modification.
				payload.IgnorePersistList = onDisk.IgnorePersistList
			default:
				// ignore_persist value matches another element of ignore_persist_list - move it to the first position.
				payload.IgnorePersistList = slices.Concat(onDisk.IgnorePersistList[found:found+1], onDisk.IgnorePersistList[:found], onDisk.IgnorePersistList[found+1:])
			}
		} else {
			// Otherwise, add ignore_persist_list with the value of the provided ignore_persist as its only element.
			payload.IgnorePersistList = []*models.IgnorePersist{{
				Cond:     payload.IgnorePersist.Cond,
				CondTest: payload.IgnorePersist.CondTest,
			}}
		}

		// Remove ignore_persist from the payload.
		payload.IgnorePersist = nil
	}
}

// Handle executing the request and returning a response
func (h *CreateBackendHandlerImpl) Handle(params backend.CreateBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return backend.NewCreateBackendDefault(int(*e.Code)).WithPayload(e)
	}

	logDeprecatedFieldsWarning(params.Data)

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewCreateBackendDefault(int(*e.Code)).WithPayload(e)
	}

	// Populate ignore_persist_list if the deprecated ignore_persist field
	// is present in the request payload.
	handleDeprecatedBackendFields(http.MethodPost, params.Data, nil)

	err = configuration.CreateBackend(params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewCreateBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend.NewCreateBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return backend.NewCreateBackendCreated().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return backend.NewCreateBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return backend.NewCreateBackendAccepted().WithPayload(params.Data)
}

// Handle executing the request and returning a response
func (h *DeleteBackendHandlerImpl) Handle(params backend.DeleteBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return backend.NewDeleteBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewDeleteBackendDefault(int(*e.Code)).WithPayload(e)
	}

	err = configuration.DeleteBackend(params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewDeleteBackendDefault(int(*e.Code)).WithPayload(e)
	}
	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend.NewDeleteBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return backend.NewDeleteBackendNoContent()
		}
		rID := h.ReloadAgent.Reload()
		return backend.NewDeleteBackendAccepted().WithReloadID(rID)
	}
	return backend.NewDeleteBackendAccepted()
}

// Handle executing the request and returning a response
func (h *GetBackendHandlerImpl) Handle(params backend.GetBackendParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewGetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	v, bck, err := configuration.GetBackend(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewGetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	// Populate deprecated ignore_persist field in returned response.
	handleDeprecatedBackendFields(http.MethodGet, bck, nil)

	return backend.NewGetBackendOK().WithPayload(&backend.GetBackendOKBody{Version: v, Data: bck})
}

// Handle executing the request and returning a response
func (h *GetBackendsHandlerImpl) Handle(params backend.GetBackendsParams, principal interface{}) middleware.Responder {
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewGetBackendsDefault(int(*e.Code)).WithPayload(e)
	}

	v, bcks, err := configuration.GetBackends(t)
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewGetBackendsDefault(int(*e.Code)).WithPayload(e)
	}

	// Populate deprecated ignore_persist field in returned response.
	for _, bck := range bcks {
		handleDeprecatedBackendFields(http.MethodGet, bck, nil)
	}

	return backend.NewGetBackendsOK().WithPayload(&backend.GetBackendsOKBody{Version: v, Data: bcks})
}

// Handle executing the request and returning a response
func (h *ReplaceBackendHandlerImpl) Handle(params backend.ReplaceBackendParams, principal interface{}) middleware.Responder {
	t := ""
	v := int64(0)
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	if params.Version != nil {
		v = *params.Version
	}

	logDeprecatedFieldsWarning(params.Data)

	if t != "" && *params.ForceReload {
		msg := "Both force_reload and transaction specified, specify only one"
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return backend.NewReplaceBackendDefault(int(*e.Code)).WithPayload(e)
	}

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewReplaceBackendDefault(int(*e.Code)).WithPayload(e)
	}

	// Populate or modify ignore_persist_list if the deprecated ignore_persist field
	// is present in the request payload.
	if params.Data.IgnorePersist != nil {
		_, onDisk, confErr := configuration.GetBackend(params.Data.Name, t)
		if confErr != nil {
			e := misc.HandleError(confErr)
			return backend.NewReplaceBackendDefault(int(*e.Code)).WithPayload(e)
		}
		handleDeprecatedBackendFields(http.MethodPut, params.Data, onDisk)
	}

	err = configuration.EditBackend(params.Name, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewReplaceBackendDefault(int(*e.Code)).WithPayload(e)
	}

	if params.TransactionID == nil {
		if *params.ForceReload {
			err := h.ReloadAgent.ForceReload()
			if err != nil {
				e := misc.HandleError(err)
				return backend.NewReplaceBackendDefault(int(*e.Code)).WithPayload(e)
			}
			return backend.NewReplaceBackendOK().WithPayload(params.Data)
		}
		rID := h.ReloadAgent.Reload()
		return backend.NewReplaceBackendAccepted().WithReloadID(rID).WithPayload(params.Data)
	}
	return backend.NewReplaceBackendAccepted().WithPayload(params.Data)
}
