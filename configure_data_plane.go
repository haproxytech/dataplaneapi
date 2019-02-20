// This file is safe to edit. Once it exists it will not be overwritten

package dataplaneapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	parser "github.com/haproxytech/config-parser"
	"github.com/haproxytech/config-parser/types"

	log "github.com/sirupsen/logrus"

	"github.com/carbocation/interpose/adaptors"

	"github.com/meatballhat/negroni-logrus"

	"github.com/haproxytech/dataplaneapi/misc"

	"github.com/go-openapi/runtime/middleware"
	"github.com/haproxytech/dataplaneapi/operations/discovery"

	"github.com/haproxytech/client-native"

	"github.com/haproxytech/client-native/configuration"
	runtime_api "github.com/haproxytech/client-native/runtime"
	"github.com/haproxytech/dataplaneapi/handlers"
	"github.com/haproxytech/dataplaneapi/haproxy"

	"github.com/dre1080/recover"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	swag "github.com/go-openapi/swag"
	"github.com/rs/cors"
	graceful "github.com/tylerb/graceful"

	"github.com/haproxytech/dataplaneapi/operations"
)

//go:generate swagger generate server --target ../../../../../../github.com/haproxytech --name controller --spec ../../../../../../../../haproxy-api/haproxy-open-api-spec/build/haproxy_spec.yaml --server-package controller --tags Stats --tags Information --tags Configuration --tags Discovery --tags Frontend --tags Backend --tags Bind --tags Server --tags TCPRequestRule --tags HTTPRequestRule --tags HTTPResponseRule --tags Acl --tags BackendSwitchingRule --tags ServerSwitchingRule --tags TCPResponseRule --skip-models --exclude-main

var haproxyOptions struct {
	ConfigFile     string `short:"c" long:"config-file" description:"Path to the haproxy configuration file" default:"/etc/haproxy/haproxy.cfg"`
	Userlist       string `short:"u" long:"userlist" description:"Userlist in HAProxy configuration to use for API Basic Authentication" default:"controller"`
	HAProxy        string `short:"b" long:"haproxy-bin" description:"Path to the haproxy binary file" default:"haproxy"`
	ReloadDelay    int    `short:"d" long:"reload-delay" description:"Minimum delay between two reloads (in s)"`
	ReloadCmd      string `short:"r" long:"reload-cmd" description:"Reload command"`
	TransactionDir string `short:"t" long:"transaction-dir" description:"Path to the transaction directory" default:"/tmp/haproxy"`
}

func configureFlags(api *operations.DataPlaneAPI) {
	haproxyOptionsGroup := swag.CommandLineOptionsGroup{
		ShortDescription: "HAProxy options",
		LongDescription:  "Options for configuring haproxy locations",
		Options:          &haproxyOptions,
	}

	api.CommandLineOptionsGroups = make([]swag.CommandLineOptionsGroup, 0, 1)
	api.CommandLineOptionsGroups = append(api.CommandLineOptionsGroups, haproxyOptionsGroup)
}

