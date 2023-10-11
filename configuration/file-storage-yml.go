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
	"os"

	"github.com/google/renameio"
	"gopkg.in/yaml.v2"
)

type StorageYML struct {
	cfg      *StorageDataplaneAPIConfiguration
	filename string
}

func (s *StorageYML) Load(filename string) error {
	s.filename = filename
	cfg := &StorageDataplaneAPIConfiguration{}
	var err error

	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return err
	}

	s.cfg = cfg
	return nil
}

func (s *StorageYML) Get() *StorageDataplaneAPIConfiguration {
	if s.cfg == nil {
		s.cfg = &StorageDataplaneAPIConfiguration{}
	}
	return s.cfg
}

func (s *StorageYML) Set(cfg *StorageDataplaneAPIConfiguration) {
	s.cfg = cfg
}

func (s *StorageYML) SaveAs(filename string) error {
	data, err := yaml.Marshal(s.cfg)
	if err != nil {
		return err
	}

	return renameio.WriteFile(filename, data, 0o644)
}

func (s *StorageYML) Save() error {
	return s.SaveAs(s.filename)
}
