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
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/client-native/v6/options"
	cn_runtime "github.com/haproxytech/client-native/v6/runtime"
	"github.com/haproxytech/client-native/v6/spoe"
	"github.com/haproxytech/client-native/v6/storage"
	"github.com/haproxytech/dataplaneapi/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/cors"

	"github.com/haproxytech/dataplaneapi/adapters"
	cn "github.com/haproxytech/dataplaneapi/client-native"
	dataplaneapi_config "github.com/haproxytech/dataplaneapi/configuration"
	service_discovery "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/handlers"
	"github.com/haproxytech/dataplaneapi/haproxy"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/operations"
	"github.com/haproxytech/dataplaneapi/operations/discovery"
	"github.com/haproxytech/dataplaneapi/operations/specification"
	"github.com/haproxytech/dataplaneapi/operations/version3"
	"github.com/haproxytech/dataplaneapi/rate"
	"github.com/haproxytech/dataplaneapi/resilient"
	socket_runtime "github.com/haproxytech/dataplaneapi/runtime"

	// import various crypting algorithms
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

//go:generate swagger generate server --target ../../../../../../github.com/haproxytech --name controller --spec ../../../../../../../../haproxy-api/haproxy-open-api-spec/build/haproxy_spec.yaml --server-package controller --tags Stats --tags Information --tags Configuration --tags Discovery --tags Frontend --tags Backend --tags Bind --tags Server --tags TCPRequestRule --tags HTTPRequestRule --tags HTTPResponseRule --tags Acl --tags BackendSwitchingRule --tags ServerSwitchingRule --tags TCPResponseRule --skip-models --exclude-main

var (
	Version               string
	BuildTime             string
	mWorker               = false
	logFile               *os.File
	AppLogger             *log.Logger
	AccLogger             *log.Logger
	serverStartedCallback func()
	clientMutex           sync.Mutex
)

func SetServerStartedCallback(callFunc func()) {
	serverStartedCallback = callFunc
}

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

