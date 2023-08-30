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
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"strings"

	"github.com/haproxytech/dataplaneapi/log"
	"github.com/hashicorp/hcl"
	"github.com/rodaine/hclencoder"
)

type StorageHCL struct {
	cfg      *StorageDataplaneAPIConfiguration
	filename string
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
	if s.cfg == nil {
		s.cfg = &StorageDataplaneAPIConfiguration{}
	}
	return s.cfg
}

func (s *StorageHCL) Set(cfg *StorageDataplaneAPIConfiguration) {
	s.cfg = cfg
}

func (s *StorageHCL) SaveAs(filename string) error {
	var err error
	var hcl []byte
	var localCopy StorageDataplaneAPIConfiguration

	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err = e.Encode(s.cfg); err != nil {
		return err
	}
	d := gob.NewDecoder(&b)
	if err = d.Decode(&localCopy); err != nil {
		return err
	}

	// check if we have cluster log targets in config file
	if localCopy.Cluster != nil && len(localCopy.Cluster.ClusterLogTargets) > 0 {
		// since this can contain " character, escape it
		for index, value := range localCopy.Cluster.ClusterLogTargets {
			localCopy.Cluster.ClusterLogTargets[index].LogFormat = strings.ReplaceAll(value.LogFormat, `"`, `\"`)
		}
	}
	if localCopy.LogTargets != nil && len(*localCopy.LogTargets) > 0 {
		var logTargets []log.Target
		for _, value := range *localCopy.LogTargets {
			value.ACLFormat = strings.ReplaceAll(value.ACLFormat, `"`, `\"`)
			logTargets = append(logTargets, value)
		}
		localCopy.LogTargets = (*log.Targets)(&logTargets)
	}
	if localCopy.Log != nil {
		if localCopy.Log.ACLFormat != nil {
			aclF := strings.ReplaceAll(*localCopy.Log.ACLFormat, `"`, `\"`)
			localCopy.Log.ACLFormat = &aclF
		}
	}

	hcl, err = hclencoder.Encode(localCopy)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, hcl, 0o644) //nolint:gosec
}

func (s *StorageHCL) Save() error {
	return s.SaveAs(s.filename)
}
