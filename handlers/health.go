package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/operations/health"
)

// GetHaproxyProcessInfoHandlerImpl implementation of the GetHaproxyProcessInfoHandler interface using client-native client
type GetHealthHandlerImpl struct {
	HAProxy haproxy.IReloadAgent
}

func (h *GetHealthHandlerImpl) Handle(health.GetHealthParams, interface{}) middleware.Responder {
	data := models.Health{}
	status, err := h.HAProxy.Status()
	if err == nil {
		if status {
			data.Haproxy = models.HealthHaproxyUp
		} else {
			data.Haproxy = models.HealthHaproxyDown
		}
	} else {
		data.Haproxy = models.HealthHaproxyUnknown
	}

	return health.NewGetHealthOK().WithPayload(&data)
}
