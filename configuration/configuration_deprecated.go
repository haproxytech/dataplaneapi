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

package configuration

import (
	"fmt"
	"path"

	"github.com/haproxytech/dataplaneapi/storage"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

func (c *Configuration) migrateDeprecatedFields() ([]string, error) {
	deprecationInfoMsg := make([]string, 0)

	// Service Discovery Consuls
	deprMsgConsuls, err := c.migrateSDConsuls()
	deprecationInfoMsg = append(deprecationInfoMsg, deprMsgConsuls...)
	if err != nil {
		return deprecationInfoMsg, err
	}

	// Service Discovery AWS Regions
	deprMsgAwsRegions, err := c.migrateSDAwsRegions()
	deprecationInfoMsg = append(deprecationInfoMsg, deprMsgAwsRegions...)
	if err != nil {
		return deprecationInfoMsg, err
	}

	// Remove deprecated fields after migration
	c.storage.Get().emptyDeprecatedSections()

	return deprecationInfoMsg, nil
}

func (c *Configuration) migrateSDConsuls() ([]string, error) {
	deprecationInfoMsg := make([]string, 0)

	dapiCfgStorage := c.storage.Get()
	dapiStorageConsuls := c.sdConsulStorage.GetConsuls()
	consulsStoragePath := path.Join(c.HAProxy.DataplaneStorageDir, storage.ConsulFileName)

	consulsToMigrate := make(storagetype.Consuls, 0)
	if dapiCfgStorage.DeprecatedServiceDiscovery == nil || dapiCfgStorage.DeprecatedServiceDiscovery.Consuls == nil {
		return deprecationInfoMsg, nil
	}
	for _, cfgConsul := range *dapiCfgStorage.DeprecatedServiceDiscovery.Consuls {
		found := false
		// Check done on ID
		for _, dapiStorageConsul := range dapiStorageConsuls {
			if dapiStorageConsul.ID != nil && cfgConsul.ID != nil &&
				*dapiStorageConsul.ID == *cfgConsul.ID {
				found = true
				break
			}
		}

		consulID := ""
		if cfgConsul.ID != nil {
			consulID = *cfgConsul.ID
		}
		// Already migrated
		if found {
			msg := fmt.Sprintf("[CFG DEPRECATED] [SKIP] [Consul] [%s]: already migrated. Old location [%s] New location [%s]. Use only new location",
				consulID,
				c.HAProxy.DataplaneConfig,
				consulsStoragePath)
			// Logging is not done here as at startup, the logger is not yet initialized
			// so it's done later on
			deprecationInfoMsg = append(deprecationInfoMsg, msg)
		} else {
			// If not already migrated, then migrate it
			msg := fmt.Sprintf("[CFG DEPRECATED] [MIGRATE] [Consul] [%s]: migrating. Old location [%s] New location [%s]. Use only new location",
				consulID,
				c.HAProxy.DataplaneConfig,
				consulsStoragePath)
			// Logging is not done here as at startup, the logger is not yet initialized
			// so it's done later on
			deprecationInfoMsg = append(deprecationInfoMsg, msg)
			consulsToMigrate = append(consulsToMigrate, cfgConsul)
		}
	}
	if err := c.sdConsulStorage.AddConsulsAndStore(consulsToMigrate); err != nil {
		return deprecationInfoMsg, err
	}
	c.ServiceDiscovery.Consuls = append(c.ServiceDiscovery.Consuls, consulsToMigrate...)
	return deprecationInfoMsg, nil
}

func (c *Configuration) migrateSDAwsRegions() ([]string, error) {
	deprecationInfoMsg := make([]string, 0)

	dapiCfgStorage := c.storage.Get()
	dapiStorageAwsRegions := c.sdAWSRegionStorage.GetAWSRegions()
	awsRegionsStoragePath := path.Join(c.HAProxy.DataplaneStorageDir, storage.AWSRegionFileName)

	awsRegionssToMigrate := make(storagetype.AWSRegions, 0)
	if dapiCfgStorage.DeprecatedServiceDiscovery == nil || dapiCfgStorage.DeprecatedServiceDiscovery.AWSRegions == nil {
		return deprecationInfoMsg, nil
	}
	for _, cfgAwsRegion := range *dapiCfgStorage.DeprecatedServiceDiscovery.AWSRegions {
		found := false
		for _, dapiStorageAwsRegion := range dapiStorageAwsRegions {
			if dapiStorageAwsRegion.Name != nil && cfgAwsRegion.Name != nil &&
				*dapiStorageAwsRegion.Name == *cfgAwsRegion.Name {
				found = true
				break
			}
		}

		awsRegionName := ""
		if cfgAwsRegion.Name != nil {
			awsRegionName = *cfgAwsRegion.Name
		}
		// Already migrated
		if found {
			msg := fmt.Sprintf("[CFG DEPRECATED] [SKIP] [AWS Region] [%s]: already migrated. Old location [%s] New location [%s]. Use only new location",
				awsRegionName,
				c.HAProxy.DataplaneConfig,
				awsRegionsStoragePath)
			// Logging is not done here as at startup, the logger is not yet initialized
			// so it's done later on
			deprecationInfoMsg = append(deprecationInfoMsg, msg)
		} else {
			// If not already migrated, then migrate it
			msg := fmt.Sprintf("[CFG DEPRECATED] [MIGRATE] [AWS Region] [%s]: migrating. Old location [%s] New location [%s]. Use only new location",
				awsRegionName,
				c.HAProxy.DataplaneConfig,
				awsRegionsStoragePath)
			// Logging is not done here as at startup, the logger is not yet initialized
			// so it's done later on
			deprecationInfoMsg = append(deprecationInfoMsg, msg)
			awsRegionssToMigrate = append(awsRegionssToMigrate, cfgAwsRegion)
		}
	}
	if err := c.sdAWSRegionStorage.AddAWSRegionsAndStore(awsRegionssToMigrate); err != nil {
		return deprecationInfoMsg, err
	}
	c.ServiceDiscovery.AWSRegions = append(c.ServiceDiscovery.AWSRegions, awsRegionssToMigrate...)
	return deprecationInfoMsg, nil
}
