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
	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/spoe"
)

// SpoeCreateSpoeGroupHandlerImpl implementation of the SpoeCreateSpoeGroupHandler interface
type SpoeCreateSpoeGroupHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeCreateSpoeGroupHandlerImpl) Handle(params spoe.CreateSpoeGroupParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.CreateGroup(params.Scope, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewCreateSpoeGroupCreated().WithPayload(spoe.NewCreateSpoeGroupCreated().Payload)
}

// SpoeDeleteSpoeGroupHandlerImpl implementation of the SpoeDeleteSpoeGroupHandler interface
type SpoeDeleteSpoeGroupHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeDeleteSpoeGroupHandlerImpl) Handle(params spoe.DeleteSpoeGroupParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.DeleteGroup(params.Scope, params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewDeleteSpoeAgentNoContent()
}

// SpoeGetSpoeGroupsHandlerImpl implementation of the SpoeGetSpoeGroupsHandler interface
type SpoeGetSpoeGroupsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// SpoeGetAllSpoeFilesHandlerImpl implementation of the SpoeGetAllSpoeFilesHandler
func (h *SpoeGetSpoeGroupsHandlerImpl) Handle(params spoe.GetSpoeGroupsParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeGroupsDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeGroupsDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	_, groups, err := ss.GetGroups(params.Scope, t)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeGroupsDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewGetSpoeGroupsOK().WithPayload(groups)
}

// SpoeGetSpoeGroupHandlerImpl implementation of the SpoeGetSpoeGroupHandler interface
type SpoeGetSpoeGroupHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeGetSpoeGroupHandlerImpl) Handle(params spoe.GetSpoeGroupParams, c interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	_, group, err := ss.GetGroup(params.Scope, params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	if group == nil {
		return spoe.NewGetSpoeGroupNotFound()
	}
	return spoe.NewGetSpoeGroupOK().WithPayload(group)
}

// SpoeReplaceSpoeGroupHandlerImpl implementation of the SpoeReplaceSpoeGroupHandler interface
type SpoeReplaceSpoeGroupHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeReplaceSpoeGroupHandlerImpl) Handle(params spoe.ReplaceSpoeGroupParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.EditGroup(params.Scope, params.Data, params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeGroupDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewReplaceSpoeGroupOK()
}
