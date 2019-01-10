package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations/stats"
	"github.com/haproxytech/models"
)

//GetStatsHandlerImpl implementation of the GetStatsHandler interface using client-native client
type GetStatsHandlerImpl struct {
	Client *client_native.HAProxyClient
}

//Handle executing the request and returning a response
func (h *GetStatsHandlerImpl) Handle(params stats.GetStatsParams, principal interface{}) middleware.Responder {
	s, err := h.Client.Runtime.GetStats()
	if err != nil {
		code := misc.ErrHTTPInternalServerError
		msg := err.Error()
		e := &models.Error{
			Code:    &code,
			Message: &msg,
		}
		return stats.NewGetStatsDefault(int(misc.ErrHTTPInternalServerError)).WithPayload(e)
	}

	nativeStats := models.NativeStats{}
	for i, nStat := range s {
		for _, item := range nStat {
			nativeStatItem := *item
			nativeStatItem.Process = int64(i)
			nativeStats = append(nativeStats, &nativeStatItem)
		}
	}

	return stats.NewGetStatsOK().WithPayload(nativeStats)
}
