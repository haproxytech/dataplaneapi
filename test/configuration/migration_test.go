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

package clustermode_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/renameio"
	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/storagetype"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var (
	fixtureDir  = "fixture"
	expectedDir = "expected"
)

func initConfiguration(dapiCfgPath, dapiStorageDir string) *configuration.Configuration {
	testCfg := &configuration.Configuration{
		HAProxy: configuration.HAProxyConfiguration{
			DataplaneStorageDir: dapiStorageDir,
			DataplaneConfig:     dapiCfgPath,
		},
	}
	return testCfg
}

func getDapiCfg(t *testing.T, dapiCfgPath string) *configuration.Configuration {
	t.Helper()
	dapiCfgJ, err := os.ReadFile(dapiCfgPath)
	require.NoError(t, err)
	var cfg configuration.Configuration
	yaml.Unmarshal(dapiCfgJ, &cfg)
	return &cfg
}

func TestConfiguration_LoadDataplaneStorageConfig_SD_Consul(t *testing.T) {
	tests := []struct {
		name            string
		fixtureDapiCfg  string
		fixtureConsul   string
		expectedDapiCfg string
		expectedConsul  string
	}{
		{
			name:            "from empty consul.json",
			fixtureDapiCfg:  "dataplaneapi-sd-consul-1.yaml",
			fixtureConsul:   "empty.json",
			expectedDapiCfg: "dataplaneapi-1.yaml",
			expectedConsul:  "consul-1.json",
		},
		{
			name:            "from non empty consul.json",
			fixtureDapiCfg:  "dataplaneapi-sd-consul-1.yaml",
			fixtureConsul:   "consul-1_2.json",
			expectedDapiCfg: "dataplaneapi-1.yaml",
			expectedConsul:  "consul-1_2.json",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dapiCfgPathSrc := path.Join(fixtureDir, tc.fixtureDapiCfg)
			dapiStorageConsulPathSrc := path.Join(fixtureDir, tc.fixtureConsul)

			dir := t.TempDir()
			dapiCfgPath := path.Join(dir, "dataplaneapi.yaml")
			dapiStoragePath := path.Join(dir, "service_discovery", "consul.json")
			os.MkdirAll(path.Join(dir, "service_discovery"), 0o755)

			// Copy dapi config file to tmp dir
			bytesRead, err := os.ReadFile(dapiCfgPathSrc)
			require.NoError(t, err)
			err = renameio.WriteFile(dapiCfgPath, bytesRead, 0o644)
			require.NoError(t, err)
			bytesRead, err = os.ReadFile(dapiStorageConsulPathSrc)
			require.NoError(t, err)
			err = renameio.WriteFile(dapiStoragePath, bytesRead, 0o644)
			require.NoError(t, err)

			// Load dapi config
			cfg := initConfiguration(dapiCfgPath, dir)
			cfg.Load()
			// override storage Dir
			cfg.HAProxy.DataplaneStorageDir = dir

			// Load and migrate
			cfg.LoadDataplaneStorageConfig()
			cfg.Save()

			// Check migrated consul.json
			// Check dataplaneapi.yaml (removed cluster...)dapiCfgMigrated, err := os.ReadFile(dapiCfgPath)
			var areEqual bool
			dapiCfg := getDapiCfg(t, dapiCfgPath)
			dapiExectedCfg := getDapiCfg(t, path.Join(expectedDir, tc.expectedDapiCfg))
			areEqual = reflect.DeepEqual(dapiExectedCfg, dapiCfg)
			require.True(t, areEqual, "migrated dataplaneapi.yaml is not equal to expected")

			dapiStorageMigrated := getConsul(t, dapiStoragePath)
			dapiExpectedStorage := getConsul(t, path.Join(expectedDir, tc.expectedConsul))
			areEqual = reflect.DeepEqual(dapiExpectedStorage, dapiStorageMigrated)
			diff := cmp.Diff(dapiExpectedStorage, dapiStorageMigrated)
			fmt.Println(diff)
			require.True(t, areEqual, fmt.Sprintf("test: %s migrated consul.json is not equal to expected", tc.name))
		})
	}
}

