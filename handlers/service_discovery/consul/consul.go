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

package consul

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/haproxytech/client-native/v6/models"

	sc "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all Consul service discovery routes onto r.
func RegisterRouter(r chi.Router, discovery sc.ServiceDiscoveries, persistCallback func([]*models.Consul) error, useValidation bool) error {
	spec, err := GetSpec()
	if err != nil {
		return err
	}
	HandlerWithOptions(&HandlerImpl{
		Discovery:       discovery,
		PersistCallback: persistCallback,
		UseValidation:   useValidation,
	}, ChiServerOptions{
		BaseRouter:       r,
		Middlewares:      []MiddlewareFunc{middleware.NewValidator(spec)},
		ErrorHandlerFunc: middleware.ErrorHandler,
	})
	return nil
}

// HandlerImpl implements ServerInterface for Consul service discovery.
type HandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.Consul) error
	UseValidation   bool
}

func (h *HandlerImpl) GetConsuls(w http.ResponseWriter, r *http.Request) {
	consuls, err := getConsuls(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, consuls)
}

func (h *HandlerImpl) CreateConsul(w http.ResponseWriter, r *http.Request) {
	var data Consul
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	data.ID = sc.NewServiceDiscoveryUUID()

	if err := sc.ValidateConsulData(&data, h.UseValidation); err != nil {
		respond.Error(w, err)
		return
	}

	if data.HealthCheckPolicy != nil && *data.HealthCheckPolicy == models.ConsulHealthCheckPolicyMin && data.HealthCheckPolicyMin <= 0 {
		respond.BadRequest(w, "health_check_policy_min is required for 'min' health_check_policy")
		return
	}

	if err := h.Discovery.AddNode("consul", *data.ID, &data); err != nil {
		respond.Error(w, err)
		return
	}

	consuls, err := getConsuls(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = h.PersistCallback(consuls); err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, &data)
}

func (h *HandlerImpl) DeleteConsul(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Discovery.RemoveNode("consul", id); err != nil {
		respond.Error(w, err)
		return
	}

	consuls, err := getConsuls(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = h.PersistCallback(consuls); err != nil {
		respond.Error(w, err)
		return
	}

	respond.NoContent(w)
}

func (h *HandlerImpl) GetConsul(w http.ResponseWriter, r *http.Request, id string) {
	nodes, err := h.Discovery.GetNode("consul", id)
	if err != nil {
		respond.Error(w, err)
		return
	}

	consul, ok := nodes.(*models.Consul)
	if !ok {
		respond.Error(w, errors.New("expected *models.Consul"))
		return
	}

	respond.JSON(w, http.StatusOK, consul)
}

func (h *HandlerImpl) ReplaceConsul(w http.ResponseWriter, r *http.Request, id string) {
	var data Consul
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	if err := sc.ValidateConsulData(&data, h.UseValidation); err != nil {
		respond.Error(w, err)
		return
	}

	if data.HealthCheckPolicy != nil && *data.HealthCheckPolicy == models.ConsulHealthCheckPolicyMin && data.HealthCheckPolicyMin <= 0 {
		respond.BadRequest(w, "health_check_policy_min is required for 'min' health_check_policy")
		return
	}

	if err := h.Discovery.UpdateNode("consul", *data.ID, &data); err != nil {
		respond.Error(w, err)
		return
	}

	consuls, err := getConsuls(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = h.PersistCallback(consuls); err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, &data)
}

func getConsuls(discovery sc.ServiceDiscoveries) (models.Consuls, error) {
	nodes, err := discovery.GetNodes("consul")
	if err != nil {
		return nil, err
	}
	consuls, ok := nodes.(models.Consuls)
	if !ok {
		return nil, errors.New("expected models.Consuls")
	}
	return consuls, nil
}
