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
	"io/ioutil"

	"github.com/hashicorp/hcl"
	"github.com/rodaine/hclencoder"
)

type StorageHCL struct {
	filename string
	cfg      *StorageDataplaneAPIConfiguration
}

func (s *StorageHCL) Load(filename string) error {
	s.filename = filename
	cfg := &StorageDataplaneAPIConfiguration{}
	var hclFile []byte
	var err error
	hclFile, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = hcl.Decode(&cfg, string(hclFile))
	if err != nil {
		return err
	}
	s.cfg = cfg
	return nil
}

func (s *StorageHCL) Get() *StorageDataplaneAPIConfiguration {
	return s.cfg
}

func (s *StorageHCL) Set(cfg *StorageDataplaneAPIConfiguration) {
	s.cfg = cfg
}

func (s *StorageHCL) SaveAs(filename string) error {
	hcl, err := hclencoder.Encode(s.cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, hcl, 0644) //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageHCL) Save() error {
	return s.SaveAs(s.filename)
}
