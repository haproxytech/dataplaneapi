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
	"strconv"

	"github.com/google/renameio"
	"github.com/haproxytech/dataplaneapi/log"
)

func HandlePIDFile(haproxyOptions HAProxyConfiguration) {
	if haproxyOptions.PIDFile == "" {
		return
	}

	if fileExists(haproxyOptions.PIDFile) {
		data, err := os.ReadFile(haproxyOptions.PIDFile)
		if err != nil {
			log.Fatalf("error while reading PID file content: %v", err)
		}
		pid, err := strconv.ParseInt(string(data), 10, 32)
		if err != nil {
			log.Fatalf("error while parsing PID file content: %v", err)
		}
		if os.Getpid() == int(pid) {
			log.Info("Stored PID matches current PID, proceeding...")
			return
		}
		if processExists(int(pid)) {
			log.Fatalf("process with PID %v already exists", pid)
		}
	}

	err := renameio.WriteFile(haproxyOptions.PIDFile, []byte(strconv.Itoa(os.Getpid())), 0o644)
	if err != nil {
		log.Fatalf("error while writing PID file: %s %s", haproxyOptions.PIDFile, err.Error())
	} else {
		log.Infof("PID %v stored in %v", os.Getpid(), haproxyOptions.PIDFile)
	}
}
