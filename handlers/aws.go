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
	"errors"

	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/misc"

	sc "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/operations/service_discovery"
)

// CreateConsulHandlerImpl implementation of the CreateConsulHandler interface using client-native client
type CreateAWSHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.AwsRegion) error
	UseValidation   bool
}

func (c CreateAWSHandlerImpl) Handle(params service_discovery.CreateAWSRegionParams, i interface{}) middleware.Responder {
	var err error
	handleError := func(err error) *service_discovery.CreateAWSRegionDefault {
		e := misc.HandleError(err)
		return service_discovery.NewCreateAWSRegionDefault(int(*e.Code)).WithPayload(e)
	}

	params.Data.ID = sc.NewServiceDiscoveryUUID()
	if err = sc.ValidateAWSData(params.Data, c.UseValidation); err != nil {
		return handleError(err)
	}
	if err = c.Discovery.AddNode("aws", *params.Data.ID, params.Data); err != nil {
		return handleError(err)
	}
	var regions models.AwsRegions
	if regions, err = getAWSRegions(c.Discovery); err != nil {
		return handleError(err)
	}
	if err = c.PersistCallback(regions); err != nil {
		return handleError(err)
	}
	return service_discovery.NewCreateAWSRegionCreated().WithPayload(params.Data)
}

func getAWSRegions(discovery sc.ServiceDiscoveries) (models.AwsRegions, error) {
	nodes, err := discovery.GetNodes("aws")
	if err != nil {
		return nil, err
	}
	instances, ok := nodes.(models.AwsRegions)
	if !ok {
		return nil, errors.New("expected *models.AwsRegion")
	}
	return instances, nil
}

type GetAWSRegionHandlerImpl struct {
	Discovery sc.ServiceDiscoveries
}

func (g GetAWSRegionHandlerImpl) Handle(params service_discovery.GetAWSRegionParams, i interface{}) middleware.Responder {
	handleError := func(err error) *service_discovery.GetAWSRegionDefault {
		e := misc.HandleError(err)
		return service_discovery.NewGetAWSRegionDefault(int(*e.Code)).WithPayload(e)
	}

	nodes, err := g.Discovery.GetNode("aws", params.ID)
	if err != nil {
		return handleError(err)
	}
	region, ok := nodes.(*models.AwsRegion)
	if !ok {
		return handleError(err)
	}
	return service_discovery.NewGetAWSRegionOK().WithPayload(region)
}

type GetAWSRegionsHandlerImpl struct {
	Discovery sc.ServiceDiscoveries
}

func (g GetAWSRegionsHandlerImpl) Handle(params service_discovery.GetAWSRegionsParams, i interface{}) middleware.Responder {
	regions, err := getAWSRegions(g.Discovery)
	if err != nil {
		e := misc.HandleError(err)
		return service_discovery.NewGetAWSRegionsDefault(int(*e.Code)).WithPayload(e)
	}
	return service_discovery.NewGetAWSRegionsOK().WithPayload(regions)
}

type ReplaceAWSRegionHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.AwsRegion) error
	UseValidation   bool
}

func (r ReplaceAWSRegionHandlerImpl) Handle(params service_discovery.ReplaceAWSRegionParams, i interface{}) middleware.Responder {
	handleError := func(err error) *service_discovery.ReplaceAWSRegionDefault {
		e := misc.HandleError(err)
		return service_discovery.NewReplaceAWSRegionDefault(int(*e.Code)).WithPayload(e)
	}
	var err error
	if err = sc.ValidateAWSData(params.Data, r.UseValidation); err != nil {
		return handleError(err)
	}
	if err = r.Discovery.UpdateNode("aws", *params.Data.ID, params.Data); err != nil {
		return handleError(err)
	}
	var regions models.AwsRegions
	regions, err = getAWSRegions(r.Discovery)
	if err != nil {
		return handleError(err)
	}
	if err = r.PersistCallback(regions); err != nil {
		return handleError(err)
	}
	return service_discovery.NewReplaceAWSRegionOK().WithPayload(params.Data)
}

type DeleteAWSRegionHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	PersistCallback func([]*models.AwsRegion) error
}

func (d DeleteAWSRegionHandlerImpl) Handle(params service_discovery.DeleteAWSRegionParams, i interface{}) middleware.Responder {
	handleError := func(err error) *service_discovery.DeleteAWSRegionDefault {
		e := misc.HandleError(err)
		return service_discovery.NewDeleteAWSRegionDefault(int(*e.Code)).WithPayload(e)
	}
	var err error
	if err = d.Discovery.RemoveNode("aws", params.ID); err != nil {
		return handleError(err)
	}
	var regions models.AwsRegions
	if regions, err = getAWSRegions(d.Discovery); err != nil {
		return handleError(err)
	}
	if err = d.PersistCallback(regions); err != nil {
		return handleError(err)
	}
	return service_discovery.NewDeleteAWSRegionNoContent()
}
