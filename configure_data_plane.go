// This file is safe to edit. Once it exists it will not be overwritten

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

package dataplaneapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	client_native "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/models"
	runtime_api "github.com/haproxytech/client-native/v2/runtime"
	"github.com/haproxytech/client-native/v2/spoe"
	"github.com/haproxytech/client-native/v2/storage"
	parser "github.com/haproxytech/config-parser/v3"
	"github.com/haproxytech/config-parser/v3/types"
	"github.com/haproxytech/dataplaneapi/syslog"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"

	"github.com/haproxytech/dataplaneapi/adapters"
	dataplaneapi_config "github.com/haproxytech/dataplaneapi/configuration"
	service_discovery "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/handlers"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations"
	"github.com/haproxytech/dataplaneapi/operations/discovery"
	"github.com/haproxytech/dataplaneapi/operations/specification"
	"github.com/haproxytech/dataplaneapi/operations/specification_openapiv3"
	"github.com/haproxytech/dataplaneapi/rate"

	// import various crypting algorithms
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

// go:generate swagger generate server --target ../../../../../../github.com/haproxytech --name controller --spec ../../../../../../../../haproxy-api/haproxy-open-api-spec/build/haproxy_spec.yaml --server-package controller --tags Stats --tags Information --tags Configuration --tags Discovery --tags Frontend --tags Backend --tags Bind --tags Server --tags TCPRequestRule --tags HTTPRequestRule --tags HTTPResponseRule --tags Acl --tags BackendSwitchingRule --tags ServerSwitchingRule --tags TCPResponseRule --skip-models --exclude-main

var (
	Version   string
	BuildTime string
	mWorker   = false
	logFile   *os.File
)

func configureFlags(api *operations.DataPlaneAPI) {
	cfg := dataplaneapi_config.Get()

	haproxyOptionsGroup := swag.CommandLineOptionsGroup{
		ShortDescription: "HAProxy options",
		LongDescription:  "Options for configuring haproxy locations.",
		Options:          &cfg.HAProxy,
	}

	loggingOptionsGroup := swag.CommandLineOptionsGroup{
		ShortDescription: "Logging options",
		LongDescription:  "Options for configuring logging.",
		Options:          &cfg.Logging,
	}

	syslogOptionsGroup := swag.CommandLineOptionsGroup{
		ShortDescription: "Syslog options",
		LongDescription:  "Options for configuring syslog logging.",
		Options:          &cfg.Syslog,
	}

	api.CommandLineOptionsGroups = make([]swag.CommandLineOptionsGroup, 0, 1)
	api.CommandLineOptionsGroups = append(api.CommandLineOptionsGroups, haproxyOptionsGroup)
	api.CommandLineOptionsGroups = append(api.CommandLineOptionsGroups, loggingOptionsGroup)
	api.CommandLineOptionsGroups = append(api.CommandLineOptionsGroups, syslogOptionsGroup)
}

