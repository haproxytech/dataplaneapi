package main

import (
	"fmt"
	"log"
	"os"

	runtimeClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/dataplaneapi/client"
	frontend "github.com/haproxytech/dataplaneapi/client/frontend"
)

func main() {
	// create the transport
	rt := runtimeClient.New(os.Getenv("HAPROXY_DATAPLANE_ADDR"), client.DefaultBasePath, []string{"http"})
	writer := runtimeClient.BasicAuth("admin", "mypassword")

	// create the API client, with the transport
	client := frontend.New(rt, strfmt.Default)

	// make the request to get all frontends
	resp, err := client.GetFrontends(nil, writer)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("response: %#v\n", resp.Payload)
}
