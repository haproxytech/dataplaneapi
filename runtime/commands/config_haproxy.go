// Copyright 2023 HAProxy Technologies
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

package commands

import (
	"errors"
	"strconv"

	client_native "github.com/haproxytech/client-native/v5"
)

const (
	ConfCmdKey = "conf"
)

type HAProxyConfiguration struct {
	Client client_native.HAProxyClient
}

func (g HAProxyConfiguration) Definition() definition {
	return definition{
		Key:  ConfCmdKey,
		Info: "show HAProxy configuration",
		Commands: []allCommands{
			{"conf", "show HAProxy current raw configuration"},
			{"conf raw version [transactionID]", "show HAProxy raw configuration for version (transactionID is optional default \"\" (for a transactionID, put '0' for version))"},
			{"conf structured", "show HAProxy current structured configuration"},
			{"conf structured version [transactionID]", "show HAProxy structured configuration for version (transactionID is optional default \"\" (for a transactionID, put '0' for version))"},
		},
	}
}

func (g HAProxyConfiguration) Command(cmd []string) (response []byte, err error) {
	configurationClient, err := g.Client.Configuration()
	if err != nil {
		return []byte{}, err
	}
	var version int64
	transactionID := ""
	configType := "raw"

	if len(cmd) >= 2 {
		configType = cmd[1]
	}

	if len(cmd) >= 3 {
		versionS := cmd[2]
		version, err = strconv.ParseInt(versionS, 10, 64)
		if err != nil {
			return []byte{}, err
		}
		if len(cmd) >= 4 {
			transactionID = cmd[3]
		}
	}

	var config string
	switch configType {
	case "raw":
		_, config, err = configurationClient.GetRawConfiguration(transactionID, version)
		if err != nil {
			return []byte{}, err
		}
	case "structured":
		return []byte{}, errors.New("structured not implemented in community edition")
	}

	return []byte(config), nil
}
