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

var cfg *Configuration

type HAProxyConfiguration struct {
	ConfigFile      string `short:"c" long:"config-file" description:"Path to the haproxy configuration file" default:"/etc/haproxy/haproxy.cfg"`
	Userlist        string `short:"u" long:"userlist" description:"Userlist in HAProxy configuration to use for API Basic Authentication" default:"controller"`
	HAProxy         string `short:"b" long:"haproxy-bin" description:"Path to the haproxy binary file" default:"haproxy"`
	ReloadDelay     int    `short:"d" long:"reload-delay" description:"Minimum delay between two reloads (in s)" default:"5"`
	ReloadCmd       string `short:"r" long:"reload-cmd" description:"Reload command"`
	RestartCmd      string `short:"s" long:"restart-cmd" description:"Restart command"`
	ReloadRetention int    `long:"reload-retention" description:"Reload retention in days, every older reload id will be deleted" default:"1"`
	TransactionDir  string `short:"t" long:"transaction-dir" description:"Path to the transaction directory" default:"/tmp/haproxy"`
	BackupsNumber   int    `short:"n" long:"backups-number" description:"Number of backup configuration files you want to keep, stored in the config dir with version number suffix" default:"0"`
	MasterRuntime   string `short:"m" long:"master-runtime" description:"Path to the master Runtime API socket"`
	ShowSystemInfo  bool   `short:"i" long:"show-system-info" description:"Show system info on info endpoint"`
	GitMode         bool   `short:"g" long:"git-mode" description:"Run dataplaneapi in git mode, without running the haproxy and ability to push to Git"`
	GitSettingsFile string `long:"git-settings-file" description:"Path to the git settings file" default:"/etc/haproxy/git.settings"`
}

type LoggingOptions struct {
	LogTo     string `long:"log-to" description:"Log target, can be stdout or file" default:"stdout" choice:"stdout" choice:"file"`
	LogFile   string `long:"log-file" description:"Location of the log file" default:"/var/log/dataplaneapi/dataplaneapi.log"`
	LogLevel  string `long:"log-level" description:"Logging level" default:"warning" choice:"trace" choice:"debug" choice:"info" choice:"warning" choice:"error"`
	LogFormat string `long:"log-format" description:"Logging format" default:"text" choice:"text" choice:"JSON"`
}

type Configuration struct {
	HAProxy HAProxyConfiguration `yaml:"haproxy"`
	Logging LoggingOptions       `yaml:"logging"`
}

//Get retuns pointer to configuration
func Get() *Configuration {

	if cfg == nil {
		cfg = &Configuration{}
	}
	return cfg
}
