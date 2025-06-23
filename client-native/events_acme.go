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
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/dataplaneapi/log"
)

const EventAcmeNewCert = "newcert"

func (h *HAProxyEventListener) handleAcmeEvent(message string) {
	name, args, ok := strings.Cut(message, " ")
	if !ok {
		log.Warningf("failed to parse ACME Event: '%s'", message)
		return
	}

	if name == EventAcmeNewCert {
		h.handleAcmeNewCertEvent(args)
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
