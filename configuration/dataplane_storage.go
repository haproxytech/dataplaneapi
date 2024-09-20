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

	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/storage"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

// loadClusterModeData loads data from storage/dataplane (and not from dapi config file)
func (c *Configuration) loadClusterModeData() error {
	// ClusterModeData
	clusterModeStorage, cmErr := storage.NewClusterModeStorage(path.Join(c.HAProxy.DataplaneStorageDir, storage.ClusterModeDataFileName))
	if cmErr != nil {
		return fmt.Errorf("cluster mode storage error: %w", cmErr)
	}
	c.clusterModeStorage = clusterModeStorage

	if err := c.clusterModeStorage.Load(); err != nil {
		return err
	}
	if c.clusterModeStorage.GetStatus() != nil {
		c.Status.Store(*c.clusterModeStorage.GetStatus())
	}

	c.copyClusterToConfiguration(c.clusterModeStorage.GetCluster())

	return nil
}

func (c *Configuration) SaveClusterModeData() error {
	cfgCluster := c.Cluster
	dPort := cfgCluster.Port.Load()
	cfgCertificateFetched := cfgCluster.CertificateFetched.Load()

	dapiStorageCluster := storagetype.Cluster{
		APINodesPath:       misc.StringP(cfgCluster.APINodesPath.Load()),
		Token:              misc.StringP(cfgCluster.Token.Load()),
		ClusterTLSCertDir:  &c.HAProxy.ClusterTLSCertDir,
		ActiveBootstrapKey: misc.StringP(cfgCluster.ActiveBootstrapKey.Load()),
		APIRegisterPath:    misc.StringP(cfgCluster.APIRegisterPath.Load()),
		URL:                misc.StringP(cfgCluster.URL.Load()),
		Port:               &dPort,
		StorageDir:         misc.StringP(cfgCluster.StorageDir.Load()),
		BootstrapKey:       misc.StringP(cfgCluster.BootstrapKey.Load()),
		ID:                 misc.StringP(cfgCluster.ID.Load()),
		APIBasePath:        misc.StringP(cfgCluster.APIBasePath.Load()),
		CertificateDir:     misc.StringP(cfgCluster.CertificateDir.Load()),
		CertificateFetched: &cfgCertificateFetched,
		Name:               misc.StringP(cfgCluster.Name.Load()),
		Description:        misc.StringP(cfgCluster.Description.Load()),
		ClusterID:          misc.StringP(cfgCluster.ClusterID.Load()),
		ClusterLogTargets:  cfgCluster.ClusterLogTargets,
	}
	cfgStatus := c.Status.Load()
	c.clusterModeStorage.SetStatusAndStore(&cfgStatus)
	if err := c.clusterModeStorage.SetClusterAndStore(&dapiStorageCluster); err != nil {
		return err
	}

	return nil
}

func (c *Configuration) copyClusterToConfiguration(dapiStorageCluster *storagetype.Cluster) {
	if dapiStorageCluster == nil {
		return
	}
	if dapiStorageCluster.ClusterTLSCertDir != nil && !misc.HasOSArg("", "cluster-tls-dir", "") {
		c.HAProxy.ClusterTLSCertDir = *dapiStorageCluster.ClusterTLSCertDir
	}
	if dapiStorageCluster.ID != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.ID.Store(*dapiStorageCluster.ID)
	}
	if dapiStorageCluster.BootstrapKey != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.BootstrapKey.Store(*dapiStorageCluster.BootstrapKey)
	}
	if dapiStorageCluster.ActiveBootstrapKey != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.ActiveBootstrapKey.Store(*dapiStorageCluster.ActiveBootstrapKey)
	}
	if dapiStorageCluster.Token != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.Token.Store(*dapiStorageCluster.Token)
	}
	if dapiStorageCluster.URL != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.URL.Store(*dapiStorageCluster.URL)
	}
	if dapiStorageCluster.Port != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.Port.Store(*dapiStorageCluster.Port)
	}
	if dapiStorageCluster.APIBasePath != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.APIBasePath.Store(*dapiStorageCluster.APIBasePath)
	}
	if dapiStorageCluster.APINodesPath != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.APINodesPath.Store(*dapiStorageCluster.APINodesPath)
	}
	if dapiStorageCluster.APIRegisterPath != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.APIRegisterPath.Store(*dapiStorageCluster.APIRegisterPath)
	}
	if dapiStorageCluster.StorageDir != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.StorageDir.Store(*dapiStorageCluster.StorageDir)
	}
	if dapiStorageCluster.CertificateDir != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.CertificateDir.Store(*dapiStorageCluster.CertificateDir)
	}
	if dapiStorageCluster.CertificateFetched != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.CertificateFetched.Store(*dapiStorageCluster.CertificateFetched)
	}
	if dapiStorageCluster.Name != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.Name.Store(*dapiStorageCluster.Name)
	}
	if dapiStorageCluster.Description != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.Description.Store(*dapiStorageCluster.Description)
	}
	if dapiStorageCluster.ClusterID != nil && !misc.HasOSArg("", "", "") {
		c.Cluster.ClusterID.Store(*dapiStorageCluster.ClusterID)
	}
	if len(dapiStorageCluster.ClusterLogTargets) > 0 {
		c.Cluster.ClusterLogTargets = dapiStorageCluster.ClusterLogTargets
	}
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