func configureAPI(api *operations.DataPlaneAPI) http.Handler {
	cfg := dataplaneapi_config.Get()

	haproxyOptions := cfg.HAProxy
	// Override options with env variables
	if os.Getenv("HAPROXY_MWORKER") == "1" {
		mWorker = true
		masterRuntime := os.Getenv("HAPROXY_MASTER_CLI")
		if misc.IsUnixSocketAddr(masterRuntime) {
			haproxyOptions.MasterRuntime = strings.Replace(masterRuntime, "unix@", "", 1)
		}
	}
	// Override options with env variables
	if os.Getenv("HAPROXY_MWORKER") == "1" {
		mWorker = true
		masterRuntime := os.Getenv("HAPROXY_MASTER_CLI")
		if misc.IsUnixSocketAddr(masterRuntime) {
			haproxyOptions.MasterRuntime = strings.Replace(masterRuntime, "unix@", "", 1)
		}
	}

	if cfgFiles := os.Getenv("HAPROXY_CFGFILES"); cfgFiles != "" {
		m := map[string]bool{"configuration": false}
		if len(haproxyOptions.UserListFile) > 0 {
			m["userlist"] = false
		}

		for _, f := range strings.Split(cfgFiles, ";") {
			var conf bool
			var user bool

			if f == haproxyOptions.ConfigFile {
				conf = true
				m["configuration"] = true
			}
			if len(haproxyOptions.UserListFile) > 0 && f == haproxyOptions.UserListFile {
				user = true
				m["userlist"] = true
			}
			if !conf && !user {
				log.Warningf("The configuration file %s in HAPROXY_CFGFILES is not defined, neither by --config-file or --userlist-file flags.", f)
			}
		}
		for f, ok := range m {
			if !ok {
				log.Fatalf("The %s file is not declared in the HAPROXY_CFGFILES environment variable, cannot start.", f)
			}
		}
	}
	// end overriding options with env variables

	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.TxtConsumer = runtime.TextConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ServerShutdown = serverShutdown

	ctx := ContextHandler.Context()
	clientCtx, cancel := context.WithCancel(ctx)

	client := configureNativeClient(clientCtx, haproxyOptions, mWorker)

	users := dataplaneapi_config.GetUsersStore()

	// Handle reload signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	go handleSignals(ctx, cancel, sigs, client, haproxyOptions, users)

	if !haproxyOptions.DisableInotify {
		if err := startWatcher(ctx, client, haproxyOptions, users); err != nil {
			haproxyOptions.DisableInotify = true
			client = configureNativeClient(clientCtx, haproxyOptions, mWorker)
		}
	}

	// Sync map physical file with runtime map entries
	if haproxyOptions.UpdateMapFiles {
		go cfg.MapSync.SyncAll(client)
	}

	// Initialize reload agent
	raParams := haproxy.ReloadAgentParams{
		Delay:      haproxyOptions.ReloadDelay,
		ReloadCmd:  haproxyOptions.ReloadCmd,
		RestartCmd: haproxyOptions.RestartCmd,
		ConfigFile: haproxyOptions.ConfigFile,
		BackupDir:  haproxyOptions.BackupsDir,
		Retention:  haproxyOptions.ReloadRetention,
		Ctx:        ctx,
	}

	ra, e := haproxy.NewReloadAgent(raParams)
	if e != nil {
		log.Fatalf("Cannot initialize reload agent: %v", e)
	}

	// setup discovery handlers
	api.DiscoveryGetAPIEndpointsHandler = discovery.GetAPIEndpointsHandlerFunc(func(params discovery.GetAPIEndpointsParams, principal interface{}) middleware.Responder {
		ends, err := misc.DiscoverChildPaths("", SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetAPIEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetAPIEndpointsOK().WithPayload(ends)
	})
	api.DiscoveryGetServicesEndpointsHandler = discovery.GetServicesEndpointsHandlerFunc(func(params discovery.GetServicesEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetServicesEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetServicesEndpointsOK().WithPayload(ends)
	})
	api.DiscoveryGetConfigurationEndpointsHandler = discovery.GetConfigurationEndpointsHandlerFunc(func(params discovery.GetConfigurationEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetConfigurationEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetConfigurationEndpointsOK().WithPayload(ends)
	})
	api.DiscoveryGetRuntimeEndpointsHandler = discovery.GetRuntimeEndpointsHandlerFunc(func(params discovery.GetRuntimeEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetRuntimeEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetRuntimeEndpointsOK().WithPayload(ends)
	})
	api.DiscoveryGetHaproxyEndpointsHandler = discovery.GetHaproxyEndpointsHandlerFunc(func(params discovery.GetHaproxyEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetHaproxyEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetHaproxyEndpointsOK().WithPayload(ends)
	})
	api.DiscoveryGetStatsEndpointsHandler = discovery.GetStatsEndpointsHandlerFunc(func(params discovery.GetStatsEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetStatsEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetStatsEndpointsOK().WithPayload(ends)
	})
	api.DiscoveryGetSpoeEndpointsHandler = discovery.GetSpoeEndpointsHandlerFunc(func(params discovery.GetSpoeEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetSpoeEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetSpoeEndpointsOK().WithPayload(ends)
	})
	api.DiscoveryGetStorageEndpointsHandler = discovery.GetStorageEndpointsHandlerFunc(func(params discovery.GetStorageEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return discovery.NewGetStorageEndpointsDefault(int(*e.Code)).WithPayload(e)
		}
		return discovery.NewGetStorageEndpointsOK().WithPayload(ends)
	})

	// setup transaction handlers
	api.TransactionsStartTransactionHandler = &handlers.StartTransactionHandlerImpl{Client: client}
	api.TransactionsDeleteTransactionHandler = &handlers.DeleteTransactionHandlerImpl{Client: client}
	api.TransactionsGetTransactionHandler = &handlers.GetTransactionHandlerImpl{Client: client}
	api.TransactionsGetTransactionsHandler = &handlers.GetTransactionsHandlerImpl{Client: client}
	api.TransactionsCommitTransactionHandler = &handlers.CommitTransactionHandlerImpl{Client: client, ReloadAgent: ra, Mutex: &sync.Mutex{}}
	if cfg.HAProxy.MaxOpenTransactions > 0 {
		// creating the threshold limit using the CLI flag as hard quota and current open transactions as starting point
		actualCount := func() uint64 {
			ts, err := client.Configuration.GetTransactions(models.TransactionStatusInProgress)
			if err != nil {
				log.Errorf("Cannot retrieve current open transactions for rate limit, default to zero (%s)", err.Error())
				return 0
			}
			return uint64(len(*ts))
		}
		transactionLimiter := rate.NewThresholdLimit(uint64(cfg.HAProxy.MaxOpenTransactions), actualCount)
		api.TransactionsStartTransactionHandler = &handlers.RateLimitedStartTransactionHandlerImpl{
			TransactionCounter: transactionLimiter,
			Handler:            api.TransactionsStartTransactionHandler,
		}
	}

	// setup transaction handlers
	api.SpoeTransactionsStartSpoeTransactionHandler = &handlers.SpoeTransactionsStartSpoeTransactionHandlerImpl{Client: client}
	api.SpoeTransactionsDeleteSpoeTransactionHandler = &handlers.SpoeTransactionsDeleteSpoeTransactionHandlerImpl{Client: client}
	api.SpoeTransactionsGetSpoeTransactionHandler = &handlers.SpoeTransactionsGetSpoeTransactionHandlerImpl{Client: client}
	api.SpoeTransactionsGetSpoeTransactionsHandler = &handlers.SpoeTransactionsGetSpoeTransactionsHandlerImpl{Client: client}
	api.SpoeTransactionsCommitSpoeTransactionHandler = &handlers.SpoeTransactionsCommitSpoeTransactionHandlerImpl{Client: client, ReloadAgent: ra}

	// setup sites handlers
	api.SitesCreateSiteHandler = &handlers.CreateSiteHandlerImpl{Client: client, ReloadAgent: ra}
	api.SitesDeleteSiteHandler = &handlers.DeleteSiteHandlerImpl{Client: client, ReloadAgent: ra}
	api.SitesGetSiteHandler = &handlers.GetSiteHandlerImpl{Client: client}
	api.SitesGetSitesHandler = &handlers.GetSitesHandlerImpl{Client: client}
	api.SitesReplaceSiteHandler = &handlers.ReplaceSiteHandlerImpl{Client: client, ReloadAgent: ra}

	// setup backend handlers
	api.BackendCreateBackendHandler = &handlers.CreateBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.BackendDeleteBackendHandler = &handlers.DeleteBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.BackendGetBackendHandler = &handlers.GetBackendHandlerImpl{Client: client}
	api.BackendGetBackendsHandler = &handlers.GetBackendsHandlerImpl{Client: client}
	api.BackendReplaceBackendHandler = &handlers.ReplaceBackendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup frontend handlers
	api.FrontendCreateFrontendHandler = &handlers.CreateFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FrontendDeleteFrontendHandler = &handlers.DeleteFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FrontendGetFrontendHandler = &handlers.GetFrontendHandlerImpl{Client: client}
	api.FrontendGetFrontendsHandler = &handlers.GetFrontendsHandlerImpl{Client: client}
	api.FrontendReplaceFrontendHandler = &handlers.ReplaceFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup server handlers
	api.ServerCreateServerHandler = &handlers.CreateServerHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerDeleteServerHandler = &handlers.DeleteServerHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerGetServerHandler = &handlers.GetServerHandlerImpl{Client: client}
	api.ServerGetServersHandler = &handlers.GetServersHandlerImpl{Client: client}
	api.ServerReplaceServerHandler = &handlers.ReplaceServerHandlerImpl{Client: client, ReloadAgent: ra}

	// setup bind handlers
	api.BindCreateBindHandler = &handlers.CreateBindHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindDeleteBindHandler = &handlers.DeleteBindHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindGetBindHandler = &handlers.GetBindHandlerImpl{Client: client}
	api.BindGetBindsHandler = &handlers.GetBindsHandlerImpl{Client: client}
	api.BindReplaceBindHandler = &handlers.ReplaceBindHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http request rule handlers
	api.HTTPRequestRuleCreateHTTPRequestRuleHandler = &handlers.CreateHTTPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleDeleteHTTPRequestRuleHandler = &handlers.DeleteHTTPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleGetHTTPRequestRuleHandler = &handlers.GetHTTPRequestRuleHandlerImpl{Client: client}
	api.HTTPRequestRuleGetHTTPRequestRulesHandler = &handlers.GetHTTPRequestRulesHandlerImpl{Client: client}
	api.HTTPRequestRuleReplaceHTTPRequestRuleHandler = &handlers.ReplaceHTTPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http response rule handlers
	api.HTTPResponseRuleCreateHTTPResponseRuleHandler = &handlers.CreateHTTPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleDeleteHTTPResponseRuleHandler = &handlers.DeleteHTTPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleGetHTTPResponseRuleHandler = &handlers.GetHTTPResponseRuleHandlerImpl{Client: client}
	api.HTTPResponseRuleGetHTTPResponseRulesHandler = &handlers.GetHTTPResponseRulesHandlerImpl{Client: client}
	api.HTTPResponseRuleReplaceHTTPResponseRuleHandler = &handlers.ReplaceHTTPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp content rule handlers
	api.TCPRequestRuleCreateTCPRequestRuleHandler = &handlers.CreateTCPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleDeleteTCPRequestRuleHandler = &handlers.DeleteTCPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleGetTCPRequestRuleHandler = &handlers.GetTCPRequestRuleHandlerImpl{Client: client}
	api.TCPRequestRuleGetTCPRequestRulesHandler = &handlers.GetTCPRequestRulesHandlerImpl{Client: client}
	api.TCPRequestRuleReplaceTCPRequestRuleHandler = &handlers.ReplaceTCPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp connection rule handlers
	api.TCPResponseRuleCreateTCPResponseRuleHandler = &handlers.CreateTCPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleDeleteTCPResponseRuleHandler = &handlers.DeleteTCPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleGetTCPResponseRuleHandler = &handlers.GetTCPResponseRuleHandlerImpl{Client: client}
	api.TCPResponseRuleGetTCPResponseRulesHandler = &handlers.GetTCPResponseRulesHandlerImpl{Client: client}
	api.TCPResponseRuleReplaceTCPResponseRuleHandler = &handlers.ReplaceTCPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}

	// setup backend switching rule handlers
	api.BackendSwitchingRuleCreateBackendSwitchingRuleHandler = &handlers.CreateBackendSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.BackendSwitchingRuleDeleteBackendSwitchingRuleHandler = &handlers.DeleteBackendSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.BackendSwitchingRuleGetBackendSwitchingRuleHandler = &handlers.GetBackendSwitchingRuleHandlerImpl{Client: client}
	api.BackendSwitchingRuleGetBackendSwitchingRulesHandler = &handlers.GetBackendSwitchingRulesHandlerImpl{Client: client}
	api.BackendSwitchingRuleReplaceBackendSwitchingRuleHandler = &handlers.ReplaceBackendSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}

	// setup server switching rule handlers
	api.ServerSwitchingRuleCreateServerSwitchingRuleHandler = &handlers.CreateServerSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerSwitchingRuleDeleteServerSwitchingRuleHandler = &handlers.DeleteServerSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerSwitchingRuleGetServerSwitchingRuleHandler = &handlers.GetServerSwitchingRuleHandlerImpl{Client: client}
	api.ServerSwitchingRuleGetServerSwitchingRulesHandler = &handlers.GetServerSwitchingRulesHandlerImpl{Client: client}
	api.ServerSwitchingRuleReplaceServerSwitchingRuleHandler = &handlers.ReplaceServerSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}

	// setup filter handlers
	api.FilterCreateFilterHandler = &handlers.CreateFilterHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterDeleteFilterHandler = &handlers.DeleteFilterHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterGetFilterHandler = &handlers.GetFilterHandlerImpl{Client: client}
	api.FilterGetFiltersHandler = &handlers.GetFiltersHandlerImpl{Client: client}
	api.FilterReplaceFilterHandler = &handlers.ReplaceFilterHandlerImpl{Client: client, ReloadAgent: ra}

	// setup stick rule handlers
	api.StickRuleCreateStickRuleHandler = &handlers.CreateStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleDeleteStickRuleHandler = &handlers.DeleteStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleGetStickRuleHandler = &handlers.GetStickRuleHandlerImpl{Client: client}
	api.StickRuleGetStickRulesHandler = &handlers.GetStickRulesHandlerImpl{Client: client}
	api.StickRuleReplaceStickRuleHandler = &handlers.ReplaceStickRuleHandlerImpl{Client: client, ReloadAgent: ra}

	// setup log target handlers
	api.LogTargetCreateLogTargetHandler = &handlers.CreateLogTargetHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetHandler = &handlers.DeleteLogTargetHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetHandler = &handlers.GetLogTargetHandlerImpl{Client: client}
	api.LogTargetGetLogTargetsHandler = &handlers.GetLogTargetsHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetHandler = &handlers.ReplaceLogTargetHandlerImpl{Client: client, ReloadAgent: ra}

	// setup acl rule handlers
	api.ACLCreateACLHandler = &handlers.CreateACLHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLDeleteACLHandler = &handlers.DeleteACLHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLGetACLHandler = &handlers.GetACLHandlerImpl{Client: client}
	api.ACLGetAclsHandler = &handlers.GetAclsHandlerImpl{Client: client}
	api.ACLReplaceACLHandler = &handlers.ReplaceACLHandlerImpl{Client: client, ReloadAgent: ra}

	// setup resolvers handlers
	api.ResolverCreateResolverHandler = &handlers.CreateResolverHandlerImpl{Client: client, ReloadAgent: ra}
	api.ResolverDeleteResolverHandler = &handlers.DeleteResolverHandlerImpl{Client: client, ReloadAgent: ra}
	api.ResolverGetResolverHandler = &handlers.GetResolverHandlerImpl{Client: client}
	api.ResolverGetResolversHandler = &handlers.GetResolversHandlerImpl{Client: client}
	api.ResolverReplaceResolverHandler = &handlers.ReplaceResolverHandlerImpl{Client: client, ReloadAgent: ra}

	// setup nameserver handlers
	api.NameserverCreateNameserverHandler = &handlers.CreateNameserverHandlerImpl{Client: client, ReloadAgent: ra}
	api.NameserverDeleteNameserverHandler = &handlers.DeleteNameserverHandlerImpl{Client: client, ReloadAgent: ra}
	api.NameserverGetNameserverHandler = &handlers.GetNameserverHandlerImpl{Client: client}
	api.NameserverGetNameserversHandler = &handlers.GetNameserversHandlerImpl{Client: client}
	api.NameserverReplaceNameserverHandler = &handlers.ReplaceNameserverHandlerImpl{Client: client, ReloadAgent: ra}

	// setup peer section handlers
	api.PeerCreatePeerHandler = &handlers.CreatePeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.PeerDeletePeerHandler = &handlers.DeletePeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.PeerGetPeerSectionHandler = &handlers.GetPeerHandlerImpl{Client: client}
	api.PeerGetPeerSectionsHandler = &handlers.GetPeersHandlerImpl{Client: client}

	// setup peer entries handlers
	api.PeerEntryCreatePeerEntryHandler = &handlers.CreatePeerEntryHandlerImpl{Client: client, ReloadAgent: ra}
	api.PeerEntryDeletePeerEntryHandler = &handlers.DeletePeerEntryHandlerImpl{Client: client, ReloadAgent: ra}
	api.PeerEntryGetPeerEntryHandler = &handlers.GetPeerEntryHandlerImpl{Client: client}
	api.PeerEntryGetPeerEntriesHandler = &handlers.GetPeerEntriesHandlerImpl{Client: client}
	api.PeerEntryReplacePeerEntryHandler = &handlers.ReplacePeerEntryHandlerImpl{Client: client, ReloadAgent: ra}

	// setup stats handler
	api.StatsGetStatsHandler = &handlers.GetStatsHandlerImpl{Client: client}

	// setup info handler
	api.InformationGetHaproxyProcessInfoHandler = &handlers.GetHaproxyProcessInfoHandlerImpl{Client: client}

	// setup raw configuration handlers
	api.ConfigurationGetHAProxyConfigurationHandler = &handlers.GetRawConfigurationHandlerImpl{Client: client}
	api.ConfigurationPostHAProxyConfigurationHandler = &handlers.PostRawConfigurationHandlerImpl{Client: client, ReloadAgent: ra}

	// setup global configuration handlers
	api.GlobalGetGlobalHandler = &handlers.GetGlobalHandlerImpl{Client: client}
	api.GlobalReplaceGlobalHandler = &handlers.ReplaceGlobalHandlerImpl{Client: client, ReloadAgent: ra}

	// setup defaults configuration handlers
	api.DefaultsGetDefaultsHandler = &handlers.GetDefaultsHandlerImpl{Client: client}
	api.DefaultsReplaceDefaultsHandler = &handlers.ReplaceDefaultsHandlerImpl{Client: client, ReloadAgent: ra}

	// setup reload handlers
	api.ReloadsGetReloadHandler = &handlers.GetReloadHandlerImpl{ReloadAgent: ra}
	api.ReloadsGetReloadsHandler = &handlers.GetReloadsHandlerImpl{ReloadAgent: ra}

	// setup runtime server handlers
	api.ServerGetRuntimeServerHandler = &handlers.GetRuntimeServerHandlerImpl{Client: client}
	api.ServerGetRuntimeServersHandler = &handlers.GetRuntimeServersHandlerImpl{Client: client}
	api.ServerReplaceRuntimeServerHandler = &handlers.ReplaceRuntimeServerHandlerImpl{Client: client}

	// setup stick table handlers
	api.StickTableGetStickTablesHandler = &handlers.GetStickTablesHandlerImpl{Client: client}
	api.StickTableGetStickTableHandler = &handlers.GetStickTableHandlerImpl{Client: client}
	api.StickTableGetStickTableEntriesHandler = &handlers.GetStickTableEntriesHandlerImpl{Client: client}

	// setup ACL runtime handlers
	api.ACLRuntimeGetServicesHaproxyRuntimeAclsHandler = &handlers.GetACLSHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeGetServicesHaproxyRuntimeAclsIDHandler = &handlers.GetACLHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeGetServicesHaproxyRuntimeACLFileEntriesHandler = &handlers.GetACLFileEntriesHandlerRuntimeImpl{Client: client}
	api.ACLRuntimePostServicesHaproxyRuntimeACLFileEntriesHandler = &handlers.PostACLFileEntryHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeGetServicesHaproxyRuntimeACLFileEntriesIDHandler = &handlers.GetACLFileEntryRuntimeImpl{Client: client}
	api.ACLRuntimeDeleteServicesHaproxyRuntimeACLFileEntriesIDHandler = &handlers.DeleteACLFileEntryHandlerRuntimeImpl{Client: client}

	// setup map handlers
	api.MapsGetAllRuntimeMapFilesHandler = &handlers.GetMapsHandlerImpl{Client: client}
	api.MapsGetOneRuntimeMapHandler = &handlers.GetMapHandlerImpl{Client: client}
	api.MapsClearRuntimeMapHandler = &handlers.ClearMapHandlerImpl{Client: client}
	api.MapsShowRuntimeMapHandler = &handlers.ShowMapHandlerImpl{Client: client}
	api.MapsAddMapEntryHandler = &handlers.AddMapEntryHandlerImpl{Client: client}
	api.MapsGetRuntimeMapEntryHandler = &handlers.GetRuntimeMapEntryHandlerImpl{Client: client}
	api.MapsReplaceRuntimeMapEntryHandler = &handlers.ReplaceRuntimeMapEntryHandlerImpl{Client: client}
	api.MapsDeleteRuntimeMapEntryHandler = &handlers.DeleteRuntimeMapEntryHandlerImpl{Client: client}

	// setup info handler
	api.InformationGetInfoHandler = &handlers.GetInfoHandlerImpl{SystemInfo: haproxyOptions.ShowSystemInfo, BuildTime: BuildTime, Version: Version}

	// setup cluster handlers
	api.DiscoveryGetClusterHandler = &handlers.GetClusterHandlerImpl{Config: cfg}
	api.ClusterPostClusterHandler = &handlers.CreateClusterHandlerImpl{Client: client, Config: cfg, ReloadAgent: ra}
	api.ClusterInitiateCertificateRefreshHandler = &handlers.ClusterInitiateCertificateRefreshHandlerImpl{Config: cfg}

	clusterSync := dataplaneapi_config.ClusterSync{ReloadAgent: ra, Context: ctx}
	go clusterSync.Monitor(cfg, client)

	// setup specification handler
	api.SpecificationGetSpecificationHandler = specification.GetSpecificationHandlerFunc(func(params specification.GetSpecificationParams, principal interface{}) middleware.Responder {
		var m map[string]interface{}
		if err := json.Unmarshal(SwaggerJSON, &m); err != nil {
			e := misc.HandleError(err)
			return specification.NewGetSpecificationDefault(int(*e.Code)).WithPayload(e)
		}
		return specification.NewGetSpecificationOK().WithPayload(&m)
	})

	// set up service discovery handlers
	discovery := service_discovery.NewServiceDiscoveries(service_discovery.ServiceDiscoveriesParams{
		Client:      client.Configuration,
		ReloadAgent: ra,
		Context:     ctx,
	})
	api.ServiceDiscoveryCreateConsulHandler = &handlers.CreateConsulHandlerImpl{Discovery: discovery, PersistCallback: cfg.SaveConsuls}
	api.ServiceDiscoveryDeleteConsulHandler = &handlers.DeleteConsulHandlerImpl{Discovery: discovery, PersistCallback: cfg.SaveConsuls}
	api.ServiceDiscoveryGetConsulHandler = &handlers.GetConsulHandlerImpl{Discovery: discovery}
	api.ServiceDiscoveryGetConsulsHandler = &handlers.GetConsulsHandlerImpl{Discovery: discovery}
	api.ServiceDiscoveryReplaceConsulHandler = &handlers.ReplaceConsulHandlerImpl{Discovery: discovery, PersistCallback: cfg.SaveConsuls}

	api.ServiceDiscoveryCreateAWSRegionHandler = &handlers.CreateAWSHandlerImpl{Discovery: discovery, PersistCallback: cfg.SaveAWS}
	api.ServiceDiscoveryGetAWSRegionHandler = &handlers.GetAWSRegionHandlerImpl{Discovery: discovery}
	api.ServiceDiscoveryGetAWSRegionsHandler = &handlers.GetAWSRegionsHandlerImpl{Discovery: discovery}
	api.ServiceDiscoveryReplaceAWSRegionHandler = &handlers.ReplaceAWSRegionHandlerImpl{Discovery: discovery, PersistCallback: cfg.SaveAWS}
	api.ServiceDiscoveryDeleteAWSRegionHandler = &handlers.DeleteAWSRegionHandlerImpl{Discovery: discovery, PersistCallback: cfg.SaveAWS}

	// create stored consul instances
	for _, data := range cfg.ServiceDiscovery.Consuls {
		var err error

		if data.ID == nil || len(*data.ID) == 0 {
			data.ID = service_discovery.NewServiceDiscoveryUUID()
		}
		if err = service_discovery.ValidateConsulData(data, true); err != nil {
			log.Fatalf("Error validating Consul instance: " + err.Error())
		}
		if err = discovery.AddNode("consul", *data.ID, data); err != nil {
			log.Warning("Error creating consul instance: " + err.Error())
		}
	}
	_ = cfg.SaveConsuls(cfg.ServiceDiscovery.Consuls)

	// create stored AWS instances
	for _, data := range cfg.ServiceDiscovery.AWSRegions {
		var err error

		if data.ID == nil || len(*data.ID) == 0 {
			data.ID = service_discovery.NewServiceDiscoveryUUID()
		}
		if err = service_discovery.ValidateAWSData(data, true); err != nil {
			log.Fatalf("Error validating AWS instance: " + err.Error())
		}
		if err = discovery.AddNode("aws", *data.ID, data); err != nil {
			log.Warning("Error creating AWS instance: " + err.Error())
		}
	}
	_ = cfg.SaveAWS(cfg.ServiceDiscovery.AWSRegions)

	api.ConfigurationGetConfigurationVersionHandler = &handlers.ConfigurationGetConfigurationVersionHandlerImpl{Client: client}

	// map file storage handlers

	api.StorageCreateStorageMapFileHandler = &handlers.StorageCreateStorageMapFileHandlerImpl{Client: client}
	api.StorageGetAllStorageMapFilesHandler = &handlers.GetAllStorageMapFilesHandlerImpl{Client: client}
	api.StorageGetOneStorageMapHandler = &handlers.GetOneStorageMapHandlerImpl{Client: client}
	api.StorageDeleteStorageMapHandler = &handlers.StorageDeleteStorageMapHandlerImpl{Client: client}
	api.StorageReplaceStorageMapFileHandler = &handlers.StorageReplaceStorageMapFileHandlerImpl{Client: client, ReloadAgent: ra}

	// SSL certs file storage handlers
	api.StorageGetAllStorageSSLCertificatesHandler = &handlers.StorageGetAllStorageSSLCertificatesHandlerImpl{Client: client}
	api.StorageGetOneStorageSSLCertificateHandler = &handlers.StorageGetOneStorageSSLCertificateHandlerImpl{Client: client}
	api.StorageDeleteStorageSSLCertificateHandler = &handlers.StorageDeleteStorageSSLCertificateHandlerImpl{Client: client, ReloadAgent: ra}
	api.StorageReplaceStorageSSLCertificateHandler = &handlers.StorageReplaceStorageSSLCertificateHandlerImpl{Client: client, ReloadAgent: ra}
	api.StorageCreateStorageSSLCertificateHandler = &handlers.StorageCreateStorageSSLCertificateHandlerImpl{Client: client, ReloadAgent: ra}

	// setup OpenAPI v3 specification handler
	api.SpecificationOpenapiv3GetOpenapiv3SpecificationHandler = specification_openapiv3.GetOpenapiv3SpecificationHandlerFunc(func(params specification_openapiv3.GetOpenapiv3SpecificationParams, principal interface{}) middleware.Responder {
		v2 := openapi2.Swagger{}
		err := v2.UnmarshalJSON(SwaggerJSON)
		if err != nil {
			e := misc.HandleError(err)
			return specification_openapiv3.NewGetOpenapiv3SpecificationDefault(int(*e.Code)).WithPayload(e)
		}

		// if host is empty(dynamic hosts), server prop is empty,
		// so we need to set it explicitly
		if v2.Host == "" {
			cfg = dataplaneapi_config.Get()
			v2.Host = cfg.RuntimeData.Host
		}

		v3, err := openapi2conv.ToV3Swagger(&v2)
		if err != nil {
			e := misc.HandleError(err)
			return specification_openapiv3.NewGetOpenapiv3SpecificationDefault(int(*e.Code)).WithPayload(e)
		}
		return specification_openapiv3.NewGetOpenapiv3SpecificationOK().WithPayload(v3)
	})

	// TODO: do we need a ReloadAgent for SPOE
	// setup SPOE handlers
	api.SpoeCreateSpoeHandler = &handlers.SpoeCreateSpoeHandlerImpl{Client: client}
	api.SpoeDeleteSpoeFileHandler = &handlers.SpoeDeleteSpoeFileHandlerImpl{Client: client}
	api.SpoeGetAllSpoeFilesHandler = &handlers.SpoeGetAllSpoeFilesHandlerImpl{Client: client}
	api.SpoeGetOneSpoeFileHandler = &handlers.SpoeGetOneSpoeFileHandlerImpl{Client: client}

	// SPOE scope
	api.SpoeGetSpoeScopesHandler = &handlers.SpoeGetSpoeScopesHandlerImpl{Client: client}
	api.SpoeGetSpoeScopeHandler = &handlers.SpoeGetSpoeScopeHandlerImpl{Client: client}
	api.SpoeCreateSpoeScopeHandler = &handlers.SpoeCreateSpoeScopeHandlerImpl{Client: client}
	api.SpoeDeleteSpoeScopeHandler = &handlers.SpoeDeleteSpoeScopeHandlerImpl{Client: client}

	// SPOE agent
	api.SpoeGetSpoeAgentsHandler = &handlers.SpoeGetSpoeAgentsHandlerImpl{Client: client}
	api.SpoeGetSpoeAgentHandler = &handlers.SpoeGetSpoeAgentHandlerImpl{Client: client}
	api.SpoeCreateSpoeAgentHandler = &handlers.SpoeCreateSpoeAgentHandlerImpl{Client: client}
	api.SpoeDeleteSpoeAgentHandler = &handlers.SpoeDeleteSpoeAgentHandlerImpl{Client: client}
	api.SpoeReplaceSpoeAgentHandler = &handlers.SpoeReplaceSpoeAgentHandlerImpl{Client: client}

	// SPOE messages
	api.SpoeGetSpoeMessagesHandler = &handlers.SpoeGetSpoeMessagesHandlerImpl{Client: client}
	api.SpoeGetSpoeMessageHandler = &handlers.SpoeGetSpoeMessageHandlerImpl{Client: client}
	api.SpoeCreateSpoeMessageHandler = &handlers.SpoeCreateSpoeMessageHandlerImpl{Client: client}
	api.SpoeDeleteSpoeMessageHandler = &handlers.SpoeDeleteSpoeMessageHandlerImpl{Client: client}
	api.SpoeReplaceSpoeMessageHandler = &handlers.SpoeReplaceSpoeMessageHandlerImpl{Client: client}

	// SPOE groups
	api.SpoeGetSpoeGroupsHandler = &handlers.SpoeGetSpoeGroupsHandlerImpl{Client: client}
	api.SpoeGetSpoeGroupHandler = &handlers.SpoeGetSpoeGroupHandlerImpl{Client: client}
	api.SpoeCreateSpoeGroupHandler = &handlers.SpoeCreateSpoeGroupHandlerImpl{Client: client}
	api.SpoeDeleteSpoeGroupHandler = &handlers.SpoeDeleteSpoeGroupHandlerImpl{Client: client}
	api.SpoeReplaceSpoeGroupHandler = &handlers.SpoeReplaceSpoeGroupHandlerImpl{Client: client}

	// SPOE version
	api.SpoeGetSpoeConfigurationVersionHandler = &handlers.SpoeGetSpoeConfigurationVersionHandlerImpl{Client: client}

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Error starting Data Plane API: %s\n Stacktrace from panic: \n%s", err, string(debug.Stack()))
		}
	}()

	al, err := cfg.ApacheLogFormat()
	if err != nil {
		println("Cannot setup custom Apache Log Format", err.Error())
	}

	appLogger := log.StandardLogger()
	configureLogging(appLogger, cfg.Logging, func(opts dataplaneapi_config.SyslogOptions) dataplaneapi_config.SyslogOptions {
		opts.SyslogMsgID = "app"
		return opts
	}(cfg.Syslog))

	accLogger := log.New()
	configureLogging(accLogger, cfg.Logging, func(opts dataplaneapi_config.SyslogOptions) dataplaneapi_config.SyslogOptions {
		opts.SyslogMsgID = "accesslog"
		return opts
	}(cfg.Syslog))

	applicationEntry := log.NewEntry(appLogger)
	accessEntry := log.NewEntry(accLogger)

	// middlewares
	var adpts []adapters.Adapter
	adpts = append(adpts,
		adapters.RecoverMiddleware(applicationEntry),
		cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Reload-ID", "Configuration-Version"},
			AllowCredentials: true,
			MaxAge:           86400,
		}).Handler,
		adapters.UniqueIDMiddleware(applicationEntry),
		adapters.LoggingMiddleware(applicationEntry),
		adapters.ApacheLogMiddleware(accessEntry, al),
	)

	return setupGlobalMiddleware(api.Serve(setupMiddlewares), adpts...)
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler, adapters ...adapters.Adapter) http.Handler {
	for _, adpt := range adapters {
		handler = adpt(handler)
	}

	return handler
}

