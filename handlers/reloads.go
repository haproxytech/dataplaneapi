package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/reloads"
	"github.com/haproxytech/models"
)

//GetReloadHandlerImpl implementation of the GetReloadHandler interface using client-native client
type GetReloadHandlerImpl struct {
	ReloadAgent *haproxy.ReloadAgent
}

//GetReloadsHandlerImpl implementation of the GetReloadsHandler interface using client-native client
type GetReloadsHandlerImpl struct {
	ReloadAgent *haproxy.ReloadAgent
}

//Handle executing the request and returning a response
func (rh *GetReloadHandlerImpl) Handle(params reloads.GetReloadParams, principal interface{}) middleware.Responder {
	r := rh.ReloadAgent.GetReload(params.ID)
	if r == nil {
		msg := fmt.Sprintf("Reload with ID %s does not exist", params.ID)
		c := misc.ErrHTTPNotFound
		e := &models.Error{
			Code:    &c,
			Message: &msg,
		}
		return reloads.NewGetReloadDefault(404).WithPayload(e)
	}
	return reloads.NewGetReloadOK().WithPayload(r)
}

//Handle executing the request and returning a response
func (rh *GetReloadsHandlerImpl) Handle(params reloads.GetReloadsParams, principal interface{}) middleware.Responder {
	rs := rh.ReloadAgent.GetReloads()
	return reloads.NewGetReloadsOK().WithPayload(rs)
}