func configureAPI(api *operations.DataPlaneAPI) http.Handler {
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

	// Initialize HAProxy native client
	confClient := &configuration.Client{}
	confParams := configuration.ClientParams{
		ConfigurationFile: haproxyOptions.ConfigFile,
		Haproxy:           haproxyOptions.HAProxy,
		UseValidation:     false,
		UseCache:          true,
		TransactionDir:    haproxyOptions.TransactionDir,
	}
	err := confClient.Init(confParams)
	if err != nil {
		fmt.Println("Error setting up configuration client, using default one")
		confClient, err = configuration.DefaultClient()
		if err != nil {
			fmt.Println("Error setting up default configuration client, exiting...")
			api.ServerShutdown()
		}
	}

	runtimeClient := &runtime_api.Client{}
	globalConf, err := confClient.GetGlobalConfiguration("")

	if err != nil {
		fmt.Println("Stats socket not configured, no runtime client initiated")
	}

	nbproc := globalConf.Data.Nbproc
	if nbproc == 0 {
		nbproc = 1
	}

	runtimeAPI := globalConf.Data.RuntimeAPI
	if runtimeAPI == "" {
		fmt.Println("Stats socket not configured, no runtime client initiated")
	} else {
		socketList := make([]string, 0, 1)
		if nbproc > 1 {
			for i := int64(0); i < nbproc; i++ {
				socketList = append(socketList, fmt.Sprintf("%v.%v", runtimeAPI, i))
			}
		} else {
			socketList = append(socketList, runtimeAPI)
		}
		err := runtimeClient.Init(socketList)
		if err != nil {
			fmt.Println("Error setting up runtime client, not using one")
			runtimeClient = nil
		}
	}

	client := &client_native.HAProxyClient{}
	client.Init(confClient, runtimeClient)

	// Initialize reload agent
	ra := &haproxy.ReloadAgent{}
	ra.Init(haproxyOptions.ReloadDelay, haproxyOptions.ReloadCmd)

	// Applies when the Authorization header is set with the Basic scheme
	api.BasicAuthAuth = func(user string, pass string) (interface{}, error) {
		return authenticateUser(user, pass, client)
	}

	// setup discovery handlers
	api.DiscoveryGetAPIEndpointsHandler = discovery.GetAPIEndpointsHandlerFunc(func(params discovery.GetAPIEndpointsParams, principal interface{}) middleware.Responder {
		rURI := "/" + strings.SplitN(params.HTTPRequest.RequestURI[1:], "/", 2)[1]
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

	// setup stats handler
	api.StatsGetStatsHandler = &handlers.GetStatsHandlerImpl{Client: client}

	// setup info handler
	api.InformationGetHaproxyProcessInfoHandler = &handlers.GetInformationHandlerImpl{Client: client}

	// setup raw configuration handlers
	api.ConfigurationGetHAProxyConfigurationHandler = &handlers.GetRawConfigurationHandlerImpl{Client: client}
	api.ConfigurationPostHAProxyConfigurationHandler = &handlers.PostRawConfigurationHandlerImpl{Client: client, ReloadAgent: ra}

	// setup global configuration handlers
	api.GlobalGetGlobalHandler = &handlers.GetGlobalHandlerImpl{Client: client}
	api.GlobalReplaceGlobalHandler = &handlers.ReplaceGlobalHandlerImpl{Client: client, ReloadAgent: ra}

	api.ServerShutdown = func() {}

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
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	handleCORS := cors.AllowAll().Handler

	recovery := recover.New(&recover.Options{
		Log: log.Print,
	})

	mw := negronilogrus.NewMiddlewareFromLogger(log.StandardLogger(), "controller")
	logViaLogrus := adaptors.FromNegroni(mw)
	return recovery(logViaLogrus(handleCORS(handler)))
}

func authenticateUser(user string, pass string, cli *client_native.HAProxyClient) (interface{}, error) {
	if err := cli.Configuration.ConfigParser.LoadData(cli.Configuration.ConfigurationFile); err != nil {
		return nil, fmt.Errorf("Error reading users from %v userlist in conf", haproxyOptions.Userlist)
	}
	data, err := cli.Configuration.ConfigParser.Get(parser.UserList, haproxyOptions.Userlist, "user")
	if err != nil {
		return nil, err
	}
	users, ok := data.([]types.User)
	if !ok {
		return nil, fmt.Errorf("Error reading users from %v userlist in conf", haproxyOptions.Userlist)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("No users configured in %v userlist in conf", haproxyOptions.Userlist)
	}

	for _, u := range users {
		if u.Name == user {
			if u.Password == pass {
				return user, nil
			}
			return nil, errors.New(401, "Invalid password")
		}
	}
	return nil, errors.New(401, "User does not exist")
}
