// Copyright 2025 HAProxy Technologies
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

package cn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/dataplaneapi/acme"
	"github.com/haproxytech/dataplaneapi/log"
	jsoniter "github.com/json-iterator/go"
)

// Supported event types.
const (
	EventAcmeNewCert = "newcert"
	EventAcmeDeploy  = "deploy"
)

// Structs used to unmarshal ACME messages in JSON.
type acmeIdentifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type acmeStatus struct {
	Identifier acmeIdentifier `json:"identifier"`
}

func (h *HAProxyEventListener) handleAcmeEvent(ctx context.Context, message string) {
	name, args, ok := strings.Cut(message, " ")
	if !ok {
		log.Warningf("failed to parse ACME Event: '%s'", message)
		return
	}

	switch name {
	case EventAcmeNewCert:
		h.handleAcmeNewCertEvent(args)
		return
	case EventAcmeDeploy:
		h.handleAcmeDeployEvent(ctx, args)
		return
	}

	log.Debugf("unknown ACME Event type: '%s'", message)
}

// HAProxy has created a new certificate and needs dpapi to write it to Storage.
// ex: acme newcert foobar.pem.rsa
func (h *HAProxyEventListener) handleAcmeNewCertEvent(args string) {
	if len(args) == 0 {
		log.Error("received HAProxy Event 'acme newcert' without a cert name")
		return
	}

	// Do not use the certificate name in args as a storage name.
	// It could be an alias, or the user could have split keys and certs storage.

	crt, err := h.rt.ShowCertificate(args)
	if err != nil {
		log.Errorf("events: acme newcert %s: %s", args, err.Error())
		return
	}

	storage, err := h.client.SSLCertStorage()
	if err != nil {
		log.Error(err)
		return
	}

	// 'dump ssl cert' can only be issued on sockets with level "admin".
	pem, err := h.rt.DumpCertificate(crt.StorageName)
	if err != nil {
		log.Errorf("events: acme newcert %s: dump cert: %s", args, err.Error())
		return
	}

	// The storage API only wants the filename, while the runtime API uses paths.
	storageName := filepath.Base(crt.StorageName)

	// Create or Replace the certificate.
	_, _, err = storage.Get(storageName)
	if err != nil {
		if errors.Is(err, configuration.ErrObjectDoesNotExist) {
			rc := io.NopCloser(strings.NewReader(pem))
			_, _, err = storage.Create(storageName, rc)
		}
	} else {
		_, err = storage.Replace(storageName, pem)
	}

	if err != nil {
		log.Errorf("events: acme newcert %s: storage: %s", args, err.Error())
		return
	}

	log.Debugf("events: OK: acme newcert %s => %s", args, crt.StorageName)
}