func configureLogging(logger *log.Logger, loggingOptions dataplaneapi_config.LoggingOptions, syslogOptions dataplaneapi_config.SyslogOptions) {
	switch loggingOptions.LogFormat {
	case "text":
		logger.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		})
	case "JSON":
		logger.SetFormatter(&log.JSONFormatter{})
	}

	switch loggingOptions.LogTo {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "file":
		dir := filepath.Dir(loggingOptions.LogFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			logger.Warning("Error opening log file, no logging implemented: " + err.Error())
		}
		//nolint:govet
		logFile, err := os.OpenFile(loggingOptions.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Warning("Error opening log file, no logging implemented: " + err.Error())
		}
		log.SetOutput(logFile)
	case "syslog":
		logger.SetOutput(ioutil.Discard)
		hook, err := syslog.NewRFC5424Hook(syslogOptions)
		if err != nil {
			logger.Warningf("Error configuring Syslog logging: %s", err.Error())
			break
		}
		logger.AddHook(hook)
	}

	switch loggingOptions.LogLevel {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warning":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	}
}

func serverShutdown() {
	cfg := dataplaneapi_config.Get()
	if logFile != nil {
		logFile.Close()
	}
	if cfg.HAProxy.UpdateMapFiles {
		cfg.MapSync.Stop()
	}
}

