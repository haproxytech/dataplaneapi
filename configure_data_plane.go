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
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"

	"github.com/Masterminds/semver/v3"
	"github.com/go-openapi/swag/cmdutils"
	client_native "github.com/haproxytech/client-native/v6"
	"github.com/haproxytech/client-native/v6/options"
	cn_runtime "github.com/haproxytech/client-native/v6/runtime"
	"github.com/haproxytech/client-native/v6/spoe"
	"github.com/haproxytech/client-native/v6/storage"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/reload_agent"
	"github.com/rs/cors"

	"github.com/haproxytech/dataplaneapi/adapters"
	cn "github.com/haproxytech/dataplaneapi/client-native"
	dataplaneapi_config "github.com/haproxytech/dataplaneapi/configuration"
	service_discovery "github.com/haproxytech/dataplaneapi/discovery"
	"github.com/haproxytech/dataplaneapi/handlers"
	"github.com/haproxytech/dataplaneapi/misc"
	"github.com/haproxytech/dataplaneapi/resilient"
	socket_runtime "github.com/haproxytech/dataplaneapi/runtime"

	// import various crypting algorithms
	_ "github.com/GehirnInc/crypt/md5_crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	_ "github.com/GehirnInc/crypt/sha512_crypt"
)

var (
	Version               string
	BuildTime             string
	mWorker               = false
	logFile               *os.File
	AppLogger             *log.Logger
	AccLogger             *log.Logger
	serverStartedCallback func()
	clientMutex           sync.Mutex
	eventListener         *cn.HAProxyEventListener
)

func SetServerStartedCallback(callFunc func()) {
	serverStartedCallback = callFunc
}

// GetCommandLineOptionsGroups returns the flag groups that should be registered
// with the CLI parser for HAProxy, logging, and syslog options.
func GetCommandLineOptionsGroups() []cmdutils.CommandLineOptionsGroup {
	cfg := dataplaneapi_config.Get()
	return []cmdutils.CommandLineOptionsGroup{
		{
			ShortDescription: "HAProxy options",
			LongDescription:  "Options for configuring haproxy locations.",
			Options:          &cfg.HAProxy,
		},
		{
			ShortDescription: "Logging options",
			LongDescription:  "Options for configuring logging.",
			Options:          &cfg.Logging,
		},
		{
			ShortDescription: "Syslog options",
			LongDescription:  "Options for configuring syslog logging.",
			Options:          &cfg.Syslog,
		},
	}
}

