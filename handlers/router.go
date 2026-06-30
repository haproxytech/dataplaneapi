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

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"

	sc "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/handlers/respond"
	"github.com/haproxytech/dataplaneapi/reload_agent"

	// configuration
	"github.com/haproxytech/dataplaneapi/handlers/configuration/acl"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/acme"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/backend"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/backend_switching_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/bind"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/cache"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/capture"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/crt_store"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/defaults"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/dgram_bind"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/fcgi_app"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/filter"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/force_be_switch"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/frontend"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/global"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/groups"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/health_check"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/http/http_after_response_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/http/http_check"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/http/http_error_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/http/http_request_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/http/http_response_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/http_errors_section"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/log_forward"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/log_profile"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/log_target"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/log_targets_global"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/mailers"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/peers"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/quic/quic_initial_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/raw"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/resolver"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/ring"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/server"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/server_switching_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/server_template"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/ssl_front_use"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/stick_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/table"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/tcp/tcp_check"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/tcp/tcp_request_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/tcp/tcp_response_rule"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/traces"
	"github.com/haproxytech/dataplaneapi/handlers/configuration/userlist"
	config_version "github.com/haproxytech/dataplaneapi/handlers/configuration/version"

	// dataplane
	"github.com/haproxytech/dataplaneapi/handlers/dataplane/discovery"
	"github.com/haproxytech/dataplaneapi/handlers/dataplane/health"
	"github.com/haproxytech/dataplaneapi/handlers/dataplane/info"
	"github.com/haproxytech/dataplaneapi/handlers/dataplane/reloads"
	"github.com/haproxytech/dataplaneapi/handlers/dataplane/sites"
	"github.com/haproxytech/dataplaneapi/handlers/dataplane/specification"
	"github.com/haproxytech/dataplaneapi/handlers/dataplane/stats"
	dataplane_transactions "github.com/haproxytech/dataplaneapi/handlers/dataplane/transactions"

	// runtime
	"github.com/haproxytech/dataplaneapi/handlers/runtime/acl_entries"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/acl_files"
	runtime_acme "github.com/haproxytech/dataplaneapi/handlers/runtime/acme"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/map_entries"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/maps"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/process_info"
	runtime_servers "github.com/haproxytech/dataplaneapi/handlers/runtime/servers"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/ssl_ca_files"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/ssl_certs"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/ssl_crl_files"
	runtime_ssl_crt_lists "github.com/haproxytech/dataplaneapi/handlers/runtime/ssl_crt_lists"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/stick_table_entries"
	"github.com/haproxytech/dataplaneapi/handlers/runtime/stick_tables"

	// service discovery
	"github.com/haproxytech/dataplaneapi/handlers/service_discovery/aws"
	"github.com/haproxytech/dataplaneapi/handlers/service_discovery/consul"

	// spoe
	"github.com/haproxytech/dataplaneapi/handlers/spoe/agents"
	spoe_files "github.com/haproxytech/dataplaneapi/handlers/spoe/files"
	spoe_groups "github.com/haproxytech/dataplaneapi/handlers/spoe/groups"
	"github.com/haproxytech/dataplaneapi/handlers/spoe/messages"
	"github.com/haproxytech/dataplaneapi/handlers/spoe/scopes"
	spoe_transactions "github.com/haproxytech/dataplaneapi/handlers/spoe/transactions"
	spoe_version "github.com/haproxytech/dataplaneapi/handlers/spoe/version"

	// storage
	"github.com/haproxytech/dataplaneapi/handlers/storage/general"
	storage_maps "github.com/haproxytech/dataplaneapi/handlers/storage/maps"
	"github.com/haproxytech/dataplaneapi/handlers/storage/ssl_certificates"
	storage_ssl_crt_lists "github.com/haproxytech/dataplaneapi/handlers/storage/ssl_crt_lists"
)

// Options holds all dependencies needed to build the router.
type Options struct {
	Client      client_native.HAProxyClient
	ReloadAgent reload_agent.IReloadAgent

	// Service discovery
	ConsulDiscovery       sc.ServiceDiscoveries
	ConsulPersistCallback func([]*models.Consul) error
	AWSDiscovery          sc.ServiceDiscoveries
	AWSPersistCallback    func([]*models.AwsRegion) error
	UseValidation         bool

	// API info
	Version    string
	BuildTime  string
	SystemInfo bool

	// MaxOpenTransactions limits concurrent in-progress transactions (0 = unlimited).
	MaxOpenTransactions int

	// SwaggerJSON is the raw OpenAPI v2 / swagger spec passed to the discovery and
	// specification handlers. It must not be nil.
	SwaggerJSON json.RawMessage
}