func configureNativeClient(cyx context.Context, haproxyOptions dataplaneapi_config.HAProxyConfiguration, mWorker bool) *client_native.HAProxyClient {

	// Initialize HAProxy native client
	confClient, err := configureConfigurationClient(haproxyOptions, mWorker)
	if err != nil {
		log.Fatalf("Error initializing configuration client: %v", err)
	}

	runtimeClient := configureRuntimeClient(cyx, confClient, haproxyOptions)
	client := &client_native.HAProxyClient{}
	if err = client.Init(confClient, runtimeClient); err != nil {
		log.Fatalf("Error setting up native client: %v", err)
	}

	if haproxyOptions.MapsDir != "" {
		client.MapStorage, err = storage.New(haproxyOptions.MapsDir, storage.MapsType)
		if err != nil {
			log.Fatalf("error initializing map storage: %v", err)
		}
	} else {
		log.Fatalf("error trying to use empty string for managed map directory")
	}

	if haproxyOptions.SSLCertsDir != "" {
		client.SSLCertStorage, err = storage.New(haproxyOptions.SSLCertsDir, storage.SSLType)
		if err != nil {
			log.Fatalf("error initializing SSL certs storage: %v", err)
		}
	} else {
		log.Fatalf("error trying to use empty string for managed map directory")
	}

	if haproxyOptions.SpoeDir != "" {
		prms := spoe.Params{
			SpoeDir:        haproxyOptions.SpoeDir,
			TransactionDir: haproxyOptions.SpoeTransactionDir,
		}
		client.Spoe, err = spoe.NewSpoe(prms)
		if err != nil {
			log.Fatalf("error setting up spoe: %v", err)
		}
	} else {
		log.Fatalf("error trying to use empty string for SPOE configuration directory")
	}

	return client
}

