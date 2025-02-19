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

package configuration

import (
	"time"

	"github.com/jessevdk/go-flags"

	"github.com/haproxytech/client-native/v6/models"
	dpapilog "github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

type configTypeDataplaneapi struct {
	WriteTimeout     *string                `yaml:"write_timeout,omitempty"`
	GracefulTimeout  *string                `yaml:"graceful_timeout,omitempty"`
	ShowSystemInfo   *bool                  `yaml:"show_system_info,omitempty"`
	MaxHeaderSize    *string                `yaml:"max_header_size,omitempty"`
	SocketPath       *flags.Filename        `yaml:"socket_path,omitempty"`
	DebugSocketPath  *string                `yaml:"debug_socket_path,omitempty"`
	Host             *string                `yaml:"host,omitempty"`
	Port             *int                   `yaml:"port,omitempty"`
	ListenLimit      *int                   `yaml:"listen_limit,omitempty"`
	DisableInotify   *bool                  `yaml:"disable_inotify,omitempty"`
	ReadTimeout      *string                `yaml:"read_timeout,omitempty"`
	Advertised       *configTypeAdvertised  `yaml:"advertised,omitempty"`
	CleanupTimeout   *string                `yaml:"cleanup_timeout,omitempty"`
	KeepAlive        *string                `yaml:"keep_alive,omitempty"`
	PIDFile          *string                `yaml:"pid_file,omitempty"`
	UID              *int                   `yaml:"uid,omitempty"`
	GID              *int                   `yaml:"gid,omitempty"`
	TLS              *configTypeTLS         `yaml:"tls,omitempty"`
	EnabledListeners *[]string              `yaml:"scheme,omitempty"`
	Userlist         *configTypeUserlist    `yaml:"userlist,omitempty"`
	Transaction      *configTypeTransaction `yaml:"transaction,omitempty"`
	Resources        *configTypeResources   `yaml:"resources,omitempty"`
	Users            []storagetype.User     `yaml:"user,omitempty"`
}

type configTypeTLS struct {
	TLSHost           *string         `yaml:"tls_host,omitempty"`
	TLSPort           *int            `yaml:"tls_port,omitempty"`
	TLSCertificate    *flags.Filename `yaml:"tls_certificate,omitempty"`
	TLSCertificateKey *flags.Filename `yaml:"tls_key,omitempty"`
	TLSCACertificate  *flags.Filename `yaml:"tls_ca,omitempty"`
	TLSListenLimit    *int            `yaml:"tls_listen_limit,omitempty"`
	TLSKeepAlive      *string         `yaml:"tls_keep_alive,omitempty"`
	TLSReadTimeout    *string         `yaml:"tls_read_timeout,omitempty"`
	TLSWriteTimeout   *string         `yaml:"tls_write_timeout,omitempty"`
}

type configTypeHaproxy struct {
	ConfigFile       *string           `yaml:"config_file,omitempty"`
	HAProxy          *string           `yaml:"haproxy_bin,omitempty"`
	MasterRuntime    *string           `yaml:"master_runtime,omitempty"`
	NodeIDFile       *string           `yaml:"fid,omitempty"`
	MasterWorkerMode *bool             `yaml:"master_worker_mode,omitempty"`
	Reload           *configTypeReload `yaml:"reload,omitempty"`
	DelayedStartMax  *string           `yaml:"delayed_start_max,omitempty"`
	DelayedStartTick *string           `yaml:"delayed_start_tick,omitempty"`
}

type configTypeUserlist struct {
	Userlist     *string `yaml:"userlist,omitempty"`
	UserListFile *string `yaml:"userlist_file,omitempty"`
}

type configTypeReload struct {
	ReloadDelay     *int    `yaml:"reload_delay,omitempty"`
	ReloadCmd       *string `yaml:"reload_cmd,omitempty"`
	RestartCmd      *string `yaml:"restart_cmd,omitempty"`
	StatusCmd       *string `yaml:"status_cmd,omitempty"`
	ServiceName     *string `yaml:"service_name,omitempty"`
	ReloadRetention *int    `yaml:"reload_retention,omitempty"`
	ReloadStrategy  *string `yaml:"reload_strategy,omitempty"`
	ValidateCmd     *string `yaml:"validate_cmd,omitempty"`
}

type configTypeTransaction struct {
	TransactionDir      *string `yaml:"transaction_dir,omitempty"`
	BackupsNumber       *int    `yaml:"backups_number,omitempty"`
	BackupsDir          *string `yaml:"backups_dir,omitempty"`
	MaxOpenTransactions *int64  `yaml:"max_open_transactions,omitempty"`
}

type configTypeResources struct {
	MapsDir              *string `yaml:"maps_dir,omitempty"`
	SSLCertsDir          *string `yaml:"ssl_certs_dir,omitempty"`
	GeneralStorageDir    *string `yaml:"general_storage_dir,omitempty"`
	DataplaneStorageDir  *string `yaml:"dataplane_storage_dir,omitempty"`
	UpdateMapFiles       *bool   `yaml:"update_map_files,omitempty"`
	UpdateMapFilesPeriod *int64  `yaml:"update_map_files_period,omitempty"`
	SpoeDir              *string `yaml:"spoe_dir,omitempty"`
	SpoeTransactionDir   *string `yaml:"spoe_transaction_dir,omitempty"`
}

type configTypeAdvertised struct {
	APIAddress *string `yaml:"api_address,omitempty"`
	APIPort    *int64  `yaml:"api_port,omitempty"`
}

type storagetypeerviceDiscovery struct {
	Consuls    *[]*models.Consul    `yaml:"consuls,omitempty"`
	AWSRegions *[]*models.AwsRegion `yaml:"aws_regions,omitempty"`
}

type configKeepalived struct {
	ConfigFile     *string `yaml:"config_file"`
	StartCmd       *string `yaml:"start_cmd"`
	ReloadCmd      *string `yaml:"reload_cmd"`
	RestartCmd     *string `yaml:"restart_cmd"`
	StopCmd        *string `yaml:"stop_cmd"`
	StatusCmd      *string `yaml:"status_cmd,omitempty"`
	ValidateCmd    *string `yaml:"validate_cmd"`
	ReloadDelay    *int    `yaml:"reload_delay"`
	TransactionDir *string `yaml:"transaction_dir"`
	BackupsDir     *string `yaml:"backups_dir"`
	BackupsNumber  *int    `yaml:"backups_number"`
}

type storagetypeyslog struct {
	SyslogAddr     *string `yaml:"syslog_address,omitempty"`
	SyslogProto    *string `yaml:"syslog_protocol,omitempty"`
	SyslogTag      *string `yaml:"syslog_tag,omitempty"`
	SyslogLevel    *string `yaml:"syslog_level,omitempty"`
	SyslogFacility *string `yaml:"syslog_facility,omitempty"`
}

type configTypeLog struct {
	LogTo     *string           `yaml:"log_to,omitempty"`
	LogFile   *string           `yaml:"log_file,omitempty"`
	LogLevel  *string           `yaml:"log_level,omitempty"`
	LogFormat *string           `yaml:"log_format,omitempty"`
	ACLFormat *string           `yaml:"apache_common_log_format,omitempty"`
	Syslog    *storagetypeyslog `yaml:"syslog,omitempty"`
}

type StorageDataplaneAPIConfiguration struct {
	Version                    *int                        `yaml:"config_version,omitempty"`
	Name                       *string                     `yaml:"name,omitempty"`
	DeprecatedMode             *string                     `yaml:"mode,omitempty"`
	DeprecatedBootstrapKey     *string                     `yaml:"bootstrap_key,omitempty"`
	DeprecatedStatus           *string                     `yaml:"status,omitempty"`
	Dataplaneapi               *configTypeDataplaneapi     `yaml:"dataplaneapi,omitempty"`
	Haproxy                    *configTypeHaproxy          `yaml:"haproxy,omitempty"`
	DeprecatedCluster          *storagetype.Cluster        `yaml:"cluster,omitempty"`
	DeprecatedServiceDiscovery *storagetypeerviceDiscovery `yaml:"service_discovery,omitempty"`
	Log                        *configTypeLog              `yaml:"log,omitempty"`
	LogTargets                 *dpapilog.Targets           `yaml:"log_targets,omitempty"`
}

func copyToConfiguration(cfg *Configuration) { //nolint:cyclop,maintidx
	cfgStorage := cfg.storage.Get()
	if cfgStorage.Name != nil {
		cfg.Name.Store(*cfgStorage.Name)
	}
	if cfgStorage.DeprecatedBootstrapKey != nil {
		cfg.DeprecatedBootstrapKey.Store(*cfgStorage.DeprecatedBootstrapKey)
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.ShowSystemInfo != nil && !misc.HasOSArg("i", "show-system-info", "") {
		cfg.HAProxy.ShowSystemInfo = *cfgStorage.Dataplaneapi.ShowSystemInfo
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.DisableInotify != nil && !misc.HasOSArg("", "disable-inotify", "") {
		cfg.HAProxy.DisableInotify = *cfgStorage.Dataplaneapi.DisableInotify
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.PIDFile != nil && !misc.HasOSArg("", "pid-file", "") {
		cfg.HAProxy.PIDFile = *cfgStorage.Dataplaneapi.PIDFile
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.DebugSocketPath != nil && !misc.HasOSArg("", "debug-socket-path", "") {
		cfg.HAProxy.DebugSocketPath = *cfgStorage.Dataplaneapi.DebugSocketPath
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.UID != nil && !misc.HasOSArg("", "uid", "") {
		cfg.HAProxy.UID = *cfgStorage.Dataplaneapi.UID
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.GID != nil && !misc.HasOSArg("", "gid", "") {
		cfg.HAProxy.GID = *cfgStorage.Dataplaneapi.GID
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Users != nil {
		cfg.Users = []User{}
		// If find users in dataplaneapi config file, then use them
		for _, item := range cfgStorage.Dataplaneapi.Users {
			itemUser := User{
				Name: item.Name,
			}
			if item.Insecure != nil {
				itemUser.Insecure = *item.Insecure
			}
			if item.Password != nil {
				itemUser.Password = *item.Password
			}
			cfg.Users = append(cfg.Users, itemUser)
		}
	}

	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.ConfigFile != nil && !misc.HasOSArg("c", "config-file", "") {
		cfg.HAProxy.ConfigFile = *cfgStorage.Haproxy.ConfigFile
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.HAProxy != nil && !misc.HasOSArg("b", "haproxy-bin", "") {
		cfg.HAProxy.HAProxy = *cfgStorage.Haproxy.HAProxy
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.MasterRuntime != nil && !misc.HasOSArg("m", "master-runtime", "") {
		cfg.HAProxy.MasterRuntime = *cfgStorage.Haproxy.MasterRuntime
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.NodeIDFile != nil && !misc.HasOSArg("", "fid", "") {
		cfg.HAProxy.NodeIDFile = *cfgStorage.Haproxy.NodeIDFile
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.MasterWorkerMode != nil && !misc.HasOSArg("", "master-worker-mode", "") {
		cfg.HAProxy.MasterWorkerMode = *cfgStorage.Haproxy.MasterWorkerMode
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Userlist != nil && cfgStorage.Dataplaneapi.Userlist.Userlist != nil && !misc.HasOSArg("u", "userlist", "") {
		cfg.HAProxy.Userlist = *cfgStorage.Dataplaneapi.Userlist.Userlist
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Userlist != nil && cfgStorage.Dataplaneapi.Userlist.UserListFile != nil && !misc.HasOSArg("", "userlist-file", "") {
		cfg.HAProxy.UserListFile = *cfgStorage.Dataplaneapi.Userlist.UserListFile
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.ReloadDelay != nil && !misc.HasOSArg("d", "reload-delay", "") {
		cfg.HAProxy.ReloadDelay = *cfgStorage.Haproxy.Reload.ReloadDelay
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.ReloadCmd != nil && !misc.HasOSArg("r", "reload-cmd", "") {
		cfg.HAProxy.ReloadCmd = *cfgStorage.Haproxy.Reload.ReloadCmd
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.RestartCmd != nil && !misc.HasOSArg("s", "restart-cmd", "") {
		cfg.HAProxy.RestartCmd = *cfgStorage.Haproxy.Reload.RestartCmd
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.StatusCmd != nil && !misc.HasOSArg("", "status-cmd", "") {
		cfg.HAProxy.StatusCmd = *cfgStorage.Haproxy.Reload.StatusCmd
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.ServiceName != nil && !misc.HasOSArg("", "service", "") {
		cfg.HAProxy.ServiceName = *cfgStorage.Haproxy.Reload.ServiceName
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.ReloadRetention != nil && !misc.HasOSArg("", "reload-retention", "") {
		cfg.HAProxy.ReloadRetention = *cfgStorage.Haproxy.Reload.ReloadRetention
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.ReloadStrategy != nil && !misc.HasOSArg("", "reload-strategy", "") {
		cfg.HAProxy.ReloadStrategy = *cfgStorage.Haproxy.Reload.ReloadStrategy
	}
	if cfgStorage.Haproxy != nil && cfgStorage.Haproxy.Reload != nil && cfgStorage.Haproxy.Reload.ValidateCmd != nil && !misc.HasOSArg("", "validate-cmd", "") {
		cfg.HAProxy.ValidateCmd = *cfgStorage.Haproxy.Reload.ValidateCmd
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Transaction != nil && cfgStorage.Dataplaneapi.Transaction.TransactionDir != nil && !misc.HasOSArg("t", "transaction-dir", "") {
		cfg.HAProxy.TransactionDir = *cfgStorage.Dataplaneapi.Transaction.TransactionDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Transaction != nil && cfgStorage.Dataplaneapi.Transaction.BackupsNumber != nil && !misc.HasOSArg("n", "backups-number", "") {
		cfg.HAProxy.BackupsNumber = *cfgStorage.Dataplaneapi.Transaction.BackupsNumber
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Transaction != nil && cfgStorage.Dataplaneapi.Transaction.BackupsDir != nil && !misc.HasOSArg("", "backups-dir", "") {
		cfg.HAProxy.BackupsDir = *cfgStorage.Dataplaneapi.Transaction.BackupsDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Transaction != nil && cfgStorage.Dataplaneapi.Transaction.MaxOpenTransactions != nil && !misc.HasOSArg("", "max-open-transactions", "") {
		cfg.HAProxy.MaxOpenTransactions = *cfgStorage.Dataplaneapi.Transaction.MaxOpenTransactions
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.MapsDir != nil && !misc.HasOSArg("p", "maps-dir", "") {
		cfg.HAProxy.MapsDir = *cfgStorage.Dataplaneapi.Resources.MapsDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.SSLCertsDir != nil && !misc.HasOSArg("", "ssl-certs-dir", "") {
		cfg.HAProxy.SSLCertsDir = *cfgStorage.Dataplaneapi.Resources.SSLCertsDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.GeneralStorageDir != nil && !misc.HasOSArg("", "general-storage-dir", "") {
		cfg.HAProxy.GeneralStorageDir = *cfgStorage.Dataplaneapi.Resources.GeneralStorageDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.DataplaneStorageDir != nil && !misc.HasOSArg("", "dataplane-storage-dir", "") {
		cfg.HAProxy.DataplaneStorageDir = *cfgStorage.Dataplaneapi.Resources.DataplaneStorageDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.UpdateMapFiles != nil && !misc.HasOSArg("", "update-map-files", "") {
		cfg.HAProxy.UpdateMapFiles = *cfgStorage.Dataplaneapi.Resources.UpdateMapFiles
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.UpdateMapFilesPeriod != nil && !misc.HasOSArg("", "update-map-files-period", "") {
		cfg.HAProxy.UpdateMapFilesPeriod = *cfgStorage.Dataplaneapi.Resources.UpdateMapFilesPeriod
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.SpoeDir != nil && !misc.HasOSArg("", "spoe-dir", "") {
		cfg.HAProxy.SpoeDir = *cfgStorage.Dataplaneapi.Resources.SpoeDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Resources != nil && cfgStorage.Dataplaneapi.Resources.SpoeTransactionDir != nil && !misc.HasOSArg("", "spoe-transaction-dir", "") {
		cfg.HAProxy.SpoeTransactionDir = *cfgStorage.Dataplaneapi.Resources.SpoeTransactionDir
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Advertised != nil && cfgStorage.Dataplaneapi.Advertised.APIAddress != nil && !misc.HasOSArg("", "api-address", "") {
		cfg.APIOptions.APIAddress = *cfgStorage.Dataplaneapi.Advertised.APIAddress
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Advertised != nil && cfgStorage.Dataplaneapi.Advertised.APIPort != nil && !misc.HasOSArg("", "api-port", "") {
		cfg.APIOptions.APIPort = *cfgStorage.Dataplaneapi.Advertised.APIPort
	}
	if cfgStorage.Log != nil && cfgStorage.Log.Syslog != nil && cfgStorage.Log.Syslog.SyslogAddr != nil && !misc.HasOSArg("", "syslog-address", "") {
		cfg.Syslog.SyslogAddr = *cfgStorage.Log.Syslog.SyslogAddr
	}
	if cfgStorage.Log != nil && cfgStorage.Log.Syslog != nil && cfgStorage.Log.Syslog.SyslogProto != nil && !misc.HasOSArg("", "syslog-protocol", "") {
		cfg.Syslog.SyslogProto = *cfgStorage.Log.Syslog.SyslogProto
	}
	if cfgStorage.Log != nil && cfgStorage.Log.Syslog != nil && cfgStorage.Log.Syslog.SyslogTag != nil && !misc.HasOSArg("", "syslog-tag", "") {
		cfg.Syslog.SyslogTag = *cfgStorage.Log.Syslog.SyslogTag
	}
	if cfgStorage.Log != nil && cfgStorage.Log.Syslog != nil && cfgStorage.Log.Syslog.SyslogLevel != nil && !misc.HasOSArg("", "syslog-level", "") {
		cfg.Syslog.SyslogLevel = *cfgStorage.Log.Syslog.SyslogLevel
	}
	if cfgStorage.Log != nil && cfgStorage.Log.Syslog != nil && cfgStorage.Log.Syslog.SyslogFacility != nil && !misc.HasOSArg("", "syslog-facility", "") {
		cfg.Syslog.SyslogFacility = *cfgStorage.Log.Syslog.SyslogFacility
	}
	if cfgStorage.Log != nil && cfgStorage.Log.LogTo != nil && !misc.HasOSArg("", "log-to", "") {
		cfg.Logging.LogTo = *cfgStorage.Log.LogTo
	}
	if cfgStorage.Log != nil && cfgStorage.Log.LogFile != nil && !misc.HasOSArg("", "log-file", "") {
		cfg.Logging.LogFile = *cfgStorage.Log.LogFile
	}
	if cfgStorage.Log != nil && cfgStorage.Log.LogLevel != nil && !misc.HasOSArg("", "log-level", "") {
		cfg.Logging.LogLevel = *cfgStorage.Log.LogLevel
	}
	if cfgStorage.Log != nil && cfgStorage.Log.LogFormat != nil && !misc.HasOSArg("", "log-format", "") {
		cfg.Logging.LogFormat = *cfgStorage.Log.LogFormat
	}
	if cfgStorage.Log != nil && cfgStorage.Log.ACLFormat != nil && !misc.HasOSArg("", "apache-common-log-format", "") {
		cfg.Logging.ACLFormat = *cfgStorage.Log.ACLFormat
	}
	if cfgStorage.LogTargets != nil {
		cfg.LogTargets = *cfgStorage.LogTargets
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Haproxy.DelayedStartMax != nil && !misc.HasOSArg("", "delayed-start-max", "") {
		if d, err := time.ParseDuration(*cfgStorage.Haproxy.DelayedStartMax); err == nil {
			cfg.HAProxy.DelayedStartMax = d
		}
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Haproxy.DelayedStartTick != nil && !misc.HasOSArg("", "delayed-start-tick", "") {
		if d, err := time.ParseDuration(*cfgStorage.Haproxy.DelayedStartTick); err == nil {
			cfg.HAProxy.DelayedStartTick = d
		}
	}
}

func copyConfigurationToStorage(cfg *Configuration) {
	cfgStorage := cfg.storage.Get()

	version := 2
	cfgStorage.Version = &version

	valueName := cfg.Name.Load()
	cfgStorage.Name = &valueName

	if cfgStorage.Dataplaneapi == nil {
		cfgStorage.Dataplaneapi = &configTypeDataplaneapi{}
	}
	if cfgStorage.Dataplaneapi.Advertised == nil {
		cfgStorage.Dataplaneapi.Advertised = &configTypeAdvertised{}
	}
	cfgStorage.Dataplaneapi.Advertised.APIAddress = &cfg.APIOptions.APIAddress
	cfgStorage.Dataplaneapi.Advertised.APIPort = &cfg.APIOptions.APIPort

	if cfgStorage.Dataplaneapi.Userlist == nil {
		cfgStorage.Dataplaneapi.Userlist = &configTypeUserlist{}
	}
	cfgStorage.Dataplaneapi.Userlist.Userlist = &cfg.HAProxy.Userlist
	cfgStorage.Dataplaneapi.Userlist.UserListFile = &cfg.HAProxy.UserListFile

	if cfgStorage.Haproxy == nil {
		cfgStorage.Haproxy = &configTypeHaproxy{
			Reload: &configTypeReload{},
		}
	}
	cfgStorage.Haproxy.Reload.ReloadStrategy = &cfg.HAProxy.ReloadStrategy

	if cfgStorage.LogTargets == nil {
		cfgStorage.LogTargets = &dpapilog.Targets{}
	}
	cfgStorage.LogTargets = &cfg.LogTargets
}

func (cfgStorage *StorageDataplaneAPIConfiguration) emptyDeprecatedSections() {
	cfgStorage.DeprecatedBootstrapKey = nil
	// Remove Cluster Users from dapi configuration file
	if cfgStorage.Dataplaneapi != nil {
		for i := 0; i < len(cfgStorage.Dataplaneapi.Users); {
			if cfgStorage.Dataplaneapi.Users[i].IsClusterUser() {
				if len(cfgStorage.Dataplaneapi.Users) > i {
					cfgStorage.Dataplaneapi.Users = append(cfgStorage.Dataplaneapi.Users[:i], cfgStorage.Dataplaneapi.Users[i+1:]...)
					continue
				}
			}
			i++
		}
	}
	// Remove Cluster
	cfgStorage.DeprecatedCluster = nil
	// Remove Status
	cfgStorage.DeprecatedStatus = nil
	// Remove Mode
	cfgStorage.DeprecatedMode = nil
	// Remove ServiceDiscovery Consuls and AWS Regions
	cfgStorage.DeprecatedServiceDiscovery = nil
}
