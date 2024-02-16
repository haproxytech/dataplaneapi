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

// SpoeCreateSpoeMessageHandlerImpl implementation of the SpoeCreateSpoeMessageHandler interface
type SpoeCreateSpoeMessageHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeCreateSpoeMessageHandlerImpl) Handle(params spoe.CreateSpoeMessageParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.CreateMessage(params.Scope, params.Data, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewCreateSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewCreateSpoeMessageCreated().WithPayload(spoe.NewCreateSpoeMessageCreated().Payload)
}

// SpoeDeleteSpoeMessageHandlerImpl implementation of the SpoeDeleteSpoeMessageHandler interface
type SpoeDeleteSpoeMessageHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeDeleteSpoeMessageHandlerImpl) Handle(params spoe.DeleteSpoeMessageParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.DeleteMessage(params.Scope, params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewDeleteSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewDeleteSpoeAgentNoContent()
}

// SpoeGetSpoeMessagesHandlerImpl implementation of the SpoeGetSpoeMessagesHandler interface
type SpoeGetSpoeMessagesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// SpoeGetAllSpoeFilesHandlerImpl implementation of the SpoeGetAllSpoeFilesHandler
func (h *SpoeGetSpoeMessagesHandlerImpl) Handle(params spoe.GetSpoeMessagesParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeMessagesDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeMessagesDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	_, messages, err := ss.GetMessages(params.Scope, t)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeMessagesDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewGetSpoeMessagesOK().WithPayload(messages)
}

// SpoeGetSpoeMessageHandlerImpl implementation of the SpoeGetSpoeMessageHandler interface
type SpoeGetSpoeMessageHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeGetSpoeMessageHandlerImpl) Handle(params spoe.GetSpoeMessageParams, c interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	_, message, err := ss.GetMessage(params.Scope, params.Name, t)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewGetSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	if message == nil {
		return spoe.NewGetSpoeMessageNotFound()
	}
	return spoe.NewGetSpoeMessageOK().WithPayload(message)
}

// SpoeReplaceSpoeMessageHandlerImpl implementation of the SpoeReplaceSpoeMessageHandler interface
type SpoeReplaceSpoeMessageHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *SpoeReplaceSpoeMessageHandlerImpl) Handle(params spoe.ReplaceSpoeMessageParams, principal interface{}) middleware.Responder {
	spoeStorage, err := h.Client.Spoe()
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}

	ss, err := spoeStorage.GetSingleSpoe(params.Spoe)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	t := ""
	if params.TransactionID != nil {
		t = *params.TransactionID
	}
	v := int64(0)
	if params.Version != nil {
		v = *params.Version
	}
	err = ss.EditMessage(params.Scope, params.Data, params.Name, t, v)
	if err != nil {
		e := misc.HandleError(err)
		return spoe.NewReplaceSpoeMessageDefault(int(*e.Code)).WithPayload(e)
	}
	return spoe.NewReplaceSpoeMessageOK()
}