func configureConfigurationClient(haproxyOptions dataplaneapi_config.HAProxyConfiguration, mWorker bool) (*configuration.Client, error) {
	confClient := &configuration.Client{}
	confParams := configuration.ClientParams{
		ConfigurationFile:         haproxyOptions.ConfigFile,
		Haproxy:                   haproxyOptions.HAProxy,
		BackupsNumber:             haproxyOptions.BackupsNumber,
		UseValidation:             false,
		PersistentTransactions:    true,
		TransactionDir:            haproxyOptions.TransactionDir,
		ValidateCmd:               haproxyOptions.ValidateCmd,
		ValidateConfigurationFile: true,
		MasterWorker:              true,
		UseMd5Hash:                !haproxyOptions.DisableInotify,
	}

	err := confClient.Init(confParams)
	if err != nil {
		return nil, fmt.Errorf("error setting up configuration client: %s", err.Error())
	}

	p := confClient.Parser
	comments, err := p.Get(parser.Comments, parser.CommentsSectionName, "#")
	insertDisclaimer := false
	if err != nil {
		insertDisclaimer = true
	}
	data, ok := comments.([]types.Comments)
	if !ok {
		insertDisclaimer = true
	} else if len(data) == 0 || data[0].Value != "Dataplaneapi managed File" {
		insertDisclaimer = true
	}
	if insertDisclaimer {
		commentsNew := types.Comments{Value: "Dataplaneapi managed File"}
		err = p.Insert(parser.Comments, parser.CommentsSectionName, "#", commentsNew, 0)
		if err != nil {
			return nil, fmt.Errorf("error setting up configuration client: %s", err.Error())
		}
		commentsNew = types.Comments{Value: "changing file directly can cause a conflict if dataplaneapi is running"}
		err = p.Insert(parser.Comments, parser.CommentsSectionName, "#", commentsNew, 1)
		if err != nil {
			return nil, fmt.Errorf("error setting up configuration client: %s", err.Error())
		}
	}

	return confClient, nil
}

