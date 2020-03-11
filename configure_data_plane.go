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
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"

	"github.com/haproxytech/dataplaneapi/adapters"
	"github.com/haproxytech/dataplaneapi/operations/specification"

	parser "github.com/haproxytech/config-parser/v2"
	"github.com/haproxytech/config-parser/v2/types"

	log "github.com/sirupsen/logrus"

	"github.com/haproxytech/dataplaneapi/misc"

	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/dataplaneapi/operations/discovery"

	client_native "github.com/haproxytech/client-native"

	"github.com/haproxytech/client-native/configuration"
	runtime_api "github.com/haproxytech/client-native/runtime"
	dataplaneapi_config "github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/handlers"
	"github.com/haproxytech/dataplaneapi/haproxy"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	swag "github.com/go-openapi/swag"
	"github.com/rs/cors"

	"github.com/haproxytech/dataplaneapi/operations"

	"github.com/GehirnInc/crypt"
	// import various crypting algorithms
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

//go:generate swagger generate server --target ../../../../../../github.com/haproxytech --name controller --spec ../../../../../../../../haproxy-api/haproxy-open-api-spec/build/haproxy_spec.yaml --server-package controller --tags Stats --tags Information --tags Configuration --tags Discovery --tags Frontend --tags Backend --tags Bind --tags Server --tags TCPRequestRule --tags HTTPRequestRule --tags HTTPResponseRule --tags Acl --tags BackendSwitchingRule --tags ServerSwitchingRule --tags TCPResponseRule --skip-models --exclude-main

var Version string
var BuildTime string

var logFile *os.File

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

	api.CommandLineOptionsGroups = make([]swag.CommandLineOptionsGroup, 0, 1)
	api.CommandLineOptionsGroups = append(api.CommandLineOptionsGroups, haproxyOptionsGroup)
	api.CommandLineOptionsGroups = append(api.CommandLineOptionsGroups, loggingOptionsGroup)
}

func configureAPI(api *operations.DataPlaneAPI) http.Handler {
	cfg := dataplaneapi_config.Get()
	haproxyOptions := cfg.HAProxy

	configureLogging(cfg.Logging)

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Error starting Data Plane API: %s\n Stacktrace from panic: \n%s", err, string(debug.Stack()))
		}
	}()
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

	api.TxtProducer = runtime.TextProducer()

	api.ServerShutdown = serverShutdown

	client := configureNativeClient(haproxyOptions)

	// Handle reload signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	go handleSignals(sigs, client, haproxyOptions)

	// Initialize reload agent
	ra := &haproxy.ReloadAgent{}
	if err := ra.Init(haproxyOptions.ReloadDelay, haproxyOptions.ReloadCmd, haproxyOptions.RestartCmd, haproxyOptions.ConfigFile, haproxyOptions.ReloadRetention); err != nil {
		log.Fatalf("Cannot initialize reload agent: %v", err)
	}

	// Applies when the Authorization header is set with the Basic scheme
	api.BasicAuthAuth = func(user string, pass string) (interface{}, error) {
		return authenticateUser(user, pass, client)
	}
	// setup discovery handlers
	api.DiscoveryGetAPIEndpointsHandler = discovery.GetAPIEndpointsHandlerFunc(func(params discovery.GetAPIEndpointsParams, principal interface{}) middleware.Responder {
		uriSlice := strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)
		rURI := ""
		if len(uriSlice) < 2 {
			rURI = "/"
		} else {
			rURI = "/" + uriSlice[1]
		}

		ends, err := misc.DiscoverChildPaths(rURI, SwaggerJSON)
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

	// setup transaction handlers
	api.TransactionsStartTransactionHandler = &handlers.StartTransactionHandlerImpl{Client: client}
	api.TransactionsDeleteTransactionHandler = &handlers.DeleteTransactionHandlerImpl{Client: client}
	api.TransactionsGetTransactionHandler = &handlers.GetTransactionHandlerImpl{Client: client}
	api.TransactionsGetTransactionsHandler = &handlers.GetTransactionsHandlerImpl{Client: client}
	api.TransactionsCommitTransactionHandler = &handlers.CommitTransactionHandlerImpl{Client: client, ReloadAgent: ra}

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

	// setup log target handlers
	api.ACLCreateACLHandler = &handlers.CreateACLHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLDeleteACLHandler = &handlers.DeleteACLHandlerImpl{Client: client, ReloadAgent: ra}
	api.ACLGetACLHandler = &handlers.GetACLHandlerImpl{Client: client}
	api.ACLGetAclsHandler = &handlers.GetAclsHandlerImpl{Client: client}
	api.ACLReplaceACLHandler = &handlers.ReplaceACLHandlerImpl{Client: client, ReloadAgent: ra}

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

	// setup info handler
	api.InformationGetInfoHandler = &handlers.GetInfoHandlerImpl{SystemInfo: haproxyOptions.ShowSystemInfo, BuildTime: BuildTime, Version: Version}

	// setup specification handler
	api.SpecificationGetSpecificationHandler = specification.GetSpecificationHandlerFunc(func(params specification.GetSpecificationParams, principal interface{}) middleware.Responder {
		var m map[string]interface{}
		if err := json.Unmarshal(SwaggerJSON, &m); err != nil {
			e := misc.HandleError(err)
			return specification.NewGetSpecificationDefault(int(*e.Code)).WithPayload(e)
		}
		return specification.NewGetSpecificationOK().WithPayload(&m)
	})

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
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
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.New(cors.Options{
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
	}).Handler
	recovery := adapters.RecoverMiddleware(log.StandardLogger())
	logViaLogrus := adapters.LoggingMiddleware(log.StandardLogger())
	return (logViaLogrus(handleCORS(recovery(handler))))
}

