// Copyright 2026 HAProxy Technologies
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

package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/dataplaneapi/storagetype"
)

// Service discovery storage holds credentials (e.g. the Consul token): owner only.
func TestFileStorage_SecretsNotWorldReadable(t *testing.T) {
	p := filepath.Join(t.TempDir(), "consul.json")
	fs := &fileStorage[storagetype.ConsulData]{p}
	if err := fs.Store(storagetype.ConsulData{
		Consuls: storagetype.Consuls{&models.Consul{Token: "dummy-consul-token"}},
	}); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("consul storage mode=%v", fi.Mode().Perm())
	if fi.Mode().Perm()&0o077 != 0 {
		t.Fatalf("consul storage with token readable by group/others: %v", fi.Mode().Perm())
	}
}
