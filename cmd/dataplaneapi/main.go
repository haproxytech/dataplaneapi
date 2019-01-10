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
	// log.SetReportCaller(true)
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
