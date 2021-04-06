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
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/haproxytech/client-native/v2/models"
	"github.com/haproxytech/dataplaneapi/misc"

	sc "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/operations/service_discovery"
)

// CreateConsulHandlerImpl implementation of the CreateConsulHandler interface using client-native client
type CreateAWSHandlerImpl struct {
	Discovery       sc.ServiceDiscoveries
	UseValidation   bool
	PersistCallback func([]*models.AwsRegion) error
}

func (c CreateAWSHandlerImpl) Handle(params service_discovery.CreateAWSRegionParams, i interface{}) middleware.Responder {
	var err error
	handleError := func(err error) *service_discovery.CreateAWSRegionDefault {
		e := misc.HandleError(err)
		return service_discovery.NewCreateAWSRegionDefault(int(*e.Code)).WithPayload(e)
	}

	id := uuid.New().String()
	params.Data.ID = &id
	if err = validateAWSData(params.Data, c.UseValidation); err != nil {
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

func validateAWSData(data *models.AwsRegion, useValidation bool) error {
	if useValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return validationErr
		}
	}
	if data.ID == nil || *data.ID == "" {
		return errors.New("missing ID")
	}
	if data.ServerSlotsBase == nil || *data.ServerSlotsBase < 10 {
		data.ServerSlotsBase = misc.Int64P(10)
	}
	if data.ServerSlotsGrowthType == nil {
		data.ServerSlotsGrowthType = misc.StringP("linear")
	}
	if *data.ServerSlotsGrowthType == "linear" && (data.ServerSlotsGrowthIncrement == 0 || data.ServerSlotsGrowthIncrement < 10) {
		data.ServerSlotsGrowthIncrement = 10
	}
	return nil
}
