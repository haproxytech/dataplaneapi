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
	"fmt"
	"os"
	"path"

	loads "github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/security"
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"

	"github.com/haproxytech/dataplaneapi"
	"github.com/haproxytech/dataplaneapi/configuration"
	"github.com/haproxytech/dataplaneapi/operations"
)

//GitRepo ...
var GitRepo = ""

//GitTag ...
var GitTag = ""

//GitCommit ...
var GitCommit = "dev"

//GitDirty ...
var GitDirty = ".dirty"

//BuildTime ...
var BuildTime = ""

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})
	log.SetOutput(os.Stdout)
}

var cliOptions struct {
	Version bool `short:"v" long:"version" description:"Version and build information"`
}

func main() {
	cfg := configuration.Get()
	for {
		restart := startServer(cfg)
		if !restart.Load() {
			break
		}
	}
}

func startServer(cfg *configuration.Configuration) (reload configuration.AtomicBool) {

	swaggerSpec, err := loads.Embedded(dataplaneapi.SwaggerJSON, dataplaneapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	dataplaneapi.BuildTime = BuildTime
	dataplaneapi.Version = fmt.Sprintf("%s %s%s", GitTag, GitCommit, GitDirty)

	api := operations.NewDataPlaneAPI(swaggerSpec)
	server := dataplaneapi.NewServer(api)
	//nolint
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "HAProxy Data Plane API"
	parser.LongDescription = "API for editing and managing haproxy instances"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err = parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	_, err = parser.AddGroup("Show version", "Show build and version information", &cliOptions)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err = parser.Parse(); err != nil {
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				log.Fatalln(err)
			}
		}
	}

	if cliOptions.Version {
		fmt.Printf("HAProxy Data Plane API %s %s%s\n\n", GitTag, GitCommit, GitDirty)
		fmt.Printf("Build from: %s\n", GitRepo)
		fmt.Printf("Build date: %s\n\n", BuildTime)
		return
	}

	err = cfg.Load(dataplaneapi.SwaggerJSON, server.Host, server.Port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Infof("HAProxy Data Plane API %s %s%s", GitTag, GitCommit, GitDirty)
	log.Infof("Build from: %s", GitRepo)
	log.Infof("Build date: %s", BuildTime)

	if cfg.Mode.Load() == "cluster" {
		if cfg.Cluster.Certificate.Fetched.Load() {
			log.Info("HAProxy Data Plane API in cluster mode")
			server.TLSCertificate = flags.Filename(path.Join(cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.crt", cfg.Name.Load())))
			server.TLSCertificateKey = flags.Filename(path.Join(cfg.GetClusterCertDir(), fmt.Sprintf("dataplane-%s.key", cfg.Name.Load())))
			server.EnabledListeners = []string{"https"}
			if server.TLSPort == 0 {
				server.TLSPort = server.Port
			}
		} else if cfg.Cluster.ActiveBootstrapKey.Load() != "" {
			cfg.Notify.BootstrapKeyChanged.NotifyWithRetry()
		}
	}
	err = cfg.Save()
	if err != nil {
		log.Fatalln(err)
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

	go func() {
		for range cfg.Notify.Reload.Subscribe("main") {
			log.Info("HAProxy Data Plane API reloading")
			reload.Store(true)
			cfg.UnSubscribeAll()
			err := server.Shutdown()
			if err != nil {
				log.Fatalln(err)
			}
		}
	}()

	go func() {
		for range cfg.Notify.Shutdown.Subscribe("main") {
			log.Info("HAProxy Data Plane API shutting down")
			err := server.Shutdown()
			if err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		}
	}()

	server.ConfigureAPI()
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
	return reload
}
