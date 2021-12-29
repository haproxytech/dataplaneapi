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

func (c *ClusterSync) getNodeVariables() map[string]string {
	variables := map[string]string{}

	// report the dataplane_cmdline if started from within haproxy
	if c.cfg.HAProxy.MasterWorkerMode || os.Getenv("HAPROXY_MWORKER") == "1" {
		variables["dataplane_cmdline"] = c.cfg.Cmdline.String()
	}

	processInfos, err := c.cli.Runtime.GetInfo()
	if err != nil || len(processInfos) < 1 {
		log.Error("unable to fetch processInfo")
	} else {
		if processInfos[0].Info != nil {
			variables["haproxy_version"] = processInfos[0].Info.Version
		} else {
			log.Error("empty process info")
		}
	}
	return variables
}
