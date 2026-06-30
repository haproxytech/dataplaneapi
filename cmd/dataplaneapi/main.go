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
	"syscall"

	_ "github.com/KimMachineGun/automemlimit"
	"github.com/haproxytech/dataplaneapi"
	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/log"
	socket_runtime "github.com/haproxytech/dataplaneapi/runtime"
	flags "github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	_ "go.uber.org/automaxprocs"
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
	_ = godotenv.Load()
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
	dataplaneapi.BuildTime = BuildTime
	dataplaneapi.Version = fmt.Sprintf("%s %s%s", GitTag, GitCommit, GitDirty)

	server := dataplaneapi.NewServer()
	server.Logger = log.Printf

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "HAProxy Data Plane API"
	parser.LongDescription = "API for editing and managing haproxy instances"

	for _, optsGroup := range dataplaneapi.GetCommandLineOptionsGroups() {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	var err error
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
		return reload
	}

	var loadMsg []string
	_, err = cfg.Load()
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

	loadMsg, err = cfg.LoadDataplaneStorageConfig()
	if err != nil {
		fmt.Println("configuration error:", err)
		os.Exit(1)
	}
	cfg.FlagLoadDapiStorageData = true

	// incorporate changes from file to global settings
	dataplaneapi.SyncWithFileSettings(server, cfg)
	err = cfg.LoadRuntimeVars(dataplaneapi.SwaggerJSON, server.Host, server.Port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configuration.HandlePIDFile(cfg.HAProxy)

	if err = log.InitWithConfiguration(cfg.LogTargets); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg.InitSignalHandler()

	log.Infof("HAProxy Data Plane API %s %s%s", GitTag, GitCommit, GitDirty)
	log.Infof("Build from: %s", GitRepo)
	log.Infof("Build date: %s", BuildTime)
	log.Infof("Reload strategy: %s", cfg.HAProxy.ReloadStrategy)

	// log deprecation message
	if len(loadMsg) > 0 {
		for _, msg := range loadMsg {
			log.Warning(msg)
		}
	}

	err = cfg.Save()
	if err != nil {
		log.Fatalf("Error saving configuration: %s", err.Error())
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