// HAProxy needs dpapi to solve a dns-01 challenge.
// example:
// acme deploy CertIdentifier thumbprint "QPFLnguBJSfyTiN2c4DWiWJvpveUB3bvY3EoC8cZC-U"\n
// provider-name "godaddy"\n
// acme-vars "var=var1,var=var2"\n
// acme-vars "var1=foobar,var2=var2"\n
//
//	{
//	  "identifier": {
//	    "type": "dns",
//	    "value": "test1.example.com"
//	  },
//	  "status": "pending",
//	  "expires": "2025-04-02T13:25:16Z",
//	  "challenges": [
//	    {
//	      "type": "dns-01",
//	      "url": "https://acme-staging-v02.api.letsencrypt.org/acme/chall/189956024/16553103724/hj-ldw",
//	      "status": "pending",
//	      "token": "Yz3R-LFz6JPr04FN6FjnArcojzyFoD9ojFXZAaG5Rmo"
//	    },
//	    ...
//	  ]
//	}\0
func (h *HAProxyEventListener) handleAcmeDeployEvent(ctx context.Context, args string) {
	if len(args) == 0 {
		log.Error("received HAProxy Event 'acme deploy' without any argument")
		return
	}

	var (
		firstLine = true
		isJSON    = false
		certID    string
		provider  string
		keyAuth   string
		acmeArgs  []string
		acmeJSON  string
		parseErr  error
	)

	// Parse the message line by line.
	strings.SplitSeq(args, "\n")(func(line string) bool {
		if firstLine {
			firstLine = false
			words := strings.Split(line, " ")
			if len(words) != 3 || words[1] != "thumbprint" || len(words[2]) < 3 {
				parseErr = fmt.Errorf("invalid acme deploy line: '%s'", line)
				return false
			}
			certID = strings.Trim(words[0], `"`)
			return true
		}
		if isJSON {
			acmeJSON += line
			return true
		}
		if strings.HasPrefix(line, "provider-name ") {
			words := strings.Split(line, " ")
			if len(words) != 2 {
				parseErr = fmt.Errorf("invalid provider-name line: '%s'", line)
				return false
			}
			provider = strings.Trim(words[1], `"`)
			return true
		}
		if strings.HasPrefix(line, "acme-vars ") {
			_, vars, found := strings.Cut(line, " ")
			if !found || len(vars) == 0 {
				parseErr = fmt.Errorf("invalid acme-vars line: '%s'", line)
				return false
			}
			// Do not trim the double-quotes here.
			acmeArgs = append(acmeArgs, vars)
			return true
		}
		if strings.HasPrefix(line, "dns-01-record ") {
			words := strings.Split(line, " ")
			if len(words) != 2 {
				parseErr = fmt.Errorf("invalid dns-01-record line: '%s'", line)
				return false
			}
			keyAuth = strings.Trim(words[1], `"`)
			return true
		}
		if strings.HasPrefix(line, "{") {
			isJSON = true
			acmeJSON += line
			return true
		}
		// Ignore anything else.
		return true
	})

	if parseErr != nil {
		log.Errorf("events: acme deploy: %s", parseErr.Error())
		return
	}

	// Parse the JSON to get the domain name.
	var status acmeStatus
	if err := jsoniter.UnmarshalFromString(acmeJSON, &status); err != nil {
		log.Errorf("events: acme deploy: json.Unmarshal: %s", err.Error())
		return
	}

	domainName := status.Identifier.Value

	// Merge the acme-vars
	vars := make(map[string]any, 8)
	for _, line := range acmeArgs {
		hmap := configuration.ParseAcmeVars(line)
		for k, v := range hmap {
			vars[k] = v
		}
	}

	// Solve the DNS challenge.
	solver, err := acme.NewDNS01Solver(provider, vars)
	if err != nil {
		log.Errorf("events: acme deploy: DNS provider: %s", err.Error())
		return
	}

	// These options will be configurable from haproxy.cfg in a future version.
	// For now use environment variables.
	solver.PropagationDelay = getEnvDuration("DPAPI_ACME_PROPAGDELAY_SEC", 0)
	solver.PropagationTimeout = getEnvDuration("DPAPI_ACME_PROPAGTIMEOUT_SEC", time.Hour)

	var zone string
	if solver.PropagationTimeout != -1 {
		zone = acme.GuessZone(domainName)
	} else {
		zone, err = acme.FindZoneByFQDN(ctx, domainName, acme.RecursiveNameservers(nil))
	}
	if err != nil {
		log.Errorf("events: acme deploy: failed to find root zone for '%s': %s", domainName, err.Error())
		return
	}
	err = solver.Present(ctx, domainName, zone, keyAuth)
	if err != nil {
		log.Errorf("events: acme deploy: DNS solver: %s", err.Error())
		return
	}
	// Wait for DNS propagation and cleanup.
	err = solver.Wait(ctx, domainName, zone, keyAuth)
	// Remove the challenge in 10m if Wait() was successful. This should be
	// more than enough for HAProxy to finish the challenge with the ACME server.
	waitBeforeCleanup := 10 * time.Minute
	if err != nil {
		waitBeforeCleanup = time.Second
	}
	go func() {
		time.Sleep(waitBeforeCleanup)
		if err := solver.CleanUp(ctx, domainName, zone, keyAuth); err != nil {
			log.Errorf("events: acme deploy: cleanup failed for %s: %v", domainName, err)
		}
	}()
	if err != nil {
		log.Errorf("events: acme deploy: DNS propagation check failed for '%s': %v", domainName, err)
		return
	}

	// Send back a response to HAProxy.
	rt, err := h.client.Runtime()
	if err != nil {
		log.Error(err)
		return
	}
	resp, err := rt.ExecuteRaw("acme challenge_ready " + certID + " domain " + domainName)
	if err != nil {
		log.Errorf("events: acme deploy: sending response: %s", err.Error())
		return
	}

	log.Debugf("events: OK: acme deploy %s => %s", domainName, resp)
}

// Parse an environment variable containing a duration in seconds, or return a default value.
func getEnvDuration(name string, def time.Duration) time.Duration {
	str := os.Getenv(name)
	if str == "" {
		return def
	}
	if str == "-1" {
		// special case to disable waiting for propagation
		return -1
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return time.Duration(num) * time.Second
}
