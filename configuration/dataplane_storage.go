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
)

// loadClusterModeData loads dataplane-managed users from the dataplane storage
// file (and not from the dapi config file).
func (c *Configuration) loadClusterModeData() error {
	clusterModeStorage, cmErr := storage.NewClusterModeStorage(path.Join(c.HAProxy.DataplaneStorageDir, storage.ClusterModeDataFileName))
	if cmErr != nil {
		return fmt.Errorf("dataplane storage error: %w", cmErr)
	}
	c.clusterModeStorage = clusterModeStorage

	return c.clusterModeStorage.Load()
}

func (c *Configuration) loadSDConsuls() error {
	// ClusterModeData
	consulStorage, cmErr := storage.NewSDConsulStorage(path.Join(c.HAProxy.DataplaneStorageDir, storage.ConsulFileName))
	if cmErr != nil {
		return fmt.Errorf("consul storage error: %w", cmErr)
	}

	if err := consulStorage.Load(); err != nil {
		return err
	}
	c.sdConsulStorage = consulStorage

	// Copy to configuration
	c.ServiceDiscovery.Consuls = c.sdConsulStorage.GetConsuls()

	return nil
}

func (c *Configuration) SaveSDConsuls() error {
	if err := c.sdConsulStorage.ReplaceAllConsulsAndStore(c.ServiceDiscovery.Consuls); err != nil {
		return err
	}
	return nil
}

func (c *Configuration) loadSDAWSRegions() error {
	// ClusterModeData
	awsRegionStorage, cmErr := storage.NewSDAWSRegionStorage(path.Join(c.HAProxy.DataplaneStorageDir, storage.AWSRegionFileName))
	if cmErr != nil {
		return fmt.Errorf("aws region storage error: %w", cmErr)
	}

	if err := awsRegionStorage.Load(); err != nil {
		return err
	}
	c.sdAWSRegionStorage = awsRegionStorage

	// Copy to configuration
	c.ServiceDiscovery.AWSRegions = c.sdAWSRegionStorage.GetAWSRegions()

	return nil
}

func (c *Configuration) SaveSDAWSRegions() error {
	if err := c.sdAWSRegionStorage.ReplaceAllAWSRegionsAndStore(c.ServiceDiscovery.AWSRegions); err != nil {
		return err
	}
	return nil
}
