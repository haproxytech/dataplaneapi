// Copyright 2022 HAProxy Technologies
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
	"errors"
	"fmt"
)

const (
	ReloadStratCustom  = "custom"
	ReloadStratS6      = "s6"
	ReloadStratSystemd = "systemd"

	// Default systemd commands.
	systemdReloadCmd      = "sudo -n systemctl reload"
	systemdRestartCmd     = "sudo -n systemctl restart"
	systemdStatusCmd      = "systemctl --quiet is-active"
	systemdDefaultService = "haproxy.service"

	// Default s6 commands.
	s6ReloadCmd      = "s6-svc -2"
	s6RestartCmd     = "s6-svc -r"
	s6StatusCmd      = "s6-svstat -u"
	s6DefaultService = "/var/run/service/haproxy"
)

// Set a default value to a string.
func setDefault(option *string, cmd, service string) {
	if *option == "" {
		*option = cmd + " " + service
	}
}

// Validate and initialize the way dataplaneapi will reload and monitor HAPRoxy.
func validateReloadConfiguration(c *HAProxyConfiguration) error {
	// By default, use the custom strategy.
	if c.ReloadStrategy == "" {
		c.ReloadStrategy = ReloadStratCustom
	}

	switch c.ReloadStrategy {
	case ReloadStratCustom:
		// The custom commands need to be set.
		if c.ReloadCmd == "" || c.RestartCmd == "" {
			return errors.New("the custom reload strategy requires these options to be set: " +
				"ReloadCmd, RestartCmd")
		}
	case ReloadStratS6:
		svc := s6DefaultService
		if c.ServiceName != "" {
			svc = c.ServiceName
		}
		setDefault(&c.ReloadCmd, s6ReloadCmd, svc)
		setDefault(&c.RestartCmd, s6RestartCmd, svc)
		setDefault(&c.StatusCmd, s6StatusCmd, svc)
	case ReloadStratSystemd:
		svc := systemdDefaultService
		if c.ServiceName != "" {
			svc = c.ServiceName
		}
		setDefault(&c.ReloadCmd, systemdReloadCmd, svc)
		setDefault(&c.RestartCmd, systemdRestartCmd, svc)
		setDefault(&c.StatusCmd, systemdStatusCmd, svc)
	default:
		return fmt.Errorf("invalid reload strategy: '%s'", c.ReloadStrategy)
	}

	return nil
}
