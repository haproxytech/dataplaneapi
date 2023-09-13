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

package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"syscall"

	loads "github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/security"
	flags "github.com/jessevdk/go-flags"

	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/client-native/v5/storage"
	"github.com/haproxytech/dataplaneapi"
	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/log"
	"github.com/haproxytech/dataplaneapi/operations"
	socket_runtime "github.com/haproxytech/dataplaneapi/runtime"
)

// GitRepo ...
var GitRepo = ""

// GitTag ...
var GitTag = ""

// GitCommit ...
var GitCommit = "dev"

// GitDirty ...
var GitDirty = ".dirty"

// BuildTime ...
var BuildTime = ""

var cliOptions struct {
	Version bool `short:"v" long:"version" description:"Version and build information"`
}

func main() {
	cancelDebugServer := startRuntimeDebugServer()

	cfg := configuration.Get()
	for {
		restart := startServer(cfg, cancelDebugServer)
		if !restart.Load() {
			break
		}
	}
}

func startRuntimeDebugServer() context.CancelFunc {
	ctx := context.Background()
	ctx, cancelDebugServer := context.WithCancel(ctx)
	debugServer := socket_runtime.GetServer()
	debugServer.DAPIVersion = fmt.Sprintf("%s %s%s", GitTag, GitCommit, GitDirty)
	go debugServer.Start(ctx, cancelDebugServer)
	return cancelDebugServer
}

