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
	"github.com/haproxytech/client-native/v6/models"
	"github.com/jessevdk/go-flags"

	dpapilog "github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/misc"
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
	User             []configTypeUser       `yaml:"user,omitempty"`
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

type configTypeUser struct {
	Name     string  `yaml:"name"`
	Insecure *bool   `yaml:"insecure,omitempty"`
	Password *string `yaml:"password,omitempty"`
}

type configTypeHaproxy struct {
	ConfigFile       *string           `yaml:"config_file,omitempty"`
	HAProxy          *string           `yaml:"haproxy_bin,omitempty"`
	MasterRuntime    *string           `yaml:"master_runtime,omitempty"`
	NodeIDFile       *string           `yaml:"fid,omitempty"`
	MasterWorkerMode *bool             `yaml:"master_worker_mode,omitempty"`
	Reload           *configTypeReload `yaml:"reload,omitempty"`
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
	UpdateMapFiles       *bool   `yaml:"update_map_files,omitempty"`
	UpdateMapFilesPeriod *int64  `yaml:"update_map_files_period,omitempty"`
	SpoeDir              *string `yaml:"spoe_dir,omitempty"`
	SpoeTransactionDir   *string `yaml:"spoe_transaction_dir,omitempty"`
}

type configTypeCluster struct {
	APINodesPath       *string                    `yaml:"api_nodes_path,omitempty"`
	Token              *string                    `yaml:"token,omitempty"`
	ClusterTLSCertDir  *string                    `yaml:"cluster_tls_dir,omitempty"`
	ActiveBootstrapKey *string                    `yaml:"active_bootstrap_key,omitempty"`
	APIRegisterPath    *string                    `yaml:"api_register_path,omitempty"`
	URL                *string                    `yaml:"url,omitempty"`
	Port               *int                       `yaml:"port,omitempty"`
	StorageDir         *string                    `yaml:"storage_dir,omitempty"`
	BootstrapKey       *string                    `yaml:"bootstrap_key,omitempty"`
	ID                 *string                    `yaml:"id,omitempty"`
	APIBasePath        *string                    `yaml:"api_base_path,omitempty"`
	CertificateDir     *string                    `yaml:"cert_path,omitempty"`
	CertificateFetched *bool                      `yaml:"cert_fetched,omitempty"`
	Name               *string                    `yaml:"name,omitempty"`
	Description        *string                    `yaml:"description,omitempty"`
	ClusterID          *string                    `yaml:"cluster_id,omitempty" group:"cluster" save:"true"`
	ClusterLogTargets  []*models.ClusterLogTarget `yaml:"cluster_log_targets,omitempty" group:"cluster" save:"true"`
}

type configTypeAdvertised struct {
	APIAddress *string `yaml:"api_address,omitempty"`
	APIPort    *int64  `yaml:"api_port,omitempty"`
}

type configTypeServiceDiscovery struct {
	Consuls    *[]*models.Consul    `yaml:"consuls,omitempty"`
	AWSRegions *[]*models.AwsRegion `yaml:"aws_regions,omitempty"`
}

type configTypeSyslog struct {
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
	Syslog    *configTypeSyslog `yaml:"syslog,omitempty"`
}

type StorageDataplaneAPIConfiguration struct {
	Version                *int                        `yaml:"config_version,omitempty"`
	Name                   *string                     `yaml:"name,omitempty"`
	Mode                   *string                     `yaml:"mode,omitempty"`
	DeprecatedBootstrapKey *string                     `yaml:"bootstrap_key,omitempty"`
	Status                 *string                     `yaml:"status,omitempty"`
	Dataplaneapi           *configTypeDataplaneapi     `yaml:"dataplaneapi,omitempty"`
	Haproxy                *configTypeHaproxy          `yaml:"haproxy,omitempty"`
	Cluster                *configTypeCluster          `yaml:"cluster,omitempty"`
	ServiceDiscovery       *configTypeServiceDiscovery `yaml:"service_discovery,omitempty"`
	Log                    *configTypeLog              `yaml:"log,omitempty"`
	LogTargets             *dpapilog.Targets           `yaml:"log_targets,omitempty"`
}

