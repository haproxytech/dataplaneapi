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
	"strconv"

	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/operations/cluster"
	"github.com/haproxytech/dataplaneapi/operations/discovery"
	"github.com/haproxytech/models"
)

//CreateClusterHandlerImpl implementation of the CreateClusterHandler interface using client-native client
type CreateClusterHandlerImpl struct {
	Config *configuration.Configuration
}

//GetClusterHandlerImpl implementation of the GetClusterHandler interface using client-native client
type GetClusterHandlerImpl struct {
	Config *configuration.Configuration
}

//Handle executing the request and returning a response
func (h *CreateClusterHandlerImpl) Handle(params cluster.PostClusterParams, principal interface{}) middleware.Responder {
	key := h.Config.Cluster.BootstrapKey.Load()
	if key != params.Data.BootstrapKey || true {
		h.Config.Cluster.Mode.Store("cluster")
		h.Config.BotstrapKeyChanged(params.Data.BootstrapKey)
	}

	err := h.Config.Save()
	if err != nil {
		return cluster.NewPostClusterDefault(500)
	}
	return cluster.NewPostClusterOK().WithPayload(params.Data)
}

//Handle executing the request and returning a response
func (h *GetClusterHandlerImpl) Handle(params discovery.GetClusterParams, principal interface{}) middleware.Responder {

	portStr := h.Config.Cluster.Port.Load()
	p, err := strconv.Atoi(portStr)
	if err != nil {
		p = 0
	}
	port := int64(p)
	clusterSettings := &models.ClusterSettingsCluster{
		Address:     h.Config.Cluster.URL.Load(),
		Port:        &port,
		APIBasePath: h.Config.Cluster.APIBasePath.Load(),
		Name:        h.Config.Cluster.Name.Load(),
		Description: "",
	}
	settings := &models.ClusterSettings{
		BootstrapKey: h.Config.Cluster.BootstrapKey.Load(),
		Cluster:      clusterSettings,
		Mode:         h.Config.Cluster.Mode.Load(),
		Status:       h.Config.Cluster.Status.Load(),
	}
	return discovery.NewGetClusterOK().WithPayload(settings)
}