func configureRuntimeClient(ctx context.Context, confClient *configuration.Client, haproxyOptions dataplaneapi_config.HAProxyConfiguration) *runtime_api.Client {
	runtimeParams := runtime_api.ClientParams{
		MapsDir: haproxyOptions.MapsDir,
	}
	runtimeClient := &runtime_api.Client{ClientParams: runtimeParams}

	_, globalConf, err := confClient.GetGlobalConfiguration("")

	// First try to setup master runtime socket
	if err == nil {
		var err error
		// If master socket is set and a valid unix socket, use only this
		if haproxyOptions.MasterRuntime != "" && misc.IsUnixSocketAddr(haproxyOptions.MasterRuntime) {
			masterSocket := haproxyOptions.MasterRuntime
			// if nbproc is set, set nbproc sockets
			if globalConf.Nbproc > 0 {
				nbproc := int(globalConf.Nbproc)
				if err = runtimeClient.InitWithMasterSocketAndContext(ctx, masterSocket, nbproc); err == nil {
					return runtimeClient
				}
				log.Warningf("Error setting up runtime client with master socket: %s : %s", masterSocket, err.Error())
			} else {
				// if nbproc is not set, use master socket with 1 process
				if err = runtimeClient.InitWithMasterSocketAndContext(ctx, masterSocket, 1); err == nil {
					return runtimeClient
				}
				log.Warningf("Error setting up runtime client with master socket: %s : %s", masterSocket, err.Error())
			}
		}
		runtimeAPIs := globalConf.RuntimeAPIs
		// if no master socket set, read from first valid socket if nbproc <= 1
		if globalConf.Nbproc <= 1 {
			socketList := make(map[int]string)
			for _, r := range runtimeAPIs {
				if misc.IsUnixSocketAddr(*r.Address) {
					socketList[1] = *r.Address
					if err = runtimeClient.InitWithSocketsAndContext(ctx, socketList); err == nil {
						return runtimeClient
					}
					log.Warningf("Error setting up runtime client with socket: %s : %s", *r.Address, err.Error())
				}
			}
		} else {
			// else try to find process specific sockets and set them up
			sockets := make(map[int]string)
			for _, r := range runtimeAPIs {
				//nolint:govet
				if misc.IsUnixSocketAddr(*r.Address) && r.Process != "" {
					process, err := strconv.ParseInt(r.Process, 10, 64)
					if err == nil {
						sockets[int(process)] = *r.Address
					}
				}
			}
			// no process specific settings found, Issue a warning and return nil
			if len(sockets) == 0 {
				log.Warning("Runtime API not configured, found multiple processes and no stats sockets bound to them.")
				return nil
				// use only found process specific sockets issue a warning if not all processes have a socket configured
			}
			if len(sockets) < int(globalConf.Nbproc) {
				log.Warning("Runtime API not configured properly, there are more processes then configured sockets")
			}
			if err = runtimeClient.InitWithSocketsAndContext(ctx, sockets); err == nil {
				return runtimeClient
			}
			log.Warningf("Error setting up runtime client with sockets: %v : %s", sockets, err.Error())

		}
		if err != nil {
			log.Warning("Runtime API not configured, not using it: " + err.Error())
		} else {
			log.Warning("Runtime API not configured, not using it")
		}
		return nil
	}
	log.Warning("Cannot read runtime API configuration, not using it")
	return nil
}