func copyToConfiguration(cfg *Configuration) { //nolint:cyclop,maintidx
	cfgStorage := cfg.storage.Get()
	if cfgStorage.Name != nil {
		cfg.Name.Store(*cfgStorage.Name)
	}
	if cfgStorage.Mode != nil {
		cfg.Mode.Store(*cfgStorage.Mode)
	}
	if cfgStorage.DeprecatedBootstrapKey != nil {
		cfg.DeprecatedBootstrapKey.Store(*cfgStorage.DeprecatedBootstrapKey)
	}
	if cfgStorage.Status != nil {
		cfg.Status.Store(*cfgStorage.Status)
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
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.User != nil {
		cfg.Users = []User{}
		for _, item := range cfgStorage.Dataplaneapi.User {
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
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.ClusterTLSCertDir != nil && !misc.HasOSArg("", "cluster-tls-dir", "") {
		cfg.HAProxy.ClusterTLSCertDir = *cfgStorage.Cluster.ClusterTLSCertDir
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.ID != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.ID.Store(*cfgStorage.Cluster.ID)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.BootstrapKey != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.BootstrapKey.Store(*cfgStorage.Cluster.BootstrapKey)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.ActiveBootstrapKey != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.ActiveBootstrapKey.Store(*cfgStorage.Cluster.ActiveBootstrapKey)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.Token != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.Token.Store(*cfgStorage.Cluster.Token)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.URL != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.URL.Store(*cfgStorage.Cluster.URL)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.Port != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.Port.Store(*cfgStorage.Cluster.Port)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.APIBasePath != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.APIBasePath.Store(*cfgStorage.Cluster.APIBasePath)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.APINodesPath != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.APINodesPath.Store(*cfgStorage.Cluster.APINodesPath)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.APIRegisterPath != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.APIRegisterPath.Store(*cfgStorage.Cluster.APIRegisterPath)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.StorageDir != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.StorageDir.Store(*cfgStorage.Cluster.StorageDir)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.CertificateDir != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.CertificateDir.Store(*cfgStorage.Cluster.CertificateDir)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.CertificateFetched != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.CertificateFetched.Store(*cfgStorage.Cluster.CertificateFetched)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.Name != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.Name.Store(*cfgStorage.Cluster.Name)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.Description != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.Description.Store(*cfgStorage.Cluster.Description)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.ClusterID != nil && !misc.HasOSArg("", "", "") {
		cfg.Cluster.ClusterID.Store(*cfgStorage.Cluster.ClusterID)
	}
	if cfgStorage.Cluster != nil && cfgStorage.Cluster.ClusterLogTargets != nil && len(cfgStorage.Cluster.ClusterLogTargets) > 0 {
		cfg.Cluster.ClusterLogTargets = cfgStorage.Cluster.ClusterLogTargets
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Advertised != nil && cfgStorage.Dataplaneapi.Advertised.APIAddress != nil && !misc.HasOSArg("", "api-address", "") {
		cfg.APIOptions.APIAddress = *cfgStorage.Dataplaneapi.Advertised.APIAddress
	}
	if cfgStorage.Dataplaneapi != nil && cfgStorage.Dataplaneapi.Advertised != nil && cfgStorage.Dataplaneapi.Advertised.APIPort != nil && !misc.HasOSArg("", "api-port", "") {
		cfg.APIOptions.APIPort = *cfgStorage.Dataplaneapi.Advertised.APIPort
	}
	if cfgStorage.ServiceDiscovery != nil && cfgStorage.ServiceDiscovery.Consuls != nil && !misc.HasOSArg("", "", "") {
		cfg.ServiceDiscovery.Consuls = *cfgStorage.ServiceDiscovery.Consuls
	}
	if cfgStorage.ServiceDiscovery != nil && cfgStorage.ServiceDiscovery.AWSRegions != nil && !misc.HasOSArg("", "", "") {
		cfg.ServiceDiscovery.AWSRegions = *cfgStorage.ServiceDiscovery.AWSRegions
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
}

func copyConfigurationToStorage(cfg *Configuration) {
	cfgStorage := cfg.storage.Get()

	version := 2
	cfgStorage.Version = &version

	valueName := cfg.Name.Load()
	cfgStorage.Name = &valueName

	valueMode := cfg.Mode.Load()
	cfgStorage.Mode = &valueMode

	cfgStorage.DeprecatedBootstrapKey = nil

	valueStatus := cfg.Status.Load()
	cfgStorage.Status = &valueStatus

	valueClusterNodeID := cfg.Cluster.ID.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.ID = &valueClusterNodeID

	valueClusterBootstrapKey := cfg.Cluster.BootstrapKey.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.BootstrapKey = &valueClusterBootstrapKey

	valueClusterActiveBootstrapKey := cfg.Cluster.ActiveBootstrapKey.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.ActiveBootstrapKey = &valueClusterActiveBootstrapKey

	valueClusterToken := cfg.Cluster.Token.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.Token = &valueClusterToken

	valueClusterURL := cfg.Cluster.URL.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.URL = &valueClusterURL

	valueClusterPort := cfg.Cluster.Port.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.Port = &valueClusterPort

	valueClusterAPIBasePath := cfg.Cluster.APIBasePath.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.APIBasePath = &valueClusterAPIBasePath

	valueClusterAPINodesPath := cfg.Cluster.APINodesPath.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.APINodesPath = &valueClusterAPINodesPath

	valueClusterAPIRegisterPath := cfg.Cluster.APIRegisterPath.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.APIRegisterPath = &valueClusterAPIRegisterPath

	valueClusterStorageDir := cfg.Cluster.StorageDir.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.StorageDir = &valueClusterStorageDir

	valueClusterCertificateDir := cfg.Cluster.CertificateDir.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.CertificateDir = &valueClusterCertificateDir

	valueClusterCertificateFetched := cfg.Cluster.CertificateFetched.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	if cfgStorage.Cluster != nil {
		cfgStorage.Cluster.CertificateFetched = &valueClusterCertificateFetched
	}

	valueClusterName := cfg.Cluster.Name.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.Name = &valueClusterName

	valueClusterDescription := cfg.Cluster.Description.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.Description = &valueClusterDescription

	valueClusterID := cfg.Cluster.ClusterID.Load()
	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.ClusterID = &valueClusterID

	if cfgStorage.Cluster == nil {
		cfgStorage.Cluster = &configTypeCluster{}
	}
	cfgStorage.Cluster.ClusterLogTargets = cfg.Cluster.ClusterLogTargets

	if cfgStorage.Dataplaneapi == nil {
		cfgStorage.Dataplaneapi = &configTypeDataplaneapi{}
	}
	if cfgStorage.Dataplaneapi.Advertised == nil {
		cfgStorage.Dataplaneapi.Advertised = &configTypeAdvertised{}
	}
	cfgStorage.Dataplaneapi.Advertised.APIAddress = &cfg.APIOptions.APIAddress

	if cfgStorage.Dataplaneapi == nil {
		cfgStorage.Dataplaneapi = &configTypeDataplaneapi{}
	}
	if cfgStorage.Dataplaneapi.Advertised == nil {
		cfgStorage.Dataplaneapi.Advertised = &configTypeAdvertised{}
	}
	cfgStorage.Dataplaneapi.Advertised.APIPort = &cfg.APIOptions.APIPort

	if cfgStorage.ServiceDiscovery == nil {
		cfgStorage.ServiceDiscovery = &configTypeServiceDiscovery{}
	}
	cfgStorage.ServiceDiscovery.Consuls = &cfg.ServiceDiscovery.Consuls

	if cfgStorage.ServiceDiscovery == nil {
		cfgStorage.ServiceDiscovery = &configTypeServiceDiscovery{}
	}
	cfgStorage.ServiceDiscovery.AWSRegions = &cfg.ServiceDiscovery.AWSRegions

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
