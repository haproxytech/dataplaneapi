// Copyright 2021 HAProxy Technologies
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

package storagetype

import (
	"github.com/haproxytech/client-native/v6/models"
)

type Cluster struct {
	APINodesPath       *string                    `json:"api_nodes_path,omitempty" yaml:"api_nodes_path,omitempty"`
	Token              *string                    `json:"token,omitempty" yaml:"token,omitempty"`
	ClusterTLSCertDir  *string                    `json:"cluster_tls_dir,omitempty" yaml:"cluster_tls_dir,omitempty"`
	ActiveBootstrapKey *string                    `json:"active_bootstrap_key,omitempty" yaml:"active_bootstrap_key,omitempty"`
	APIRegisterPath    *string                    `json:"api_register_path,omitempty" yaml:"api_register_path,omitempty"`
	URL                *string                    `json:"url,omitempty" yaml:"url,omitempty"`
	Port               *int                       `json:"port,omitempty" yaml:"port,omitempty"`
	StorageDir         *string                    `json:"storage_dir,omitempty" yaml:"storage_dir,omitempty"`
	BootstrapKey       *string                    `json:"bootstrap_key,omitempty" yaml:"bootstrap_key,omitempty"`
	ID                 *string                    `json:"id,omitempty" yaml:"id,omitempty"`
	APIBasePath        *string                    `json:"api_base_path,omitempty" yaml:"api_base_path,omitempty"`
	CertificateDir     *string                    `json:"cert_path,omitempty" yaml:"cert_path,omitempty"`
	CertificateFetched *bool                      `json:"cert_fetched,omitempty" yaml:"cert_fetched,omitempty"`
	Name               *string                    `json:"name,omitempty" yaml:"name,omitempty"`
	Description        *string                    `json:"description,omitempty" yaml:"description,omitempty"`
	ClusterID          *string                    `json:"cluster_id,omitempty" yaml:"cluster_id,omitempty" group:"cluster" save:"true"`
	ClusterLogTargets  []*models.ClusterLogTarget `json:"cluster_log_targets,omitempty" yaml:"cluster_log_targets,omitempty" group:"cluster" save:"true"`
}

func (c Cluster) LogDisplayName() string {
	if c.ClusterID != nil {
		return *c.ClusterID
	}
	if c.BootstrapKey != nil {
		return *c.BootstrapKey
	}
	return ""
}

func (c Cluster) IsEmpty() bool {
	res := isEmptyOrNilS(c.APINodesPath) &&
		isEmptyOrNilS(c.Token) &&
		isEmptyOrNilS(c.ClusterTLSCertDir) &&
		isEmptyOrNilS(c.ActiveBootstrapKey) &&
		isEmptyOrNilS(c.APIRegisterPath) &&
		isEmptyOrNilS(c.URL) &&
		isEmptyOrNilI(c.Port) &&
		isEmptyOrNilS(c.StorageDir) &&
		isEmptyOrNilS(c.BootstrapKey) &&
		isEmptyOrNilS(c.ID) &&
		isEmptyOrNilS(c.APIBasePath) &&
		isEmptyOrNilS(c.CertificateDir) &&
		isEmptyOrNilB(c.CertificateFetched) &&
		isEmptyOrNilS(c.Name) &&
		isEmptyOrNilS(c.Description)

	return res
}

func isEmptyOrNilS(s *string) bool {
	return (s == nil || *s == "")
}

func isEmptyOrNilI(i *int) bool {
	return (i == nil || *i == 0)
}

func isEmptyOrNilB(b *bool) bool {
	return (b == nil || !*b)
}
