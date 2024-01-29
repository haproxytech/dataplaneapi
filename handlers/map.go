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
	"path/filepath"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	config "github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/maps"
)

// GetMapsHandlerImpl implementation of the GetAllRuntimeMapFilesHandler interface using client-native client
type GetMapsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *GetMapsHandlerImpl) Handle(params maps.GetAllRuntimeMapFilesParams, principal interface{}) middleware.Responder {
	mapList := []*models.Map{}

	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewShowRuntimeMapDefault(int(*e.Code)).WithPayload(e)
	}

	runtimeMaps, err := runtime.ShowMaps()
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewShowRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}

	mapsDir, err := runtime.GetMapsDir()
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewShowRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}

	for _, m := range runtimeMaps {
		if *params.IncludeUnmanaged || strings.HasPrefix(filepath.Dir(m.File), mapsDir) {
			if strings.HasPrefix(filepath.Dir(m.File), mapsDir) {
				m.StorageName = filepath.Base(m.File)
			}
			mapList = append(mapList, m)
		}
	}
	return maps.NewGetAllRuntimeMapFilesOK().WithPayload(mapList)
}

// GetMapHandlerImpl implementation of the MapsGetOneRuntimeMapHandler interface using client-native client
type GetMapHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *GetMapHandlerImpl) Handle(params maps.GetOneRuntimeMapParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewGetOneRuntimeMapDefault(int(*e.Code)).WithPayload(e)
	}

	m, err := runtime.GetMap(params.Name)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewGetOneRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if m == nil {
		return maps.NewGetOneRuntimeMapNotFound()
	}
	return maps.NewGetOneRuntimeMapOK().WithPayload(m)
}

// ClearMapHandlerImpl implementation of the ClearRuntimeMapHandler interface using client-native client
type ClearMapHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *ClearMapHandlerImpl) Handle(params maps.ClearRuntimeMapParams, principal interface{}) middleware.Responder {
	forceDelete := false
	if params.ForceDelete != nil {
		forceDelete = *params.ForceDelete
	}

	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewClearRuntimeMapDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.ClearMap(params.Name, forceDelete)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewClearRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if *params.ForceSync {
		m, err := runtime.GetMap(params.Name)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewClearRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
		ms := config.NewMapSync()
		_, err = ms.Sync(m, h.Client)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewClearRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
	}
	return maps.NewClearRuntimeMapNoContent()
}

// ShowMapHandlerImpl implementation of the ShowMapHandlerImpl interface using client-native client
type ShowMapHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *ShowMapHandlerImpl) Handle(params maps.ShowRuntimeMapParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewShowRuntimeMapDefault(int(*e.Code)).WithPayload(e)
	}
	m, err := runtime.ShowMapEntries(params.Map)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewShowRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if m == nil {
		return maps.NewShowRuntimeMapNotFound()
	}
	return maps.NewShowRuntimeMapOK().WithPayload(m)
}

// AddMapEntryHandlerImpl implementation of the AddMapEntryHandler interface using client-native client
type AddMapEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *AddMapEntryHandlerImpl) Handle(params maps.AddMapEntryParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewAddMapEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.AddMapEntry(params.Map, params.Data.Key, params.Data.Value)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewAddMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if *params.ForceSync {
		m, err := runtime.GetMap(params.Map)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewAddMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
		ms := config.NewMapSync()
		_, err = ms.Sync(m, h.Client)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewAddMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
	}
	return maps.NewAddMapEntryCreated().WithPayload(params.Data)
}

type MapsAddPayloadRuntimeMapHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *MapsAddPayloadRuntimeMapHandlerImpl) Handle(params maps.AddPayloadRuntimeMapParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewAddPayloadRuntimeMapDefault(int(*e.Code)).WithPayload(e)
	}
	err = runtime.AddMapPayloadVersioned(params.Name, params.Data)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewAddPayloadRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if *params.ForceSync {
		m, err := runtime.GetMap(params.Name)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewAddPayloadRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
		ms := config.NewMapSync()
		_, err = ms.Sync(m, h.Client)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewAddPayloadRuntimeMapDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
	}
	return maps.NewAddPayloadRuntimeMapCreated().WithPayload(params.Data)
}

// GetRuntimeMapEntryHandlerImpl implementation of the GetRuntimeMapEntryHandler interface using client-native client
type GetRuntimeMapEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *GetRuntimeMapEntryHandlerImpl) Handle(params maps.GetRuntimeMapEntryParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewReplaceRuntimeMapEntryDefault(int(*e.Code)).WithPayload(e)
	}
	m, err := runtime.GetMapEntry(params.Map, params.ID)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewGetRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if m == nil {
		return maps.NewGetRuntimeMapEntryNotFound()
	}
	return maps.NewGetRuntimeMapEntryOK().WithPayload(m)
}

// ReplaceRuntimeMapEntryHandlerImpl implementation of the ReplaceRuntimeMapEntryHandler interface using client-native client
type ReplaceRuntimeMapEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *ReplaceRuntimeMapEntryHandlerImpl) Handle(params maps.ReplaceRuntimeMapEntryParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewReplaceRuntimeMapEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.SetMapEntry(params.Map, params.ID, *params.Data.Value)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewGetRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	e, err := runtime.GetMapEntry(params.Map, params.ID)
	if err != nil {
		return maps.NewReplaceRuntimeMapEntryNotFound()
	}
	if *params.ForceSync {
		m, err := runtime.GetMap(params.Map)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewGetRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
		ms := config.NewMapSync()
		_, err = ms.Sync(m, h.Client)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewGetRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
	}
	return maps.NewGetRuntimeMapEntryOK().WithPayload(e)
}

// DeleteRuntimeMapEntryHandlerImpl implementation of the DeleteRuntimeMapEntryHandler interface using client-native client
type DeleteRuntimeMapEntryHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *DeleteRuntimeMapEntryHandlerImpl) Handle(params maps.DeleteRuntimeMapEntryParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return maps.NewDeleteRuntimeMapEntryDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.DeleteMapEntry(params.Map, params.ID)
	if err != nil {
		status := misc.GetHTTPStatusFromErr(err)
		return maps.NewDeleteRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
	}
	if *params.ForceSync {
		m, err := runtime.GetMap(params.Map)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewDeleteRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
		ms := config.NewMapSync()
		_, err = ms.Sync(m, h.Client)
		if err != nil {
			status := misc.GetHTTPStatusFromErr(err)
			return maps.NewDeleteRuntimeMapEntryDefault(status).WithPayload(misc.SetError(status, err.Error()))
		}
	}
	return maps.NewDeleteRuntimeMapEntryNoContent()
}
