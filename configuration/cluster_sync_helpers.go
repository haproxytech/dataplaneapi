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

	"github.com/haproxytech/dataplaneapi/log"
)

func (c *ClusterSync) getNodeFacts() map[string]string {
	facts := map[string]string{}

	// report the dataplane_cmdline if started from within haproxy
	if c.cfg.HAProxy.MasterWorkerMode || os.Getenv("HAPROXY_MWORKER") == "1" {
		facts["dataplane_cmdline"] = c.cfg.Cmdline.String()
	}

	runtime, err := c.cli.Runtime()
	if err != nil {
		log.Errorf("unable to fetch processInfo: %s", err.Error())
		return facts
	}
	processInfo, err := runtime.GetInfo()
	if err != nil {
		log.Error("unable to fetch processInfo")
	} else {
		if processInfo.Info != nil {
			facts["haproxy_version"] = processInfo.Info.Version
		} else {
			log.Error("empty process info")
		}
	}
	return facts
}
