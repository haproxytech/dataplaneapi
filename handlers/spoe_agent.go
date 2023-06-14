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
	client_native "github.com/haproxytech/client-native/v5"

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/spoe"
)

// SpoeCreateSpoeAgentHandlerImpl implementation of the SpoeCreateSpoeAgentHandler interface
type SpoeCreateSpoeAgentHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeCreateSpoeAgentHandlerImpl) Handle(params spoe.CreateSpoeAgentParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.CreateAgent(params.Scope, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewCreateSpoeAgentCreated().WithPayload(spoe.NewCreateSpoeAgentCreated().Payload)
}

// SpoeDeleteSpoeAgentHandlerImpl implementation of the SpoeDeleteSpoeAgentHandler interface
type SpoeDeleteSpoeAgentHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeDeleteSpoeAgentHandlerImpl) Handle(params spoe.DeleteSpoeAgentParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.DeleteAgent(params.Scope, params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewDeleteSpoeAgentNoContent()
}

// SpoeGetSpoeAgentsHandlerImpl implementation of the SpoeGetSpoeAgentsHandler interface
type SpoeGetSpoeAgentsHandlerImpl struct {
	Client client_native.HAProxyClient
}

// SpoeGetAllSpoeFilesHandlerImpl implementation of the SpoeGetAllSpoeFilesHandler
func (h *SpoeGetSpoeAgentsHandlerImpl) Handle(params spoe.GetSpoeAgentsParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetAllSpoeFilesDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetAllSpoeFilesDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	_, agents, err := ss.GetAgents(params.Scope, t)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetAllSpoeFilesDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewGetSpoeAgentsOK().WithPayload(&spoe.GetSpoeAgentsOKBody{Data: agents})
}

// SpoeGetSpoeAgentHandlerImpl implementation of the SpoeGetSpoeAgentHandler interface
type SpoeGetSpoeAgentHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeGetSpoeAgentHandlerImpl) Handle(params spoe.GetSpoeAgentParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v, agent, err := ss.GetAgent(params.Scope, params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	if agent == nil {
		return spoe.NewGetSpoeAgentNotFound()
	}
	return spoe.NewGetSpoeAgentOK().WithPayload(&spoe.GetSpoeAgentOKBody{Version: v, Data: agent})
}

// SpoeReplaceSpoeAgentHandlerImpl implementation of the SpoeReplaceSpoeAgentHandler interface
type SpoeReplaceSpoeAgentHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeReplaceSpoeAgentHandlerImpl) Handle(params spoe.ReplaceSpoeAgentParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.EditAgent(params.Scope, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeAgentDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewReplaceSpoeAgentOK()
}
