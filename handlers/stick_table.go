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
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	client_native "github.com/haproxytech/client-native/v5"
	"github.com/haproxytech/client-native/v5/models"

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/stick_table"
)

// GetStickTablesHandlerImpl implementation of the GetStickTablesHandler interface using client-native client
type GetStickTablesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetStickTableHandlerImpl implementation of the GetStickTableHandler interface using client-native client
type GetStickTableHandlerImpl struct {
	Client client_native.HAProxyClient
}

// GetStickTableEntriesHandlerImpl implementation of the GetStickTableEntriesHandler interface using client-native client
type GetStickTableEntriesHandlerImpl struct {
	Client client_native.HAProxyClient
}

type SetStickTableEntriesHandlerImpl struct {
	Client client_native.HAProxyClient
}

// Handle executing the request and returning a response
func (h *GetStickTablesHandlerImpl) Handle(params stick_table.GetStickTablesParams, principal interface{}) middleware.Responder {
	process := 0
	if params.Process != nil {
		process = int(*params.Process)
	}

	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewGetStickTablesDefault(int(*e.Code)).WithPayload(e)
	}

	stkTS, err := runtime.ShowTables(process)
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewGetStickTablesDefault(int(*e.Code)).WithPayload(e)
	}

	for _, table := range stkTS {
		table.Fields = findTableFields(table.Name, h.Client)
	}

	return stick_table.NewGetStickTablesOK().WithPayload(stkTS)
}

// Handle executing the request and returning a response
func (h *GetStickTableHandlerImpl) Handle(params stick_table.GetStickTableParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewGetStickTableDefault(int(*e.Code)).WithPayload(e)
	}

	stkT, err := runtime.ShowTable(params.Name, int(params.Process))
	if stkT == nil {
		msg := fmt.Sprintf("Stick table %s not found in process %d", params.Name, params.Process)
		c := misc.ErrHTTPNotFound
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return stick_table.NewGetStickTableDefault(int(*e.Code)).WithPayload(e)
	}

	stkT.Fields = findTableFields(stkT.Name, h.Client)
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewGetStickTableDefault(int(*e.Code)).WithPayload(e)
	}

	return stick_table.NewGetStickTableOK().WithPayload(stkT)
}

// Handle executing the request and returning a response
func (h *SetStickTableEntriesHandlerImpl) Handle(params stick_table.SetStickTableEntriesParams, principal interface{}) middleware.Responder {
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewSetStickTableEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	err = runtime.SetTableEntry(params.StickTable, *params.StickTableEntry.Key, *params.StickTableEntry.DataType, int(params.Process))
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewSetStickTableEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	return stick_table.NewSetStickTableEntriesNoContent()
}

// Handle executing the request and returning a response
func (h *GetStickTableEntriesHandlerImpl) Handle(params stick_table.GetStickTableEntriesParams, principal interface{}) middleware.Responder {
	filter := make([]string, 0)
	if params.Filter != nil {
		filter = strings.Split(*params.Filter, ",")
	}

	key := ""
	if params.Key != nil {
		key = *params.Key
	}
	runtime, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewGetStickTableEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	stkEntries, err := runtime.GetTableEntries(params.StickTable, int(params.Process), filter, key)
	if err != nil {
		e := misc.HandleError(err)
		return stick_table.NewGetStickTableEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	// if no entries return empty array
	if len(stkEntries) == 0 {
		return stick_table.NewGetStickTableEntriesOK().WithPayload(stkEntries)
	}

	// else check for pagination
	offset := int64(0)
	if params.Offset != nil {
		offset = *params.Offset
	}

	if int(offset) >= len(stkEntries) {
		msg := fmt.Sprintf("Offset %d is larger than the slice size %d", offset, len(stkEntries))
		c := misc.ErrHTTPBadRequest
		e := &models.Error{
			Message: &msg,
			Code:    &c,
		}
		return stick_table.NewGetStickTableEntriesDefault(int(*e.Code)).WithPayload(e)
	}

	if params.Count != nil {
		if int(offset+*params.Count) >= len(stkEntries) {
			stkEntries = stkEntries[offset:]
		} else {
			stkEntries = stkEntries[offset : offset+*params.Count]
		}
	} else {
		stkEntries = stkEntries[offset:]
	}
	return stick_table.NewGetStickTableEntriesOK().WithPayload(stkEntries)
}

func findTableFields(name string, client client_native.HAProxyClient) []*models.StickTableField {
	configuration, err := client.Configuration()
	if err != nil {
		return nil
	}

	_, bck, err := configuration.GetBackend(name, "")
	if err != nil {
		return nil
	}

	if bck.StickTable == nil {
		return nil
	}

	data := strings.Split(bck.StickTable.Store, ",")
	fields := make([]*models.StickTableField, 0)
	for _, d := range data {
		f := &models.StickTableField{}
		spl := strings.Split(d, "(")
		if len(spl) == 1 {
			f.Field = d
			f.Type = "counter"
			fields = append(fields, f)
		} else if len(spl) == 2 {
			p := misc.ParseTimeout(spl[1][:len(spl[1])-1])
			if p != nil {
				f.Field = spl[0]
				f.Period = *p
				f.Type = "rate"
				fields = append(fields, f)
			}
		}
	}

	return fields
}
