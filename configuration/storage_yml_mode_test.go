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

package configuration

import (
	"os"
	"path/filepath"
	"testing"
)

// dataplaneapi.yml holds API users (and plaintext passwords for insecure users): owner only.
func TestStorageYML_NotWorldReadable(t *testing.T) {
	p := filepath.Join(t.TempDir(), "dataplaneapi.yml")
	s := &StorageYML{}
	s.Set(&StorageDataplaneAPIConfiguration{})
	if err := s.SaveAs(p); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("dataplaneapi.yml mode=%v", fi.Mode().Perm())
	if fi.Mode().Perm()&0o077 != 0 {
		t.Fatalf("dataplaneapi.yml readable by group/others: %v", fi.Mode().Perm())
	}
}