func TestConfiguration_LoadDataplaneStorageConfig_SD_AWS(t *testing.T) {
	tests := []struct {
		name            string
		fixtureDapiCfg  string
		fixtureAWS      string
		expectedDapiCfg string
		expectedAWS     string
	}{
		{
			name:            "from empty aws.json",
			fixtureDapiCfg:  "dataplaneapi-sd-aws-1.yaml",
			fixtureAWS:      "empty.json",
			expectedDapiCfg: "dataplaneapi-1.yaml",
			expectedAWS:     "aws-1.json",
		},
		{
			name:            "from non empty aws.json",
			fixtureDapiCfg:  "dataplaneapi-sd-aws-1.yaml",
			fixtureAWS:      "aws-1_2.json",
			expectedDapiCfg: "dataplaneapi-1.yaml",
			expectedAWS:     "aws-2.json",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dapiCfgPathSrc := path.Join(fixtureDir, tc.fixtureDapiCfg)
			dapiStorageConsulPathSrc := path.Join(fixtureDir, tc.fixtureAWS)

			dir := t.TempDir()
			dapiCfgPath := path.Join(dir, "dataplaneapi.yaml")
			dapiStoragePath := path.Join(dir, "service_discovery", "aws.json")
			os.MkdirAll(path.Join(dir, "service_discovery"), 0o755)

			// Copy dapi config file to tmp dir
			bytesRead, err := os.ReadFile(dapiCfgPathSrc)
			require.NoError(t, err)
			err = renameio.WriteFile(dapiCfgPath, bytesRead, 0o644)
			require.NoError(t, err)
			bytesRead, err = os.ReadFile(dapiStorageConsulPathSrc)
			require.NoError(t, err)
			err = renameio.WriteFile(dapiStoragePath, bytesRead, 0o644)
			require.NoError(t, err)

			// Load dapi config
			cfg := initConfiguration(dapiCfgPath, dir)
			cfg.Load()
			// override storage Dir
			cfg.HAProxy.DataplaneStorageDir = dir

			// Load and migrate
			cfg.LoadDataplaneStorageConfig()
			cfg.Save()

			// Check aws.json
			// Check dataplaneapi.yaml (removed cluster...)dapiCfgMigrated, err := os.ReadFile(dapiCfgPath)
			var areEqual bool
			dapiCfg := getDapiCfg(t, dapiCfgPath)
			dapiExectedCfg := getDapiCfg(t, path.Join(expectedDir, tc.expectedDapiCfg))
			areEqual = reflect.DeepEqual(dapiExectedCfg, dapiCfg)
			require.True(t, areEqual, "migrated dataplaneapi.yaml is not equal to expected")

			dapiStorageMigrated := getAWS(t, dapiStoragePath)
			dapiExpectedStorage := getAWS(t, path.Join(expectedDir, tc.expectedAWS))
			areEqual = reflect.DeepEqual(dapiExpectedStorage, dapiStorageMigrated)
			diff := cmp.Diff(dapiExpectedStorage, dapiStorageMigrated)
			fmt.Println(diff)
			require.True(t, areEqual, fmt.Sprintf("test: %s migrated aws.json is not equal to expected", tc.name))
		})
	}
}

func getConsul(t *testing.T, dapiStoragePath string) *storagetype.ConsulData {
	t.Helper()
	dapiStorageJ, err := os.ReadFile(dapiStoragePath)
	require.NoError(t, err)
	var data storagetype.ConsulData
	json.Unmarshal(dapiStorageJ, &data)
	return &data
}

func getAWS(t *testing.T, dapiStoragePath string) *storagetype.AWSRegionData {
	t.Helper()
	dapiStorageJ, err := os.ReadFile(dapiStoragePath)
	require.NoError(t, err)
	var data storagetype.AWSRegionData
	json.Unmarshal(dapiStorageJ, &data)
	return &data
}