// NewRouter assembles all handler packages onto a single chi router and returns
// it as an http.Handler. All routes are served under the /v3 prefix to match
// the OpenAPI spec's server URL. The first error encountered during registration
// is returned; all RegisterRouter calls are otherwise non-failing (the embedded
// OpenAPI spec is always valid).
func NewRouter(opts Options) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	// Emit JSON instead of chi's default plain-text bodies for unmatched routes
	// and disallowed methods, matching the previous go-swagger server.
	r.NotFound(respond.NotFound)
	r.MethodNotAllowed(respond.MethodNotAllowed(r))

	registrations := []func() error{
		// — configuration —
		func() error { return acl.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return acme.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return backend.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return backend_switching_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return bind.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return cache.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return capture.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return crt_store.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return defaults.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return dgram_bind.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return fcgi_app.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return filter.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return force_be_switch.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return frontend.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return global.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return groups.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return health_check.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return http_after_response_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return http_check.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return http_error_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return http_request_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return http_response_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return http_errors_section.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return log_forward.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return log_profile.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return log_target.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return log_targets_global.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return mailers.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return peers.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return quic_initial_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return raw.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return resolver.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return ring.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return server.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return server_switching_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return server_template.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return ssl_front_use.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return stick_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return table.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return tcp_check.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return tcp_request_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return tcp_response_rule.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return traces.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return userlist.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return config_version.RegisterRouter(r, opts.Client) },

		// — runtime —
		func() error { return acl_entries.RegisterRouter(r, opts.Client) },
		func() error { return acl_files.RegisterRouter(r, opts.Client) },
		func() error { return runtime_acme.RegisterRouter(r, opts.Client) },
		func() error { return map_entries.RegisterRouter(r, opts.Client) },
		func() error { return maps.RegisterRouter(r, opts.Client) },
		func() error { return process_info.RegisterRouter(r, opts.Client) },
		func() error { return runtime_servers.RegisterRouter(r, opts.Client) },
		func() error { return ssl_ca_files.RegisterRouter(r, opts.Client) },
		func() error { return ssl_certs.RegisterRouter(r, opts.Client) },
		func() error { return ssl_crl_files.RegisterRouter(r, opts.Client) },
		func() error { return runtime_ssl_crt_lists.RegisterRouter(r, opts.Client) },
		func() error { return stick_table_entries.RegisterRouter(r, opts.Client) },
		func() error { return stick_tables.RegisterRouter(r, opts.Client) },

		// — spoe —
		func() error { return agents.RegisterRouter(r, opts.Client) },
		func() error { return spoe_files.RegisterRouter(r, opts.Client) },
		func() error { return spoe_groups.RegisterRouter(r, opts.Client) },
		func() error { return messages.RegisterRouter(r, opts.Client) },
		func() error { return scopes.RegisterRouter(r, opts.Client) },
		func() error { return spoe_transactions.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return spoe_version.RegisterRouter(r, opts.Client) },

		// — dataplane —
		func() error { return discovery.RegisterRouter(r, opts.SwaggerJSON) },
		func() error { return health.RegisterRouter(r, opts.ReloadAgent) },
		func() error { return info.RegisterRouter(r, opts.Version, opts.BuildTime, opts.SystemInfo) },
		func() error { return reloads.RegisterRouter(r, opts.ReloadAgent) },
		func() error { return sites.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return specification.RegisterRouter(r, opts.SwaggerJSON) },
		func() error { return stats.RegisterRouter(r, opts.Client) },
		func() error {
			return dataplane_transactions.RegisterRouter(r, opts.Client, opts.ReloadAgent, opts.MaxOpenTransactions)
		},

		// — storage —
		func() error { return general.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return storage_maps.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return ssl_certificates.RegisterRouter(r, opts.Client, opts.ReloadAgent) },
		func() error { return storage_ssl_crt_lists.RegisterRouter(r, opts.Client, opts.ReloadAgent) },

		// — service discovery —
		func() error {
			return consul.RegisterRouter(r, opts.ConsulDiscovery, opts.ConsulPersistCallback, opts.UseValidation)
		},
		func() error {
			return aws.RegisterRouter(r, opts.AWSDiscovery, opts.AWSPersistCallback, opts.UseValidation)
		},
	}

	for _, reg := range registrations {
		if err := reg(); err != nil {
			return nil, err
		}
	}

	root := chi.NewRouter()
	root.NotFound(respond.NotFound)
	root.MethodNotAllowed(respond.MethodNotAllowed(root))
	root.Mount("/v3", r)
	return root, nil
}
