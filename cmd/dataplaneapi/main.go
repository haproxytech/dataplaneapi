// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this files except in compliance with the License.
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
	"os"

	log "github.com/sirupsen/logrus"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	"github.com/haproxytech/dataplaneapi"
	"github.com/haproxytech/dataplaneapi/operations"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})
	log.SetOutput(os.Stdout)
}

func main() {
	swaggerSpec, err := loads.Embedded(dataplaneapi.SwaggerJSON, dataplaneapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewDataPlaneAPI(swaggerSpec)
	server := dataplaneapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "HAProxy API"
	parser.LongDescription = "API for editing and managing HAPEE instances"

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				os.Exit(0)
			} else {
				log.Fatalln(err)
			}
		}
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