func handleSignals(ctx context.Context, cancel context.CancelFunc, sigs chan os.Signal, client *client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users) {
	//nolint:gosimple
	for {
		select {
		case sig := <-sigs:
			if sig == syscall.SIGUSR1 {
				var clientCtx context.Context
				cancel()
				clientCtx, cancel = context.WithCancel(ctx)
				client.Runtime = configureRuntimeClient(clientCtx, client.Configuration, haproxyOptions)
				log.Info("Reloaded Data Plane API")
			} else if sig == syscall.SIGUSR2 {
				reloadConfigurationFile(client, haproxyOptions, users)
			}
		case <-ctx.Done():
			return
		}
	}
}

func reloadConfigurationFile(client *client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users) {
	confClient, err := configureConfigurationClient(haproxyOptions, mWorker)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if err := users.Init(); err != nil {
		log.Fatalf(err.Error())
	}
	log.Info("Rereading Configuration Files")
	client.Configuration = confClient
}

func startWatcher(ctx context.Context, client *client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users) error {
	cb := func() {
		reloadConfigurationFile(client, haproxyOptions, users)
		if err := client.Configuration.IncrementVersion(); err != nil {
			log.Warningf("Failed to increment configuration version: %v", err)
		}
	}

	watcher, err := dataplaneapi_config.NewConfigWatcher(dataplaneapi_config.ConfigWatcherParams{
		FilePath: haproxyOptions.ConfigFile,
		Callback: cb,
		Ctx:      ctx,
	})
	if err != nil {
		return err
	}
	go watcher.Listen()
	return nil
}