func startServer(cfg *configuration.Configuration, cancelDebugServer context.CancelFunc) (reload configuration.AtomicBool) {
	swaggerSpec, err := loads.Embedded(dataplaneapi.SwaggerJSON, dataplaneapi.FlatSwaggerJSON)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataplaneapi.BuildTime = BuildTime
	dataplaneapi.Version = fmt.Sprintf("%s %s%s", GitTag, GitCommit, GitDirty)

	api := operations.NewDataPlaneAPI(swaggerSpec)
	server := dataplaneapi.NewServer(api)

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "HAProxy Data Plane API"
	parser.LongDescription = "API for editing and managing haproxy instances"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err = parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	_, err = parser.AddGroup("Show version", "Show build and version information", &cliOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err = parser.Parse(); err != nil {
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}

	if cliOptions.Version {
		fmt.Printf("HAProxy Data Plane API %s %s%s\n\n", GitTag, GitCommit, GitDirty)
		fmt.Printf("Build from: %s\n", GitRepo)
		fmt.Printf("Build date: %s\n\n", BuildTime)
		return
	}

	err = cfg.Load()
	if err != nil {
		fmt.Println("configuration error:", err)
		os.Exit(1)
	}

	if cfg.HAProxy.UID != 0 {
		if err = syscall.Setuid(cfg.HAProxy.UID); err != nil {
			fmt.Println("set uid:", err)
			os.Exit(1)
		}
	}

	if cfg.HAProxy.GID != 0 {
		if err = syscall.Setgid(cfg.HAProxy.GID); err != nil {
			fmt.Println("set gid:", err)
			os.Exit(1)
		}
	}

	// incorporate changes from file to global settings
	dataplaneapi.SyncWithFileSettings(server, cfg)
	err = cfg.LoadRuntimeVars(dataplaneapi.SwaggerJSON, server.Host, server.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configuration.HandlePIDFile(cfg.HAProxy)

	if cfg.Mode.Load() == configuration.ModeCluster {
		if cfg.Cluster.CertificateFetched.Load() {
			log.Info("HAProxy Data Plane API in cluster mode")
			server.TLSCertificate = flags.Filename(path.Join(cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.crt", cfg.Name.Load())))
			server.TLSCertificateKey = flags.Filename(path.Join(cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.key", cfg.Name.Load())))
			server.EnabledListeners = []string{"https"}
			if server.TLSPort == 0 {
				server.TLSPort = server.Port
			}
			// override storage dir location
			storageDir := cfg.Cluster.StorageDir.Load()
			if storageDir != "" {
				cfg.HAProxy.MapsDir = path.Join(storageDir, string(storage.MapsType))
				cfg.HAProxy.SSLCertsDir = path.Join(storageDir, string(storage.SSLType))
				cfg.HAProxy.GeneralStorageDir = path.Join(storageDir, string(storage.GeneralType))
				cfg.HAProxy.SpoeDir = path.Join(storageDir, string(storage.SpoeType))
				cfg.HAProxy.SpoeTransactionDir = path.Join(storageDir, string(storage.SpoeTransactionsType))
				cfg.HAProxy.BackupsDir = path.Join(storageDir, string(storage.BackupsType))
				cfg.HAProxy.TransactionDir = path.Join(storageDir, string(storage.TransactionsType))
				// dataplane internal
				cfg.HAProxy.ClusterTLSCertDir = path.Join(storageDir, "certs-cluster")
				cfg.Cluster.CertificateDir.Store(path.Join(storageDir, "certs-cluster"))
			}
		} else if cfg.Cluster.ActiveBootstrapKey.Load() != "" {
			cfg.Notify.BootstrapKeyChanged.NotifyWithRetry()
		}
	}

	clusterLogTargets := parseClusterLogTargets(cfg)
	if err = log.InitWithConfiguration(cfg.LogTargets, cfg.Logging, cfg.Syslog, clusterLogTargets, cfg.Cluster.ID.Load()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg.InitSignalHandler()

	log.Infof("HAProxy Data Plane API %s %s%s", GitTag, GitCommit, GitDirty)
	log.Infof("Build from: %s", GitRepo)
	log.Infof("Build date: %s", BuildTime)
	log.Infof("Reload strategy: %s", cfg.HAProxy.ReloadStrategy)

	err = cfg.Save()
	if err != nil {
		log.Fatalf("Error saving configuration: %s", err.Error())
	}

	// Applies when the Authorization header is set with the Basic scheme
	api.BasicAuthAuth = configuration.AuthenticateUser
	api.BasicAuthenticator = func(authentication security.UserPassAuthentication) runtime.Authenticator {
		// if mTLS is enabled with backing Certificate Authority, skipping basic authentication
		if len(server.TLSCACertificate) > 0 && server.TLSPort > 0 {
			return runtime.AuthenticatorFunc(func(i interface{}) (bool, interface{}, error) {
				return true, "", nil
			})
		}
		return security.BasicAuthRealm("", authentication)
	}

	dataplaneapi.ContextHandler.Init()
	go func() {
		<-cfg.Notify.Reload.Subscribe("main")
		log.Info("HAProxy Data Plane API reloading")
		reload.Store(true)
		cfg.UnSubscribeAll()
		cfg.StopSignalHandler()
		dataplaneapi.ContextHandler.Cancel()
		err = server.Shutdown()
		if err != nil {
			log.Fatalf("Error reloading HAProxy Data Plane API: %s", err.Error())
		}
	}()

	go func() {
		select {
		case <-cfg.Notify.Shutdown.Subscribe("main"):
			log.Info("HAProxy Data Plane API shutting down")
			err = server.Shutdown()
			if err != nil {
				log.Fatalf("Error shutting down HAProxy Data Plane API: %s", err.Error())
			}
			cancelDebugServer()
			os.Exit(0)
		case <-dataplaneapi.ContextHandler.Context().Done():
			return
		}
	}()

	server.ConfigureAPI()
	dataplaneapi.SetServerStartedCallback(cfg.Notify.ServerStarted.Notify)
	if err := server.Serve(); err != nil {
		log.Fatalf("Error running HAProxy Data Plane API: %s", err.Error())
	}

	defer server.Shutdown() //nolint:errcheck

	return reload
}

func parseClusterLogTargets(cfg *configuration.Configuration) []*models.ClusterLogTarget {
	if cfg.Mode.Load() == "cluster" && cfg.Cluster.ClusterLogTargets != nil && len(cfg.Cluster.ClusterLogTargets) > 0 {
		return cfg.Cluster.ClusterLogTargets
	}
	return []*models.ClusterLogTarget{}
}
