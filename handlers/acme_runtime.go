// Copyright 2025 HAProxy Technologies
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
	"github.com/haproxytech/dataplaneapi/operations/acme_runtime"
)

type GetAcmeStatusHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *GetAcmeStatusHandlerImpl) Handle(params acme_runtime.GetAcmeStatusParams, principal interface{}) middleware.Responder {
	rt, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acme_runtime.NewGetAcmeStatusDefault(int(*e.Code)).WithPayload(e)
	}

	status, err := rt.AcmeStatus()
	if err != nil {
		e := misc.HandleError(err)
		return acme_runtime.NewGetAcmeStatusDefault(int(*e.Code)).WithPayload(e)
	}

	return acme_runtime.NewGetAcmeStatusOK().WithPayload(status)
}

type RenewAcmeCertificateHandlerImpl struct {
	Client client_native.HAProxyClient
}

func (h *RenewAcmeCertificateHandlerImpl) Handle(params acme_runtime.RenewAcmeCertificateParams, principal interface{}) middleware.Responder {
	rt, err := h.Client.Runtime()
	if err != nil {
		e := misc.HandleError(err)
		return acme_runtime.NewRenewAcmeCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	if err := rt.AcmeRenew(params.Certificate); err != nil {
		e := misc.HandleError(err)
		return acme_runtime.NewRenewAcmeCertificateDefault(int(*e.Code)).WithPayload(e)
	}

	return acme_runtime.NewRenewAcmeCertificateOK()
}