func authenticateUser(user string, pass string, cli *client_native.HAProxyClient) (interface{}, error) {
	cfg := dataplaneapi_config.Get()
	haproxyOptions := cfg.HAProxy

	data, err := cli.Configuration.Parser.Get(parser.UserList, haproxyOptions.Userlist, "user")
	if err != nil {
		return nil, fmt.Errorf("error reading userlist %v userlist in conf: %s", haproxyOptions.Userlist, err.Error())
	}
	users, ok := data.([]types.User)
	if !ok {
		return nil, fmt.Errorf("error reading users from %v userlist in conf", haproxyOptions.Userlist)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no users configured in %v userlist in conf", haproxyOptions.Userlist)
	}

	for _, u := range users {
		if u.Name == user {
			if u.IsInsecure {
				if u.Password == pass {
					return user, nil
				}
			} else {
				if checkPassword(pass, u.Password) {
					return user, nil
				}
			}
			return nil, errors.New(401, "Invalid username/password")
		}
	}
	return nil, errors.New(401, "Invalid username/password")
}

func configureLogging(loggingOptions dataplaneapi_config.LoggingOptions) {
	switch loggingOptions.LogFormat {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		})
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
	}

	switch loggingOptions.LogTo {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "file":
		dir := filepath.Dir(loggingOptions.LogFile)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Warning("Error opening log file, no logging implemented: " + err.Error())
		}
		//nolint:govet
		logFile, err := os.OpenFile(loggingOptions.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Warning("Error opening log file, no logging implemented: " + err.Error())
		}
		log.SetOutput(logFile)
	}

	switch loggingOptions.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	}
}

func checkPassword(pass, storedPass string) bool {
	parts := strings.Split(storedPass, "$")
	if len(parts) == 4 {
		var c crypt.Crypter
		switch parts[1] {
		case "1":
			c = crypt.MD5.New()
		case "5":
			c = crypt.SHA256.New()
		case "6":
			c = crypt.SHA512.New()
		default:
			return false
		}
		if err := c.Verify(storedPass, []byte(pass)); err == nil {
			return true
		}
	}

	return false
}

func serverShutdown() {
	if logFile != nil {
		logFile.Close()
	}
}

