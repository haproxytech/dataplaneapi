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

package dataplaneapi

import (
	"time"

	"github.com/docker/go-units"
	"github.com/go-openapi/runtime/flagext"

	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/misc"
)

func SyncWithFileSettings(server *Server, cfg *configuration.Configuration) { //nolint: cyclop
	configStorage := cfg.GetStorageData()
	// This is added to allow backward compatibility if no scheme is defined it defaults to specification schemes
	// and recently we added https to specification for clarity, and then old configs without scheme in them would
	// start to fail.
	if len(server.EnabledListeners) == 0 && configStorage.Dataplaneapi != nil {
		if configStorage.Dataplaneapi.EnabledListeners == nil {
			configStorage.Dataplaneapi.EnabledListeners = &[]string{"http"}
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.EnabledListeners != nil && !misc.HasOSArg("", "scheme", "") {
		server.EnabledListeners = *configStorage.Dataplaneapi.EnabledListeners
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.CleanupTimeout != nil && !misc.HasOSArg("", "cleanup-timeout", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.CleanupTimeout); err == nil {
			server.CleanupTimeout = d
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.GracefulTimeout != nil && !misc.HasOSArg("", "graceful-timeout", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.GracefulTimeout); err == nil {
			server.GracefulTimeout = d
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.MaxHeaderSize != nil && !misc.HasOSArg("", "max-header-size", "") {
		s, err := units.FromHumanSize(*configStorage.Dataplaneapi.MaxHeaderSize)
		if err == nil {
			server.MaxHeaderSize = flagext.ByteSize(s)
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.SocketPath != nil && !misc.HasOSArg("", "socket-path", "") {
		server.SocketPath = *configStorage.Dataplaneapi.SocketPath
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.Host != nil && !misc.HasOSArg("", "host", "HOST") {
		server.Host = *configStorage.Dataplaneapi.Host
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.Port != nil && !misc.HasOSArg("", "port", "PORT") {
		server.Port = *configStorage.Dataplaneapi.Port
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.ListenLimit != nil && !misc.HasOSArg("", "listen-limit", "") {
		server.ListenLimit = *configStorage.Dataplaneapi.ListenLimit
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.KeepAlive != nil && !misc.HasOSArg("", "keep-alive", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.KeepAlive); err == nil {
			server.KeepAlive = d
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.ReadTimeout != nil && !misc.HasOSArg("", "read-timeout", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.ReadTimeout); err == nil {
			server.ReadTimeout = d
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.WriteTimeout != nil && !misc.HasOSArg("", "write-timeout", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.WriteTimeout); err == nil {
			server.WriteTimeout = d
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSHost != nil && !misc.HasOSArg("", "tls-host", "TLS_HOST") {
		server.TLSHost = *configStorage.Dataplaneapi.TLS.TLSHost
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSPort != nil && !misc.HasOSArg("", "tls-port", "TLS_PORT") {
		server.TLSPort = *configStorage.Dataplaneapi.TLS.TLSPort
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSCertificate != nil && !misc.HasOSArg("", "tls-certificate", "TLS_CERTIFICATE") {
		server.TLSCertificate = *configStorage.Dataplaneapi.TLS.TLSCertificate
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSCertificateKey != nil && !misc.HasOSArg("", "tls-key", "TLS_PRIVATE_KEY") {
		server.TLSCertificateKey = *configStorage.Dataplaneapi.TLS.TLSCertificateKey
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSCACertificate != nil && !misc.HasOSArg("", "tls-ca", "TLS_CA_CERTIFICATE") {
		server.TLSCACertificate = *configStorage.Dataplaneapi.TLS.TLSCACertificate
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSListenLimit != nil && !misc.HasOSArg("", "tls-listen-limit", "") {
		server.TLSListenLimit = *configStorage.Dataplaneapi.TLS.TLSListenLimit
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSKeepAlive != nil && !misc.HasOSArg("", "tls-keep-alive", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.TLS.TLSKeepAlive); err == nil {
			server.TLSKeepAlive = d
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSReadTimeout != nil && !misc.HasOSArg("", "tls-read-timeout", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.TLS.TLSReadTimeout); err == nil {
			server.TLSReadTimeout = d
		}
	}
	if configStorage.Dataplaneapi != nil && configStorage.Dataplaneapi.TLS != nil && configStorage.Dataplaneapi.TLS.TLSWriteTimeout != nil && !misc.HasOSArg("", "tls-write-timeout", "") {
		if d, err := time.ParseDuration(*configStorage.Dataplaneapi.TLS.TLSWriteTimeout); err == nil {
			server.TLSWriteTimeout = d
		}
	}
}
