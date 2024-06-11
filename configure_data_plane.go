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
	"github.com/haproxytech/dataplaneapi/operations/specification_openapiv3"
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
	api.ServerCreateServerHandler = &handlers.CreateServerHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerDeleteServerHandler = &handlers.DeleteServerHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerGetServerHandler = &handlers.GetServerHandlerImpl{Client: client}
	api.ServerGetServersHandler = &handlers.GetServersHandlerImpl{Client: client}
	api.ServerReplaceServerHandler = &handlers.ReplaceServerHandlerImpl{Client: client, ReloadAgent: ra}

	// setup server template handlers
	api.ServerTemplateCreateServerTemplateHandler = &handlers.CreateServerTemplateHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerTemplateDeleteServerTemplateHandler = &handlers.DeleteServerTemplateHandlerImpl{Client: client, ReloadAgent: ra}
	api.ServerTemplateGetServerTemplateHandler = &handlers.GetServerTemplateHandlerImpl{Client: client}
	api.ServerTemplateGetServerTemplatesHandler = &handlers.GetServerTemplatesHandlerImpl{Client: client}
	api.ServerTemplateReplaceServerTemplateHandler = &handlers.ReplaceServerTemplateHandlerImpl{Client: client, ReloadAgent: ra}

	// setup bind handlers
	api.BindCreateBindHandler = &handlers.CreateBindHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindDeleteBindHandler = &handlers.DeleteBindHandlerImpl{Client: client, ReloadAgent: ra}
	api.BindGetBindHandler = &handlers.GetBindHandlerImpl{Client: client}
	api.BindGetBindsHandler = &handlers.GetBindsHandlerImpl{Client: client}
	api.BindReplaceBindHandler = &handlers.ReplaceBindHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http check handlers
	api.HTTPCheckGetHTTPCheckHandler = &handlers.GetHTTPCheckHandlerImpl{Client: client}
	api.HTTPCheckGetHTTPChecksHandler = &handlers.GetHTTPChecksHandlerImpl{Client: client}
	api.HTTPCheckCreateHTTPCheckHandler = &handlers.CreateHTTPCheckHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckReplaceHTTPCheckHandler = &handlers.ReplaceHTTPCheckHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckDeleteHTTPCheckHandler = &handlers.DeleteHTTPCheckHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPCheckReplaceHTTPChecksHandler = &handlers.ReplaceHTTPChecksHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http request rule handlers
	api.HTTPRequestRuleCreateHTTPRequestRuleHandler = &handlers.CreateHTTPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleDeleteHTTPRequestRuleHandler = &handlers.DeleteHTTPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleGetHTTPRequestRuleHandler = &handlers.GetHTTPRequestRuleHandlerImpl{Client: client}
	api.HTTPRequestRuleGetHTTPRequestRulesHandler = &handlers.GetHTTPRequestRulesHandlerImpl{Client: client}
	api.HTTPRequestRuleReplaceHTTPRequestRuleHandler = &handlers.ReplaceHTTPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPRequestRuleReplaceHTTPRequestRulesHandler = &handlers.ReplaceHTTPRequestRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http after response rule handlers
	api.HTTPAfterResponseRuleCreateHTTPAfterResponseRuleHandler = &handlers.CreateHTTPAfterResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleDeleteHTTPAfterResponseRuleHandler = &handlers.DeleteHTTPAfterResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleGetHTTPAfterResponseRuleHandler = &handlers.GetHTTPAfterResponseRuleHandlerImpl{Client: client}
	api.HTTPAfterResponseRuleGetHTTPAfterResponseRulesHandler = &handlers.GetHTTPAfterResponseRulesHandlerImpl{Client: client}
	api.HTTPAfterResponseRuleReplaceHTTPAfterResponseRuleHandler = &handlers.ReplaceHTTPAfterResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPAfterResponseRuleReplaceHTTPAfterResponseRulesHandler = &handlers.ReplaceHTTPAfterResponseRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http response rule handlers
	api.HTTPResponseRuleCreateHTTPResponseRuleHandler = &handlers.CreateHTTPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleDeleteHTTPResponseRuleHandler = &handlers.DeleteHTTPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleGetHTTPResponseRuleHandler = &handlers.GetHTTPResponseRuleHandlerImpl{Client: client}
	api.HTTPResponseRuleGetHTTPResponseRulesHandler = &handlers.GetHTTPResponseRulesHandlerImpl{Client: client}
	api.HTTPResponseRuleReplaceHTTPResponseRuleHandler = &handlers.ReplaceHTTPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPResponseRuleReplaceHTTPResponseRulesHandler = &handlers.ReplaceHTTPResponseRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup http error rule handlers
	api.HTTPErrorRuleCreateHTTPErrorRuleHandler = &handlers.CreateHTTPErrorRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleDeleteHTTPErrorRuleHandler = &handlers.DeleteHTTPErrorRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleGetHTTPErrorRuleHandler = &handlers.GetHTTPErrorRuleHandlerImpl{Client: client}
	api.HTTPErrorRuleGetHTTPErrorRulesHandler = &handlers.GetHTTPErrorRulesHandlerImpl{Client: client}
	api.HTTPErrorRuleReplaceHTTPErrorRuleHandler = &handlers.ReplaceHTTPErrorRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.HTTPErrorRuleReplaceHTTPErrorRulesHandler = &handlers.ReplaceHTTPErrorRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp content rule handlers
	api.TCPRequestRuleCreateTCPRequestRuleHandler = &handlers.CreateTCPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleDeleteTCPRequestRuleHandler = &handlers.DeleteTCPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleGetTCPRequestRuleHandler = &handlers.GetTCPRequestRuleHandlerImpl{Client: client}
	api.TCPRequestRuleGetTCPRequestRulesHandler = &handlers.GetTCPRequestRulesHandlerImpl{Client: client}
	api.TCPRequestRuleReplaceTCPRequestRuleHandler = &handlers.ReplaceTCPRequestRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPRequestRuleReplaceTCPRequestRulesHandler = &handlers.ReplaceTCPRequestRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp connection rule handlers
	api.TCPResponseRuleCreateTCPResponseRuleHandler = &handlers.CreateTCPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleDeleteTCPResponseRuleHandler = &handlers.DeleteTCPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleGetTCPResponseRuleHandler = &handlers.GetTCPResponseRuleHandlerImpl{Client: client}
	api.TCPResponseRuleGetTCPResponseRulesHandler = &handlers.GetTCPResponseRulesHandlerImpl{Client: client}
	api.TCPResponseRuleReplaceTCPResponseRuleHandler = &handlers.ReplaceTCPResponseRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPResponseRuleReplaceTCPResponseRulesHandler = &handlers.ReplaceTCPResponseRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup tcp check handlers
	api.TCPCheckCreateTCPCheckHandler = &handlers.CreateTCPCheckHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckDeleteTCPCheckHandler = &handlers.DeleteTCPCheckHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckGetTCPCheckHandler = &handlers.GetTCPCheckHandlerImpl{Client: client}
	api.TCPCheckGetTCPChecksHandler = &handlers.GetTCPChecksHandlerImpl{Client: client}
	api.TCPCheckReplaceTCPCheckHandler = &handlers.ReplaceTCPCheckHandlerImpl{Client: client, ReloadAgent: ra}
	api.TCPCheckReplaceTCPChecksHandler = &handlers.ReplaceTCPChecksHandlerImpl{Client: client, ReloadAgent: ra}

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
	api.FilterCreateFilterHandler = &handlers.CreateFilterHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterDeleteFilterHandler = &handlers.DeleteFilterHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterGetFilterHandler = &handlers.GetFilterHandlerImpl{Client: client}
	api.FilterGetFiltersHandler = &handlers.GetFiltersHandlerImpl{Client: client}
	api.FilterReplaceFilterHandler = &handlers.ReplaceFilterHandlerImpl{Client: client, ReloadAgent: ra}
	api.FilterReplaceFiltersHandler = &handlers.ReplaceFiltersHandlerImpl{Client: client, ReloadAgent: ra}

	// setup stick rule handlers
	api.StickRuleCreateStickRuleHandler = &handlers.CreateStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleDeleteStickRuleHandler = &handlers.DeleteStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleGetStickRuleHandler = &handlers.GetStickRuleHandlerImpl{Client: client}
	api.StickRuleGetStickRulesHandler = &handlers.GetStickRulesHandlerImpl{Client: client}
	api.StickRuleReplaceStickRuleHandler = &handlers.ReplaceStickRuleHandlerImpl{Client: client, ReloadAgent: ra}
	api.StickRuleReplaceStickRulesHandler = &handlers.ReplaceStickRulesHandlerImpl{Client: client, ReloadAgent: ra}

	// setup log target handlers
	api.LogTargetCreateLogTargetHandler = &handlers.CreateLogTargetHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetDeleteLogTargetHandler = &handlers.DeleteLogTargetHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetGetLogTargetHandler = &handlers.GetLogTargetHandlerImpl{Client: client}
	api.LogTargetGetLogTargetsHandler = &handlers.GetLogTargetsHandlerImpl{Client: client}
	api.LogTargetReplaceLogTargetHandler = &handlers.ReplaceLogTargetHandlerImpl{Client: client, ReloadAgent: ra}
	api.LogTargetReplaceLogTargetsHandler = &handlers.ReplaceLogTargetsHandlerImpl{Client: client, ReloadAgent: ra}

	// setup acl rule handlers
	api.ACLCreateACLHandler = &handlers.CreateACLHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLDeleteACLHandler = &handlers.DeleteACLHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLGetACLHandler = &handlers.GetACLHandlerImpl{Client: client}
	api.ACLGetAclsHandler = &handlers.GetAclsHandlerImpl{Client: client}
	api.ACLReplaceACLHandler = &handlers.ReplaceACLHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLReplaceAclsHandler = &handlers.ReplaceAclsHandlerImpl{Client: client, ReloadAgent: ra}

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
	api.DefaultsDeleteDefaultsSectionHandler = &handlers.DeleteDefaultsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.DefaultsReplaceDefaultsSectionHandler = &handlers.ReplaceDefaultsSectionHandlerImpl{Client: client, ReloadAgent: ra}
	api.DefaultsGetDefaultsSectionHandler = &handlers.GetDefaultsSectionHandlerImpl{Client: client}
	api.DefaultsGetDefaultsSectionsHandler = &handlers.GetDefaultsSectionsHandlerImpl{Client: client}

	// setup reload handlers
	api.ReloadsGetReloadHandler = &handlers.GetReloadHandlerImpl{ReloadAgent: ra}
	api.ReloadsGetReloadsHandler = &handlers.GetReloadsHandlerImpl{ReloadAgent: ra}

	// setup runtime server handlers
	api.ServerGetRuntimeServerHandler = &handlers.GetRuntimeServerHandlerImpl{Client: client}
	api.ServerGetRuntimeServersHandler = &handlers.GetRuntimeServersHandlerImpl{Client: client}
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
	api.ACLRuntimeGetServicesHaproxyRuntimeACLFileEntriesHandler = &handlers.GetACLFileEntriesHandlerRuntimeImpl{Client: client}
	api.ACLRuntimePostServicesHaproxyRuntimeACLFileEntriesHandler = &handlers.PostACLFileEntryHandlerRuntimeImpl{Client: client}
	api.ACLRuntimeGetServicesHaproxyRuntimeACLFileEntriesIDHandler = &handlers.GetACLFileEntryRuntimeImpl{Client: client}
	api.ACLRuntimeDeleteServicesHaproxyRuntimeACLFileEntriesIDHandler = &handlers.DeleteACLFileEntryHandlerRuntimeImpl{Client: client}
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
			log.Fatalf("Error validating Consul instance: " + errSD.Error())
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
			log.Fatalf("Error validating AWS instance: " + errSD.Error())
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
	api.SpecificationOpenapiv3GetOpenapiv3SpecificationHandler = specification_openapiv3.GetOpenapiv3SpecificationHandlerFunc(func(params specification_openapiv3.GetOpenapiv3SpecificationParams, principal interface{}) middleware.Responder {
		v2 := openapi2.T{}
		err = v2.UnmarshalJSON(SwaggerJSON)
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

		var v3 *openapi3.T
		v3, err = openapi2conv.ToV3(&v2)
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
		log.Fatalf(err.Error())
	}
	if err := users.Init(); err != nil {
		log.Fatalf(err.Error())
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
