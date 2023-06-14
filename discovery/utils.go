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

package discovery

import (
	"errors"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/dataplaneapi/misc"
)

const (
	minimumServerSlotsBase = 10
)

func ValidateAWSData(data *models.AwsRegion, useValidation bool) error {
	if useValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return validationErr
		}
	}
	if data.ID == nil || *data.ID == "" {
		return errors.New("missing ID")
	}
	if _, err := uuid.Parse(*data.ID); err != nil {
		return err
	}
	if data.ServerSlotsBase == nil || *data.ServerSlotsBase < minimumServerSlotsBase {
		data.ServerSlotsBase = misc.Int64P(10)
	}
	if data.ServerSlotsGrowthType == nil {
		data.ServerSlotsGrowthType = misc.StringP(models.AwsRegionServerSlotsGrowthTypeExponential)
	}
	if *data.ServerSlotsGrowthType == models.AwsRegionServerSlotsGrowthTypeLinear && (data.ServerSlotsGrowthIncrement == 0 || data.ServerSlotsGrowthIncrement < minimumServerSlotsBase) {
		data.ServerSlotsGrowthIncrement = minimumServerSlotsBase
	}
	return nil
}

func ValidateConsulData(data *models.Consul, useValidation bool) error {
	if useValidation {
		validationErr := data.Validate(strfmt.Default)
		if validationErr != nil {
			return validationErr
		}
	}
	if data.ID == nil || *data.ID == "" {
		return errors.New("missing ID")
	}
	if _, err := uuid.Parse(*data.ID); err != nil {
		return err
	}
	if data.ServerSlotsBase == nil || *data.ServerSlotsBase < minimumServerSlotsBase {
		data.ServerSlotsBase = misc.Int64P(minimumServerSlotsBase)
	}
	if data.ServerSlotsGrowthType == nil {
		data.ServerSlotsGrowthType = misc.StringP(models.ConsulServerSlotsGrowthTypeLinear)
	}
	if *data.ServerSlotsGrowthType == models.ConsulServerSlotsGrowthTypeLinear && (data.ServerSlotsGrowthIncrement == 0 || data.ServerSlotsGrowthIncrement < minimumServerSlotsBase) {
		data.ServerSlotsGrowthIncrement = minimumServerSlotsBase
	}
	return nil
}

func NewServiceDiscoveryUUID() *string {
	id := uuid.New().String()
	return &id
}