func configureAPI(skipBasicAuth bool, maxBodySize int64) (http.Handler, func()) { //nolint:maintidx
	clientMutex.Lock()
	defer clientMutex.Unlock()

	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Error starting Data Plane API: %s\n Stacktrace from panic: \n%s", err, string(debug.Stack()))
		}
	}()

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

		for f := range strings.SplitSeq(cfgFiles, ";") {
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

	ctx := ContextHandler.Context()
	clientCtx, cancel := context.WithCancel(ctx)

	client := configureNativeClient(clientCtx, haproxyOptions, mWorker)

	initDataplaneStorage(haproxyOptions.DataplaneStorageDir, client)

	configureEventListener(clientCtx, client)

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

	// Cluster sync
	clusterSync := dataplaneapi_config.ClusterSync{ReloadAgent: ra, Context: ctx}
	go clusterSync.Monitor(cfg, client)

	// Service discovery setup — create the shared ServiceDiscoveries object and
	// populate it with any persisted Consul/AWS instances from the config file.
	configurationClient, err := client.Configuration()
	if err != nil {
		log.Fatal(err)
	}
	discovery := service_discovery.NewServiceDiscoveries(service_discovery.ServiceDiscoveriesParams{
		Client:      configurationClient,
		ReloadAgent: ra,
		Context:     ctx,
	})
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

	client = resilient.NewClient(client)

	// Build the chi router — all handler packages register their routes here.
	opts := handlers.Options{
		Client:                client,
		ReloadAgent:           ra,
		Config:                cfg,
		Users:                 users,
		ConsulDiscovery:       discovery,
		ConsulPersistCallback: cfg.SaveConsuls,
		AWSDiscovery:          discovery,
		AWSPersistCallback:    cfg.SaveAWS,
		UseValidation:         true,
		Version:               Version,
		BuildTime:             BuildTime,
		SystemInfo:            haproxyOptions.ShowSystemInfo,
		MaxOpenTransactions:   int(cfg.HAProxy.MaxOpenTransactions),
		SwaggerJSON:           SwaggerJSON,
	}
	chiRouter, err := handlers.NewRouter(opts)
	if err != nil {
		log.Fatalf("Error building router: %s", err.Error())
	}

	// Configure/Re-configure DebugServer on runtime socket
	debugServer := socket_runtime.GetServer()
	select {
	case debugServer.CnChannel <- client:
	default:
		// ... do not block dataplane
		log.Warning("-- command socket failed to update cn client")
	}

	// Request flow: RecoverMiddleware → ApacheLog → ConfigVersion → CORS → SpecDocs → BasicAuth → MaxBodySize → chi router
	// This matches the go-swagger server's ordering: ConfigVersion sits outside
	// CORS and auth so the Configuration-Version header is present on 401 and
	// CORS preflight responses too, and SpecDocs sits in front of BasicAuth so
	// /swagger.json and /v3/docs stay unauthenticated.
	var adpts []adapters.Adapter
	adpts = append(
		adpts,
		adapters.MaxBodySizeMiddleware(maxBodySize),
		adapters.BasicAuthMiddleware(skipBasicAuth),
		adapters.SpecDocsMiddleware(SwaggerJSON, "/v3", "HAProxy Data Plane API"),
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

	return setupGlobalMiddleware(chiRouter, adpts...), serverShutdown
}

func configureReloadAgent(haproxyOptions dataplaneapi_config.HAProxyConfiguration, client client_native.HAProxyClient, ctx context.Context) *reload_agent.ReloadAgent {
	// Initialize reload agent
	msReload := canUseMasterSocketReload(&haproxyOptions, client)
	raParams := reload_agent.ReloadAgentParams{
		Delay:      haproxyOptions.ReloadDelay,
		ReloadCmd:  haproxyOptions.ReloadCmd,
		UseRuntime: msReload,
		RestartCmd: haproxyOptions.RestartCmd,
		StatusCmd:  haproxyOptions.StatusCmd,
		ConfigFile: haproxyOptions.ConfigFile,
		BackupDir:  haproxyOptions.BackupsDir,
		Retention:  haproxyOptions.ReloadRetention,
		Ctx:        ctx,
	}
	if msReload {
		rt, err := client.Runtime()
		if err != nil {
			log.Fatalf("Cannot initialize reload agent: %v", err)
		}
		raParams.Runtime = rt
	}

	ra, e := reload_agent.NewReloadAgent(raParams)
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

	return cn_runtime.IsBiggerOrEqual(&cn_runtime.HAProxyVersion{Version: semver.New(2, 7, 0, "", "")}, &currVersion)
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
		// crt-lists use the same directory
		var crtListStorage storage.Storage
		crtListStorage, err = storage.New(haproxyOptions.SSLCertsDir, storage.CrtListType)
		if err != nil {
			log.Fatalf("error initializing CRT Lists storage: %v", err)
		}
		opt = append(opt, options.CrtListStorage(crtListStorage))
	} else {
		log.Fatalf("error trying to use empty string for managed SSL certificates directory")
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

func configureEventListener(ctx context.Context, client client_native.HAProxyClient) {
	rt, err := client.Runtime()
	if err != nil {
		return
	}

	if eventListener != nil {
		err = eventListener.Reconfigure(ctx, rt)
		if err != nil {
			// Stop the listener if the new conf has no master socket.
			log.Info("Stopping the EventListener:", err.Error())
			_ = eventListener.Stop()
		}
	} else {
		// First start.
		eventListener, err = cn.ListenHAProxyEvents(ctx, client)
		if err != nil && err != cn.ErrNoMasterSocket && err != cn.ErrOldVersion {
			log.Error("Failed to start HAProxy's event listener:", err.Error())
		}
	}
}

func handleSignals(ctx context.Context, cancel context.CancelFunc, sigs chan os.Signal, client client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users) {
	for {
		select {
		case sig := <-sigs:
			switch sig {
			case syscall.SIGUSR1:
				var clientCtx context.Context
				cancel()
				clientCtx, cancel = context.WithCancel(ctx)
				configuration, err := client.Configuration()
				if err != nil {
					log.Infof("Unable to reload Data Plane API: %s", err.Error())
				} else {
					client.ReplaceRuntime(cn.ConfigureRuntimeClient(clientCtx, configuration, haproxyOptions))
					configureEventListener(clientCtx, client)
					log.Info("Reloaded Data Plane API")
				}
			case syscall.SIGUSR2:
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

func startWatcher(ctx context.Context, client client_native.HAProxyClient, haproxyOptions dataplaneapi_config.HAProxyConfiguration, users *dataplaneapi_config.Users, reloadAgent *reload_agent.ReloadAgent) error {
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

		// Reconfigure the event listener if needed.
		configureEventListener(ctx, client)

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