func configureAPI(api *operations.DataPlaneAPI) http.Handler { //nolint:cyclop,maintidx
	clientMutex.Lock()
	defer clientMutex.Unlock()

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

	api.JSONConsumer = runtime.ConsumerFunc(func(reader io.Reader, data interface{}) error {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		dec := json.NewDecoder(reader)
		dec.UseNumber() // preserve number formats
		return dec.Decode(data)
	})

	api.TxtConsumer = runtime.TextConsumer()

	api.JSONProducer = runtime.ProducerFunc(func(writer io.Writer, data interface{}) error {
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		enc := json.NewEncoder(writer)
		enc.SetEscapeHTML(false)
		return enc.Encode(data)
	})

	api.ServerShutdown = serverShutdown

	ctx := ContextHandler.Context()
	clientCtx, cancel := context.WithCancel(ctx)

	client := configureNativeClient(clientCtx, haproxyOptions, mWorker)

	initDataplaneStorage(haproxyOptions.DataplaneStorageDir, client)

	users := dataplaneapi_config.GetUsersStore()
	// this is not part of GetUsersStore(),
	// in case of reload we need to reread users
	// mode might have changed from single to cluster one
	if err := users.Init(); err != nil {
		log.Fatalf("Error initiating users: %s", err.Error())
	}

	// Handle reload signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	go handleSignals(ctx, cancel, sigs, client, haproxyOptions, users)

	ra := configureReloadAgent(haproxyOptions, client, ctx)

	if !haproxyOptions.DisableInotify {
		if err := startWatcher(ctx, client, haproxyOptions, users, ra); err != nil {
			haproxyOptions.DisableInotify = true
			client = configureNativeClient(clientCtx, haproxyOptions, mWorker)
			ra = configureReloadAgent(haproxyOptions, client, ctx)
		}
	}

	// Sync map physical file with runtime map entries
	if haproxyOptions.UpdateMapFiles {
		go cfg.MapSync.SyncAll(client)
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
	client = resilient.NewClient(client)
	api.TransactionsStartTransactionHandler = &handlers.StartTransactionHandlerImpl{Client: client}
	api.TransactionsDeleteTransactionHandler = &handlers.DeleteTransactionHandlerImpl{Client: client}
	api.TransactionsGetTransactionHandler = &handlers.GetTransactionHandlerImpl{Client: client}
	api.TransactionsGetTransactionsHandler = &handlers.GetTransactionsHandlerImpl{Client: client}
	api.TransactionsCommitTransactionHandler = &handlers.CommitTransactionHandlerImpl{Client: client, ReloadAgent: ra, Mutex: &sync.Mutex{}}
	if cfg.HAProxy.MaxOpenTransactions > 0 {
		// creating the threshold limit using the CLI flag as hard quota and current open transactions as starting point
		actualCount := func() uint64 {
			configuration, err := client.Configuration()
			if err != nil {
				log.Errorf("Cannot retrieve current open transactions for rate limit, default to zero (%s)", err.Error())
				return 0
			}
			ts, err := configuration.GetTransactions(models.TransactionStatusInProgress)
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
	api.SpoeTransactionsGetAllSpoeTransactionHandler = &handlers.SpoeTransactionsGetAllSpoeTransactionHandlerImpl{Client: client}
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

	// setup ring handlers
	api.RingCreateRingHandler = &handlers.CreateRingHandlerImpl{Client: client, ReloadAgent: ra}
	api.RingDeleteRingHandler = &handlers.DeleteRingHandlerImpl{Client: client, ReloadAgent: ra}
	api.RingGetRingHandler = &handlers.GetRingHandlerImpl{Client: client}
	api.RingGetRingsHandler = &handlers.GetRingsHandlerImpl{Client: client}
	api.RingReplaceRingHandler = &handlers.ReplaceRingHandlerImpl{Client: client, ReloadAgent: ra}

	// setup log forward handlers
	api.LogForwardCreateLogForwardHandler = &handlers.CreateLogForwardHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogForwardDeleteLogForwardHandler = &handlers.DeleteLogForwardHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogForwardGetLogForwardHandler = &handlers.GetLogForwardHandlerImpl{Client: client}
	api.LogForwardGetLogForwardsHandler = &handlers.GetLogForwardsHandlerImpl{Client: client}
	api.LogForwardReplaceLogForwardHandler = &handlers.ReplaceLogForwardHandlerImpl{Client: client, ReloadAgent: ra}

	// setup dgram bind handlers
	api.DgramBindCreateDgramBindHandler = &handlers.CreateDgramBindHandlerImpl{Client: client, ReloadAgent: ra}
	api.DgramBindDeleteDgramBindHandler = &handlers.DeleteDgramBindHandlerImpl{Client: client, ReloadAgent: ra}
	api.DgramBindGetDgramBindHandler = &handlers.GetDgramBindHandlerImpl{Client: client}
	api.DgramBindGetDgramBindsHandler = &handlers.GetDgramBindsHandlerImpl{Client: client}
	api.DgramBindReplaceDgramBindHandler = &handlers.ReplaceDgramBindHandlerImpl{Client: client, ReloadAgent: ra}

	// setup frontend handlers
	api.FrontendCreateFrontendHandler = &handlers.CreateFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FrontendDeleteFrontendHandler = &handlers.DeleteFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FrontendGetFrontendHandler = &handlers.GetFrontendHandlerImpl{Client: client}
	api.FrontendGetFrontendsHandler = &handlers.GetFrontendsHandlerImpl{Client: client}
	api.FrontendReplaceFrontendHandler = &handlers.ReplaceFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup userlist handlers
	api.UserlistCreateUserlistHandler = &handlers.CreateUserListHandlerImpl{Client: client, ReloadAgent: ra}
	api.UserlistDeleteUserlistHandler = &handlers.DeleteUserListHandlerImpl{Client: client, ReloadAgent: ra}
	api.UserlistGetUserlistHandler = &handlers.GetUserListHandlerImpl{Client: client}
	api.UserlistGetUserlistsHandler = &handlers.GetUserListsHandlerImpl{Client: client}

	// setup user handlers
	api.UserCreateUserHandler = &handlers.CreateUserHandlerImpl{Client: client, ReloadAgent: ra}
	api.UserDeleteUserHandler = &handlers.DeleteUserHandlerImpl{Client: client, ReloadAgent: ra}
	api.UserGetUserHandler = &handlers.GetUserHandlerImpl{Client: client}
	api.UserGetUsersHandler = &handlers.GetUsersHandlerImpl{Client: client}
	api.UserReplaceUserHandler = &handlers.ReplaceUserHandlerImpl{Client: client, ReloadAgent: ra}

	// setup group handlers
	api.GroupCreateGroupHandler = &handlers.CreateGroupHandlerImpl{Client: client, ReloadAgent: ra}
	api.GroupDeleteGroupHandler = &handlers.DeleteGroupHandlerImpl{Client: client, ReloadAgent: ra}
	api.GroupGetGroupHandler = &handlers.GetGroupHandlerImpl{Client: client}
	api.GroupGetGroupsHandler = &handlers.GetGroupsHandlerImpl{Client: client}
	api.GroupReplaceGroupHandler = &handlers.ReplaceGroupHandlerImpl{Client: client, ReloadAgent: ra}

	// setup server handlers
	// Create
	api.ServerCreateServerBackendHandler = &handlers.CreateServerBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerCreateServerPeerHandler = &handlers.CreateServerPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerCreateServerRingHandler = &handlers.CreateServerRingHandlerImpl{Client: client, ReloadAgent: ra}
	// Get all
	api.ServerGetAllServerBackendHandler = &handlers.GetAllServerBackendHandlerImpl{Client: client}
	api.ServerGetAllServerPeerHandler = &handlers.GetAllServerPeerHandlerImpl{Client: client}
	api.ServerGetAllServerRingHandler = &handlers.GetAllServerRingHandlerImpl{Client: client}
	// Delete one
	api.ServerDeleteServerBackendHandler = &handlers.DeleteServerBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerDeleteServerPeerHandler = &handlers.DeleteServerPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerDeleteServerRingHandler = &handlers.DeleteServerRingHandlerImpl{Client: client, ReloadAgent: ra}
	// Get one
	api.ServerGetServerBackendHandler = &handlers.GetServerBackendHandlerImpl{Client: client}
	api.ServerGetServerPeerHandler = &handlers.GetServerPeerHandlerImpl{Client: client}
	api.ServerGetServerRingHandler = &handlers.GetServerRingHandlerImpl{Client: client}
	// Replace one
	api.ServerReplaceServerBackendHandler = &handlers.ReplaceServerBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerReplaceServerPeerHandler = &handlers.ReplaceServerPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerReplaceServerRingHandler = &handlers.ReplaceServerRingHandlerImpl{Client: client, ReloadAgent: ra}

	// setup server template handlers
	api.ServerTemplateCreateServerTemplateHandler = &handlers.CreateServerTemplateHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerTemplateDeleteServerTemplateHandler = &handlers.DeleteServerTemplateHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerTemplateGetServerTemplateHandler = &handlers.GetServerTemplateHandlerImpl{Client: client}
	api.ServerTemplateGetServerTemplatesHandler = &handlers.GetServerTemplatesHandlerImpl{Client: client}
	api.ServerTemplateReplaceServerTemplateHandler = &handlers.ReplaceServerTemplateHandlerImpl{Client: client, ReloadAgent: ra}

	// setup bind handlers
	api.BindCreateBindFrontendHandler = &handlers.CreateBindFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindDeleteBindFrontendHandler = &handlers.DeleteBindFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindGetBindFrontendHandler = &handlers.GetBindFrontendHandlerImpl{Client: client}
	api.BindGetAllBindFrontendHandler = &handlers.GetAllBindFrontendHandlerImpl{Client: client}
	api.BindReplaceBindFrontendHandler = &handlers.ReplaceBindFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	api.BindCreateBindPeerHandler = &handlers.CreateBindPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindDeleteBindPeerHandler = &handlers.DeleteBindPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindGetBindPeerHandler = &handlers.GetBindPeerHandlerImpl{Client: client}
	api.BindGetAllBindPeerHandler = &handlers.GetAllBindPeerHandlerImpl{Client: client}
	api.BindReplaceBindPeerHandler = &handlers.ReplaceBindPeerHandlerImpl{Client: client, ReloadAgent: ra}

	api.BindCreateBindLogForwardHandler = &handlers.CreateBindLogForwardHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindDeleteBindLogForwardHandler = &handlers.DeleteBindLogForwardHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindGetBindLogForwardHandler = &handlers.GetBindLogForwardHandlerImpl{Client: client}
	api.BindGetAllBindLogForwardHandler = &handlers.GetAllBindLogForwardHandlerImpl{Client: client}
	api.BindReplaceBindLogForwardHandler = &handlers.ReplaceBindLogForwardHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http check handlers
	api.HTTPCheckGetHTTPCheckBackendHandler = &handlers.GetHTTPCheckBackendHandlerImpl{Client: client}
	api.HTTPCheckGetAllHTTPCheckBackendHandler = &handlers.GetAllHTTPCheckBackendHandlerImpl{Client: client}
	api.HTTPCheckCreateHTTPCheckBackendHandler = &handlers.CreateHTTPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckReplaceHTTPCheckBackendHandler = &handlers.ReplaceHTTPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckDeleteHTTPCheckBackendHandler = &handlers.DeleteHTTPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckReplaceAllHTTPCheckBackendHandler = &handlers.ReplaceAllHTTPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckGetHTTPCheckDefaultsHandler = &handlers.GetHTTPCheckDefaultsHandlerImpl{Client: client}
	api.HTTPCheckGetAllHTTPCheckDefaultsHandler = &handlers.GetAllHTTPCheckDefaultsHandlerImpl{Client: client}
	api.HTTPCheckCreateHTTPCheckDefaultsHandler = &handlers.CreateHTTPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckReplaceHTTPCheckDefaultsHandler = &handlers.ReplaceHTTPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckDeleteHTTPCheckDefaultsHandler = &handlers.DeleteHTTPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckReplaceAllHTTPCheckDefaultsHandler = &handlers.ReplaceAllHTTPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http request rule handlers
	api.HTTPRequestRuleCreateHTTPRequestRuleBackendHandler = &handlers.CreateHTTPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleDeleteHTTPRequestRuleBackendHandler = &handlers.DeleteHTTPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleGetHTTPRequestRuleBackendHandler = &handlers.GetHTTPRequestRuleBackendHandlerImpl{Client: client}
	api.HTTPRequestRuleGetAllHTTPRequestRuleBackendHandler = &handlers.GetAllHTTPRequestRuleBackendHandlerImpl{Client: client}
	api.HTTPRequestRuleReplaceHTTPRequestRuleBackendHandler = &handlers.ReplaceHTTPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleReplaceAllHTTPRequestRuleBackendHandler = &handlers.ReplaceAllHTTPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleCreateHTTPRequestRuleFrontendHandler = &handlers.CreateHTTPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleDeleteHTTPRequestRuleFrontendHandler = &handlers.DeleteHTTPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleGetHTTPRequestRuleFrontendHandler = &handlers.GetHTTPRequestRuleFrontendHandlerImpl{Client: client}
	api.HTTPRequestRuleGetAllHTTPRequestRuleFrontendHandler = &handlers.GetAllHTTPRequestRuleFrontendHandlerImpl{Client: client}
	api.HTTPRequestRuleReplaceHTTPRequestRuleFrontendHandler = &handlers.ReplaceHTTPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleReplaceAllHTTPRequestRuleFrontendHandler = &handlers.ReplaceAllHTTPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http after response rule handlers
	api.HTTPAfterResponseRuleCreateHTTPAfterResponseRuleBackendHandler = &handlers.CreateHTTPAfterResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleDeleteHTTPAfterResponseRuleBackendHandler = &handlers.DeleteHTTPAfterResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleGetHTTPAfterResponseRuleBackendHandler = &handlers.GetHTTPAfterResponseRuleBackendHandlerImpl{Client: client}
	api.HTTPAfterResponseRuleGetAllHTTPAfterResponseRuleBackendHandler = &handlers.GetAllHTTPAfterResponseRuleBackendHandlerImpl{Client: client}
	api.HTTPAfterResponseRuleReplaceHTTPAfterResponseRuleBackendHandler = &handlers.ReplaceHTTPAfterResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleReplaceAllHTTPAfterResponseRuleBackendHandler = &handlers.ReplaceAllHTTPAfterResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.HTTPAfterResponseRuleCreateHTTPAfterResponseRuleFrontendHandler = &handlers.CreateHTTPAfterResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleDeleteHTTPAfterResponseRuleFrontendHandler = &handlers.DeleteHTTPAfterResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleGetHTTPAfterResponseRuleFrontendHandler = &handlers.GetHTTPAfterResponseRuleFrontendHandlerImpl{Client: client}
	api.HTTPAfterResponseRuleGetAllHTTPAfterResponseRuleFrontendHandler = &handlers.GetAllHTTPAfterResponseRuleFrontendHandlerImpl{Client: client}
	api.HTTPAfterResponseRuleReplaceHTTPAfterResponseRuleFrontendHandler = &handlers.ReplaceHTTPAfterResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleReplaceAllHTTPAfterResponseRuleFrontendHandler = &handlers.ReplaceAllHTTPAfterResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http response rule handlers
	api.HTTPResponseRuleCreateHTTPResponseRuleBackendHandler = &handlers.CreateHTTPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleDeleteHTTPResponseRuleBackendHandler = &handlers.DeleteHTTPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleGetHTTPResponseRuleBackendHandler = &handlers.GetHTTPResponseRuleBackendHandlerImpl{Client: client}
	api.HTTPResponseRuleGetAllHTTPResponseRuleBackendHandler = &handlers.GetAllHTTPResponseRuleBackendHandlerImpl{Client: client}
	api.HTTPResponseRuleReplaceHTTPResponseRuleBackendHandler = &handlers.ReplaceHTTPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleReplaceAllHTTPResponseRuleBackendHandler = &handlers.ReplaceAllHTTPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.HTTPResponseRuleCreateHTTPResponseRuleFrontendHandler = &handlers.CreateHTTPResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleDeleteHTTPResponseRuleFrontendHandler = &handlers.DeleteHTTPResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleGetHTTPResponseRuleFrontendHandler = &handlers.GetHTTPResponseRuleFrontendHandlerImpl{Client: client}
	api.HTTPResponseRuleGetAllHTTPResponseRuleFrontendHandler = &handlers.GetAllHTTPResponseRuleFrontendHandlerImpl{Client: client}
	api.HTTPResponseRuleReplaceHTTPResponseRuleFrontendHandler = &handlers.ReplaceHTTPResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleReplaceAllHTTPResponseRuleFrontendHandler = &handlers.ReplaceAllHTTPResponseRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http error rule handlers
	api.HTTPErrorRuleCreateHTTPErrorRuleBackendHandler = &handlers.CreateHTTPErrorRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleDeleteHTTPErrorRuleBackendHandler = &handlers.DeleteHTTPErrorRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleGetHTTPErrorRuleBackendHandler = &handlers.GetHTTPErrorRuleBackendHandlerImpl{Client: client}
	api.HTTPErrorRuleGetAllHTTPErrorRuleBackendHandler = &handlers.GetAllHTTPErrorRuleBackendHandlerImpl{Client: client}
	api.HTTPErrorRuleReplaceHTTPErrorRuleBackendHandler = &handlers.ReplaceHTTPErrorRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleReplaceAllHTTPErrorRuleBackendHandler = &handlers.ReplaceAllHTTPErrorRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.HTTPErrorRuleCreateHTTPErrorRuleFrontendHandler = &handlers.CreateHTTPErrorRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleDeleteHTTPErrorRuleFrontendHandler = &handlers.DeleteHTTPErrorRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleGetHTTPErrorRuleFrontendHandler = &handlers.GetHTTPErrorRuleFrontendHandlerImpl{Client: client}
	api.HTTPErrorRuleGetAllHTTPErrorRuleFrontendHandler = &handlers.GetAllHTTPErrorRuleFrontendHandlerImpl{Client: client}
	api.HTTPErrorRuleReplaceHTTPErrorRuleFrontendHandler = &handlers.ReplaceHTTPErrorRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleReplaceAllHTTPErrorRuleFrontendHandler = &handlers.ReplaceAllHTTPErrorRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	api.HTTPErrorRuleCreateHTTPErrorRuleDefaultsHandler = &handlers.CreateHTTPErrorRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleDeleteHTTPErrorRuleDefaultsHandler = &handlers.DeleteHTTPErrorRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleGetHTTPErrorRuleDefaultsHandler = &handlers.GetHTTPErrorRuleDefaultsHandlerImpl{Client: client}
	api.HTTPErrorRuleGetAllHTTPErrorRuleDefaultsHandler = &handlers.GetAllHTTPErrorRuleDefaultsHandlerImpl{Client: client}
	api.HTTPErrorRuleReplaceHTTPErrorRuleDefaultsHandler = &handlers.ReplaceHTTPErrorRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleReplaceAllHTTPErrorRuleDefaultsHandler = &handlers.ReplaceAllHTTPErrorRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp content rule handlers
	api.TCPRequestRuleCreateTCPRequestRuleBackendHandler = &handlers.CreateTCPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleDeleteTCPRequestRuleBackendHandler = &handlers.DeleteTCPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleGetTCPRequestRuleBackendHandler = &handlers.GetTCPRequestRuleBackendHandlerImpl{Client: client}
	api.TCPRequestRuleGetAllTCPRequestRuleBackendHandler = &handlers.GetAllTCPRequestRuleBackendHandlerImpl{Client: client}
	api.TCPRequestRuleReplaceTCPRequestRuleBackendHandler = &handlers.ReplaceTCPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleReplaceAllTCPRequestRuleBackendHandler = &handlers.ReplaceAllTCPRequestRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.TCPRequestRuleCreateTCPRequestRuleFrontendHandler = &handlers.CreateTCPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleDeleteTCPRequestRuleFrontendHandler = &handlers.DeleteTCPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleGetTCPRequestRuleFrontendHandler = &handlers.GetTCPRequestRuleFrontendHandlerImpl{Client: client}
	api.TCPRequestRuleGetAllTCPRequestRuleFrontendHandler = &handlers.GetAllTCPRequestRuleFrontendHandlerImpl{Client: client}
	api.TCPRequestRuleReplaceTCPRequestRuleFrontendHandler = &handlers.ReplaceTCPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleReplaceAllTCPRequestRuleFrontendHandler = &handlers.ReplaceAllTCPRequestRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp connection rule handlers
	api.TCPResponseRuleCreateTCPResponseRuleBackendHandler = &handlers.CreateTCPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleDeleteTCPResponseRuleBackendHandler = &handlers.DeleteTCPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleGetTCPResponseRuleBackendHandler = &handlers.GetTCPResponseRuleBackendHandlerImpl{Client: client}
	api.TCPResponseRuleGetAllTCPResponseRuleBackendHandler = &handlers.GetAllTCPResponseRuleBackendHandlerImpl{Client: client}
	api.TCPResponseRuleReplaceTCPResponseRuleBackendHandler = &handlers.ReplaceTCPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleReplaceAllTCPResponseRuleBackendHandler = &handlers.ReplaceAllTCPResponseRuleBackendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp check handlers
	api.TCPCheckCreateTCPCheckBackendHandler = &handlers.CreateTCPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckDeleteTCPCheckBackendHandler = &handlers.DeleteTCPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckGetTCPCheckBackendHandler = &handlers.GetTCPCheckBackendHandlerImpl{Client: client}
	api.TCPCheckGetAllTCPCheckBackendHandler = &handlers.GetAllTCPCheckBackendHandlerImpl{Client: client}
	api.TCPCheckReplaceTCPCheckBackendHandler = &handlers.ReplaceTCPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckReplaceAllTCPCheckBackendHandler = &handlers.ReplaceAllTCPCheckBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.TCPCheckCreateTCPCheckDefaultsHandler = &handlers.CreateTCPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckDeleteTCPCheckDefaultsHandler = &handlers.DeleteTCPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckGetTCPCheckDefaultsHandler = &handlers.GetTCPCheckDefaultsHandlerImpl{Client: client}
	api.TCPCheckGetAllTCPCheckDefaultsHandler = &handlers.GetAllTCPCheckDefaultsHandlerImpl{Client: client}
	api.TCPCheckReplaceTCPCheckDefaultsHandler = &handlers.ReplaceTCPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckReplaceAllTCPCheckDefaultsHandler = &handlers.ReplaceAllTCPCheckDefaultsHandlerImpl{Client: client, ReloadAgent: ra}

	// setup quic initia; rule handlers
	api.QUICInitialRuleCreateQUICInitialRuleDefaultsHandler = &handlers.CreateQUICInitialRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.QUICInitialRuleDeleteQUICInitialRuleDefaultsHandler = &handlers.DeleteQUICInitialRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.QUICInitialRuleGetQUICInitialRuleDefaultsHandler = &handlers.GetQUICInitialRuleDefaultsHandlerImpl{Client: client}
	api.QUICInitialRuleGetAllQUICInitialRuleDefaultsHandler = &handlers.GetAllQUICInitialRuleDefaultsHandlerImpl{Client: client}
	api.QUICInitialRuleReplaceQUICInitialRuleDefaultsHandler = &handlers.ReplaceQUICInitialRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.QUICInitialRuleReplaceAllQUICInitialRuleDefaultsHandler = &handlers.ReplaceAllQUICInitialRuleDefaultsHandlerImpl{Client: client, ReloadAgent: ra}

	api.QUICInitialRuleCreateQUICInitialRuleFrontendHandler = &handlers.CreateQUICInitialRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.QUICInitialRuleDeleteQUICInitialRuleFrontendHandler = &handlers.DeleteQUICInitialRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.QUICInitialRuleGetQUICInitialRuleFrontendHandler = &handlers.GetQUICInitialRuleFrontendHandlerImpl{Client: client}
	api.QUICInitialRuleGetAllQUICInitialRuleFrontendHandler = &handlers.GetAllQUICInitialRuleFrontendHandlerImpl{Client: client}
	api.QUICInitialRuleReplaceQUICInitialRuleFrontendHandler = &handlers.ReplaceQUICInitialRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.QUICInitialRuleReplaceAllQUICInitialRuleFrontendHandler = &handlers.ReplaceAllQUICInitialRuleFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup declare capture handlers
	api.DeclareCaptureCreateDeclareCaptureHandler = &handlers.CreateDeclareCaptureHandlerImpl{Client: client, ReloadAgent: ra}
	api.DeclareCaptureDeleteDeclareCaptureHandler = &handlers.DeleteDeclareCaptureHandlerImpl{Client: client, ReloadAgent: ra}
	api.DeclareCaptureGetDeclareCaptureHandler = &handlers.GetDeclareCaptureHandlerImpl{Client: client}
	api.DeclareCaptureGetDeclareCapturesHandler = &handlers.GetDeclareCapturesHandlerImpl{Client: client}
	api.DeclareCaptureReplaceDeclareCaptureHandler = &handlers.ReplaceDeclareCaptureHandlerImpl{Client: client, ReloadAgent: ra}
	api.DeclareCaptureReplaceDeclareCapturesHandler = &handlers.ReplaceDeclareCapturesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup backend switching rule handlers
	api.BackendSwitchingRuleCreateBackendSwitchingRuleHandler = &handlers.CreateBackendSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.BackendSwitchingRuleDeleteBackendSwitchingRuleHandler = &handlers.DeleteBackendSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.BackendSwitchingRuleGetBackendSwitchingRuleHandler = &handlers.GetBackendSwitchingRuleHandlerImpl{Client: client}
	api.BackendSwitchingRuleGetBackendSwitchingRulesHandler = &handlers.GetBackendSwitchingRulesHandlerImpl{Client: client}
	api.BackendSwitchingRuleReplaceBackendSwitchingRuleHandler = &handlers.ReplaceBackendSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.BackendSwitchingRuleReplaceBackendSwitchingRulesHandler = &handlers.ReplaceBackendSwitchingRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup server switching rule handlers
	api.ServerSwitchingRuleCreateServerSwitchingRuleHandler = &handlers.CreateServerSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerSwitchingRuleDeleteServerSwitchingRuleHandler = &handlers.DeleteServerSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerSwitchingRuleGetServerSwitchingRuleHandler = &handlers.GetServerSwitchingRuleHandlerImpl{Client: client}
	api.ServerSwitchingRuleGetServerSwitchingRulesHandler = &handlers.GetServerSwitchingRulesHandlerImpl{Client: client}
	api.ServerSwitchingRuleReplaceServerSwitchingRuleHandler = &handlers.ReplaceServerSwitchingRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerSwitchingRuleReplaceServerSwitchingRulesHandler = &handlers.ReplaceServerSwitchingRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup filter handlers
	api.FilterCreateFilterBackendHandler = &handlers.CreateFilterBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterDeleteFilterBackendHandler = &handlers.DeleteFilterBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterGetFilterBackendHandler = &handlers.GetFilterBackendHandlerImpl{Client: client}
	api.FilterGetAllFilterBackendHandler = &handlers.GetAllFilterBackendHandlerImpl{Client: client}
	api.FilterReplaceFilterBackendHandler = &handlers.ReplaceFilterBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterReplaceAllFilterBackendHandler = &handlers.ReplaceAllFilterBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.FilterCreateFilterFrontendHandler = &handlers.CreateFilterFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterDeleteFilterFrontendHandler = &handlers.DeleteFilterFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterGetFilterFrontendHandler = &handlers.GetFilterFrontendHandlerImpl{Client: client}
	api.FilterGetAllFilterFrontendHandler = &handlers.GetAllFilterFrontendHandlerImpl{Client: client}
	api.FilterReplaceFilterFrontendHandler = &handlers.ReplaceFilterFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterReplaceAllFilterFrontendHandler = &handlers.ReplaceAllFilterFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	// setup stick rule handlers
	api.StickRuleCreateStickRuleHandler = &handlers.CreateStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleDeleteStickRuleHandler = &handlers.DeleteStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleGetStickRuleHandler = &handlers.GetStickRuleHandlerImpl{Client: client}
	api.StickRuleGetStickRulesHandler = &handlers.GetStickRulesHandlerImpl{Client: client}
	api.StickRuleReplaceStickRuleHandler = &handlers.ReplaceStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleReplaceStickRulesHandler = &handlers.ReplaceStickRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup log target handlers
	api.LogTargetCreateLogTargetBackendHandler = &handlers.CreateLogTargetBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetBackendHandler = &handlers.DeleteLogTargetBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetBackendHandler = &handlers.GetLogTargetBackendHandlerImpl{Client: client}
	api.LogTargetGetAllLogTargetBackendHandler = &handlers.GetAllLogTargetBackendHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetBackendHandler = &handlers.ReplaceLogTargetBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetReplaceAllLogTargetBackendHandler = &handlers.ReplaceAllLogTargetBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.LogTargetCreateLogTargetFrontendHandler = &handlers.CreateLogTargetFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetFrontendHandler = &handlers.DeleteLogTargetFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetFrontendHandler = &handlers.GetLogTargetFrontendHandlerImpl{Client: client}
	api.LogTargetGetAllLogTargetFrontendHandler = &handlers.GetAllLogTargetFrontendHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetFrontendHandler = &handlers.ReplaceLogTargetFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetReplaceAllLogTargetFrontendHandler = &handlers.ReplaceAllLogTargetFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	api.LogTargetCreateLogTargetDefaultsHandler = &handlers.CreateLogTargetDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetDefaultsHandler = &handlers.DeleteLogTargetDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetDefaultsHandler = &handlers.GetLogTargetDefaultsHandlerImpl{Client: client}
	api.LogTargetGetAllLogTargetDefaultsHandler = &handlers.GetAllLogTargetDefaultsHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetDefaultsHandler = &handlers.ReplaceLogTargetDefaultsHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetReplaceAllLogTargetDefaultsHandler = &handlers.ReplaceAllLogTargetDefaultsHandlerImpl{Client: client, ReloadAgent: ra}

	api.LogTargetCreateLogTargetLogForwardHandler = &handlers.CreateLogTargetLogForwardHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetLogForwardHandler = &handlers.DeleteLogTargetLogForwardHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetLogForwardHandler = &handlers.GetLogTargetLogForwardHandlerImpl{Client: client}
	api.LogTargetGetAllLogTargetLogForwardHandler = &handlers.GetAllLogTargetLogForwardHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetLogForwardHandler = &handlers.ReplaceLogTargetLogForwardHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetReplaceAllLogTargetLogForwardHandler = &handlers.ReplaceAllLogTargetLogForwardHandlerImpl{Client: client, ReloadAgent: ra}

	api.LogTargetCreateLogTargetPeerHandler = &handlers.CreateLogTargetPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetPeerHandler = &handlers.DeleteLogTargetPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetPeerHandler = &handlers.GetLogTargetPeerHandlerImpl{Client: client}
	api.LogTargetGetAllLogTargetPeerHandler = &handlers.GetAllLogTargetPeerHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetPeerHandler = &handlers.ReplaceLogTargetPeerHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetReplaceAllLogTargetPeerHandler = &handlers.ReplaceAllLogTargetPeerHandlerImpl{Client: client, ReloadAgent: ra}

	api.LogTargetCreateLogTargetGlobalHandler = &handlers.CreateLogTargetGlobalHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetGlobalHandler = &handlers.DeleteLogTargetGlobalHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetGlobalHandler = &handlers.GetLogTargetGlobalHandlerImpl{Client: client}
	api.LogTargetGetAllLogTargetGlobalHandler = &handlers.GetAllLogTargetGlobalHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetGlobalHandler = &handlers.ReplaceLogTargetGlobalHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetReplaceAllLogTargetGlobalHandler = &handlers.ReplaceAllLogTargetGlobalHandlerImpl{Client: client, ReloadAgent: ra}

	// setup acl rule handlers
	api.ACLCreateACLBackendHandler = &handlers.CreateACLBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLDeleteACLBackendHandler = &handlers.DeleteACLBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLGetACLBackendHandler = &handlers.GetACLBackendHandlerImpl{Client: client}
	api.ACLGetAllACLBackendHandler = &handlers.GetAllACLBackendHandlerImpl{Client: client}
	api.ACLReplaceACLBackendHandler = &handlers.ReplaceACLBackendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLReplaceAllACLBackendHandler = &handlers.ReplaceAllACLBackendHandlerImpl{Client: client, ReloadAgent: ra}

	api.ACLCreateACLFrontendHandler = &handlers.CreateACLFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLDeleteACLFrontendHandler = &handlers.DeleteACLFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLGetACLFrontendHandler = &handlers.GetACLFrontendHandlerImpl{Client: client}
	api.ACLGetAllACLFrontendHandler = &handlers.GetAllACLFrontendHandlerImpl{Client: client}
	api.ACLReplaceACLFrontendHandler = &handlers.ReplaceACLFrontendHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLReplaceAllACLFrontendHandler = &handlers.ReplaceAllACLFrontendHandlerImpl{Client: client, ReloadAgent: ra}

	api.ACLCreateACLFCGIAppHandler = &handlers.CreateACLFCGIAppHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLDeleteACLFCGIAppHandler = &handlers.DeleteACLFCGIAppHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLGetACLFCGIAppHandler = &handlers.GetACLFCGIAppHandlerImpl{Client: client}
	api.ACLGetAllACLFCGIAppHandler = &handlers.GetAllACLFCGIAppHandlerImpl{Client: client}
	api.ACLReplaceACLFCGIAppHandler = &handlers.ReplaceACLFCGIAppHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLReplaceAllACLFCGIAppHandler = &handlers.ReplaceAllACLFCGIAppHandlerImpl{Client: client, ReloadAgent: ra}

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

	// setup mailers sections handlers
	api.MailersCreateMailersSectionHandler = &handlers.CreateMailersSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.MailersDeleteMailersSectionHandler = &handlers.DeleteMailersSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.MailersGetMailersSectionHandler = &handlers.GetMailersSectionHandlerImpl{Client: client}
	api.MailersGetMailersSectionsHandler = &handlers.GetMailersSectionsHandlerImpl{Client: client}
	api.MailersEditMailersSectionHandler = &handlers.EditMailersSectionHandlerImpl{Client: client, ReloadAgent: ra}

	// setup mailer entry handlers
	api.MailerEntryCreateMailerEntryHandler = &handlers.CreateMailerEntryHandlerImpl{Client: client, ReloadAgent: ra}
	api.MailerEntryDeleteMailerEntryHandler = &handlers.DeleteMailerEntryHandlerImpl{Client: client, ReloadAgent: ra}
	api.MailerEntryGetMailerEntryHandler = &handlers.GetMailerEntryHandlerImpl{Client: client}
	api.MailerEntryGetMailerEntriesHandler = &handlers.GetMailerEntriesHandlerImpl{Client: client}
	api.MailerEntryReplaceMailerEntryHandler = &handlers.ReplaceMailerEntryHandlerImpl{Client: client, ReloadAgent: ra}

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

	// setup tables handlers
	api.TableCreateTableHandler = &handlers.CreateTableHandlerImpl{Client: client, ReloadAgent: ra}
	api.TableDeleteTableHandler = &handlers.DeleteTableHandlerImpl{Client: client, ReloadAgent: ra}
	api.TableGetTableHandler = &handlers.GetTableHandlerImpl{Client: client}
	api.TableGetTablesHandler = &handlers.GetTablesHandlerImpl{Client: client}
	api.TableReplaceTableHandler = &handlers.ReplaceTableHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http-errors sections handlers
	api.HTTPErrorsCreateHTTPErrorsSectionHandler = &handlers.CreateHTTPErrorsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorsDeleteHTTPErrorsSectionHandler = &handlers.DeleteHTTPErrorsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorsGetHTTPErrorsSectionHandler = &handlers.GetHTTPErrorsSectionHandlerImpl{Client: client}
	api.HTTPErrorsGetHTTPErrorsSectionsHandler = &handlers.GetHTTPErrorsSectionsHandlerImpl{Client: client}
	api.HTTPErrorsReplaceHTTPErrorsSectionHandler = &handlers.ReplaceHTTPErrorsSectionHandlerImpl{Client: client, ReloadAgent: ra}

	// setup cache handlers
	api.CacheCreateCacheHandler = &handlers.CreateCacheHandlerImpl{Client: client, ReloadAgent: ra}
	api.CacheDeleteCacheHandler = &handlers.DeleteCacheHandlerImpl{Client: client, ReloadAgent: ra}
	api.CacheGetCacheHandler = &handlers.GetCacheHandlerImpl{Client: client}
	api.CacheGetCachesHandler = &handlers.GetCachesHandlerImpl{Client: client}
	api.CacheReplaceCacheHandler = &handlers.ReplaceCacheHandlerImpl{Client: client, ReloadAgent: ra}

	// setup program handlers
	api.ProcessManagerCreateProgramHandler = &handlers.CreateProgramHandlerImpl{Client: client, ReloadAgent: ra}
	api.ProcessManagerDeleteProgramHandler = &handlers.DeleteProgramHandlerImpl{Client: client, ReloadAgent: ra}
	api.ProcessManagerGetProgramHandler = &handlers.GetProgramHandlerImpl{Client: client}
	api.ProcessManagerGetProgramsHandler = &handlers.GetProgramsHandlerImpl{Client: client}
	api.ProcessManagerReplaceProgramHandler = &handlers.ReplaceProgramHandlerImpl{Client: client, ReloadAgent: ra}

	// setup fcgi handlers
	api.FCGIAppCreateFCGIAppHandler = &handlers.CreateFCGIAppHandlerImpl{Client: client, ReloadAgent: ra}
	api.FCGIAppDeleteFCGIAppHandler = &handlers.DeleteFCGIAppHandlerImpl{Client: client, ReloadAgent: ra}
	api.FCGIAppGetFCGIAppHandler = &handlers.GetFCGIAppHandlerImpl{Client: client}
	api.FCGIAppGetFCGIAppsHandler = &handlers.GetFCGIAppsHandlerImpl{Client: client}
	api.FCGIAppReplaceFCGIAppHandler = &handlers.ReplaceFCGIAppHandlerImpl{Client: client, ReloadAgent: ra}

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
	api.DefaultsCreateDefaultsSectionHandler = &handlers.CreateDefaultsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.DefaultsAddDefaultsSectionHandler = &handlers.AddDefaultsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.DefaultsDeleteDefaultsSectionHandler = &handlers.DeleteDefaultsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.DefaultsReplaceDefaultsSectionHandler = &handlers.ReplaceDefaultsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.DefaultsGetDefaultsSectionHandler = &handlers.GetDefaultsSectionHandlerImpl{Client: client}
	api.DefaultsGetDefaultsSectionsHandler = &handlers.GetDefaultsSectionsHandlerImpl{Client: client}

	// setup reload handlers
	api.ReloadsGetReloadHandler = &handlers.GetReloadHandlerImpl{ReloadAgent: ra}
	api.ReloadsGetReloadsHandler = &handlers.GetReloadsHandlerImpl{ReloadAgent: ra}

	// setup runtime server handlers
	api.ServerGetRuntimeServerHandler = &handlers.GetRuntimeServerHandlerImpl{Client: client}
	api.ServerGetAllRuntimeServerHandler = &handlers.GetAllRuntimeServerHandlerImpl{Client: client}
	api.ServerReplaceRuntimeServerHandler = &handlers.ReplaceRuntimeServerHandlerImpl{Client: client}
	api.ServerAddRuntimeServerHandler = &handlers.AddRuntimeServerHandlerImpl{Client: client}
	api.ServerDeleteRuntimeServerHandler = &handlers.DeleteRuntimeServerHandlerImpl{Client: client}

	// setup stick table handlers
	api.StickTableGetStickTablesHandler = &handlers.GetStickTablesHandlerImpl{Client: client}
	api.StickTableGetStickTableHandler = &handlers.GetStickTableHandlerImpl{Client: client}
	api.StickTableGetStickTableEntriesHandler = &handlers.GetStickTableEntriesHandlerImpl{Client: client}
	api.StickTableSetStickTableEntriesHandler = &handlers.SetStickTableEntriesHandlerImpl{Client: client}

	// setup ACL runtime handlers
	api.ACLRuntimeGetServicesHaproxyRuntimeAclsHandler = &handlers.GetACLSHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeGetServicesHaproxyRuntimeAclsIDHandler = &handlers.GetACLHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeGetServicesHaproxyRuntimeAclsParentNameEntriesHandler = &handlers.GetACLFileEntriesHandlerRuntimeImpl{Client: client}
	api.ACLRuntimePostServicesHaproxyRuntimeAclsParentNameEntriesHandler = &handlers.PostACLFileEntryHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeGetServicesHaproxyRuntimeAclsParentNameEntriesIDHandler = &handlers.GetACLFileEntryRuntimeImpl{Client: client}
	api.ACLRuntimeDeleteServicesHaproxyRuntimeAclsParentNameEntriesIDHandler = &handlers.DeleteACLFileEntryHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeAddPayloadRuntimeACLHandler = &handlers.ACLRuntimeAddPayloadRuntimeACLHandlerImpl{Client: client}

	// setup map handlers
	api.MapsGetAllRuntimeMapFilesHandler = &handlers.GetMapsHandlerImpl{Client: client}
	api.MapsGetOneRuntimeMapHandler = &handlers.GetMapHandlerImpl{Client: client}
	api.MapsClearRuntimeMapHandler = &handlers.ClearMapHandlerImpl{Client: client}
	api.MapsShowRuntimeMapHandler = &handlers.ShowMapHandlerImpl{Client: client}
	api.MapsAddMapEntryHandler = &handlers.AddMapEntryHandlerImpl{Client: client}
	api.MapsAddPayloadRuntimeMapHandler = &handlers.MapsAddPayloadRuntimeMapHandlerImpl{Client: client}
	api.MapsGetRuntimeMapEntryHandler = &handlers.GetRuntimeMapEntryHandlerImpl{Client: client}
	api.MapsReplaceRuntimeMapEntryHandler = &handlers.ReplaceRuntimeMapEntryHandlerImpl{Client: client}
	api.MapsDeleteRuntimeMapEntryHandler = &handlers.DeleteRuntimeMapEntryHandlerImpl{Client: client}

	// crt-store handlers
	api.CrtStoreGetCrtStoresHandler = &handlers.GetCrtStoresHandlerImpl{Client: client}
	api.CrtStoreGetCrtStoreHandler = &handlers.GetCrtStoreHandlerImpl{Client: client}
	api.CrtStoreCreateCrtStoreHandler = &handlers.CreateCrtStoreHandlerImpl{Client: client, ReloadAgent: ra}
	api.CrtStoreEditCrtStoreHandler = &handlers.EditCrtStoreHandler{Client: client, ReloadAgent: ra}
	api.CrtStoreDeleteCrtStoreHandler = &handlers.DeleteCrtStoreHandlerImpl{Client: client, ReloadAgent: ra}
	// crt-store load handlers
	api.CrtLoadGetCrtLoadsHandler = &handlers.GetCrtLoadsHandlerImpl{Client: client}
	api.CrtLoadGetCrtLoadHandler = &handlers.GetCrtLoadHandlerImpl{Client: client}
	api.CrtLoadCreateCrtLoadHandler = &handlers.CreateCrtLoadHandlerImpl{Client: client, ReloadAgent: ra}
	api.CrtLoadReplaceCrtLoadHandler = &handlers.ReplaceCrtLoadHandler{Client: client, ReloadAgent: ra}
	api.CrtLoadDeleteCrtLoadHandler = &handlers.DeleteCrtLoadHandlerImpl{Client: client, ReloadAgent: ra}

	// traces handlers
	api.TracesGetTracesHandler = &handlers.GetTracesHandlerImpl{Client: client}
	api.TracesCreateTracesHandler = &handlers.CreateTracesHandlerImpl{Client: client, ReloadAgent: ra}
	api.TracesReplaceTracesHandler = &handlers.ReplaceTracesHandler{Client: client, ReloadAgent: ra}
	api.TracesDeleteTracesHandler = &handlers.DeleteTracesHandlerImpl{Client: client, ReloadAgent: ra}
	api.TracesCreateTraceEntryHandler = &handlers.CreateTraceEntryHandlerImpl{Client: client, ReloadAgent: ra}
	api.TracesDeleteTraceEntryHandler = &handlers.DeleteTraceEntryHandlerImpl{Client: client, ReloadAgent: ra}

	// log-profile handlers
	api.LogProfileGetLogProfilesHandler = &handlers.GetLogProfilesHandlerImpl{Client: client}
	api.LogProfileGetLogProfileHandler = &handlers.GetLogProfileHandlerImpl{Client: client}
	api.LogProfileCreateLogProfileHandler = &handlers.CreateLogProfileHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogProfileEditLogProfileHandler = &handlers.EditLogProfileHandler{Client: client, ReloadAgent: ra}
	api.LogProfileDeleteLogProfileHandler = &handlers.DeleteLogProfileHandlerImpl{Client: client, ReloadAgent: ra}

	// setup info handler
	api.InformationGetInfoHandler = &handlers.GetInfoHandlerImpl{SystemInfo: haproxyOptions.ShowSystemInfo, BuildTime: BuildTime, Version: Version}

	// setup cluster handlers
	api.ClusterGetClusterHandler = &handlers.GetClusterHandlerImpl{Config: cfg}
	api.ClusterPostClusterHandler = &handlers.CreateClusterHandlerImpl{Client: client, Config: cfg, ReloadAgent: ra}
	api.ClusterDeleteClusterHandler = &handlers.DeleteClusterHandlerImpl{Client: client, Config: cfg, Users: dataplaneapi_config.GetUsersStore(), ReloadAgent: ra}
	api.ClusterEditClusterHandler = &handlers.EditClusterHandlerImpl{Config: cfg}
	api.ClusterInitiateCertificateRefreshHandler = &handlers.ClusterInitiateCertificateRefreshHandlerImpl{Config: cfg}

	clusterSync := dataplaneapi_config.ClusterSync{ReloadAgent: ra, Context: ctx}
	go clusterSync.Monitor(cfg, client)

	// setup specification handler
	api.SpecificationGetSpecificationHandler = specification.GetSpecificationHandlerFunc(func(params specification.GetSpecificationParams, principal interface{}) middleware.Responder {
		var m map[string]interface{}
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		if err := json.Unmarshal(SwaggerJSON, &m); err != nil {
			e := misc.HandleError(err)
			return specification.NewGetSpecificationDefault(int(*e.Code)).WithPayload(e)
		}
		return specification.NewGetSpecificationOK().WithPayload(&m)
	})

	configurationClient, err := client.Configuration()
	if err != nil {
		log.Fatal(err)
	}
	// set up service discovery handlers
	discovery := service_discovery.NewServiceDiscoveries(service_discovery.ServiceDiscoveriesParams{
		Client:      configurationClient,
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
		var errSD error
		if data.ID == nil || len(*data.ID) == 0 {
			data.ID = service_discovery.NewServiceDiscoveryUUID()
		}
		if errSD = service_discovery.ValidateConsulData(data, true); errSD != nil {
			log.Fatal("Error validating Consul instance: " + errSD.Error())
		}
		if errSD = discovery.AddNode("consul", *data.ID, data); errSD != nil {
			log.Warning("Error creating consul instance: " + errSD.Error())
		}
	}
	_ = cfg.SaveConsuls(cfg.ServiceDiscovery.Consuls)

	// create stored AWS instances
	for _, data := range cfg.ServiceDiscovery.AWSRegions {
		var errSD error

		if data.ID == nil || len(*data.ID) == 0 {
			data.ID = service_discovery.NewServiceDiscoveryUUID()
		}
		if errSD = service_discovery.ValidateAWSData(data, true); errSD != nil {
			log.Fatal("Error validating AWS instance: " + errSD.Error())
		}
		if errSD = discovery.AddNode("aws", *data.ID, data); errSD != nil {
			log.Warning("Error creating AWS instance: " + errSD.Error())
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

	// general file storage handlers
	api.StorageCreateStorageGeneralFileHandler = &handlers.StorageCreateStorageGeneralFileHandlerImpl{Client: client}
	api.StorageGetAllStorageGeneralFilesHandler = &handlers.StorageGetAllStorageGeneralFilesHandlerImpl{Client: client}
	api.StorageGetOneStorageGeneralFileHandler = &handlers.StorageGetOneStorageGeneralFileHandlerImpl{Client: client}
	api.StorageDeleteStorageGeneralFileHandler = &handlers.StorageDeleteStorageGeneralFileHandlerImpl{Client: client}
	api.StorageReplaceStorageGeneralFileHandler = &handlers.StorageReplaceStorageGeneralFileHandlerImpl{Client: client, ReloadAgent: ra}

	// setup OpenAPI v3 specification handler
	api.Version3GetOpenapiv3SpecificationHandler = version3.GetOpenapiv3SpecificationHandlerFunc(func(params version3.GetOpenapiv3SpecificationParams, principal interface{}) middleware.Responder {
		v2 := openapi2.T{}
		v2JSONString := string(SwaggerJSON)
		v2JSONString = strings.ReplaceAll(v2JSONString, "#/definitions", "#/components/schemas")
		curatedV2 := json.RawMessage([]byte(v2JSONString))

		err = v2.UnmarshalJSON(curatedV2)
		if err != nil {
			e := misc.HandleError(err)
			return version3.NewGetOpenapiv3SpecificationDefault(int(*e.Code)).WithPayload(e)
		}

		// if host is empty(dynamic hosts), server prop is empty,
		// so we need to set it explicitly
		if v2.Host == "" {
			cfg = dataplaneapi_config.Get()
			v2.Host = cfg.RuntimeData.Host
		}

		var v3 *openapi3.T
		v3, err = openapi2conv.ToV3(&v2)
		if err != nil {
			e := misc.HandleError(err)
			return version3.NewGetOpenapiv3SpecificationDefault(int(*e.Code)).WithPayload(e)
		}
		return version3.NewGetOpenapiv3SpecificationOK().WithPayload(v3)
	})

	// TODO: do we need a ReloadAgent for SPOE
	// setup SPOE handlers
	api.SpoeCreateSpoeHandler = &handlers.SpoeCreateSpoeHandlerImpl{Client: client}
	api.SpoeDeleteSpoeFileHandler = &handlers.SpoeDeleteSpoeFileHandlerImpl{Client: client}
	api.SpoeGetAllSpoeFilesHandler = &handlers.SpoeGetAllSpoeFilesHandlerImpl{Client: client}
	api.SpoeGetOneSpoeFileHandler = &handlers.SpoeGetOneSpoeFileHandlerImpl{Client: client}

	// SPOE scope
	api.SpoeGetAllSpoeScopeHandler = &handlers.SpoeGetAllSpoeScopeHandlerImpl{Client: client}
	api.SpoeGetSpoeScopeHandler = &handlers.SpoeGetSpoeScopeHandlerImpl{Client: client}
	api.SpoeCreateSpoeScopeHandler = &handlers.SpoeCreateSpoeScopeHandlerImpl{Client: client}
	api.SpoeDeleteSpoeScopeHandler = &handlers.SpoeDeleteSpoeScopeHandlerImpl{Client: client}

	// SPOE agent
	api.SpoeGetAllSpoeAgentHandler = &handlers.SpoeGetAllSpoeAgentHandlerImpl{Client: client}
	api.SpoeGetSpoeAgentHandler = &handlers.SpoeGetSpoeAgentHandlerImpl{Client: client}
	api.SpoeCreateSpoeAgentHandler = &handlers.SpoeCreateSpoeAgentHandlerImpl{Client: client}
	api.SpoeDeleteSpoeAgentHandler = &handlers.SpoeDeleteSpoeAgentHandlerImpl{Client: client}
	api.SpoeReplaceSpoeAgentHandler = &handlers.SpoeReplaceSpoeAgentHandlerImpl{Client: client}

	// SPOE messages
	api.SpoeGetAllSpoeMessageHandler = &handlers.SpoeGetAllSpoeMessageHandlerImpl{Client: client}
	api.SpoeGetSpoeMessageHandler = &handlers.SpoeGetSpoeMessageHandlerImpl{Client: client}
	api.SpoeCreateSpoeMessageHandler = &handlers.SpoeCreateSpoeMessageHandlerImpl{Client: client}
	api.SpoeDeleteSpoeMessageHandler = &handlers.SpoeDeleteSpoeMessageHandlerImpl{Client: client}
	api.SpoeReplaceSpoeMessageHandler = &handlers.SpoeReplaceSpoeMessageHandlerImpl{Client: client}

	// SPOE groups
	api.SpoeGetAllSpoeGroupHandler = &handlers.SpoeGetAllSpoeGroupHandlerImpl{Client: client}
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

	// Health
	api.HealthGetHealthHandler = &handlers.GetHealthHandlerImpl{HAProxy: ra}

	// middlewares
	var adpts []adapters.Adapter
	adpts = append(adpts,
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
		adapters.ConfigVersionMiddleware(client),
	)

	// setup logging middlewares
	accessLogger, err := log.AccessLogger()
	if err != nil {
		log.Warningf("Error getting access loggers: %s", err.Error())
	}

	if accessLogger != nil {
		for _, logger := range accessLogger.Loggers() {
			adpts = append(adpts, adapters.ApacheLogMiddleware(logger))
		}
	}

	appLogger, err := log.AppLogger()
	if err != nil {
		log.Warningf("Error getting app loggers: %s", err.Error())
	}
	if appLogger != nil {
		adpts = append(adpts, adapters.RecoverMiddleware(appLogger))
	}

	// Configure/Re-configure DebugServer on runtime socket
	debugServer := socket_runtime.GetServer()
	select {
	case debugServer.CnChannel <- client:
	default:
		// ... do not block dataplane
		log.Warning("-- command socket failed to update cn client")
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares), adpts...)
}

func configureReloadAgent(haproxyOptions dataplaneapi_config.HAProxyConfiguration, client client_native.HAProxyClient, ctx context.Context) *haproxy.ReloadAgent {
	// Initialize reload agent
	raParams := haproxy.ReloadAgentParams{
		Delay:           haproxyOptions.ReloadDelay,
		ReloadCmd:       haproxyOptions.ReloadCmd,
		UseMasterSocket: canUseMasterSocketReload(&haproxyOptions, client),
		RestartCmd:      haproxyOptions.RestartCmd,
		StatusCmd:       haproxyOptions.StatusCmd,
		ConfigFile:      haproxyOptions.ConfigFile,
		BackupDir:       haproxyOptions.BackupsDir,
		Retention:       haproxyOptions.ReloadRetention,
		Client:          client,
		Ctx:             ctx,
	}

	ra, e := haproxy.NewReloadAgent(raParams)
	if e != nil {
		log.Fatalf("Cannot initialize reload agent: %v", e)
	}
	return ra
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

func serverShutdown() {
	cfg := dataplaneapi_config.Get()
	if logFile != nil {
		logFile.Close()
	}
	if cfg.HAProxy.UpdateMapFiles {
		cfg.MapSync.Stop()
	}
}

// Determine if we can reload HAProxy's configuration using the master socket 'reload' command.
// This requires at least HAProxy 2.7 configured in master/worker mode, with a master socket.
func canUseMasterSocketReload(conf *dataplaneapi_config.HAProxyConfiguration, client client_native.HAProxyClient) bool {
	useMasterSocket := conf.MasterWorkerMode && conf.MasterRuntime != "" && misc.IsUnixSocketAddr(conf.MasterRuntime)
	if !useMasterSocket {
		return false
	}

	rt, err := client.Runtime()
	if err != nil {
		return false
	}
	currVersion, err := rt.GetVersion()
	if err != nil {
		return false
	}

	return cn_runtime.IsBiggerOrEqual(&cn_runtime.HAProxyVersion{Major: 2, Minor: 7}, &currVersion)
}

func configureNativeClient(cyx context.Context, haproxyOptions dataplaneapi_config.HAProxyConfiguration, mWorker bool) client_native.HAProxyClient {
	// Initialize HAProxy native client
	confClient, err := cn.ConfigureConfigurationClient(haproxyOptions, mWorker)
	if err != nil {
		log.Fatalf("Error initializing configuration client: %v", err)
	}

	runtimeClient := cn.ConfigureRuntimeClient(cyx, confClient, haproxyOptions)

	opt := []options.Option{
		options.Configuration(confClient),
		options.Runtime(runtimeClient),
	}
	if haproxyOptions.MapsDir != "" {
		var mapStorage storage.Storage
		mapStorage, err = storage.New(haproxyOptions.MapsDir, storage.MapsType)
		if err != nil {
			log.Fatalf("error initializing map storage: %v", err)
		}
		opt = append(opt, options.MapStorage(mapStorage))
	} else {
		log.Fatalf("error trying to use empty string for managed map directory")
	}

	if haproxyOptions.SSLCertsDir != "" {
		var sslCertStorage storage.Storage
		sslCertStorage, err = storage.New(haproxyOptions.SSLCertsDir, storage.SSLType)
		if err != nil {
			log.Fatalf("error initializing SSL certs storage: %v", err)
		}
		opt = append(opt, options.SSLCertStorage(sslCertStorage))
	} else {
		log.Fatalf("error trying to use empty string for managed map directory")
	}

	if haproxyOptions.GeneralStorageDir != "" {
		var generalStorage storage.Storage
		generalStorage, err = storage.New(haproxyOptions.GeneralStorageDir, storage.GeneralType)
		if err != nil {
			log.Fatalf("error initializing General storage: %v", err)
		}
		opt = append(opt, options.GeneralStorage(generalStorage))
	} else {
		log.Fatalf("error trying to use empty string for managed general files directory")
	}

	if haproxyOptions.SpoeDir != "" {
		prms := spoe.Params{
			SpoeDir:        haproxyOptions.SpoeDir,
			TransactionDir: haproxyOptions.SpoeTransactionDir,
		}
		var spoeClient spoe.Spoe
		spoeClient, err = spoe.NewSpoe(prms)
		if err != nil {
			log.Fatalf("error setting up spoe: %v", err)
		}
		opt = append(opt, options.Spoe(spoeClient))
	} else {
		log.Fatalf("error trying to use empty string for SPOE configuration directory")
	}

	client, err := client_native.New(cyx, opt...)
	if err != nil {
		log.Fatalf("Error initializing configuration client: %v", err)
	}

	return client
}

func handleSignals(ctx context.Context, cancel context.CancelFunc, sigs chan os.Signal, client client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users) {
	for {
		select {
		case sig := <-sigs:
			if sig == syscall.SIGUSR1 {
				var clientCtx context.Context
				cancel()
				clientCtx, cancel = context.WithCancel(ctx)
				configuration, err := client.Configuration()
				if err != nil {
					log.Infof("Unable to reload Data Plane API: %s", err.Error())
				} else {
					client.ReplaceRuntime(cn.ConfigureRuntimeClient(clientCtx, configuration, haproxyOptions))
					log.Info("Reloaded Data Plane API")
				}
			} else if sig == syscall.SIGUSR2 {
				reloadConfigurationFile(client, haproxyOptions, users)
			}
		case <-ctx.Done():
			cancel()
			return
		}
	}
}

func reloadConfigurationFile(client client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users) {
	confClient, err := cn.ConfigureConfigurationClient(haproxyOptions, mWorker)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := users.Init(); err != nil {
		log.Fatal(err.Error())
	}
	log.Info("Rereading Configuration Files")
	clientMutex.Lock()
	defer clientMutex.Unlock()
	client.ReplaceConfiguration(confClient)
}

func startWatcher(ctx context.Context, client client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users, reloadAgent *haproxy.ReloadAgent) error {
	cb := func() {
		// reload configuration from config file.
		reloadConfigurationFile(client, haproxyOptions, users)

		// reload runtime client if necessary.
		callbackNeeded, reconfigureFunc, err := cn.ReconfigureRuntime(client)
		if err != nil {
			log.Warningf("Failed to check if native client need to be reloaded: %s", err)
			return
		}
		if callbackNeeded {
			reloadAgent.ReloadWithCallback(reconfigureFunc)
		}

		// get the last configuration which has been updated by reloadConfigurationFile and increment version in config file.
		configuration, err := client.Configuration()
		if err != nil {
			log.Warningf("Failed to get configuration: %s", err)
			return
		}
		if err := configuration.IncrementVersion(); err != nil {
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
