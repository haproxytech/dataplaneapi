package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/controller/misc"
	"github.com/haproxytech/controller/operations/information"
	"github.com/haproxytech/models"
)

//GetInformationHandlerImpl implementation of the GetInformationHandler interface using client-native client
type GetInformationHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//Handle executing the request and returning a response
func (h *GetInformationHandlerImpl) Handle(params information.GetHaproxyProcessInfoParams, principal interface{}) middleware.Responder {
	info, err := h.Client.Runtime.GetInfo()
	if err != nil || len(info) == 0 {
		code := misc.ErrHTTPInternalServerError
		msg := err.Error()
		e := &models.Error{
			Code:    &code,
			Message: &msg,
		}
		return information.NewGetHaproxyProcessInfoDefault(int(misc.ErrHTTPInternalServerError)).WithPayload(e)
	}

	data := models.ProcessInfo{}
	data.Haproxy = &info[0]

	return information.NewGetHaproxyProcessInfoOK().WithPayload(&data)
}
