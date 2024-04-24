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
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	"github.com/haproxytech/dataplaneapi/haproxy"
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

// handleDeprecatedBackendFields adds backward compatibility support for the fields
// force_persist and ignore_persist that are deprecated in favour of force_persist_list
// and ignore_persist_list.
func handleDeprecatedBackendFields(method string, payload *models.Backend, onDisk *models.Backend) {
	// A fair amount of code duplication in this function is tolerated because this code is expected to be
	// short-lived - it should be removed at the end of the sunset period for the deprecated fields.

	if method == http.MethodGet {
		// Populate force_persist with the first element of force_persist_list if it is present.
		if len(payload.ForcePersistList) > 0 {
			payload.ForcePersist = &models.BackendForcePersist{
				Cond:     payload.ForcePersistList[0].Cond,
				CondTest: payload.ForcePersistList[0].CondTest,
			}
		}
		// Populate ignore_persist with the first element of ignore_persist_list if it is present.
		if len(payload.IgnorePersistList) > 0 {
			payload.IgnorePersist = &models.BackendIgnorePersist{
				Cond:     payload.IgnorePersistList[0].Cond,
				CondTest: payload.IgnorePersistList[0].CondTest,
			}
		}
		return
	}

	if payload.ForcePersist != nil && len(payload.ForcePersistList) == 0 {
		if method == http.MethodPost || (method == http.MethodPut && (onDisk == nil || len(onDisk.ForcePersistList) == 0)) {
			// Deprecated force_persist is present in a POST payload, or in a PUT payload when force_persist_list does not yet exist in the backend.
			// Transform it into force_persist_list with the only element.
			payload.ForcePersistList = []*models.ForcePersist{{
				Cond:     payload.ForcePersist.Cond,
				CondTest: payload.ForcePersist.CondTest,
			}}
		} else {
			// Deprecated force_persist is present in a PUT payload, and force_persist_list already exists in the backend.
			// Preserve the existing force_persist_list, and add or reposition the submitted force_persist to be its first element.
			found := -1
			for i, item := range onDisk.ForcePersistList {
				if *item.Cond == *payload.ForcePersist.Cond && *item.CondTest == *payload.ForcePersist.CondTest {
					found = i
					break
				}
			}
			switch found {
			case -1:
				// force_persist value is not part of existing force_persist_list - insert it in the first position.
				payload.ForcePersistList = slices.Insert(onDisk.ForcePersistList, 0, &models.ForcePersist{
					Cond:     payload.ForcePersist.Cond,
					CondTest: payload.ForcePersist.CondTest,
				})
			case 0:
				// force_persist value matches the first element of force_persist_list - preserve it without modification.
				payload.ForcePersistList = onDisk.ForcePersistList
			default:
				// force_persist value matches another element of force_persist_list - move it to the first position.
				payload.ForcePersistList = slices.Concat(onDisk.ForcePersistList[found:found+1], onDisk.ForcePersistList[:found], onDisk.ForcePersistList[found+1:])
			}
		}
	}

	if payload.IgnorePersist != nil && len(payload.IgnorePersistList) == 0 {
		if method == http.MethodPost || (method == http.MethodPut && (onDisk == nil || len(onDisk.IgnorePersistList) == 0)) {
			// Deprecated ignore_persist is present in a POST payload, or in a PUT payload when ignore_persist_list does not yet exist in the backend.
			// Transform it into ignore_persist_list with the only element.
			payload.IgnorePersistList = []*models.IgnorePersist{{
				Cond:     payload.IgnorePersist.Cond,
				CondTest: payload.IgnorePersist.CondTest,
			}}
		} else {
			// Deprecated ignore_persist is present in a PUT payload, and ignore_persist_list already exists in the backend.
			// Preserve the existing ignore_persist_list, and add or reposition the submitted ignore_persist to be its first element.
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
		}
	}

	// Remove force_persist and ignore_persist from the payload - at this point, they were either processed,
	// or not present in the payload, or will be ignored because non-deprecated variants were submitted at the same time.
	payload.ForcePersist = nil
	payload.IgnorePersist = nil
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

	configuration, err := h.Client.Configuration()
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewCreateBackendDefault(int(*e.Code)).WithPayload(e)
	}

	// Populate force_persist_list and ignore_persist_list if the corresponding
	// deprecated fields force_persist or ignore_persist are present in the request payload.
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

	_, bck, err := configuration.GetBackend(params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewGetBackendDefault(int(*e.Code)).WithPayload(e)
	}

	// Populate deprecated force_persist and ignore_persist fields in returned response.
	handleDeprecatedBackendFields(http.MethodGet, bck, nil)

	return backend.NewGetBackendOK().WithPayload(bck)
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

	_, bcks, err := configuration.GetBackends(t)
	if err != nil {
		e := misc.HandleError(err)
		return backend.NewGetBackendsDefault(int(*e.Code)).WithPayload(e)
	}

	// Populate deprecated force_persist and ignore_persist fields in returned response.
	for _, bck := range bcks {
		handleDeprecatedBackendFields(http.MethodGet, bck, nil)
	}

	return backend.NewGetBackendsOK().WithPayload(bcks)
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

	// Populate or modify force_persist_list and ignore_persist_list if the corresponding
	// deprecated fields force_persist or ignore_persist are present in the request payload.
	if params.Data.ForcePersist != nil || params.Data.IgnorePersist != nil {
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
