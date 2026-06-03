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

package aws

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/haproxytech/client-native/v6/models"

	sc "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/handlers/middleware"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
)

// RegisterRouter registers all AWS service discovery routes onto r.
func RegisterRouter(r chi.Router, discovery sc.ServiceDiscoveries, persistCallback func([]*models.AwsRegion) error, useValidation bool) error {
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

// HandlerImpl implements ServerInterface for AWS service discovery.
type HandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.AwsRegion) error
	UseValidation   bool
}

func (h *HandlerImpl) GetAWSRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := getAWSRegions(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, regions)
}

func (h *HandlerImpl) CreateAWSRegion(w http.ResponseWriter, r *http.Request) {
	var data AwsRegion
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	data.ID = sc.NewServiceDiscoveryUUID()

	if err := sc.ValidateAWSData(&data, h.UseValidation); err != nil {
		respond.Error(w, err)
		return
	}

	if err := h.Discovery.AddNode("aws", *data.ID, &data); err != nil {
		respond.Error(w, err)
		return
	}

	regions, err := getAWSRegions(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err = h.PersistCallback(regions); err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusCreated, &data)
}

func (h *HandlerImpl) DeleteAWSRegion(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.Discovery.RemoveNode("aws", id); err != nil {
		respond.Error(w, err)
		return
	}

	regions, err := getAWSRegions(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err := h.PersistCallback(regions); err != nil {
		respond.Error(w, err)
		return
	}

	respond.NoContent(w)
}

func (h *HandlerImpl) GetAWSRegion(w http.ResponseWriter, r *http.Request, id string) {
	nodes, err := h.Discovery.GetNode("aws", id)
	if err != nil {
		respond.Error(w, err)
		return
	}

	region, ok := nodes.(*models.AwsRegion)
	if !ok {
		respond.Error(w, errors.New("expected *models.AwsRegion"))
		return
	}

	respond.JSON(w, http.StatusOK, region)
}

func (h *HandlerImpl) ReplaceAWSRegion(w http.ResponseWriter, r *http.Request, id string) {
	var data AwsRegion
	if !respond.DecodeBody(r, w, &data) {
		return
	}

	if err := sc.ValidateAWSData(&data, h.UseValidation); err != nil {
		respond.Error(w, err)
		return
	}

	if err := h.Discovery.UpdateNode("aws", *data.ID, &data); err != nil {
		respond.Error(w, err)
		return
	}

	regions, err := getAWSRegions(h.Discovery)
	if err != nil {
		respond.Error(w, err)
		return
	}

	if err := h.PersistCallback(regions); err != nil {
		respond.Error(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, &data)
}

func getAWSRegions(discovery sc.ServiceDiscoveries) (models.AwsRegions, error) {
	nodes, err := discovery.GetNodes("aws")
	if err != nil {
		return nil, err
	}
	regions, ok := nodes.(models.AwsRegions)
	if !ok {
		return nil, errors.New("expected models.AwsRegions")
	}
	return regions, nil
}