func configureNativeClient(haproxyOptions dataplaneapi_config.HAProxyConfiguration) *client_native.HAProxyClient {
	// Override options with env variables
	if os.Getenv("HAPROXY_MWORKER") == "1" {
		masterRuntime := os.Getenv("HAPROXY_MASTER_CLI")
		if misc.IsUnixSocketAddr(masterRuntime) {
			haproxyOptions.MasterRuntime = masterRuntime
		}
	}
	cfgFiles := os.Getenv("HAPROXY_CFGFILES")
	if cfgFiles != "" {
		cfg := strings.Split(cfgFiles, ";")
		haproxyOptions.ConfigFile = cfg[0]
	}
	// Initialize HAProxy native client
	confClient, err := configureConfigurationClient(haproxyOptions)
	if err != nil {
		log.Fatalf("Error initializing configuration client: %v", err)
	}

	runtimeClient := configureRuntimeClient(confClient, haproxyOptions)
	client := &client_native.HAProxyClient{}
	if err := client.Init(confClient, runtimeClient); err != nil {
		log.Fatalf("Error setting up native client: %v", err)
	}

	return client
}

func configureConfigurationClient(haproxyOptions dataplaneapi_config.HAProxyConfiguration) (*configuration.Client, error) {
	confClient := &configuration.Client{}
	confParams := configuration.ClientParams{
		ConfigurationFile:      haproxyOptions.ConfigFile,
		Haproxy:                haproxyOptions.HAProxy,
		BackupsNumber:          haproxyOptions.BackupsNumber,
		UseValidation:          false,
		PersistentTransactions: true,
		TransactionDir:         haproxyOptions.TransactionDir,
	}
	err := confClient.Init(confParams)
	if err != nil {
		return nil, fmt.Errorf("error setting up configuration client: %s", err.Error())
	}
	return confClient, nil
}

func configureRuntimeClient(confClient *configuration.Client, haproxyOptions dataplaneapi_config.HAProxyConfiguration) *runtime_api.Client {
	runtimeClient := &runtime_api.Client{}
	_, globalConf, err := confClient.GetGlobalConfiguration("")

	// First try to setup master runtime socket
	if err == nil {
		var err error
		// If master socket is set and a valid unix socket, use only this
		if haproxyOptions.MasterRuntime != "" && misc.IsUnixSocketAddr(haproxyOptions.MasterRuntime) {
			masterSocket := haproxyOptions.MasterRuntime
			// if nbproc is set set nbproc sockets
			if globalConf.Nbproc > 0 {
				nbproc := int(globalConf.Nbproc)
				if err = runtimeClient.InitWithMasterSocket(masterSocket, nbproc); err == nil {
					return runtimeClient
				}
				log.Warningf("Error setting up runtime client with master socket: %s : %s", masterSocket, err.Error())
			} else {
				// if nbproc is not set, use master socket with 1 process
				if err = runtimeClient.InitWithMasterSocket(masterSocket, 1); err == nil {
					return runtimeClient
				}
				log.Warningf("Error setting up runtime client with master socket: %s : %s", masterSocket, err.Error())
			}
		}
		runtimeAPIs := globalConf.RuntimeApis
		// if no master socket set, read from first valid socket if nbproc <= 1
		if globalConf.Nbproc <= 1 {
			socketList := make(map[int]string)
			for _, r := range runtimeAPIs {
				if misc.IsUnixSocketAddr(*r.Address) {
					socketList[1] = *r.Address
					if err = runtimeClient.InitWithSockets(socketList); err == nil {
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
			if err = runtimeClient.InitWithSockets(sockets); err == nil {
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

func handleSignals(sigs chan os.Signal, client *client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration) {
	//nolint:gosimple
	for {
		select {
		case sig := <-sigs:
			if sig == syscall.SIGUSR1 {
				client.Runtime = configureRuntimeClient(client.Configuration, haproxyOptions)
				log.Info("Reloaded Data Plane API")
			} else if sig == syscall.SIGUSR2 {
				confClient, err := configureConfigurationClient(haproxyOptions)
				if err != nil {
					log.Fatalf(err.Error())
				}
				log.Info("Rereading Configuration Files")
				client.Configuration = confClient
			}
		}
	}
}
